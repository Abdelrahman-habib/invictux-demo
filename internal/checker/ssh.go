package checker

import (
	"fmt"
	"time"

	"qwin/internal/device"

	"golang.org/x/crypto/ssh"
)

// SSHClient handles SSH connections to network devices
type SSHClient struct {
	timeout time.Duration
	retries int
}

// NewSSHClient creates a new SSH client
func NewSSHClient() *SSHClient {
	return &SSHClient{
		timeout: 30 * time.Second,
		retries: 3,
	}
}

// Connect establishes an SSH connection to a device
func (c *SSHClient) Connect(device *device.Device) (*ssh.Session, error) {
	// TODO: Implement proper credential decryption
	config := &ssh.ClientConfig{
		User: device.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password("placeholder"), // Will decrypt device.PasswordEncrypted
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For demo - use proper verification in production
		Timeout:         c.timeout,
	}

	address := fmt.Sprintf("%s:%d", device.IPAddress, device.SSHPort)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// ExecuteCommand executes a command on the SSH session
func (c *SSHClient) ExecuteCommand(session *ssh.Session, command string) (string, error) {
	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w", err)
	}

	return string(output), nil
}

// Disconnect closes the SSH session
func (c *SSHClient) Disconnect(session *ssh.Session) {
	if session != nil {
		session.Close()
	}
}

// SetTimeout sets the connection timeout
func (c *SSHClient) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// SetRetries sets the number of connection retries
func (c *SSHClient) SetRetries(retries int) {
	c.retries = retries
}
