package app

import (
	"context"
	"log"
	"time"

	"invictux-demo/internal/checker"
	"invictux-demo/internal/database"
	"invictux-demo/internal/device"
	"invictux-demo/internal/security"
)

// App struct represents the main application
type App struct {
	ctx               context.Context
	db                *database.DB
	deviceManager     *device.Manager
	checkEngine       *checker.Engine
	scanner           *device.ConnectivityScanner
	encryptionManager *security.EncryptionManager
	sessionManager    *security.SessionManager
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

	// Initialize security components
	// TODO: In production, this should be configurable or derived from user input
	a.encryptionManager = security.NewEncryptionManager("default-app-key-change-in-production")
	a.sessionManager = security.NewSessionManager(30 * time.Minute) // 30 minute session timeout

	// Initialize components
	a.deviceManager = device.NewManager(a.db.DB)

	// Initialize rule manager and load predefined rules
	ruleManager := checker.NewRuleManager(a.db.DB)
	if err := ruleManager.LoadPredefinedRules(); err != nil {
		log.Printf("Failed to load predefined rules: %v", err)
		// Continue anyway, rules can be loaded later
	}

	a.checkEngine = checker.NewEngine(ruleManager)
	a.scanner = device.NewConnectivityScanner()

	log.Println("Network Configuration Checker initialized successfully")
}

// GetEnvironment returns the current application environment (production, staging, etc.)
func (a *App) GetEnvironment() string {
	return a.environment
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
	if result, err := a.scanner.TestConnectivity(&dev); err != nil {
		log.Printf("Connectivity test failed for device %s: %v", dev.Name, err)
		// Don't fail the add operation, just log the warning
	} else if result.Error != nil {
		log.Printf("Connectivity issues for device %s: %v", dev.Name, result.Error)
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

	result, err := a.scanner.TestConnectivity(dev)
	if err != nil {
		return err
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
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

// Security and Settings Methods

// EncryptPassword encrypts a password for secure storage
func (a *App) EncryptPassword(password string) ([]byte, error) {
	if a.encryptionManager == nil {
		return nil, nil
	}
	return a.encryptionManager.Encrypt(password)
}

// DecryptPassword decrypts a stored password
func (a *App) DecryptPassword(encryptedPassword []byte) (string, error) {
	if a.encryptionManager == nil {
		return "", nil
	}
	return a.encryptionManager.Decrypt(encryptedPassword)
}

// CreateSession creates a new user session
func (a *App) CreateSession(userID string) (*security.Session, error) {
	if a.sessionManager == nil {
		return nil, nil
	}
	return a.sessionManager.CreateSession(userID)
}

// ValidateSession validates an existing session
func (a *App) ValidateSession(sessionID string) (*security.Session, error) {
	if a.sessionManager == nil {
		return nil, nil
	}
	return a.sessionManager.ValidateSession(sessionID)
}

// DestroySession destroys a user session
func (a *App) DestroySession(sessionID string) {
	if a.sessionManager != nil {
		a.sessionManager.DestroySession(sessionID)
	}
}

// GetDatabaseStats returns database statistics
func (a *App) GetDatabaseStats() map[string]interface{} {
	if a.db == nil {
		return make(map[string]interface{})
	}

	stats := a.db.GetStats()
	return map[string]interface{}{
		"maxOpenConnections": stats.MaxOpenConnections,
		"openConnections":    stats.OpenConnections,
		"inUse":              stats.InUse,
		"idle":               stats.Idle,
		"waitCount":          stats.WaitCount,
		"waitDuration":       stats.WaitDuration.String(),
		"maxIdleClosed":      stats.MaxIdleClosed,
		"maxIdleTimeClosed":  stats.MaxIdleTimeClosed,
		"maxLifetimeClosed":  stats.MaxLifetimeClosed,
	}
}

// PerformDatabaseHealthCheck performs a database health check
func (a *App) PerformDatabaseHealthCheck() error {
	if a.db == nil {
		return nil
	}
	return a.db.HealthCheck()
}

// BackupDatabase creates a backup of the database
func (a *App) BackupDatabase(backupPath string) error {
	if a.db == nil {
		return nil
	}
	return a.db.Backup(backupPath)
}
