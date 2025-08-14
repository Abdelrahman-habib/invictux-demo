package ssh

import (
	"context"
	"fmt"
	"time"
)

// DeviceSSHManager provides SSH operations for network devices
type DeviceSSHManager struct {
	client *SSHClient
}

// DeviceConnection represents connection information for a network device
type DeviceConnection struct {
	ID       string
	Name     string
	Host     string
	Port     int
	Username string
	Password string
}

// DeviceSSHManagerInterface defines the interface for device SSH operations
type DeviceSSHManagerInterface interface {
	ConnectToDevice(ctx context.Context, device *DeviceConnection) (*SSHConnection, error)
	ExecuteDeviceCommand(ctx context.Context, conn *SSHConnection, command string) (*CommandResult, error)
	ExecuteDeviceCommands(ctx context.Context, conn *SSHConnection, commands []string) ([]*CommandResult, error)
	TestDeviceConnectivity(ctx context.Context, device *DeviceConnection) error
	DisconnectFromDevice(conn *SSHConnection) error
	Close() error
}

// NewDeviceSSHManager creates a new device SSH manager
func NewDeviceSSHManager(config *ClientConfig) *DeviceSSHManager {
	return &DeviceSSHManager{
		client: NewSSHClient(config),
	}
}

// NewDeviceSSHManagerWithDefaults creates a new device SSH manager with default configuration
func NewDeviceSSHManagerWithDefaults() *DeviceSSHManager {
	return &DeviceSSHManager{
		client: NewSSHClient(DefaultClientConfig()),
	}
}

// ConnectToDevice establishes an SSH connection to a network device
func (m *DeviceSSHManager) ConnectToDevice(ctx context.Context, device *DeviceConnection) (*SSHConnection, error) {
	if device == nil {
		return nil, fmt.Errorf("device connection info cannot be nil")
	}

	connInfo := &ConnectionInfo{
		Host:       device.Host,
		Port:       device.Port,
		Username:   device.Username,
		Password:   device.Password,
		AuthMethod: AuthPassword, // Default to password authentication
	}

	return m.client.Connect(ctx, connInfo)
}

// ExecuteDeviceCommand executes a single command on a network device
func (m *DeviceSSHManager) ExecuteDeviceCommand(ctx context.Context, conn *SSHConnection, command string) (*CommandResult, error) {
	return m.client.ExecuteCommand(ctx, conn, command)
}

// ExecuteDeviceCommands executes multiple commands on a network device
func (m *DeviceSSHManager) ExecuteDeviceCommands(ctx context.Context, conn *SSHConnection, commands []string) ([]*CommandResult, error) {
	return m.client.ExecuteCommands(ctx, conn, commands)
}

// TestDeviceConnectivity tests SSH connectivity to a network device
func (m *DeviceSSHManager) TestDeviceConnectivity(ctx context.Context, device *DeviceConnection) error {
	conn, err := m.ConnectToDevice(ctx, device)
	if err != nil {
		return fmt.Errorf("failed to connect to device %s (%s): %w", device.Name, device.Host, err)
	}
	defer m.DisconnectFromDevice(conn)

	// Execute a simple command to verify the connection works
	result, err := m.ExecuteDeviceCommand(ctx, conn, "show version")
	if err != nil {
		return fmt.Errorf("failed to execute test command on device %s: %w", device.Name, err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("test command failed on device %s with exit code %d: %s", device.Name, result.ExitCode, result.Error)
	}

	return nil
}

// DisconnectFromDevice closes the SSH connection to a network device
func (m *DeviceSSHManager) DisconnectFromDevice(conn *SSHConnection) error {
	return m.client.Disconnect(conn)
}

// Close closes all SSH connections and cleans up resources
func (m *DeviceSSHManager) Close() error {
	return m.client.Close()
}

// GetConnectionStats returns connection statistics
func (m *DeviceSSHManager) GetConnectionStats() map[string]ConnectionStats {
	return m.client.GetConnectionStats()
}

// ExecuteCommandWithTimeout executes a command with a specific timeout
func (m *DeviceSSHManager) ExecuteCommandWithTimeout(ctx context.Context, conn *SSHConnection, command string, timeout time.Duration) (*CommandResult, error) {
	cmdCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return m.client.ExecuteCommand(cmdCtx, conn, command)
}

// BatchExecuteOnDevices executes commands on multiple devices concurrently
func (m *DeviceSSHManager) BatchExecuteOnDevices(ctx context.Context, devices []*DeviceConnection, commands []string) (map[string][]*CommandResult, error) {
	if len(devices) == 0 {
		return nil, fmt.Errorf("devices list cannot be empty")
	}

	if len(commands) == 0 {
		return nil, fmt.Errorf("commands list cannot be empty")
	}

	results := make(map[string][]*CommandResult)
	resultChan := make(chan struct {
		deviceID string
		results  []*CommandResult
		err      error
	}, len(devices))

	// Execute commands on each device concurrently
	for _, device := range devices {
		go func(dev *DeviceConnection) {
			deviceResults, err := m.executeCommandsOnDevice(ctx, dev, commands)
			resultChan <- struct {
				deviceID string
				results  []*CommandResult
				err      error
			}{dev.ID, deviceResults, err}
		}(device)
	}

	// Collect results
	for i := 0; i < len(devices); i++ {
		select {
		case result := <-resultChan:
			if result.err != nil {
				// Log error but continue with other devices
				// In a production system, you might want to handle this differently
				results[result.deviceID] = []*CommandResult{
					{
						Command:    "connection_error",
						Error:      result.err.Error(),
						ExitCode:   -1,
						ExecutedAt: time.Now(),
					},
				}
			} else {
				results[result.deviceID] = result.results
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("batch execution cancelled: %w", ctx.Err())
		}
	}

	return results, nil
}

// executeCommandsOnDevice executes commands on a single device
func (m *DeviceSSHManager) executeCommandsOnDevice(ctx context.Context, device *DeviceConnection, commands []string) ([]*CommandResult, error) {
	conn, err := m.ConnectToDevice(ctx, device)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to device %s: %w", device.Name, err)
	}
	defer m.DisconnectFromDevice(conn)

	return m.ExecuteDeviceCommands(ctx, conn, commands)
}

// ValidateDeviceConnection validates device connection parameters
func ValidateDeviceConnection(device *DeviceConnection) error {
	if device == nil {
		return fmt.Errorf("device connection cannot be nil")
	}

	if device.Host == "" {
		return fmt.Errorf("device host cannot be empty")
	}

	if device.Port <= 0 || device.Port > 65535 {
		return fmt.Errorf("device port must be between 1 and 65535")
	}

	if device.Username == "" {
		return fmt.Errorf("device username cannot be empty")
	}

	if device.Password == "" {
		return fmt.Errorf("device password cannot be empty")
	}

	return nil
}

// CreateDeviceConnectionFromDevice creates a DeviceConnection from device information
func CreateDeviceConnectionFromDevice(id, name, host string, port int, username, password string) *DeviceConnection {
	return &DeviceConnection{
		ID:       id,
		Name:     name,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}
