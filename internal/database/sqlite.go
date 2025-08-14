package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps the sql.DB with additional functionality
type DB struct {
	*sql.DB
	dataDir string
}

// ConnectionConfig holds database connection configuration
type ConnectionConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultConnectionConfig returns default connection configuration
func DefaultConnectionConfig() *ConnectionConfig {
	return &ConnectionConfig{
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}
}

// NewSQLiteDB creates a new SQLite database connection with proper configuration
func NewSQLiteDB(dataDir string) (*DB, error) {
	return NewSQLiteDBWithConfig(dataDir, DefaultConnectionConfig())
}

// NewSQLiteDBWithConfig creates a new SQLite database connection with custom configuration
func NewSQLiteDBWithConfig(dataDir string, config *ConnectionConfig) (*DB, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "network_checker.db")

	// SQLite connection string with optimizations
	connectionString := fmt.Sprintf("%s?_journal_mode=WAL&_synchronous=NORMAL&_cache_size=1000&_foreign_keys=ON", dbPath)

	db, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set additional SQLite pragmas for performance and reliability
	pragmas := []string{
		"PRAGMA busy_timeout = 30000",  // 30 second timeout for busy database
		"PRAGMA temp_store = MEMORY",   // Store temporary tables in memory
		"PRAGMA mmap_size = 268435456", // 256MB memory-mapped I/O
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to set pragma %s: %w", pragma, err)
		}
	}

	return &DB{
		DB:      db,
		dataDir: dataDir,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// GetDataDir returns the data directory path
func (db *DB) GetDataDir() string {
	return db.dataDir
}

// HealthCheck performs a database health check
func (db *DB) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Test a simple query
	var result int
	err := db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("database query test failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("database query returned unexpected result: %d", result)
	}

	return nil
}

// GetStats returns database statistics
func (db *DB) GetStats() sql.DBStats {
	return db.DB.Stats()
}

// Backup creates a backup of the database
func (db *DB) Backup(backupPath string) error {
	// Ensure backup directory exists
	backupDir := filepath.Dir(backupPath)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Use SQLite's VACUUM INTO command for backup
	query := "VACUUM INTO ?"
	if _, err := db.Exec(query, backupPath); err != nil {
		return fmt.Errorf("failed to backup database: %w", err)
	}

	return nil
}

// GetDefaultDataDir returns the default data directory
func GetDefaultDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".network-config-checker"), nil
}

// GetDataDir returns the default data directory (deprecated, use GetDefaultDataDir)
func GetDataDir() (string, error) {
	return GetDefaultDataDir()
}
