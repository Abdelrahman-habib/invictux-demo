package app

import (
	"context"
	"log"

	"qwin/internal/checker"
	"qwin/internal/database"
	"qwin/internal/device"
)

// App struct represents the main application
type App struct {
	ctx           context.Context
	db            *database.DB
	deviceManager *device.Manager
	checkEngine   *checker.Engine
	scanner       *device.Scanner
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called at application startup
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	dataDir, err := database.GetDataDir()
	if err != nil {
		log.Printf("Failed to get data directory: %v", err)
		return
	}

	a.db, err = database.NewSQLiteDB(dataDir)
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
		return
	}

	// Run database migrations
	if err := database.RunMigrations(a.db.DB); err != nil {
		log.Printf("Failed to run migrations: %v", err)
		return
	}

	// Initialize components
	a.deviceManager = device.NewManager(a.db.DB)
	a.checkEngine = checker.NewEngine()
	a.scanner = device.NewScanner()

	log.Println("Network Configuration Checker initialized successfully")
}

// DomReady is called after front-end resources have been loaded
func (a *App) DomReady(ctx context.Context) {
	// Add your action here
}

// BeforeClose is called when the application is about to quit
func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	return false
}

// Shutdown is called at application termination
func (a *App) Shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
	log.Println("Network Configuration Checker shutdown complete")
}

// Device Management Methods

// GetDevices returns all network devices
func (a *App) GetDevices() ([]device.Device, error) {
	if a.deviceManager == nil {
		return []device.Device{}, nil
	}
	return a.deviceManager.GetAllDevices()
}

// AddDevice adds a new network device
func (a *App) AddDevice(dev device.Device) error {
	if a.deviceManager == nil {
		return nil
	}

	// Test connectivity before adding
	if err := a.scanner.TestConnectivity(&dev); err != nil {
		log.Printf("Connectivity test failed for device %s: %v", dev.Name, err)
		// Don't fail the add operation, just log the warning
	}

	return a.deviceManager.AddDevice(&dev)
}

// UpdateDevice updates an existing device
func (a *App) UpdateDevice(dev device.Device) error {
	if a.deviceManager == nil {
		return nil
	}
	return a.deviceManager.UpdateDevice(&dev)
}

// DeleteDevice removes a device
func (a *App) DeleteDevice(deviceID string) error {
	if a.deviceManager == nil {
		return nil
	}
	return a.deviceManager.DeleteDevice(deviceID)
}

// TestDeviceConnectivity tests if a device is reachable
func (a *App) TestDeviceConnectivity(deviceID string) error {
	if a.deviceManager == nil || a.scanner == nil {
		return nil
	}

	dev, err := a.deviceManager.GetDevice(deviceID)
	if err != nil {
		return err
	}

	return a.scanner.TestConnectivity(dev)
}

// Security Check Methods

// RunSecurityCheck runs security checks on a device
func (a *App) RunSecurityCheck(deviceID string) ([]checker.CheckResult, error) {
	if a.deviceManager == nil || a.checkEngine == nil {
		return []checker.CheckResult{}, nil
	}

	dev, err := a.deviceManager.GetDevice(deviceID)
	if err != nil {
		return nil, err
	}

	return a.checkEngine.RunChecks(dev)
}

// RunBulkSecurityChecks runs security checks on all devices
func (a *App) RunBulkSecurityChecks() (map[string][]checker.CheckResult, error) {
	if a.deviceManager == nil || a.checkEngine == nil {
		return make(map[string][]checker.CheckResult), nil
	}

	devices, err := a.deviceManager.GetAllDevices()
	if err != nil {
		return nil, err
	}

	return a.checkEngine.RunBulkChecks(devices)
}
