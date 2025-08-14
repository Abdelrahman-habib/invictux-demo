package device

import (
	"strings"
	"testing"
	"time"
)

func TestDevice_Validate(t *testing.T) {
	tests := []struct {
		name    string
		device  Device
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid device",
			device: Device{
				Name:       "Test Router",
				IPAddress:  "192.168.1.1",
				DeviceType: string(TypeRouter),
				Vendor:     string(VendorCisco),
				Username:   "admin",
				SSHPort:    22,
				Tags:       "production,core",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			device: Device{
				Name:       "",
				IPAddress:  "192.168.1.1",
				DeviceType: string(TypeRouter),
				Vendor:     string(VendorCisco),
				Username:   "admin",
				SSHPort:    22,
			},
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
		{
			name: "invalid IP address",
			device: Device{
				Name:       "Test Router",
				IPAddress:  "invalid-ip",
				DeviceType: string(TypeRouter),
				Vendor:     string(VendorCisco),
				Username:   "admin",
				SSHPort:    22,
			},
			wantErr: true,
			errMsg:  "invalid IP address format",
		},
		{
			name: "invalid device type",
			device: Device{
				Name:       "Test Router",
				IPAddress:  "192.168.1.1",
				DeviceType: "invalid-type",
				Vendor:     string(VendorCisco),
				Username:   "admin",
				SSHPort:    22,
			},
			wantErr: true,
			errMsg:  "invalid device type",
		},
		{
			name: "invalid vendor",
			device: Device{
				Name:       "Test Router",
				IPAddress:  "192.168.1.1",
				DeviceType: string(TypeRouter),
				Vendor:     "invalid-vendor",
				Username:   "admin",
				SSHPort:    22,
			},
			wantErr: true,
			errMsg:  "invalid vendor",
		},
		{
			name: "invalid SSH port",
			device: Device{
				Name:       "Test Router",
				IPAddress:  "192.168.1.1",
				DeviceType: string(TypeRouter),
				Vendor:     string(VendorCisco),
				Username:   "admin",
				SSHPort:    70000,
			},
			wantErr: true,
			errMsg:  "SSH port must be between 1 and 65535",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.device.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Device.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Device.Validate() error = %v, expected to contain %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid name", "Test Router", false, ""},
		{"valid name with numbers", "Router-01", false, ""},
		{"valid name with underscores", "Core_Switch_01", false, ""},
		{"valid name with dots", "router.example.com", false, ""},
		{"empty name", "", true, "name cannot be empty"},
		{"whitespace only", "   ", true, "name cannot be empty"},
		{"too long name", strings.Repeat("a", 101), true, "name cannot exceed 100 characters"},
		{"invalid characters", "Router@#$", true, "name contains invalid characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateName() error = %v, expected to contain %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateIPAddress(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid IPv4", "192.168.1.1", false, ""},
		{"valid IPv4 standard", "10.0.0.1", false, ""},
		{"valid IPv6", "2001:db8::1", false, ""},
		{"empty IP", "", true, "IP address cannot be empty"},
		{"whitespace only", "   ", true, "IP address cannot be empty"},
		{"invalid format", "invalid-ip", true, "invalid IP address format"},
		{"invalid IPv4", "256.256.256.256", true, "invalid IP address format"},
		{"loopback address", "127.0.0.1", true, "loopback addresses are not allowed"},
		{"IPv6 loopback", "::1", true, "loopback addresses are not allowed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIPAddress(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIPAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateIPAddress() error = %v, expected to contain %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateDeviceType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid router", string(TypeRouter), false, ""},
		{"valid switch", string(TypeSwitch), false, ""},
		{"valid firewall", string(TypeFirewall), false, ""},
		{"valid load balancer", string(TypeLoadBalancer), false, ""},
		{"valid access point", string(TypeAccessPoint), false, ""},
		{"valid other", string(TypeOther), false, ""},
		{"empty type", "", true, "device type cannot be empty"},
		{"whitespace only", "   ", true, "device type cannot be empty"},
		{"invalid type", "invalid-type", true, "invalid device type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDeviceType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDeviceType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateDeviceType() error = %v, expected to contain %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateVendor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid cisco", string(VendorCisco), false, ""},
		{"valid juniper", string(VendorJuniper), false, ""},
		{"valid hp", string(VendorHP), false, ""},
		{"valid arista", string(VendorArista), false, ""},
		{"valid fortinet", string(VendorFortinet), false, ""},
		{"valid other", string(VendorOther), false, ""},
		{"empty vendor", "", true, "vendor cannot be empty"},
		{"whitespace only", "   ", true, "vendor cannot be empty"},
		{"invalid vendor", "invalid-vendor", true, "invalid vendor"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVendor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVendor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateVendor() error = %v, expected to contain %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid username", "admin", false, ""},
		{"valid username with numbers", "user123", false, ""},
		{"valid username with underscore", "test_user", false, ""},
		{"valid username with hyphen", "test-user", false, ""},
		{"valid username with dot", "user.name", false, ""},
		{"empty username", "", true, "username cannot be empty"},
		{"whitespace only", "   ", true, "username cannot be empty"},
		{"too long username", strings.Repeat("a", 51), true, "username cannot exceed 50 characters"},
		{"invalid characters", "user@domain", true, "username contains invalid characters"},
		{"invalid characters with space", "user name", true, "username contains invalid characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateUsername() error = %v, expected to contain %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateSSHPort(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
		errMsg  string
	}{
		{"valid port 22", 22, false, ""},
		{"valid port 2222", 2222, false, ""},
		{"valid port 65535", 65535, false, ""},
		{"invalid port 0", 0, true, "SSH port must be between 1 and 65535"},
		{"invalid port negative", -1, true, "SSH port must be between 1 and 65535"},
		{"invalid port too high", 65536, true, "SSH port must be between 1 and 65535"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSSHPort(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSSHPort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateSSHPort() error = %v, expected to contain %v", err, tt.errMsg)
			}
		})
	}
}

func TestValidateTags(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"empty tags", "", false, ""},
		{"valid single tag", "production", false, ""},
		{"valid multiple tags", "production,core,critical", false, ""},
		{"valid tags with spaces", "production, core, critical", false, ""},
		{"valid tags with underscores", "prod_env,core_switch", false, ""},
		{"valid tags with hyphens", "prod-env,core-switch", false, ""},
		{"too long tags", strings.Repeat("a", 501), true, "tags cannot exceed 500 characters"},
		{"invalid tag characters", "prod@env,core", true, "tag 'prod@env' contains invalid characters"},
		{"tag too long", "production," + strings.Repeat("a", 51), true, "exceeds 50 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTags(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateTags() error = %v, expected to contain %v", err, tt.errMsg)
			}
		})
	}
}

func TestIsValidDeviceType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid router", string(TypeRouter), true},
		{"valid switch", string(TypeSwitch), true},
		{"valid firewall", string(TypeFirewall), true},
		{"valid load balancer", string(TypeLoadBalancer), true},
		{"valid access point", string(TypeAccessPoint), true},
		{"valid wireless controller", string(TypeWirelessController), true},
		{"valid vpn", string(TypeVPN), true},
		{"valid proxy", string(TypeProxy), true},
		{"valid other", string(TypeOther), true},
		{"invalid type", "invalid", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidDeviceType(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidDeviceType() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsValidVendor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid cisco", string(VendorCisco), true},
		{"valid juniper", string(VendorJuniper), true},
		{"valid hp", string(VendorHP), true},
		{"valid arista", string(VendorArista), true},
		{"valid fortinet", string(VendorFortinet), true},
		{"valid palo alto", string(VendorPaloAlto), true},
		{"valid checkpoint", string(VendorCheckPoint), true},
		{"valid f5", string(VendorF5), true},
		{"valid brocade", string(VendorBrocade), true},
		{"valid dell", string(VendorDell), true},
		{"valid huawei", string(VendorHuawei), true},
		{"valid mikrotik", string(VendorMikroTik), true},
		{"valid ubiquiti", string(VendorUbiquiti), true},
		{"valid other", string(VendorOther), true},
		{"invalid vendor", "invalid", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidVendor(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidVendor() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestDevice_SetDefaults(t *testing.T) {
	device := &Device{
		Name:       "Test Device",
		IPAddress:  "192.168.1.1",
		DeviceType: string(TypeRouter),
		Vendor:     string(VendorCisco),
		Username:   "admin",
	}

	device.SetDefaults()

	if device.SSHPort != 22 {
		t.Errorf("SetDefaults() SSHPort = %v, expected 22", device.SSHPort)
	}

	if device.Status != string(StatusOffline) {
		t.Errorf("SetDefaults() Status = %v, expected %v", device.Status, string(StatusOffline))
	}

	if device.CreatedAt.IsZero() {
		t.Error("SetDefaults() CreatedAt should not be zero")
	}

	if device.UpdatedAt.IsZero() {
		t.Error("SetDefaults() UpdatedAt should not be zero")
	}
}

func TestDevice_UpdateTimestamp(t *testing.T) {
	device := &Device{
		UpdatedAt: time.Now().Add(-time.Hour), // Set to 1 hour ago
	}

	oldTime := device.UpdatedAt
	time.Sleep(time.Millisecond) // Ensure time difference
	device.UpdateTimestamp()

	if !device.UpdatedAt.After(oldTime) {
		t.Error("UpdateTimestamp() should update UpdatedAt to a more recent time")
	}
}

func TestValidDeviceTypes(t *testing.T) {
	types := ValidDeviceTypes()

	expectedCount := 9 // Update this if you add more device types
	if len(types) != expectedCount {
		t.Errorf("ValidDeviceTypes() returned %d types, expected %d", len(types), expectedCount)
	}

	// Check that all expected types are present
	expectedTypes := []DeviceType{
		TypeRouter, TypeSwitch, TypeFirewall, TypeLoadBalancer,
		TypeAccessPoint, TypeWirelessController, TypeVPN, TypeProxy, TypeOther,
	}

	for _, expectedType := range expectedTypes {
		found := false
		for _, actualType := range types {
			if actualType == expectedType {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ValidDeviceTypes() missing expected type: %s", expectedType)
		}
	}
}

func TestValidVendors(t *testing.T) {
	vendors := ValidVendors()

	expectedCount := 14 // Update this if you add more vendors
	if len(vendors) != expectedCount {
		t.Errorf("ValidVendors() returned %d vendors, expected %d", len(vendors), expectedCount)
	}

	// Check that all expected vendors are present
	expectedVendors := []Vendor{
		VendorCisco, VendorJuniper, VendorHP, VendorArista, VendorFortinet,
		VendorPaloAlto, VendorCheckPoint, VendorF5, VendorBrocade, VendorDell,
		VendorHuawei, VendorMikroTik, VendorUbiquiti, VendorOther,
	}

	for _, expectedVendor := range expectedVendors {
		found := false
		for _, actualVendor := range vendors {
			if actualVendor == expectedVendor {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ValidVendors() missing expected vendor: %s", expectedVendor)
		}
	}
}
