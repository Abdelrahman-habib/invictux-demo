package ssh

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestNewDeviceSSHManager(t *testing.T) {
	config := DefaultClientConfig()
	manager := NewDeviceSSHManager(config)

	if manager == nil {
		t.Fatal("NewDeviceSSHManager returned nil")
	}

	if manager.client == nil {
		t.Error("Manager client should not be nil")
	}
}

func TestNewDeviceSSHManagerWithDefaults(t *testing.T) {
	manager := NewDeviceSSHManagerWithDefaults()

	if manager == nil {
		t.Fatal("NewDeviceSSHManagerWithDefaults returned nil")
	}

	if manager.client == nil {
		t.Error("Manager client should not be nil")
	}
}

func TestDeviceSSHManager_ConnectToDevice_Success(t *testing.T) {
	server, err := NewMockSSHServer()
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer server.Close()

	manager := NewDeviceSSHManagerWithDefaults()
	defer manager.Close()

	device := &DeviceConnection{
		ID:       "test-device-1",
		Name:     "Test Router",
		Host:     server.GetAddress(),
		Port:     server.GetPort(),
		Username: "testuser",
		Password: "testpass",
	}

	ctx := context.Background()
	conn, err := manager.ConnectToDevice(ctx, device)

	if err != nil {
		t.Errorf("Expected successful connection, got error: %v", err)
	}

	if conn == nil {
		t.Error("Expected connection, got nil")
	}

	if conn != nil {
		manager.DisconnectFromDevice(conn)
	}
}

func TestDeviceSSHManager_ConnectToDevice_NilDevice(t *testing.T) {
	manager := NewDeviceSSHManagerWithDefaults()
	defer manager.Close()

	ctx := context.Background()
	conn, err := manager.ConnectToDevice(ctx, nil)

	if err == nil {
		t.Error("Expected error for nil device")
	}

	if conn != nil {
		t.Error("Expected nil connection for nil device")
		manager.DisconnectFromDevice(conn)
	}

	expectedError := "device connection info cannot be nil"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDeviceSSHManager_ExecuteDeviceCommand_Success(t *testing.T) {
	server, err := NewMockSSHServer()
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer server.Close()

	server.SetCommandResponse("show version", "Cisco IOS Version 15.1")

	manager := NewDeviceSSHManagerWithDefaults()
	defer manager.Close()

	device := &DeviceConnection{
		ID:       "test-device-1",
		Name:     "Test Router",
		Host:     server.GetAddress(),
		Port:     server.GetPort(),
		Username: "testuser",
		Password: "testpass",
	}

	ctx := context.Background()
	conn, err := manager.ConnectToDevice(ctx, device)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer manager.DisconnectFromDevice(conn)

	result, err := manager.ExecuteDeviceCommand(ctx, conn, "show version")

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
}

func TestDeviceSSHManager_ExecuteDeviceCommands_Success(t *testing.T) {
	server, err := NewMockSSHServer()
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer server.Close()

	server.SetCommandResponse("show version", "Cisco IOS Version 15.1")
	server.SetCommandResponse("show interfaces", "Interface status")

	manager := NewDeviceSSHManagerWithDefaults()
	defer manager.Close()

	device := &DeviceConnection{
		ID:       "test-device-1",
		Name:     "Test Router",
		Host:     server.GetAddress(),
		Port:     server.GetPort(),
		Username: "testuser",
		Password: "testpass",
	}

	ctx := context.Background()
	conn, err := manager.ConnectToDevice(ctx, device)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer manager.DisconnectFromDevice(conn)

	commands := []string{"show version", "show interfaces"}
	results, err := manager.ExecuteDeviceCommands(ctx, conn, commands)

	if err != nil {
		t.Errorf("Expected successful commands execution, got error: %v", err)
	}

	if len(results) != len(commands) {
		t.Errorf("Expected %d results, got %d", len(commands), len(results))
	}

	expectedOutputs := []string{"Cisco IOS Version 15.1", "Interface status"}
	for i, result := range results {
		if result.Output != expectedOutputs[i] {
			t.Errorf("Result %d: expected output '%s', got '%s'", i, expectedOutputs[i], result.Output)
		}
	}
}

func TestDeviceSSHManager_TestDeviceConnectivity_Success(t *testing.T) {
	server, err := NewMockSSHServer()
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer server.Close()

	server.SetCommandResponse("show version", "Cisco IOS Version 15.1")

	manager := NewDeviceSSHManagerWithDefaults()
	defer manager.Close()

	device := &DeviceConnection{
		ID:       "test-device-1",
		Name:     "Test Router",
		Host:     server.GetAddress(),
		Port:     server.GetPort(),
		Username: "testuser",
		Password: "testpass",
	}

	ctx := context.Background()
	err = manager.TestDeviceConnectivity(ctx, device)

	if err != nil {
		t.Errorf("Expected successful connectivity test, got error: %v", err)
	}
}

func TestDeviceSSHManager_TestDeviceConnectivity_ConnectionFailure(t *testing.T) {
	manager := NewDeviceSSHManagerWithDefaults()
	defer manager.Close()

	device := &DeviceConnection{
		ID:       "test-device-1",
		Name:     "Test Router",
		Host:     "192.0.2.1", // RFC5737 test address - should be unreachable
		Port:     22,
		Username: "testuser",
		Password: "testpass",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := manager.TestDeviceConnectivity(ctx, device)

	if err == nil {
		t.Error("Expected connectivity test to fail for unreachable host")
	}

	if !strings.Contains(err.Error(), "failed to connect to device") {
		t.Errorf("Expected connection failure error, got: %v", err)
	}
}

func TestDeviceSSHManager_BatchExecuteOnDevices_Success(t *testing.T) {
	server1, err := NewMockSSHServer()
	if err != nil {
		t.Fatalf("Failed to create mock server 1: %v", err)
	}
	defer server1.Close()

	server2, err := NewMockSSHServer()
	if err != nil {
		t.Fatalf("Failed to create mock server 2: %v", err)
	}
	defer server2.Close()

	server1.SetCommandResponse("show version", "Device 1 Version")
	server1.SetCommandResponse("show interfaces", "Device 1 Interfaces")
	server2.SetCommandResponse("show version", "Device 2 Version")
	server2.SetCommandResponse("show interfaces", "Device 2 Interfaces")

	manager := NewDeviceSSHManagerWithDefaults()
	defer manager.Close()

	devices := []*DeviceConnection{
		{
			ID:       "device-1",
			Name:     "Router 1",
			Host:     server1.GetAddress(),
			Port:     server1.GetPort(),
			Username: "testuser",
			Password: "testpass",
		},
		{
			ID:       "device-2",
			Name:     "Router 2",
			Host:     server2.GetAddress(),
			Port:     server2.GetPort(),
			Username: "testuser",
			Password: "testpass",
		},
	}

	commands := []string{"show version", "show interfaces"}

	ctx := context.Background()
	results, err := manager.BatchExecuteOnDevices(ctx, devices, commands)

	if err != nil {
		t.Errorf("Expected successful batch execution, got error: %v", err)
	}

	if len(results) != len(devices) {
		t.Errorf("Expected results for %d devices, got %d", len(devices), len(results))
	}

	// Check results for device 1
	device1Results, exists := results["device-1"]
	if !exists {
		t.Error("Expected results for device-1")
	} else {
		if len(device1Results) != len(commands) {
			t.Errorf("Expected %d results for device-1, got %d", len(commands), len(device1Results))
		}
		if device1Results[0].Output != "Device 1 Version" {
			t.Errorf("Expected 'Device 1 Version', got '%s'", device1Results[0].Output)
		}
	}

	// Check results for device 2
	device2Results, exists := results["device-2"]
	if !exists {
		t.Error("Expected results for device-2")
	} else {
		if len(device2Results) != len(commands) {
			t.Errorf("Expected %d results for device-2, got %d", len(commands), len(device2Results))
		}
		if device2Results[0].Output != "Device 2 Version" {
			t.Errorf("Expected 'Device 2 Version', got '%s'", device2Results[0].Output)
		}
	}
}

func TestDeviceSSHManager_BatchExecuteOnDevices_EmptyDevices(t *testing.T) {
	manager := NewDeviceSSHManagerWithDefaults()
	defer manager.Close()

	ctx := context.Background()
	results, err := manager.BatchExecuteOnDevices(ctx, []*DeviceConnection{}, []string{"show version"})

	if err == nil {
		t.Error("Expected error for empty devices list")
	}

	if results != nil {
		t.Error("Expected nil results for empty devices list")
	}

	expectedError := "devices list cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestDeviceSSHManager_BatchExecuteOnDevices_EmptyCommands(t *testing.T) {
	manager := NewDeviceSSHManagerWithDefaults()
	defer manager.Close()

	device := &DeviceConnection{
		ID:       "device-1",
		Name:     "Router 1",
		Host:     "localhost",
		Port:     22,
		Username: "testuser",
		Password: "testpass",
	}

	ctx := context.Background()
	results, err := manager.BatchExecuteOnDevices(ctx, []*DeviceConnection{device}, []string{})

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

func TestDeviceSSHManager_ExecuteCommandWithTimeout(t *testing.T) {
	server, err := NewMockSSHServer()
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer server.Close()

	server.SetCommandResponse("show version", "Cisco IOS Version 15.1")

	manager := NewDeviceSSHManagerWithDefaults()
	defer manager.Close()

	device := &DeviceConnection{
		ID:       "test-device-1",
		Name:     "Test Router",
		Host:     server.GetAddress(),
		Port:     server.GetPort(),
		Username: "testuser",
		Password: "testpass",
	}

	ctx := context.Background()
	conn, err := manager.ConnectToDevice(ctx, device)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer manager.DisconnectFromDevice(conn)

	result, err := manager.ExecuteCommandWithTimeout(ctx, conn, "show version", 5*time.Second)

	if err != nil {
		t.Errorf("Expected successful command execution, got error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected command result, got nil")
	}

	if result.Output != "Cisco IOS Version 15.1" {
		t.Errorf("Expected output 'Cisco IOS Version 15.1', got '%s'", result.Output)
	}
}

func TestValidateDeviceConnection(t *testing.T) {
	testCases := []struct {
		name     string
		device   *DeviceConnection
		expected string
	}{
		{
			name:     "nil device",
			device:   nil,
			expected: "device connection cannot be nil",
		},
		{
			name: "empty host",
			device: &DeviceConnection{
				ID:       "device-1",
				Name:     "Router 1",
				Host:     "",
				Port:     22,
				Username: "user",
				Password: "pass",
			},
			expected: "device host cannot be empty",
		},
		{
			name: "invalid port",
			device: &DeviceConnection{
				ID:       "device-1",
				Name:     "Router 1",
				Host:     "localhost",
				Port:     0,
				Username: "user",
				Password: "pass",
			},
			expected: "device port must be between 1 and 65535",
		},
		{
			name: "empty username",
			device: &DeviceConnection{
				ID:       "device-1",
				Name:     "Router 1",
				Host:     "localhost",
				Port:     22,
				Username: "",
				Password: "pass",
			},
			expected: "device username cannot be empty",
		},
		{
			name: "empty password",
			device: &DeviceConnection{
				ID:       "device-1",
				Name:     "Router 1",
				Host:     "localhost",
				Port:     22,
				Username: "user",
				Password: "",
			},
			expected: "device password cannot be empty",
		},
		{
			name: "valid device",
			device: &DeviceConnection{
				ID:       "device-1",
				Name:     "Router 1",
				Host:     "localhost",
				Port:     22,
				Username: "user",
				Password: "pass",
			},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateDeviceConnection(tc.device)

			if tc.expected == "" {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tc.expected {
					t.Errorf("Expected error '%s', got '%s'", tc.expected, err.Error())
				}
			}
		})
	}
}

func TestCreateDeviceConnectionFromDevice(t *testing.T) {
	id := "device-1"
	name := "Test Router"
	host := "192.168.1.1"
	port := 22
	username := "admin"
	password := "secret"

	device := CreateDeviceConnectionFromDevice(id, name, host, port, username, password)

	if device == nil {
		t.Fatal("Expected device connection, got nil")
	}

	if device.ID != id {
		t.Errorf("Expected ID '%s', got '%s'", id, device.ID)
	}

	if device.Name != name {
		t.Errorf("Expected Name '%s', got '%s'", name, device.Name)
	}

	if device.Host != host {
		t.Errorf("Expected Host '%s', got '%s'", host, device.Host)
	}

	if device.Port != port {
		t.Errorf("Expected Port %d, got %d", port, device.Port)
	}

	if device.Username != username {
		t.Errorf("Expected Username '%s', got '%s'", username, device.Username)
	}

	if device.Password != password {
		t.Errorf("Expected Password '%s', got '%s'", password, device.Password)
	}
}

func TestDeviceSSHManager_GetConnectionStats(t *testing.T) {
	manager := NewDeviceSSHManagerWithDefaults()
	defer manager.Close()

	stats := manager.GetConnectionStats()

	if stats == nil {
		t.Error("Expected connection stats, got nil")
	}

	if len(stats) != 0 {
		t.Errorf("Expected empty stats for new manager, got %d entries", len(stats))
	}
}

func TestDeviceSSHManager_Close(t *testing.T) {
	manager := NewDeviceSSHManagerWithDefaults()

	err := manager.Close()

	if err != nil {
		t.Errorf("Expected no error on close, got: %v", err)
	}
}
