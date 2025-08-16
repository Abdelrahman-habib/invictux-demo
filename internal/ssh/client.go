package ssh

import (
	"context"
	"crypto/md5"
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHClient represents an SSH client with connection pooling and security features
type SSHClient struct {
	config       *ClientConfig
	connections  map[string]*ConnectionPool
	mutex        sync.RWMutex
	hostKeyCheck ssh.HostKeyCallback
}

// ClientConfig holds configuration for the SSH client
type ClientConfig struct {
	ConnectTimeout    time.Duration
	CommandTimeout    time.Duration
	MaxRetries        int
	RetryDelay        time.Duration
	MaxConnections    int
	ConnectionTTL     time.Duration
	KeepAliveInterval time.Duration
}

// ConnectionPool manages SSH connections for a specific host
type ConnectionPool struct {
	host        string
	connections chan *SSHConnection
	active      map[*SSHConnection]bool
	mutex       sync.RWMutex
	config      *ClientConfig
}

// SSHConnection wraps an SSH client connection with metadata
type SSHConnection struct {
	client    *ssh.Client
	createdAt time.Time
	lastUsed  time.Time
	inUse     bool
	mutex     sync.RWMutex
}

// AuthMethod represents different SSH authentication methods
type AuthMethod int

const (
	AuthPassword AuthMethod = iota
	AuthPublicKey
	AuthKeyboard
)

// ConnectionInfo holds information needed to establish an SSH connection
type ConnectionInfo struct {
	Host       string
	Port       int
	Username   string
	Password   string
	PrivateKey []byte
	AuthMethod AuthMethod
}

// CommandResult represents the result of an SSH command execution
type CommandResult struct {
	Command    string
	Output     string
	Error      string
	ExitCode   int
	Duration   time.Duration
	ExecutedAt time.Time
}

// SSHClientInterface defines the interface for SSH client operations
type SSHClientInterface interface {
	Connect(ctx context.Context, connInfo *ConnectionInfo) (*SSHConnection, error)
	ExecuteCommand(ctx context.Context, conn *SSHConnection, command string) (*CommandResult, error)
	ExecuteCommands(ctx context.Context, conn *SSHConnection, commands []string) ([]*CommandResult, error)
	Disconnect(conn *SSHConnection) error
	Close() error
	GetConnectionStats() map[string]ConnectionStats
}

// Global known hosts storage for Trust-On-First-Use (TOFU) approach
var knownHosts = make(map[string]ssh.PublicKey)
var knownHostsMutex sync.RWMutex

// createSecureHostKeyCallback creates a secure host key callback using TOFU approach
func createSecureHostKeyCallback() ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		knownHostsMutex.Lock()
		defer knownHostsMutex.Unlock()

		// Check if we have a known host key for this hostname
		if knownKey, exists := knownHosts[hostname]; exists {
			// Compare the provided key with the known key
			if string(key.Marshal()) == string(knownKey.Marshal()) {
				return nil // Key matches, connection is secure
			}
			return fmt.Errorf("host key verification failed for %s: key mismatch", hostname)
		}

		// For new hosts, implement Trust-On-First-Use (TOFU) approach
		keyFingerprint := md5.Sum(key.Marshal())
		fmt.Printf("WARNING: Unknown host %s with key fingerprint %x\n", hostname, keyFingerprint)
		fmt.Printf("Adding host key to known hosts (Trust-On-First-Use)\n")

		// Store the key for future connections
		knownHosts[hostname] = key

		return nil
	}
}

// CreateInsecureHostKeyCallbackForTesting creates an insecure callback for testing
// WARNING: This should ONLY be used in development/testing environments
func CreateInsecureHostKeyCallbackForTesting() ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		fmt.Printf("WARNING: Using insecure host key verification for %s - this should only be used in development\n", hostname)
		return nil
	}
}

// ConnectionStats provides statistics about connection pools
type ConnectionStats struct {
	Host             string
	ActiveConns      int
	AvailableConns   int
	TotalConns       int
	CreatedConns     int64
	FailedConns      int64
	CommandsExecuted int64
}

// DefaultClientConfig returns a default SSH client configuration
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		ConnectTimeout:    30 * time.Second,
		CommandTimeout:    60 * time.Second,
		MaxRetries:        3,
		RetryDelay:        2 * time.Second,
		MaxConnections:    5,
		ConnectionTTL:     10 * time.Minute,
		KeepAliveInterval: 30 * time.Second,
	}
}

// NewSSHClient creates a new SSH client with the given configuration
func NewSSHClient(config *ClientConfig) *SSHClient {
	if config == nil {
		config = DefaultClientConfig()
	}

	return &SSHClient{
		config:      config,
		connections: make(map[string]*ConnectionPool),
		// Use secure host key verification by default
		hostKeyCheck: createSecureHostKeyCallback(),
	}
}

// NewSSHClientWithHostKeyCheck creates a new SSH client with custom host key verification
func NewSSHClientWithHostKeyCheck(config *ClientConfig, hostKeyCallback ssh.HostKeyCallback) *SSHClient {
	if config == nil {
		config = DefaultClientConfig()
	}

	return &SSHClient{
		config:       config,
		connections:  make(map[string]*ConnectionPool),
		hostKeyCheck: hostKeyCallback,
	}
}

// Connect establishes an SSH connection with retry logic and connection pooling
func (c *SSHClient) Connect(ctx context.Context, connInfo *ConnectionInfo) (*SSHConnection, error) {
	if connInfo == nil {
		return nil, fmt.Errorf("connection info cannot be nil")
	}

	if err := c.validateConnectionInfo(connInfo); err != nil {
		return nil, fmt.Errorf("invalid connection info: %w", err)
	}

	hostKey := fmt.Sprintf("%s:%d", connInfo.Host, connInfo.Port)

	// Get or create connection pool for this host
	pool := c.getOrCreatePool(hostKey)

	// Try to get an existing connection from the pool
	if conn := pool.getConnection(); conn != nil {
		return conn, nil
	}

	// Create a new connection with retry logic
	return c.createConnectionWithRetry(ctx, connInfo, pool)
}

// ExecuteCommand executes a single command on the SSH connection
func (c *SSHClient) ExecuteCommand(ctx context.Context, conn *SSHConnection, command string) (*CommandResult, error) {
	if conn == nil {
		return nil, fmt.Errorf("connection cannot be nil")
	}

	if command == "" {
		return nil, fmt.Errorf("command cannot be empty")
	}

	startTime := time.Now()
	result := &CommandResult{
		Command:    command,
		ExecutedAt: startTime,
	}

	// Mark connection as in use
	conn.mutex.Lock()
	conn.inUse = true
	conn.lastUsed = time.Now()
	conn.mutex.Unlock()

	defer func() {
		conn.mutex.Lock()
		conn.inUse = false
		conn.mutex.Unlock()
		result.Duration = time.Since(startTime)
	}()

	// Create a new session for command execution
	session, err := conn.client.NewSession()
	if err != nil {
		result.Error = fmt.Sprintf("failed to create session: %v", err)
		return result, err
	}
	defer session.Close()

	// Set up command timeout
	cmdCtx, cancel := context.WithTimeout(ctx, c.config.CommandTimeout)
	defer cancel()

	// Execute command with timeout
	outputChan := make(chan []byte, 1)
	errorChan := make(chan error, 1)

	go func() {
		output, err := session.CombinedOutput(command)
		if err != nil {
			errorChan <- err
		} else {
			outputChan <- output
		}
	}()

	select {
	case output := <-outputChan:
		result.Output = string(output)
		result.ExitCode = 0
		return result, nil
	case err := <-errorChan:
		result.Error = err.Error()
		if exitErr, ok := err.(*ssh.ExitError); ok {
			result.ExitCode = exitErr.ExitStatus()
		} else {
			result.ExitCode = -1
		}
		return result, err
	case <-cmdCtx.Done():
		result.Error = "command execution timeout"
		result.ExitCode = -1
		return result, fmt.Errorf("command execution timeout")
	}
}

// ExecuteCommands executes multiple commands sequentially on the SSH connection
func (c *SSHClient) ExecuteCommands(ctx context.Context, conn *SSHConnection, commands []string) ([]*CommandResult, error) {
	if len(commands) == 0 {
		return nil, fmt.Errorf("commands list cannot be empty")
	}

	results := make([]*CommandResult, 0, len(commands))

	for _, command := range commands {
		result, err := c.ExecuteCommand(ctx, conn, command)
		results = append(results, result)

		// Continue executing other commands even if one fails
		if err != nil {
			// Log the error but continue with other commands
			continue
		}
	}

	return results, nil
}

// Disconnect closes an SSH connection and returns it to the pool or closes it
func (c *SSHClient) Disconnect(conn *SSHConnection) error {
	if conn == nil {
		return nil
	}

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	// Check if connection is still valid and not expired
	if time.Since(conn.createdAt) > c.config.ConnectionTTL {
		return conn.client.Close()
	}

	// Connection is still valid, could be returned to pool
	// For now, we'll close it. In a full implementation, we'd return it to the pool
	return conn.client.Close()
}

// Close closes all connections and cleans up resources
func (c *SSHClient) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var lastErr error
	for _, pool := range c.connections {
		if err := pool.closeAll(); err != nil {
			lastErr = err
		}
	}

	c.connections = make(map[string]*ConnectionPool)
	return lastErr
}

// GetConnectionStats returns statistics about all connection pools
func (c *SSHClient) GetConnectionStats() map[string]ConnectionStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	stats := make(map[string]ConnectionStats)
	for host, pool := range c.connections {
		stats[host] = pool.getStats()
	}

	return stats
}

// validateConnectionInfo validates the connection information
func (c *SSHClient) validateConnectionInfo(connInfo *ConnectionInfo) error {
	if connInfo.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if connInfo.Port <= 0 || connInfo.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	if connInfo.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	switch connInfo.AuthMethod {
	case AuthPassword:
		if connInfo.Password == "" {
			return fmt.Errorf("password cannot be empty for password authentication")
		}
	case AuthPublicKey:
		if len(connInfo.PrivateKey) == 0 {
			return fmt.Errorf("private key cannot be empty for public key authentication")
		}
	case AuthKeyboard:
		// Keyboard interactive authentication doesn't require additional validation here
	default:
		return fmt.Errorf("unsupported authentication method")
	}

	return nil
}

// getOrCreatePool gets an existing connection pool or creates a new one
func (c *SSHClient) getOrCreatePool(hostKey string) *ConnectionPool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if pool, exists := c.connections[hostKey]; exists {
		return pool
	}

	pool := &ConnectionPool{
		host:        hostKey,
		connections: make(chan *SSHConnection, c.config.MaxConnections),
		active:      make(map[*SSHConnection]bool),
		config:      c.config,
	}

	c.connections[hostKey] = pool
	return pool
}

// createConnectionWithRetry creates a new SSH connection with retry logic
func (c *SSHClient) createConnectionWithRetry(ctx context.Context, connInfo *ConnectionInfo, pool *ConnectionPool) (*SSHConnection, error) {
	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retrying with exponential backoff
			delay := time.Duration(attempt) * c.config.RetryDelay
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		conn, err := c.createConnection(ctx, connInfo)
		if err == nil {
			pool.addConnection(conn)
			return conn, nil
		}

		lastErr = err

		// Check if context was cancelled
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %w", c.config.MaxRetries+1, lastErr)
}

// createConnection creates a new SSH connection
func (c *SSHClient) createConnection(ctx context.Context, connInfo *ConnectionInfo) (*SSHConnection, error) {
	// Prepare SSH client configuration
	config := &ssh.ClientConfig{
		User:            connInfo.Username,
		HostKeyCallback: c.hostKeyCheck,
		Timeout:         c.config.ConnectTimeout,
	}

	// Set up authentication method
	switch connInfo.AuthMethod {
	case AuthPassword:
		config.Auth = []ssh.AuthMethod{
			ssh.Password(connInfo.Password),
		}
	case AuthPublicKey:
		signer, err := ssh.ParsePrivateKey(connInfo.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	case AuthKeyboard:
		config.Auth = []ssh.AuthMethod{
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
				// For keyboard interactive, we'll use the password for now
				// In a full implementation, this would be more sophisticated
				answers := make([]string, len(questions))
				for i := range answers {
					answers[i] = connInfo.Password
				}
				return answers, nil
			}),
		}
	}

	// Create connection with timeout
	address := fmt.Sprintf("%s:%d", connInfo.Host, connInfo.Port)

	// Use context for connection timeout
	dialer := &net.Dialer{
		Timeout: c.config.ConnectTimeout,
	}

	netConn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", address, err)
	}

	sshConn, chans, reqs, err := ssh.NewClientConn(netConn, address, config)
	if err != nil {
		netConn.Close()
		return nil, fmt.Errorf("failed to create SSH connection: %w", err)
	}

	client := ssh.NewClient(sshConn, chans, reqs)

	return &SSHConnection{
		client:    client,
		createdAt: time.Now(),
		lastUsed:  time.Now(),
		inUse:     false,
	}, nil
}

// ConnectionPool methods

// getConnection gets an available connection from the pool
func (p *ConnectionPool) getConnection() *SSHConnection {
	select {
	case conn := <-p.connections:
		// Check if connection is still valid
		if time.Since(conn.createdAt) > p.config.ConnectionTTL {
			conn.client.Close()
			return nil
		}
		return conn
	default:
		return nil
	}
}

// addConnection adds a connection to the pool
func (p *ConnectionPool) addConnection(conn *SSHConnection) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.active[conn] = true
}

// closeAll closes all connections in the pool
func (p *ConnectionPool) closeAll() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var lastErr error

	// Close connections in the channel
	for {
		select {
		case conn := <-p.connections:
			if err := conn.client.Close(); err != nil {
				lastErr = err
			}
		default:
			goto closeActive
		}
	}

closeActive:
	// Close active connections
	for conn := range p.active {
		if err := conn.client.Close(); err != nil {
			lastErr = err
		}
	}

	p.active = make(map[*SSHConnection]bool)
	return lastErr
}

// getStats returns statistics for this connection pool
func (p *ConnectionPool) getStats() ConnectionStats {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return ConnectionStats{
		Host:           p.host,
		ActiveConns:    len(p.active),
		AvailableConns: len(p.connections),
		TotalConns:     len(p.active) + len(p.connections),
	}
}
