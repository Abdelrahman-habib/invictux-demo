package device

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

// Manager handles device CRUD operations
type Manager struct {
	db *sql.DB
}

// ManagerInterface defines the interface for device management operations
type ManagerInterface interface {
	AddDevice(device *Device) error
	GetAllDevices() ([]Device, error)
	GetDevice(id string) (*Device, error)
	GetDeviceByIP(ipAddress string) (*Device, error)
	UpdateDevice(device *Device) error
	DeleteDevice(id string) error
	TestConnectivity(device *Device) error
}

// DeviceError represents device-specific errors
type DeviceError struct {
	Type    string
	Message string
	Field   string
}

func (e *DeviceError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s error for field '%s': %s", e.Type, e.Field, e.Message)
	}
	return fmt.Sprintf("%s error: %s", e.Type, e.Message)
}

// Error types
const (
	ErrorTypeValidation = "validation"
	ErrorTypeDuplicate  = "duplicate"
	ErrorTypeNotFound   = "not_found"
	ErrorTypeDatabase   = "database"
)

// NewManager creates a new device manager
func NewManager(db *sql.DB) *Manager {
	return &Manager{db: db}
}

// AddDevice adds a new network device with proper validation and duplicate checking
func (m *Manager) AddDevice(device *Device) error {
	// Validate the device
	if err := device.Validate(); err != nil {
		return &DeviceError{
			Type:    ErrorTypeValidation,
			Message: err.Error(),
		}
	}

	// Set defaults and generate ID
	device.SetDefaults()
	device.ID = uuid.New().String()
	device.CreatedAt = time.Now()
	device.UpdatedAt = time.Now()

	// Start transaction for atomic operation
	tx, err := m.db.Begin()
	if err != nil {
		return &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to begin transaction: %v", err),
		}
	}
	defer tx.Rollback()

	// Check for duplicate IP address
	var existingID string
	checkQuery := `SELECT id FROM devices WHERE ip_address = ?`
	err = tx.QueryRow(checkQuery, device.IPAddress).Scan(&existingID)
	if err == nil {
		return &DeviceError{
			Type:    ErrorTypeDuplicate,
			Field:   "ipAddress",
			Message: fmt.Sprintf("device with IP address %s already exists", device.IPAddress),
		}
	} else if err != sql.ErrNoRows {
		return &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to check for duplicate IP: %v", err),
		}
	}

	// Insert the device
	insertQuery := `
		INSERT INTO devices (id, name, ip_address, device_type, vendor, username, 
			password_encrypted, ssh_port, snmp_community, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = tx.Exec(insertQuery, device.ID, device.Name, device.IPAddress,
		device.DeviceType, device.Vendor, device.Username, device.PasswordEncrypted,
		device.SSHPort, device.SNMPCommunity, device.Tags, device.CreatedAt, device.UpdatedAt)

	if err != nil {
		// Check if it's a SQLite constraint error
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return &DeviceError{
					Type:    ErrorTypeDuplicate,
					Field:   "ipAddress",
					Message: fmt.Sprintf("device with IP address %s already exists", device.IPAddress),
				}
			}
		}
		return &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to insert device: %v", err),
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to commit transaction: %v", err),
		}
	}

	return nil
}

// GetAllDevices retrieves all devices with proper error handling
func (m *Manager) GetAllDevices() ([]Device, error) {
	query := `
		SELECT id, name, ip_address, device_type, vendor, username, 
			password_encrypted, ssh_port, snmp_community, tags, created_at, updated_at
		FROM devices
		ORDER BY created_at DESC
	`

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to query devices: %v", err),
		}
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
			return nil, &DeviceError{
				Type:    ErrorTypeDatabase,
				Message: fmt.Sprintf("failed to scan device row: %v", err),
			}
		}
		devices = append(devices, device)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("error iterating over device rows: %v", err),
		}
	}

	return devices, nil
}

// GetDevice retrieves a device by ID with proper error handling
func (m *Manager) GetDevice(id string) (*Device, error) {
	if strings.TrimSpace(id) == "" {
		return nil, &DeviceError{
			Type:    ErrorTypeValidation,
			Field:   "id",
			Message: "device ID cannot be empty",
		}
	}

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
		if err == sql.ErrNoRows {
			return nil, &DeviceError{
				Type:    ErrorTypeNotFound,
				Message: fmt.Sprintf("device with ID %s not found", id),
			}
		}
		return nil, &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to get device: %v", err),
		}
	}

	return &device, nil
}

// GetDeviceByIP retrieves a device by IP address
func (m *Manager) GetDeviceByIP(ipAddress string) (*Device, error) {
	if strings.TrimSpace(ipAddress) == "" {
		return nil, &DeviceError{
			Type:    ErrorTypeValidation,
			Field:   "ipAddress",
			Message: "IP address cannot be empty",
		}
	}

	query := `
		SELECT id, name, ip_address, device_type, vendor, username, 
			password_encrypted, ssh_port, snmp_community, tags, created_at, updated_at
		FROM devices
		WHERE ip_address = ?
	`

	var device Device
	err := m.db.QueryRow(query, ipAddress).Scan(&device.ID, &device.Name, &device.IPAddress,
		&device.DeviceType, &device.Vendor, &device.Username,
		&device.PasswordEncrypted, &device.SSHPort, &device.SNMPCommunity,
		&device.Tags, &device.CreatedAt, &device.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &DeviceError{
				Type:    ErrorTypeNotFound,
				Message: fmt.Sprintf("device with IP address %s not found", ipAddress),
			}
		}
		return nil, &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to get device by IP: %v", err),
		}
	}

	return &device, nil
}

// UpdateDevice updates an existing device with proper validation and duplicate checking
func (m *Manager) UpdateDevice(device *Device) error {
	if strings.TrimSpace(device.ID) == "" {
		return &DeviceError{
			Type:    ErrorTypeValidation,
			Field:   "id",
			Message: "device ID cannot be empty",
		}
	}

	// Validate the device
	if err := device.Validate(); err != nil {
		return &DeviceError{
			Type:    ErrorTypeValidation,
			Message: err.Error(),
		}
	}

	device.UpdateTimestamp()

	// Start transaction for atomic operation
	tx, err := m.db.Begin()
	if err != nil {
		return &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to begin transaction: %v", err),
		}
	}
	defer tx.Rollback()

	// Check if device exists
	var existingID string
	checkExistsQuery := `SELECT id FROM devices WHERE id = ?`
	err = tx.QueryRow(checkExistsQuery, device.ID).Scan(&existingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &DeviceError{
				Type:    ErrorTypeNotFound,
				Message: fmt.Sprintf("device with ID %s not found", device.ID),
			}
		}
		return &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to check device existence: %v", err),
		}
	}

	// Check for duplicate IP address (excluding current device)
	var duplicateID string
	checkDuplicateQuery := `SELECT id FROM devices WHERE ip_address = ? AND id != ?`
	err = tx.QueryRow(checkDuplicateQuery, device.IPAddress, device.ID).Scan(&duplicateID)
	if err == nil {
		return &DeviceError{
			Type:    ErrorTypeDuplicate,
			Field:   "ipAddress",
			Message: fmt.Sprintf("another device with IP address %s already exists", device.IPAddress),
		}
	} else if err != sql.ErrNoRows {
		return &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to check for duplicate IP: %v", err),
		}
	}

	// Update the device
	updateQuery := `
		UPDATE devices 
		SET name = ?, ip_address = ?, device_type = ?, vendor = ?, username = ?,
			password_encrypted = ?, ssh_port = ?, snmp_community = ?, tags = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := tx.Exec(updateQuery, device.Name, device.IPAddress, device.DeviceType,
		device.Vendor, device.Username, device.PasswordEncrypted, device.SSHPort,
		device.SNMPCommunity, device.Tags, device.UpdatedAt, device.ID)

	if err != nil {
		// Check if it's a SQLite constraint error
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return &DeviceError{
					Type:    ErrorTypeDuplicate,
					Field:   "ipAddress",
					Message: fmt.Sprintf("device with IP address %s already exists", device.IPAddress),
				}
			}
		}
		return &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to update device: %v", err),
		}
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to get rows affected: %v", err),
		}
	}

	if rowsAffected == 0 {
		return &DeviceError{
			Type:    ErrorTypeNotFound,
			Message: fmt.Sprintf("device with ID %s not found", device.ID),
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to commit transaction: %v", err),
		}
	}

	return nil
}

// DeleteDevice removes a device with proper error handling and transaction support
func (m *Manager) DeleteDevice(id string) error {
	if strings.TrimSpace(id) == "" {
		return &DeviceError{
			Type:    ErrorTypeValidation,
			Field:   "id",
			Message: "device ID cannot be empty",
		}
	}

	// Start transaction for atomic operation
	tx, err := m.db.Begin()
	if err != nil {
		return &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to begin transaction: %v", err),
		}
	}
	defer tx.Rollback()

	// Delete the device (CASCADE will handle related records)
	deleteQuery := `DELETE FROM devices WHERE id = ?`
	result, err := tx.Exec(deleteQuery, id)
	if err != nil {
		return &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to delete device: %v", err),
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to get rows affected: %v", err),
		}
	}

	if rowsAffected == 0 {
		return &DeviceError{
			Type:    ErrorTypeNotFound,
			Message: fmt.Sprintf("device with ID %s not found", id),
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return &DeviceError{
			Type:    ErrorTypeDatabase,
			Message: fmt.Sprintf("failed to commit transaction: %v", err),
		}
	}

	return nil
}

// TestConnectivity tests the connectivity to a device (placeholder implementation)
func (m *Manager) TestConnectivity(device *Device) error {
	// This is a placeholder implementation for the interface
	// The actual connectivity testing will be implemented in the scanner component
	if device == nil {
		return &DeviceError{
			Type:    ErrorTypeValidation,
			Message: "device cannot be nil",
		}
	}

	if err := device.Validate(); err != nil {
		return &DeviceError{
			Type:    ErrorTypeValidation,
			Message: err.Error(),
		}
	}

	// TODO: Implement actual connectivity testing
	// This will be done in task 2.3 "Build device connectivity scanner"
	return nil
}
