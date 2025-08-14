package ssh

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

// MockSSHServer represents a mock SSH server for testing
type MockSSHServer struct {
	listener   net.Listener
	config     *ssh.ServerConfig
	address    string
	port       int
	commands   map[string]string // command -> response mapping
	shouldFail bool
	delay      time.Duration
}

// NewMockSSHServer creates a new mock SSH server
func NewMockSSHServer() (*MockSSHServer, error) {
	// Generate a test host key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == "testuser" && string(pass) == "testpass" {
				return nil, nil
			}
			return nil, fmt.Errorf("invalid credentials")
		},
	}
	config.AddHostKey(signer)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	address := listener.Addr().String()
	host, portStr, _ := net.SplitHostPort(address)
	port := 0
	fmt.Sscanf(portStr, "%d", &port)

	server := &MockSSHServer{
		listener: listener,
		config:   config,
		address:  host,
		port:     port,
		commands: make(map[string]string),
	}

	go server.serve()
	return server, nil
}

// SetCommandResponse sets the response for a specific command
func (s *MockSSHServer) SetCommandResponse(command, response string) {
	s.commands[command] = response
}

// SetShouldFail sets whether the server should fail connections
func (s *MockSSHServer) SetShouldFail(shouldFail bool) {
	s.shouldFail = shouldFail
}

// SetDelay sets a delay for command execution
func (s *MockSSHServer) SetDelay(delay time.Duration) {
	s.delay = delay
}

// GetAddress returns the server address
func (s *MockSSHServer) GetAddress() string {
	return s.address
}

// GetPort returns the server port
func (s *MockSSHServer) GetPort() int {
	return s.port
}

// Close closes the mock server
func (s *MockSSHServer) Close() error {
	return s.listener.Close()
}

// serve handles incoming connections
func (s *MockSSHServer) serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}

		go s.handleConnection(conn)
	}
}

// handleConnection handles a single SSH connection
func (s *MockSSHServer) handleConnection(netConn net.Conn) {
	defer netConn.Close()

	if s.shouldFail {
		return
	}

	sshConn, chans, reqs, err := ssh.NewServerConn(netConn, s.config)
	if err != nil {
		return
	}
	defer sshConn.Close()

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			continue
		}

		go s.handleSession(channel, requests)
	}
}

// handleSession handles a single SSH session
func (s *MockSSHServer) handleSession(channel ssh.Channel, requests <-chan *ssh.Request) {
	defer channel.Close()

	for req := range requests {
		switch req.Type {
		case "exec":
			if s.delay > 0 {
				time.Sleep(s.delay)
			}

			command := string(req.Payload[4:]) // Skip the length prefix
			response, exists := s.commands[command]
			if !exists {
				response = fmt.Sprintf("Command not found: %s", command)
			}

			channel.Write([]byte(response))
			channel.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
			req.Reply(true, nil)
			return
		default:
			req.Reply(false, nil)
		}
	}
}

// Test helper functions

func generateTestPrivateKey() ([]byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	return pem.EncodeToMemory(privateKeyPEM), nil
}

// Unit Tests

func TestNewSSHClient(t *testing.T) {
	client := NewSSHClient(nil)

	if client == nil {
		t.Fatal("NewSSHClient returned nil")
	}

	if client.config == nil {
		t.Error("Client config should not be nil")
	}

	if client.connections == nil {
		t.Error("Client connections map should not be nil")
	}
}

func TestNewSSHClientWithConfig(t *testing.T) {
	config := &ClientConfig{
		ConnectTimeout: 5 * time.Second,
		CommandTimeout: 10 * time.Second,
		MaxRetries:     2,
	}

	client := NewSSHClient(config)

	if client.config.ConnectTimeout != 5*time.Second {
		t.Errorf("Expected ConnectTimeout 5s, got %v", client.config.ConnectTimeout)
	}

	if client.config.CommandTimeout != 10*time.Second {
		t.Errorf("Expected CommandTimeout 10s, got %v", client.config.CommandTimeout)
	}

	if client.config.MaxRetries != 2 {
		t.Errorf("Expected MaxRetries 2, got %d", client.config.MaxRetries)
	}
}

func TestDefaultClientConfig(t *testing.T) {
	config := DefaultClientConfig()

	if config.ConnectTimeout != 30*time.Second {
		t.Errorf("Expected default ConnectTimeout 30s, got %v", config.ConnectTimeout)
	}

	if config.CommandTimeout != 60*time.Second {
		t.Errorf("Expected default CommandTimeout 60s, got %v", config.CommandTimeout)
	}

	if config.MaxRetries != 3 {
		t.Errorf("Expected default MaxRetries 3, got %d", config.MaxRetries)
	}
}

func TestSSHClient_Connect_Success(t *testing.T) {
	server, err := NewMockSSHServer()
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer server.Close()

	client := NewSSHClient(nil)
	defer client.Close()

	connInfo := &ConnectionInfo{
		Host:       server.GetAddress(),
		Port:       server.GetPort(),
		Username:   "testuser",
		Password:   "testpass",
		AuthMethod: AuthPassword,
	}

	ctx := context.Background()
	conn, err := client.Connect(ctx, connInfo)

	if err != nil {
		t.Errorf("Expected successful connection, got error: %v", err)
	}

	if conn == nil {
		t.Error("Expected connection, got nil")
	}

	if conn != nil {
		client.Disconnect(conn)
	}
}

func TestSSHClient_Connect_InvalidCredentials(t *testing.T) {
	server, err := NewMockSSHServer()
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer server.Close()

	client := NewSSHClient(nil)
	defer client.Close()

	connInfo := &ConnectionInfo{
		Host:       server.GetAddress(),
		Port:       server.GetPort(),
		Username:   "wronguser",
		Password:   "wrongpass",
		AuthMethod: AuthPassword,
	}

	ctx := context.Background()
	conn, err := client.Connect(ctx, connInfo)

	if err == nil {
		t.Error("Expected connection to fail with invalid credentials")
	}

	if conn != nil {
		t.Error("Expected nil connection for failed authentication")
		client.Disconnect(conn)
	}
}

func TestSSHClient_Connect_NilConnectionInfo(t *testing.T) {
	client := NewSSHClient(nil)
	defer client.Close()

	ctx := context.Background()
	conn, err := client.Connect(ctx, nil)

	if err == nil {
		t.Error("Expected error for nil connection info")
	}

	if conn != nil {
		t.Error("Expected nil connection for nil connection info")
	}

	expectedError := "connection info cannot be nil"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSSHClient_Connect_InvalidConnectionInfo(t *testing.T) {
	client := NewSSHClient(nil)
	defer client.Close()

	testCases := []struct {
		name     string
		connInfo *ConnectionInfo
		expected string
	}{
		{
			name: "empty host",
			connInfo: &ConnectionInfo{
				Host:       "",
				Port:       22,
				Username:   "user",
				Password:   "pass",
				AuthMethod: AuthPassword,
			},
			expected: "host cannot be empty",
		},
		{
			name: "invalid port",
			connInfo: &ConnectionInfo{
				Host:       "localhost",
				Port:       0,
				Username:   "user",
				Password:   "pass",
				AuthMethod: AuthPassword,
			},
			expected: "port must be between 1 and 65535",
		},
		{
			name: "empty username",
			connInfo: &ConnectionInfo{
				Host:       "localhost",
				Port:       22,
				Username:   "",
				Password:   "pass",
				AuthMethod: AuthPassword,
			},
			expected: "username cannot be empty",
		},
		{
			name: "empty password for password auth",
			connInfo: &ConnectionInfo{
				Host:       "localhost",
				Port:       22,
				Username:   "user",
				Password:   "",
				AuthMethod: AuthPassword,
			},
			expected: "password cannot be empty for password authentication",
		},
	}

	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conn, err := client.Connect(ctx, tc.connInfo)

			if err == nil {
				t.Error("Expected error for invalid connection info")
			}

			if conn != nil {
				t.Error("Expected nil connection for invalid connection info")
			}

			if !strings.Contains(err.Error(), tc.expected) {
				t.Errorf("Expected error containing '%s', got '%s'", tc.expected, err.Error())
			}
		})
	}
}

func TestSSHClient_ExecuteCommand_Success(t *testing.T) {
	server, err := NewMockSSHServer()
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer server.Close()

	// Set up command responses
	server.SetCommandResponse("show version", "Cisco IOS Version 15.1")
	server.SetCommandResponse("show running-config", "Current configuration")

	client := NewSSHClient(nil)
	defer client.Close()

	connInfo := &ConnectionInfo{
		Host:       server.GetAddress(),
		Port:       server.GetPort(),
		Username:   "testuser",
		Password:   "testpass",
		AuthMethod: AuthPassword,
	}

	ctx := context.Background()
	conn, err := client.Connect(ctx, connInfo)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect(conn)

	result, err := client.ExecuteCommand(ctx, conn, "show version")

	if err != nil {
		t.Errorf("Expected successful command execution, got error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected command result, got nil")
	}

	if result.Command != "show version" {
		t.Errorf("Expected command 'show version', got '%s'", result.Command)
	}

	if result.Output != "Cisco IOS Version 15.1" {
		t.Errorf("Expected output 'Cisco IOS Version 15.1', got '%s'", result.Output)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	if result.Duration < 0 {
		t.Error("Expected non-negative duration")
	}

	if result.ExecutedAt.IsZero() {
		t.Error("Expected ExecutedAt to be set")
	}
}

func TestSSHClient_ExecuteCommand_NilConnection(t *testing.T) {
	client := NewSSHClient(nil)
	defer client.Close()

	ctx := context.Background()
	result, err := client.ExecuteCommand(ctx, nil, "show version")

	if err == nil {
		t.Error("Expected error for nil connection")
	}

	if result != nil {
		t.Error("Expected nil result for nil connection")
	}

	expectedError := "connection cannot be nil"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSSHClient_ExecuteCommand_EmptyCommand(t *testing.T) {
	server, err := NewMockSSHServer()
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer server.Close()

	client := NewSSHClient(nil)
	defer client.Close()

	connInfo := &ConnectionInfo{
		Host:       server.GetAddress(),
		Port:       server.GetPort(),
		Username:   "testuser",
		Password:   "testpass",
		AuthMethod: AuthPassword,
	}

	ctx := context.Background()
	conn, err := client.Connect(ctx, connInfo)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect(conn)

	result, err := client.ExecuteCommand(ctx, conn, "")

	if err == nil {
		t.Error("Expected error for empty command")
	}

	if result != nil {
		t.Error("Expected nil result for empty command")
	}

	expectedError := "command cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSSHClient_ExecuteCommands_Success(t *testing.T) {
	server, err := NewMockSSHServer()
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer server.Close()

	// Set up command responses
	server.SetCommandResponse("show version", "Cisco IOS Version 15.1")
	server.SetCommandResponse("show interfaces", "Interface status")
	server.SetCommandResponse("show running-config", "Current configuration")

	client := NewSSHClient(nil)
	defer client.Close()

	connInfo := &ConnectionInfo{
		Host:       server.GetAddress(),
		Port:       server.GetPort(),
		Username:   "testuser",
		Password:   "testpass",
		AuthMethod: AuthPassword,
	}

	ctx := context.Background()
	conn, err := client.Connect(ctx, connInfo)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect(conn)

	commands := []string{"show version", "show interfaces", "show running-config"}
	results, err := client.ExecuteCommands(ctx, conn, commands)

	if err != nil {
		t.Errorf("Expected successful commands execution, got error: %v", err)
	}

	if len(results) != len(commands) {
		t.Errorf("Expected %d results, got %d", len(commands), len(results))
	}

	expectedOutputs := []string{
		"Cisco IOS Version 15.1",
		"Interface status",
		"Current configuration",
	}

	for i, result := range results {
		if result.Command != commands[i] {
			t.Errorf("Result %d: expected command '%s', got '%s'", i, commands[i], result.Command)
		}

		if result.Output != expectedOutputs[i] {
			t.Errorf("Result %d: expected output '%s', got '%s'", i, expectedOutputs[i], result.Output)
		}
	}
}

func TestSSHClient_ExecuteCommands_EmptyList(t *testing.T) {
	client := NewSSHClient(nil)
	defer client.Close()

	ctx := context.Background()
	results, err := client.ExecuteCommands(ctx, nil, []string{})

	if err == nil {
		t.Error("Expected error for empty commands list")
	}

	if results != nil {
		t.Error("Expected nil results for empty commands list")
	}

	expectedError := "commands list cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSSHClient_Connect_WithRetry(t *testing.T) {
	// Test that connection fails after exhausting retries
	config := &ClientConfig{
		ConnectTimeout: 1 * time.Second,
		CommandTimeout: 5 * time.Second,
		MaxRetries:     2,
		RetryDelay:     100 * time.Millisecond,
		MaxConnections: 5,
		ConnectionTTL:  5 * time.Minute,
	}

	client := NewSSHClient(config)
	defer client.Close()

	// Use a non-existent host to force connection failure
	connInfo := &ConnectionInfo{
		Host:       "192.0.2.1", // RFC5737 test address - should be unreachable
		Port:       22,
		Username:   "testuser",
		Password:   "testpass",
		AuthMethod: AuthPassword,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	startTime := time.Now()
	conn, err := client.Connect(ctx, connInfo)
	duration := time.Since(startTime)

	// Connection should fail
	if err == nil {
		t.Error("Expected connection to fail for unreachable host")
		if conn != nil {
			client.Disconnect(conn)
		}
	}

	// Should have taken some time due to retries
	expectedMinDuration := time.Duration(config.MaxRetries) * config.RetryDelay
	if duration < expectedMinDuration {
		t.Errorf("Expected connection attempt to take at least %v due to retries, took %v", expectedMinDuration, duration)
	}

	t.Logf("Connection failed as expected after %v: %v", duration, err)
}

func TestSSHClient_CommandTimeout(t *testing.T) {
	server, err := NewMockSSHServer()
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer server.Close()

	// Set server to delay command execution
	server.SetDelay(2 * time.Second)
	server.SetCommandResponse("slow command", "This is slow")

	config := &ClientConfig{
		ConnectTimeout: 5 * time.Second,
		CommandTimeout: 500 * time.Millisecond, // Short timeout
		MaxRetries:     1,
		RetryDelay:     100 * time.Millisecond,
	}

	client := NewSSHClient(config)
	defer client.Close()

	connInfo := &ConnectionInfo{
		Host:       server.GetAddress(),
		Port:       server.GetPort(),
		Username:   "testuser",
		Password:   "testpass",
		AuthMethod: AuthPassword,
	}

	ctx := context.Background()
	conn, err := client.Connect(ctx, connInfo)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect(conn)

	result, err := client.ExecuteCommand(ctx, conn, "slow command")

	if err == nil {
		t.Error("Expected timeout error")
	}

	if result == nil {
		t.Fatal("Expected result even on timeout")
	}

	if result.ExitCode != -1 {
		t.Errorf("Expected exit code -1 for timeout, got %d", result.ExitCode)
	}

	if !strings.Contains(result.Error, "timeout") {
		t.Errorf("Expected timeout error, got '%s'", result.Error)
	}
}

func TestSSHClient_GetConnectionStats(t *testing.T) {
	client := NewSSHClient(nil)
	defer client.Close()

	stats := client.GetConnectionStats()

	if stats == nil {
		t.Error("Expected connection stats, got nil")
	}

	if len(stats) != 0 {
		t.Errorf("Expected empty stats for new client, got %d entries", len(stats))
	}
}

func TestSSHClient_Close(t *testing.T) {
	client := NewSSHClient(nil)

	err := client.Close()

	if err != nil {
		t.Errorf("Expected no error on close, got: %v", err)
	}

	// Verify connections map is reset
	if len(client.connections) != 0 {
		t.Error("Expected connections map to be empty after close")
	}
}

// Benchmark tests

func BenchmarkSSHClient_Connect(b *testing.B) {
	server, err := NewMockSSHServer()
	if err != nil {
		b.Fatalf("Failed to create mock server: %v", err)
	}
	defer server.Close()

	client := NewSSHClient(nil)
	defer client.Close()

	connInfo := &ConnectionInfo{
		Host:       server.GetAddress(),
		Port:       server.GetPort(),
		Username:   "testuser",
		Password:   "testpass",
		AuthMethod: AuthPassword,
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn, err := client.Connect(ctx, connInfo)
		if err != nil {
			b.Errorf("Connection failed: %v", err)
		}
		if conn != nil {
			client.Disconnect(conn)
		}
	}
}

func BenchmarkSSHClient_ExecuteCommand(b *testing.B) {
	server, err := NewMockSSHServer()
	if err != nil {
		b.Fatalf("Failed to create mock server: %v", err)
	}
	defer server.Close()

	server.SetCommandResponse("test command", "test response")

	client := NewSSHClient(nil)
	defer client.Close()

	connInfo := &ConnectionInfo{
		Host:       server.GetAddress(),
		Port:       server.GetPort(),
		Username:   "testuser",
		Password:   "testpass",
		AuthMethod: AuthPassword,
	}

	ctx := context.Background()
	conn, err := client.Connect(ctx, connInfo)
	if err != nil {
		b.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect(conn)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.ExecuteCommand(ctx, conn, "test command")
		if err != nil {
			b.Errorf("Command execution failed: %v", err)
		}
	}
}
