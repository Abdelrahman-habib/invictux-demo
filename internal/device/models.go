package device

import "time"

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
	TypeRouter   DeviceType = "router"
	TypeSwitch   DeviceType = "switch"
	TypeFirewall DeviceType = "firewall"
	TypeOther    DeviceType = "other"
)

// Vendor represents supported device vendors
type Vendor string

const (
	VendorCisco   Vendor = "cisco"
	VendorJuniper Vendor = "juniper"
	VendorHP      Vendor = "hp"
	VendorOther   Vendor = "other"
)
