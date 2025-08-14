package device

import (
	"context"
	"testing"
	"time"
)

// MockNetworkInterface for testing network operations
type MockNetworkInterface struct {
	shouldFail        bool
	shouldTimeout     bool
	responseDelay     time.Duration
	failAfterAttempts int
	attemptCount      int
}

func TestNewConnectivityScanner(t *testing.T) {
	scanner := NewConnectivityScanner()

	if scanner == nil {
		t.Fatal("NewConnectivityScanner returned nil")
	}

	if scanner.timeout != 10*time.Second {
		t.Errorf("Expected default timeout of 10s, got %v", scanner.timeout)
	}

	if scanner.maxRetries != 3 {
		t.Errorf("Expected default maxRetries of 3, got %d", scanner.maxRetries)
	}

	if scanner.baseRetryDelay != 1*time.Second {
		t.Errorf("Expected default baseRetryDelay of 1s, got %v", scanner.baseRetryDelay)
	}
}

func TestNewConnectivityScannerWithConfig(t *testing.T) {
	timeout := 5 * time.Second
	maxRetries := 2
	baseRetryDelay := 500 * time.Millisecond

	scanner := NewConnectivityScannerWithConfig(timeout, maxRetries, baseRetryDelay)

	if scanner == nil {
		t.Fatal("NewConnectivityScannerWithConfig returned nil")
	}

	if scanner.timeout != timeout {
		t.Errorf("Expected timeout of %v, got %v", timeout, scanner.timeout)
	}

	if scanner.maxRetries != maxRetries {
		t.Errorf("Expected maxRetries of %d, got %d", maxRetries, scanner.maxRetries)
	}

	if scanner.baseRetryDelay != baseRetryDelay {
		t.Errorf("Expected baseRetryDelay of %v, got %v", baseRetryDelay, scanner.baseRetryDelay)
	}
}

func TestConnectivityScanner_TestConnectivity_NilDevice(t *testing.T) {
	scanner := NewConnectivityScanner()

	result, err := scanner.TestConnectivity(nil)

	if err == nil {
		t.Error("Expected error for nil device, got nil")
	}

	if result != nil {
		t.Error("Expected nil result for nil device")
	}

	expectedError := "device cannot be nil"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestConnectivityScanner_TestConnectivity_InvalidDevice(t *testing.T) {
	scanner := NewConnectivityScanner()

	// Create device with invalid IP
	device := &Device{
		Name:       "Test Device",
		IPAddress:  "invalid-ip",
		DeviceType: string(TypeRouter),
		Vendor:     string(VendorCisco),
		Username:   "admin",
		SSHPort:    22,
	}

	result, err := scanner.TestConnectivity(device)

	if err == nil {
		t.Error("Expected error for invalid device, got nil")
	}

	if result != nil {
		t.Error("Expected nil result for invalid device")
	}
}

func TestConnectivityScanner_TestConnectivity_ValidDevice(t *testing.T) {
	scanner := NewConnectivityScanner()

	// Create a valid device (using a valid IP for testing)
	device := &Device{
		Name:       "Test Device",
		IPAddress:  "192.168.1.1",
		DeviceType: string(TypeRouter),
		Vendor:     string(VendorCisco),
		Username:   "admin",
		SSHPort:    22,
	}

	result, err := scanner.TestConnectivity(device)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Device != device {
		t.Error("Result device doesn't match input device")
	}

	if result.TestedAt.IsZero() {
		t.Error("TestedAt should be set")
	}

	if result.ResponseTime <= 0 {
		t.Error("ResponseTime should be positive")
	}
}

func TestConnectivityScanner_TestConnectivityWithContext_Timeout(t *testing.T) {
	scanner := NewConnectivityScanner()

	// Create a context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	device := &Device{
		Name:       "Test Device",
		IPAddress:  "192.168.1.1", // Non-routable IP for testing
		DeviceType: string(TypeRouter),
		Vendor:     string(VendorCisco),
		Username:   "admin",
		SSHPort:    22,
	}

	result, err := scanner.TestConnectivityWithContext(ctx, device)

	// Should not return error immediately, but result should contain timeout info
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// The result should contain error information about the timeout
	if result.Error == nil {
		t.Error("Expected error in result due to timeout")
	}
}

func TestConnectivityScanner_BulkTestConnectivity_EmptySlice(t *testing.T) {
	scanner := NewConnectivityScanner()

	results, err := scanner.BulkTestConnectivity([]*Device{})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected empty results, got %d results", len(results))
	}
}

func TestConnectivityScanner_BulkTestConnectivity_MultipleDevices(t *testing.T) {
	scanner := NewConnectivityScannerWithConfig(30*time.Second, 1, 100*time.Millisecond)

	devices := []*Device{
		{
			Name:       "Device 1",
			IPAddress:  "192.168.1.1",
			DeviceType: string(TypeRouter),
			Vendor:     string(VendorCisco),
			Username:   "admin",
			SSHPort:    22,
		},
		{
			Name:       "Device 2",
			IPAddress:  "192.168.1.2",
			DeviceType: string(TypeSwitch),
			Vendor:     string(VendorCisco),
			Username:   "admin",
			SSHPort:    23,
		},
	}

	results, err := scanner.BulkTestConnectivity(devices)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(results) != len(devices) {
		t.Errorf("Expected %d results, got %d", len(devices), len(results))
	}

	for i, result := range results {
		if result == nil {
			t.Errorf("Result %d is nil", i)
			continue
		}

		if result.Device != devices[i] {
			t.Errorf("Result %d device doesn't match input device", i)
		}

		if result.TestedAt.IsZero() {
			t.Errorf("Result %d TestedAt should be set", i)
		}
	}
}

func TestConnectivityScanner_BulkTestConnectivityWithContext_Cancelled(t *testing.T) {
	scanner := NewConnectivityScanner()

	// Create a context that we'll cancel immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	devices := []*Device{
		{
			Name:       "Device 1",
			IPAddress:  "192.168.1.1",
			DeviceType: string(TypeRouter),
			Vendor:     string(VendorCisco),
			Username:   "admin",
			SSHPort:    22,
		},
	}

	results, err := scanner.BulkTestConnectivityWithContext(ctx, devices)

	if err == nil {
		t.Error("Expected error due to cancelled context")
	}

	if results != nil {
		t.Error("Expected nil results due to cancelled context")
	}
}

func TestConnectivityScanner_SettersAndGetters(t *testing.T) {
	scanner := NewConnectivityScanner()

	// Test timeout
	newTimeout := 15 * time.Second
	scanner.SetTimeout(newTimeout)
	if scanner.GetTimeout() != newTimeout {
		t.Errorf("Expected timeout %v, got %v", newTimeout, scanner.GetTimeout())
	}

	// Test max retries
	newMaxRetries := 5
	scanner.SetMaxRetries(newMaxRetries)
	if scanner.GetMaxRetries() != newMaxRetries {
		t.Errorf("Expected maxRetries %d, got %d", newMaxRetries, scanner.GetMaxRetries())
	}

	// Test base retry delay
	newBaseRetryDelay := 2 * time.Second
	scanner.SetBaseRetryDelay(newBaseRetryDelay)
	if scanner.GetBaseRetryDelay() != newBaseRetryDelay {
		t.Errorf("Expected baseRetryDelay %v, got %v", newBaseRetryDelay, scanner.GetBaseRetryDelay())
	}
}

func TestConnectivityScanner_testNetworkReachability_ReachableHost(t *testing.T) {
	scanner := NewConnectivityScanner()
	ctx := context.Background()

	// Test with Google DNS - should be reachable
	reachable, err := scanner.testNetworkReachability(ctx, "8.8.8.8")

	if err != nil {
		t.Errorf("Unexpected error testing Google DNS: %v", err)
	}

	if !reachable {
		t.Error("Google DNS should be reachable")
	}
}

func TestConnectivityScanner_testNetworkReachability_UnreachableHost(t *testing.T) {
	scanner := NewConnectivityScanner()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Test with non-routable IP - should be unreachable
	reachable, err := scanner.testNetworkReachability(ctx, "192.0.2.1") // RFC5737 test address

	// We expect either an error or false reachability for this test address
	if err == nil && reachable {
		t.Error("Non-routable IP should not be reachable without error")
	}
}

func TestConnectivityScanner_testSSHPort_InvalidPort(t *testing.T) {
	scanner := NewConnectivityScanner()
	ctx := context.Background()

	// Test with invalid port on Google DNS
	accessible, err := scanner.testSSHPort(ctx, "8.8.8.8", 99999)

	if err == nil {
		t.Error("Expected error for invalid port")
	}

	if accessible {
		t.Error("Invalid port should not be accessible")
	}
}

// TestConnectivityScanner_RetryLogic tests the retry mechanism
func TestConnectivityScanner_RetryLogic(t *testing.T) {
	scanner := NewConnectivityScannerWithConfig(5*time.Second, 2, 100*time.Millisecond)

	device := &Device{
		Name:       "Test Device",
		IPAddress:  "192.0.2.1", // RFC5737 test address - should be unreachable
		DeviceType: string(TypeRouter),
		Vendor:     string(VendorCisco),
		Username:   "admin",
		SSHPort:    22,
	}

	startTime := time.Now()
	result, err := scanner.TestConnectivity(device)
	duration := time.Since(startTime)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Should have taken some time due to retries
	expectedMinDuration := 100 * time.Millisecond // At least one retry delay
	if duration < expectedMinDuration {
		t.Errorf("Expected test to take at least %v due to retries, took %v", expectedMinDuration, duration)
	}

	// Result should contain error information
	if result.Error == nil {
		t.Error("Expected error in result for unreachable host")
	}
}

// TestConnectivityResult_Structure tests the ConnectivityResult structure
func TestConnectivityResult_Structure(t *testing.T) {
	device := &Device{
		Name:       "Test Device",
		IPAddress:  "192.168.1.1",
		DeviceType: string(TypeRouter),
		Vendor:     string(VendorCisco),
		Username:   "admin",
		SSHPort:    22,
	}

	result := &ConnectivityResult{
		Device:           device,
		NetworkReachable: true,
		SSHPortOpen:      false,
		ResponseTime:     100 * time.Millisecond,
		TestedAt:         time.Now(),
	}

	if result.Device != device {
		t.Error("Device field not set correctly")
	}

	if !result.NetworkReachable {
		t.Error("NetworkReachable field not set correctly")
	}

	if result.SSHPortOpen {
		t.Error("SSHPortOpen field not set correctly")
	}

	if result.ResponseTime != 100*time.Millisecond {
		t.Error("ResponseTime field not set correctly")
	}

	if result.TestedAt.IsZero() {
		t.Error("TestedAt field should be set")
	}
}

// Benchmark tests for performance
func BenchmarkConnectivityScanner_TestConnectivity(b *testing.B) {
	scanner := NewConnectivityScanner()
	device := &Device{
		Name:       "Test Device",
		IPAddress:  "192.168.1.1",
		DeviceType: string(TypeRouter),
		Vendor:     string(VendorCisco),
		Username:   "admin",
		SSHPort:    22,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := scanner.TestConnectivity(device)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkConnectivityScanner_BulkTestConnectivity(b *testing.B) {
	scanner := NewConnectivityScanner()
	devices := make([]*Device, 10)

	for i := 0; i < 10; i++ {
		devices[i] = &Device{
			Name:       "Test Device",
			IPAddress:  "192.168.1.1",
			DeviceType: string(TypeRouter),
			Vendor:     string(VendorCisco),
			Username:   "admin",
			SSHPort:    22,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := scanner.BulkTestConnectivity(devices)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}
