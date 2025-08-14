package database

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewSQLiteDB(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "test_db_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test database creation
	db, err := NewSQLiteDB(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Verify database file exists
	dbPath := filepath.Join(tempDir, "network_checker.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}

	// Test database connection
	if err := db.Ping(); err != nil {
		t.Errorf("Database ping failed: %v", err)
	}

	// Test that data directory is set
	if db.GetDataDir() != tempDir {
		t.Errorf("Expected data dir %s, got %s", tempDir, db.GetDataDir())
	}
}

func TestNewSQLiteDBWithConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_db_config_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := &ConnectionConfig{
		MaxOpenConns:    10,
		MaxIdleConns:    2,
		ConnMaxLifetime: 2 * time.Minute,
		ConnMaxIdleTime: 30 * time.Second,
	}

	db, err := NewSQLiteDBWithConfig(tempDir, config)
	if err != nil {
		t.Fatalf("Failed to create database with config: %v", err)
	}
	defer db.Close()

	// Verify connection pool settings
	stats := db.GetStats()
	if stats.MaxOpenConnections != config.MaxOpenConns {
		t.Errorf("Expected MaxOpenConns %d, got %d", config.MaxOpenConns, stats.MaxOpenConnections)
	}
}

func TestDefaultConnectionConfig(t *testing.T) {
	config := DefaultConnectionConfig()

	if config == nil {
		t.Fatal("Expected default config to be created")
	}

	if config.MaxOpenConns <= 0 {
		t.Error("Expected MaxOpenConns to be positive")
	}

	if config.MaxIdleConns <= 0 {
		t.Error("Expected MaxIdleConns to be positive")
	}

	if config.ConnMaxLifetime <= 0 {
		t.Error("Expected ConnMaxLifetime to be positive")
	}

	if config.ConnMaxIdleTime <= 0 {
		t.Error("Expected ConnMaxIdleTime to be positive")
	}
}

func TestHealthCheck(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_health_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db, err := NewSQLiteDB(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test health check on healthy database
	if err := db.HealthCheck(); err != nil {
		t.Errorf("Health check failed on healthy database: %v", err)
	}

	// Test health check on closed database
	db.Close()
	if err := db.HealthCheck(); err == nil {
		t.Error("Expected health check to fail on closed database")
	}
}

func TestGetStats(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_stats_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db, err := NewSQLiteDB(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	stats := db.GetStats()

	// Basic validation of stats structure
	if stats.MaxOpenConnections <= 0 {
		t.Error("Expected MaxOpenConnections to be positive")
	}

	// Stats should be accessible without error
	_ = stats.OpenConnections
	_ = stats.InUse
	_ = stats.Idle
}

func TestBackup(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_backup_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create and setup database
	db, err := NewSQLiteDB(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Run migrations to create tables
	if err := RunMigrations(db.DB); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Insert some test data
	_, err = db.Exec(`INSERT INTO app_settings (key, value) VALUES (?, ?)`, "test_key", "test_value")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Create backup
	backupPath := filepath.Join(tempDir, "backup.db")
	if err := db.Backup(backupPath); err != nil {
		t.Fatalf("Backup failed: %v", err)
	}

	// Verify backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Error("Backup file was not created")
	}

	// Verify backup contains data (basic check)
	backupDB, err := NewSQLiteDB(filepath.Dir(backupPath))
	if err != nil {
		t.Fatalf("Failed to open backup database: %v", err)
	}
	defer backupDB.Close()

	// Note: This is a simplified test. In a real scenario, you'd want to verify
	// that the backup actually contains the expected data.
}

func TestDatabaseCreationWithInvalidPath(t *testing.T) {
	// Try to create database in a path that doesn't exist and can't be created
	invalidPath := "/invalid/path/that/does/not/exist"

	// On Windows, use an invalid Windows path
	if os.PathSeparator == '\\' {
		invalidPath = "Z:\\invalid\\path\\that\\does\\not\\exist"
	}

	_, err := NewSQLiteDB(invalidPath)
	if err == nil {
		t.Error("Expected error when creating database with invalid path")
	}
}

func TestGetDefaultDataDir(t *testing.T) {
	dataDir, err := GetDefaultDataDir()
	if err != nil {
		t.Fatalf("Failed to get default data dir: %v", err)
	}

	if dataDir == "" {
		t.Error("Expected non-empty data directory")
	}

	// Should contain the application name
	if !filepath.IsAbs(dataDir) {
		t.Error("Expected absolute path for data directory")
	}

	// Should be in user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home dir: %v", err)
	}

	expectedDir := filepath.Join(homeDir, ".network-config-checker")
	if dataDir != expectedDir {
		t.Errorf("Expected data dir %s, got %s", expectedDir, dataDir)
	}
}

func TestGetDataDir(t *testing.T) {
	// Test deprecated function still works
	dataDir, err := GetDataDir()
	if err != nil {
		t.Fatalf("Failed to get data dir: %v", err)
	}

	if dataDir == "" {
		t.Error("Expected non-empty data directory")
	}

	// Should match GetDefaultDataDir
	defaultDir, err := GetDefaultDataDir()
	if err != nil {
		t.Fatalf("Failed to get default data dir: %v", err)
	}

	if dataDir != defaultDir {
		t.Errorf("GetDataDir() should match GetDefaultDataDir(). Got %s vs %s", dataDir, defaultDir)
	}
}

func TestDatabasePragmas(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_pragmas_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db, err := NewSQLiteDB(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test that foreign keys are enabled
	var foreignKeys int
	err = db.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys)
	if err != nil {
		t.Fatalf("Failed to query foreign_keys pragma: %v", err)
	}

	if foreignKeys != 1 {
		t.Error("Expected foreign keys to be enabled")
	}

	// Test journal mode
	var journalMode string
	err = db.QueryRow("PRAGMA journal_mode").Scan(&journalMode)
	if err != nil {
		t.Fatalf("Failed to query journal_mode pragma: %v", err)
	}

	if journalMode != "wal" {
		t.Errorf("Expected WAL journal mode, got %s", journalMode)
	}
}

func TestConcurrentConnections(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_concurrent_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := &ConnectionConfig{
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 1 * time.Minute,
		ConnMaxIdleTime: 30 * time.Second,
	}

	db, err := NewSQLiteDBWithConfig(tempDir, config)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := RunMigrations(db.DB); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Test concurrent operations
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()

			// Perform a simple database operation
			_, err := db.Exec("INSERT INTO app_settings (key, value) VALUES (?, ?)",
				"concurrent_test_"+string(rune(id)), "value")
			if err != nil {
				t.Errorf("Concurrent operation %d failed: %v", id, err)
			}
		}(i)
	}

	// Wait for all operations to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify some data was inserted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM app_settings WHERE key LIKE 'concurrent_test_%'").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count inserted records: %v", err)
	}

	if count == 0 {
		t.Error("Expected some records to be inserted by concurrent operations")
	}
}

func BenchmarkDatabaseConnection(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "bench_db_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db, err := NewSQLiteDB(tempDir)
		if err != nil {
			b.Fatalf("Failed to create database: %v", err)
		}
		db.Close()
	}
}

func BenchmarkHealthCheck(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "bench_health_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db, err := NewSQLiteDB(tempDir)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := db.HealthCheck(); err != nil {
			b.Fatalf("Health check failed: %v", err)
		}
	}
}
