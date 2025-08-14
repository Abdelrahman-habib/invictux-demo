package device

import (
	"context"
	"fmt"
	"net"
	"time"
)

// ConnectivityResult represents the result of a connectivity test
type ConnectivityResult struct {
	Device           *Device       `json:"device"`
	NetworkReachable bool          `json:"networkReachable"`
	SSHPortOpen      bool          `json:"sshPortOpen"`
	ResponseTime     time.Duration `json:"responseTime"`
	Error            error         `json:"error,omitempty"`
	TestedAt         time.Time     `json:"testedAt"`
}

// ConnectivityScanner handles device connectivity testing
type ConnectivityScanner struct {
	timeout        time.Duration
	maxRetries     int
	baseRetryDelay time.Duration
}

// ScannerInterface defines the interface for connectivity scanning
type ScannerInterface interface {
	TestConnectivity(device *Device) (*ConnectivityResult, error)
	TestConnectivityWithContext(ctx context.Context, device *Device) (*ConnectivityResult, error)
	BulkTestConnectivity(devices []*Device) ([]*ConnectivityResult, error)
	BulkTestConnectivityWithContext(ctx context.Context, devices []*Device) ([]*ConnectivityResult, error)
}

// NewConnectivityScanner creates a new connectivity scanner with default settings
func NewConnectivityScanner() *ConnectivityScanner {
	return &ConnectivityScanner{
		timeout:        10 * time.Second,
		maxRetries:     3,
		baseRetryDelay: 1 * time.Second,
	}
}

// NewConnectivityScannerWithConfig creates a new connectivity scanner with custom configuration
func NewConnectivityScannerWithConfig(timeout time.Duration, maxRetries int, baseRetryDelay time.Duration) *ConnectivityScanner {
	return &ConnectivityScanner{
		timeout:        timeout,
		maxRetries:     maxRetries,
		baseRetryDelay: baseRetryDelay,
	}
}

// TestConnectivity tests connectivity to a device with default context
func (s *ConnectivityScanner) TestConnectivity(device *Device) (*ConnectivityResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	return s.TestConnectivityWithContext(ctx, device)
}

// TestConnectivityWithContext tests connectivity to a device with custom context
func (s *ConnectivityScanner) TestConnectivityWithContext(ctx context.Context, device *Device) (*ConnectivityResult, error) {
	if device == nil {
		return nil, fmt.Errorf("device cannot be nil")
	}

	if err := device.Validate(); err != nil {
		return nil, fmt.Errorf("device validation failed: %w", err)
	}

	result := &ConnectivityResult{
		Device:   device,
		TestedAt: time.Now(),
	}

	startTime := time.Now()

	// Test network reachability with retry logic
	networkReachable, err := s.testNetworkReachabilityWithRetry(ctx, device.IPAddress)
	result.NetworkReachable = networkReachable

	if err != nil {
		result.Error = fmt.Errorf("network reachability test failed: %w", err)
		result.ResponseTime = time.Since(startTime)
		return result, nil
	}

	// If network is reachable, test SSH port accessibility
	if networkReachable {
		sshPortOpen, err := s.testSSHPortWithRetry(ctx, device.IPAddress, device.SSHPort)
		result.SSHPortOpen = sshPortOpen

		if err != nil {
			result.Error = fmt.Errorf("SSH port test failed: %w", err)
		}
	}

	result.ResponseTime = time.Since(startTime)
	return result, nil
}

// BulkTestConnectivity tests connectivity for multiple devices concurrently
func (s *ConnectivityScanner) BulkTestConnectivity(devices []*Device) ([]*ConnectivityResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout*time.Duration(len(devices)))
	defer cancel()

	return s.BulkTestConnectivityWithContext(ctx, devices)
}

// BulkTestConnectivityWithContext tests connectivity for multiple devices concurrently with custom context
func (s *ConnectivityScanner) BulkTestConnectivityWithContext(ctx context.Context, devices []*Device) ([]*ConnectivityResult, error) {
	if len(devices) == 0 {
		return []*ConnectivityResult{}, nil
	}

	results := make([]*ConnectivityResult, len(devices))
	resultChan := make(chan struct {
		index  int
		result *ConnectivityResult
		err    error
	}, len(devices))

	// Start goroutines for each device
	for i, device := range devices {
		go func(index int, dev *Device) {
			result, err := s.TestConnectivityWithContext(ctx, dev)
			resultChan <- struct {
				index  int
				result *ConnectivityResult
				err    error
			}{index, result, err}
		}(i, device)
	}

	// Collect results
	for i := 0; i < len(devices); i++ {
		select {
		case res := <-resultChan:
			if res.err != nil {
				// Create error result for failed tests
				results[res.index] = &ConnectivityResult{
					Device:   devices[res.index],
					Error:    res.err,
					TestedAt: time.Now(),
				}
			} else {
				results[res.index] = res.result
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("bulk connectivity test cancelled: %w", ctx.Err())
		}
	}

	return results, nil
}

// testNetworkReachabilityWithRetry tests basic network reachability with retry logic
func (s *ConnectivityScanner) testNetworkReachabilityWithRetry(ctx context.Context, ipAddress string) (bool, error) {
	var lastErr error

	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		if attempt > 0 {
			// Calculate exponential backoff delay
			delay := time.Duration(attempt) * s.baseRetryDelay
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return false, ctx.Err()
			}
		}

		reachable, err := s.testNetworkReachability(ctx, ipAddress)
		if err == nil {
			return reachable, nil
		}

		lastErr = err

		// Check if context was cancelled
		if ctx.Err() != nil {
			return false, ctx.Err()
		}
	}

	return false, fmt.Errorf("network reachability test failed after %d attempts: %w", s.maxRetries+1, lastErr)
}

// testNetworkReachability tests basic network reachability using ICMP ping simulation
func (s *ConnectivityScanner) testNetworkReachability(ctx context.Context, ipAddress string) (bool, error) {
	// Use TCP connection attempt to port 80 or 443 as a basic reachability test
	// This is more reliable than ICMP ping in many network environments
	ports := []int{80, 443, 22, 23, 53} // Common ports that are often open

	for _, port := range ports {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ipAddress, port), 3*time.Second)
		if err == nil {
			conn.Close()
			return true, nil
		}

		// Check if context was cancelled
		if ctx.Err() != nil {
			return false, ctx.Err()
		}
	}

	// If no common ports are open, the device might still be reachable but firewalled
	// Try a direct connection test with a very short timeout
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:22", ipAddress), 1*time.Second)
	if err == nil {
		conn.Close()
		return true, nil
	}

	// Check for specific network errors that indicate the host is reachable but port is closed
	if netErr, ok := err.(net.Error); ok {
		if netErr.Timeout() {
			// Timeout could mean host is reachable but port is filtered
			return true, nil
		}
	}

	return false, fmt.Errorf("host appears to be unreachable: %w", err)
}

// testSSHPortWithRetry tests SSH port accessibility with retry logic
func (s *ConnectivityScanner) testSSHPortWithRetry(ctx context.Context, ipAddress string, port int) (bool, error) {
	var lastErr error

	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		if attempt > 0 {
			// Calculate exponential backoff delay
			delay := time.Duration(attempt) * s.baseRetryDelay
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return false, ctx.Err()
			}
		}

		accessible, err := s.testSSHPort(ctx, ipAddress, port)
		if err == nil {
			return accessible, nil
		}

		lastErr = err

		// Check if context was cancelled
		if ctx.Err() != nil {
			return false, ctx.Err()
		}
	}

	return false, fmt.Errorf("SSH port test failed after %d attempts: %w", s.maxRetries+1, lastErr)
}

// testSSHPort tests SSH port accessibility
func (s *ConnectivityScanner) testSSHPort(ctx context.Context, ipAddress string, port int) (bool, error) {
	address := fmt.Sprintf("%s:%d", ipAddress, port)

	// Create a dialer with timeout
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		// Check for specific error types
		if netErr, ok := err.(net.Error); ok {
			if netErr.Timeout() {
				return false, fmt.Errorf("SSH port connection timeout")
			}
		}
		return false, fmt.Errorf("SSH port connection failed: %w", err)
	}

	conn.Close()
	return true, nil
}

// SetTimeout sets the default timeout for connectivity tests
func (s *ConnectivityScanner) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

// SetMaxRetries sets the maximum number of retry attempts
func (s *ConnectivityScanner) SetMaxRetries(maxRetries int) {
	s.maxRetries = maxRetries
}

// SetBaseRetryDelay sets the base delay for exponential backoff
func (s *ConnectivityScanner) SetBaseRetryDelay(delay time.Duration) {
	s.baseRetryDelay = delay
}

// GetTimeout returns the current timeout setting
func (s *ConnectivityScanner) GetTimeout() time.Duration {
	return s.timeout
}

// GetMaxRetries returns the current max retries setting
func (s *ConnectivityScanner) GetMaxRetries() int {
	return s.maxRetries
}

// GetBaseRetryDelay returns the current base retry delay setting
func (s *ConnectivityScanner) GetBaseRetryDelay() time.Duration {
	return s.baseRetryDelay
}
