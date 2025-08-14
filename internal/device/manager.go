package device

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Manager handles device CRUD operations
type Manager struct {
	db *sql.DB
}

// NewManager creates a new device manager
func NewManager(db *sql.DB) *Manager {
	return &Manager{db: db}
}

// AddDevice adds a new network device
func (m *Manager) AddDevice(device *Device) error {
	device.ID = uuid.New().String()
	device.CreatedAt = time.Now()
	device.UpdatedAt = time.Now()

	query := `
		INSERT INTO devices (id, name, ip_address, device_type, vendor, username, 
			password_encrypted, ssh_port, snmp_community, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := m.db.Exec(query, device.ID, device.Name, device.IPAddress,
		device.DeviceType, device.Vendor, device.Username, device.PasswordEncrypted,
		device.SSHPort, device.SNMPCommunity, device.Tags, device.CreatedAt, device.UpdatedAt)

	return err
}

// GetAllDevices retrieves all devices
func (m *Manager) GetAllDevices() ([]Device, error) {
	query := `
		SELECT id, name, ip_address, device_type, vendor, username, 
			password_encrypted, ssh_port, snmp_community, tags, created_at, updated_at
		FROM devices
		ORDER BY created_at DESC
	`

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var device Device
		err := rows.Scan(&device.ID, &device.Name, &device.IPAddress,
			&device.DeviceType, &device.Vendor, &device.Username,
			&device.PasswordEncrypted, &device.SSHPort, &device.SNMPCommunity,
			&device.Tags, &device.CreatedAt, &device.UpdatedAt)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// GetDevice retrieves a device by ID
func (m *Manager) GetDevice(id string) (*Device, error) {
	query := `
		SELECT id, name, ip_address, device_type, vendor, username, 
			password_encrypted, ssh_port, snmp_community, tags, created_at, updated_at
		FROM devices
		WHERE id = ?
	`

	var device Device
	err := m.db.QueryRow(query, id).Scan(&device.ID, &device.Name, &device.IPAddress,
		&device.DeviceType, &device.Vendor, &device.Username,
		&device.PasswordEncrypted, &device.SSHPort, &device.SNMPCommunity,
		&device.Tags, &device.CreatedAt, &device.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &device, nil
}

// UpdateDevice updates an existing device
func (m *Manager) UpdateDevice(device *Device) error {
	device.UpdatedAt = time.Now()

	query := `
		UPDATE devices 
		SET name = ?, ip_address = ?, device_type = ?, vendor = ?, username = ?,
			password_encrypted = ?, ssh_port = ?, snmp_community = ?, tags = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := m.db.Exec(query, device.Name, device.IPAddress, device.DeviceType,
		device.Vendor, device.Username, device.PasswordEncrypted, device.SSHPort,
		device.SNMPCommunity, device.Tags, device.UpdatedAt, device.ID)

	return err
}

// DeleteDevice removes a device
func (m *Manager) DeleteDevice(id string) error {
	query := `DELETE FROM devices WHERE id = ?`
	result, err := m.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("device with id %s not found", id)
	}

	return nil
}
