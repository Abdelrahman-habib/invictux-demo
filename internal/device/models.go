package device

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

// Device represents a network device
type Device struct {
	ID                string     `json:"id" db:"id"`
	Name              string     `json:"name" db:"name"`
	IPAddress         string     `json:"ipAddress" db:"ip_address"`
	DeviceType        string     `json:"deviceType" db:"device_type"`
	Vendor            string     `json:"vendor" db:"vendor"`
	Username          string     `json:"username" db:"username"`
	PasswordEncrypted []byte     `json:"-" db:"password_encrypted"`
	SSHPort           int        `json:"sshPort" db:"ssh_port"`
	SNMPCommunity     string     `json:"snmpCommunity" db:"snmp_community"`
	Tags              string     `json:"tags" db:"tags"`
	Status            string     `json:"status"`
	LastChecked       *time.Time `json:"lastChecked"`
	CreatedAt         time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time  `json:"updatedAt" db:"updated_at"`
}

// DeviceStatus represents the status of a device
type DeviceStatus string

const (
	StatusOnline  DeviceStatus = "online"
	StatusOffline DeviceStatus = "offline"
	StatusWarning DeviceStatus = "warning"
	StatusError   DeviceStatus = "error"
)

// DeviceType represents the type of network device
type DeviceType string

const (
	TypeRouter             DeviceType = "router"
	TypeSwitch             DeviceType = "switch"
	TypeFirewall           DeviceType = "firewall"
	TypeLoadBalancer       DeviceType = "load_balancer"
	TypeAccessPoint        DeviceType = "access_point"
	TypeWirelessController DeviceType = "wireless_controller"
	TypeVPN                DeviceType = "vpn"
	TypeProxy              DeviceType = "proxy"
	TypeOther              DeviceType = "other"
)

// Vendor represents supported device vendors
type Vendor string

const (
	VendorCisco      Vendor = "cisco"
	VendorJuniper    Vendor = "juniper"
	VendorHP         Vendor = "hp"
	VendorArista     Vendor = "arista"
	VendorFortinet   Vendor = "fortinet"
	VendorPaloAlto   Vendor = "palo_alto"
	VendorCheckPoint Vendor = "checkpoint"
	VendorF5         Vendor = "f5"
	VendorBrocade    Vendor = "brocade"
	VendorDell       Vendor = "dell"
	VendorHuawei     Vendor = "huawei"
	VendorMikroTik   Vendor = "mikrotik"
	VendorUbiquiti   Vendor = "ubiquiti"
	VendorOther      Vendor = "other"
)

// ValidDeviceTypes returns all valid device types
func ValidDeviceTypes() []DeviceType {
	return []DeviceType{
		TypeRouter,
		TypeSwitch,
		TypeFirewall,
		TypeLoadBalancer,
		TypeAccessPoint,
		TypeWirelessController,
		TypeVPN,
		TypeProxy,
		TypeOther,
	}
}

// ValidVendors returns all valid vendors
func ValidVendors() []Vendor {
	return []Vendor{
		VendorCisco,
		VendorJuniper,
		VendorHP,
		VendorArista,
		VendorFortinet,
		VendorPaloAlto,
		VendorCheckPoint,
		VendorF5,
		VendorBrocade,
		VendorDell,
		VendorHuawei,
		VendorMikroTik,
		VendorUbiquiti,
		VendorOther,
	}
}

// IsValidDeviceType checks if the given device type is valid
func IsValidDeviceType(deviceType string) bool {
	for _, validType := range ValidDeviceTypes() {
		if string(validType) == deviceType {
			return true
		}
	}
	return false
}

// IsValidVendor checks if the given vendor is valid
func IsValidVendor(vendor string) bool {
	for _, validVendor := range ValidVendors() {
		if string(validVendor) == vendor {
			return true
		}
	}
	return false
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// Validate validates the device struct
func (d *Device) Validate() error {
	// Validate name
	if err := ValidateName(d.Name); err != nil {
		return err
	}

	// Validate IP address
	if err := ValidateIPAddress(d.IPAddress); err != nil {
		return err
	}

	// Validate device type
	if err := ValidateDeviceType(d.DeviceType); err != nil {
		return err
	}

	// Validate vendor
	if err := ValidateVendor(d.Vendor); err != nil {
		return err
	}

	// Validate username
	if err := ValidateUsername(d.Username); err != nil {
		return err
	}

	// Validate SSH port
	if err := ValidateSSHPort(d.SSHPort); err != nil {
		return err
	}

	// Validate tags
	if err := ValidateTags(d.Tags); err != nil {
		return err
	}

	return nil
}

// ValidateName validates the device name
func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ValidationError{Field: "name", Message: "name cannot be empty"}
	}
	if len(name) > 100 {
		return ValidationError{Field: "name", Message: "name cannot exceed 100 characters"}
	}

	// Check for valid characters (alphanumeric, spaces, hyphens, underscores, dots)
	validNameRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-_.]+$`)
	if !validNameRegex.MatchString(name) {
		return ValidationError{Field: "name", Message: "name contains invalid characters"}
	}

	return nil
}

// ValidateIPAddress validates the IP address format
func ValidateIPAddress(ipAddress string) error {
	ipAddress = strings.TrimSpace(ipAddress)
	if ipAddress == "" {
		return ValidationError{Field: "ipAddress", Message: "IP address cannot be empty"}
	}

	// Parse the IP address
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return ValidationError{Field: "ipAddress", Message: "invalid IP address format"}
	}

	// Check if it's a valid IPv4 or IPv6 address
	if ip.To4() == nil && ip.To16() == nil {
		return ValidationError{Field: "ipAddress", Message: "IP address must be valid IPv4 or IPv6"}
	}

	// Reject loopback addresses for network devices
	if ip.IsLoopback() {
		return ValidationError{Field: "ipAddress", Message: "loopback addresses are not allowed for network devices"}
	}

	return nil
}

// ValidateDeviceType validates the device type
func ValidateDeviceType(deviceType string) error {
	deviceType = strings.TrimSpace(deviceType)
	if deviceType == "" {
		return ValidationError{Field: "deviceType", Message: "device type cannot be empty"}
	}

	if !IsValidDeviceType(deviceType) {
		return ValidationError{Field: "deviceType", Message: fmt.Sprintf("invalid device type: %s", deviceType)}
	}

	return nil
}

// ValidateVendor validates the vendor
func ValidateVendor(vendor string) error {
	vendor = strings.TrimSpace(vendor)
	if vendor == "" {
		return ValidationError{Field: "vendor", Message: "vendor cannot be empty"}
	}

	if !IsValidVendor(vendor) {
		return ValidationError{Field: "vendor", Message: fmt.Sprintf("invalid vendor: %s", vendor)}
	}

	return nil
}

// ValidateUsername validates the username
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return ValidationError{Field: "username", Message: "username cannot be empty"}
	}
	if len(username) > 50 {
		return ValidationError{Field: "username", Message: "username cannot exceed 50 characters"}
	}

	// Check for valid username characters (alphanumeric, hyphens, underscores, dots)
	validUsernameRegex := regexp.MustCompile(`^[a-zA-Z0-9\-_.]+$`)
	if !validUsernameRegex.MatchString(username) {
		return ValidationError{Field: "username", Message: "username contains invalid characters"}
	}

	return nil
}

// ValidateSSHPort validates the SSH port number
func ValidateSSHPort(port int) error {
	if port <= 0 || port > 65535 {
		return ValidationError{Field: "sshPort", Message: "SSH port must be between 1 and 65535"}
	}
	return nil
}

// ValidateTags validates the tags field
func ValidateTags(tags string) error {
	// Tags are optional, so empty is valid
	if tags == "" {
		return nil
	}

	if len(tags) > 500 {
		return ValidationError{Field: "tags", Message: "tags cannot exceed 500 characters"}
	}

	// If tags are provided, validate the format (comma-separated values)
	tagList := strings.Split(tags, ",")
	for _, tag := range tagList {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue // Skip empty tags
		}

		// Check for valid tag characters (alphanumeric, hyphens, underscores)
		validTagRegex := regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
		if !validTagRegex.MatchString(tag) {
			return ValidationError{Field: "tags", Message: fmt.Sprintf("tag '%s' contains invalid characters", tag)}
		}

		if len(tag) > 50 {
			return ValidationError{Field: "tags", Message: fmt.Sprintf("tag '%s' exceeds 50 characters", tag)}
		}
	}

	return nil
}

// SetDefaults sets default values for optional fields
func (d *Device) SetDefaults() {
	if d.SSHPort == 0 {
		d.SSHPort = 22
	}
	if d.Status == "" {
		d.Status = string(StatusOffline)
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now()
	}
	if d.UpdatedAt.IsZero() {
		d.UpdatedAt = time.Now()
	}
}

// UpdateTimestamp updates the UpdatedAt field to current time
func (d *Device) UpdateTimestamp() {
	d.UpdatedAt = time.Now()
}
