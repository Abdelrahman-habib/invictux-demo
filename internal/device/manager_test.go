package device

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates a test database for testing
func setupTestDB(t *testing.T) *sql.DB {
	// Create temporary directory for test database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=ON")
	require.NoError(t, err)

	// Create devices table
	createTableSQL := `
		CREATE TABLE devices (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			ip_address TEXT NOT NULL UNIQUE,
			device_type TEXT NOT NULL,
			vendor TEXT NOT NULL,
			username TEXT NOT NULL,
			password_encrypted BLOB NOT NULL,
			ssh_port INTEGER DEFAULT 22,
			snmp_community TEXT,
			tags TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err = db.Exec(createTableSQL)
	require.NoError(t, err)

	return db
}

// createTestDevice creates a valid test device
func createTestDevice() *Device {
	return &Device{
		Name:              "Test Router",
		IPAddress:         "192.168.1.1",
		DeviceType:        string(TypeRouter),
		Vendor:            string(VendorCisco),
		Username:          "admin",
		PasswordEncrypted: []byte("encrypted_password"),
		SSHPort:           22,
		SNMPCommunity:     "public",
		Tags:              "test,router",
	}
}

func TestNewManager(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	manager := NewManager(db)
	assert.NotNil(t, manager)
	assert.Equal(t, db, manager.db)
}

func TestManager_AddDevice(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	manager := NewManager(db)

	t.Run("successful add", func(t *testing.T) {
		device := createTestDevice()

		err := manager.AddDevice(device)
		assert.NoError(t, err)
		assert.NotEmpty(t, device.ID)
		assert.False(t, device.CreatedAt.IsZero())
		assert.False(t, device.UpdatedAt.IsZero())
	})

	t.Run("duplicate IP address", func(t *testing.T) {
		device1 := createTestDevice()
		device1.IPAddress = "192.168.1.10"

		err := manager.AddDevice(device1)
		require.NoError(t, err)

		// Try to add another device with same IP
		device2 := createTestDevice()
		device2.IPAddress = "192.168.1.10"
		device2.Name = "Another Router"

		err = manager.AddDevice(device2)
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeDuplicate, deviceErr.Type)
		assert.Equal(t, "ipAddress", deviceErr.Field)
	})

	t.Run("invalid device validation", func(t *testing.T) {
		device := createTestDevice()
		device.Name = "" // Invalid name

		err := manager.AddDevice(device)
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeValidation, deviceErr.Type)
	})

	t.Run("invalid IP address", func(t *testing.T) {
		device := createTestDevice()
		device.IPAddress = "invalid-ip"

		err := manager.AddDevice(device)
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeValidation, deviceErr.Type)
	})

	t.Run("invalid device type", func(t *testing.T) {
		device := createTestDevice()
		device.DeviceType = "invalid_type"

		err := manager.AddDevice(device)
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeValidation, deviceErr.Type)
	})

	t.Run("invalid vendor", func(t *testing.T) {
		device := createTestDevice()
		device.Vendor = "invalid_vendor"

		err := manager.AddDevice(device)
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeValidation, deviceErr.Type)
	})
}

func TestManager_GetAllDevices(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	manager := NewManager(db)

	t.Run("empty database", func(t *testing.T) {
		devices, err := manager.GetAllDevices()
		assert.NoError(t, err)
		assert.Empty(t, devices)
	})

	t.Run("multiple devices", func(t *testing.T) {
		// Add test devices
		device1 := createTestDevice()
		device1.IPAddress = "192.168.1.1"
		device1.Name = "Router 1"

		device2 := createTestDevice()
		device2.IPAddress = "192.168.1.2"
		device2.Name = "Router 2"

		err := manager.AddDevice(device1)
		require.NoError(t, err)

		// Add small delay to ensure different timestamps
		time.Sleep(time.Millisecond)

		err = manager.AddDevice(device2)
		require.NoError(t, err)

		devices, err := manager.GetAllDevices()
		assert.NoError(t, err)
		assert.Len(t, devices, 2)

		// Should be ordered by created_at DESC (newest first)
		assert.Equal(t, "Router 2", devices[0].Name)
		assert.Equal(t, "Router 1", devices[1].Name)
	})
}

func TestManager_GetDevice(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	manager := NewManager(db)

	t.Run("existing device", func(t *testing.T) {
		device := createTestDevice()
		err := manager.AddDevice(device)
		require.NoError(t, err)

		retrieved, err := manager.GetDevice(device.ID)
		assert.NoError(t, err)
		assert.Equal(t, device.ID, retrieved.ID)
		assert.Equal(t, device.Name, retrieved.Name)
		assert.Equal(t, device.IPAddress, retrieved.IPAddress)
	})

	t.Run("non-existent device", func(t *testing.T) {
		_, err := manager.GetDevice("non-existent-id")
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeNotFound, deviceErr.Type)
	})

	t.Run("empty ID", func(t *testing.T) {
		_, err := manager.GetDevice("")
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeValidation, deviceErr.Type)
		assert.Equal(t, "id", deviceErr.Field)
	})

	t.Run("whitespace ID", func(t *testing.T) {
		_, err := manager.GetDevice("   ")
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeValidation, deviceErr.Type)
		assert.Equal(t, "id", deviceErr.Field)
	})
}

func TestManager_GetDeviceByIP(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	manager := NewManager(db)

	t.Run("existing device", func(t *testing.T) {
		device := createTestDevice()
		device.IPAddress = "192.168.1.100"
		err := manager.AddDevice(device)
		require.NoError(t, err)

		retrieved, err := manager.GetDeviceByIP("192.168.1.100")
		assert.NoError(t, err)
		assert.Equal(t, device.ID, retrieved.ID)
		assert.Equal(t, device.IPAddress, retrieved.IPAddress)
	})

	t.Run("non-existent IP", func(t *testing.T) {
		_, err := manager.GetDeviceByIP("192.168.1.200")
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeNotFound, deviceErr.Type)
	})

	t.Run("empty IP", func(t *testing.T) {
		_, err := manager.GetDeviceByIP("")
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeValidation, deviceErr.Type)
		assert.Equal(t, "ipAddress", deviceErr.Field)
	})
}

func TestManager_UpdateDevice(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	manager := NewManager(db)

	t.Run("successful update", func(t *testing.T) {
		device := createTestDevice()
		err := manager.AddDevice(device)
		require.NoError(t, err)

		originalUpdatedAt := device.UpdatedAt
		time.Sleep(time.Millisecond) // Ensure timestamp difference

		// Update device
		device.Name = "Updated Router"
		device.Tags = "updated,test"

		err = manager.UpdateDevice(device)
		assert.NoError(t, err)
		assert.True(t, device.UpdatedAt.After(originalUpdatedAt))

		// Verify update
		retrieved, err := manager.GetDevice(device.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Router", retrieved.Name)
		assert.Equal(t, "updated,test", retrieved.Tags)
	})

	t.Run("non-existent device", func(t *testing.T) {
		device := createTestDevice()
		device.ID = "non-existent-id"

		err := manager.UpdateDevice(device)
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeNotFound, deviceErr.Type)
	})

	t.Run("empty ID", func(t *testing.T) {
		device := createTestDevice()
		device.ID = ""

		err := manager.UpdateDevice(device)
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeValidation, deviceErr.Type)
		assert.Equal(t, "id", deviceErr.Field)
	})

	t.Run("duplicate IP address", func(t *testing.T) {
		// Add first device
		device1 := createTestDevice()
		device1.IPAddress = "192.168.1.10"
		err := manager.AddDevice(device1)
		require.NoError(t, err)

		// Add second device
		device2 := createTestDevice()
		device2.IPAddress = "192.168.1.11"
		device2.Name = "Second Router"
		err = manager.AddDevice(device2)
		require.NoError(t, err)

		// Try to update second device with first device's IP
		device2.IPAddress = "192.168.1.10"
		err = manager.UpdateDevice(device2)
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeDuplicate, deviceErr.Type)
		assert.Equal(t, "ipAddress", deviceErr.Field)
	})

	t.Run("invalid device validation", func(t *testing.T) {
		device := createTestDevice()
		device.IPAddress = "192.168.1.99" // Use unique IP
		err := manager.AddDevice(device)
		require.NoError(t, err)

		// Make device invalid
		device.Name = ""

		err = manager.UpdateDevice(device)
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeValidation, deviceErr.Type)
	})
}

func TestManager_DeleteDevice(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	manager := NewManager(db)

	t.Run("successful delete", func(t *testing.T) {
		device := createTestDevice()
		err := manager.AddDevice(device)
		require.NoError(t, err)

		err = manager.DeleteDevice(device.ID)
		assert.NoError(t, err)

		// Verify device is deleted
		_, err = manager.GetDevice(device.ID)
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeNotFound, deviceErr.Type)
	})

	t.Run("non-existent device", func(t *testing.T) {
		err := manager.DeleteDevice("non-existent-id")
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeNotFound, deviceErr.Type)
	})

	t.Run("empty ID", func(t *testing.T) {
		err := manager.DeleteDevice("")
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeValidation, deviceErr.Type)
		assert.Equal(t, "id", deviceErr.Field)
	})

	t.Run("whitespace ID", func(t *testing.T) {
		err := manager.DeleteDevice("   ")
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeValidation, deviceErr.Type)
		assert.Equal(t, "id", deviceErr.Field)
	})
}

func TestManager_TestConnectivity(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	manager := NewManager(db)

	t.Run("valid device", func(t *testing.T) {
		device := createTestDevice()

		err := manager.TestConnectivity(device)
		// Since we're testing with a real IP that may not be reachable,
		// we expect either no error (if reachable) or a connectivity error
		if err != nil {
			deviceErr, ok := err.(*DeviceError)
			require.True(t, ok)
			assert.Equal(t, "connectivity", deviceErr.Type)
		}
		// The device status should be updated regardless
		assert.NotEmpty(t, device.Status)
		assert.NotNil(t, device.LastChecked)
	})

	t.Run("nil device", func(t *testing.T) {
		err := manager.TestConnectivity(nil)
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeValidation, deviceErr.Type)
	})

	t.Run("invalid device", func(t *testing.T) {
		device := createTestDevice()
		device.Name = "" // Make invalid

		err := manager.TestConnectivity(device)
		assert.Error(t, err)

		deviceErr, ok := err.(*DeviceError)
		require.True(t, ok)
		assert.Equal(t, ErrorTypeValidation, deviceErr.Type)
	})
}

// Test transaction rollback behavior
func TestManager_TransactionRollback(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	manager := NewManager(db)

	t.Run("add device transaction rollback on constraint violation", func(t *testing.T) {
		// Add a device first
		device1 := createTestDevice()
		device1.IPAddress = "192.168.1.50"
		err := manager.AddDevice(device1)
		require.NoError(t, err)

		// Count devices before failed add
		devicesBefore, err := manager.GetAllDevices()
		require.NoError(t, err)
		countBefore := len(devicesBefore)

		// Try to add device with duplicate IP (should fail and rollback)
		device2 := createTestDevice()
		device2.IPAddress = "192.168.1.50" // Same IP
		device2.Name = "Duplicate IP Device"

		err = manager.AddDevice(device2)
		assert.Error(t, err)

		// Count devices after failed add (should be same)
		devicesAfter, err := manager.GetAllDevices()
		require.NoError(t, err)
		countAfter := len(devicesAfter)

		assert.Equal(t, countBefore, countAfter, "Transaction should have been rolled back")
	})
}

// Benchmark tests
func BenchmarkManager_AddDevice(b *testing.B) {
	db := setupTestDB(&testing.T{})
	defer db.Close()
	manager := NewManager(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		device := createTestDevice()
		device.IPAddress = fmt.Sprintf("192.168.1.%d", i+1)
		device.Name = fmt.Sprintf("Device %d", i+1)

		err := manager.AddDevice(device)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkManager_GetAllDevices(b *testing.B) {
	db := setupTestDB(&testing.T{})
	defer db.Close()
	manager := NewManager(db)

	// Add some test devices
	for i := 0; i < 100; i++ {
		device := createTestDevice()
		device.IPAddress = fmt.Sprintf("192.168.1.%d", i+1)
		device.Name = fmt.Sprintf("Device %d", i+1)
		manager.AddDevice(device)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.GetAllDevices()
		if err != nil {
			b.Fatal(err)
		}
	}
}
