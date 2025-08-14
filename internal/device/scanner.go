package device

import (
	"fmt"
	"net"
	"time"
)

// Scanner handles device connectivity testing
type Scanner struct {
	timeout time.Duration
}

// NewScanner creates a new device scanner
func NewScanner() *Scanner {
	return &Scanner{
		timeout: 5 * time.Second,
	}
}

// TestConnectivity tests if a device is reachable
func (s *Scanner) TestConnectivity(device *Device) error {
	// Test basic network connectivity first
	if err := s.testPing(device.IPAddress); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	// Test SSH port connectivity
	if err := s.testPort(device.IPAddress, device.SSHPort); err != nil {
		return fmt.Errorf("SSH port %d not accessible: %w", device.SSHPort, err)
	}

	return nil
}

// testPing tests basic network connectivity
func (s *Scanner) testPing(ipAddress string) error {
	// For now, we'll use TCP connection test instead of ICMP ping
	// as ICMP requires elevated privileges
	return s.testPort(ipAddress, 22) // Default SSH port for basic connectivity
}

// testPort tests if a specific port is accessible
func (s *Scanner) testPort(ipAddress string, port int) error {
	address := fmt.Sprintf("%s:%d", ipAddress, port)
	conn, err := net.DialTimeout("tcp", address, s.timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

// SetTimeout sets the connection timeout
func (s *Scanner) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}
