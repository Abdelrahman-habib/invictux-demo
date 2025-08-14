package database

import (
	"database/sql"
	"fmt"
)

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	SQL     string
}

// GetMigrations returns all database migrations
func GetMigrations() []Migration {
	return []Migration{
		{
			Version: 1,
			Name:    "create_devices_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS devices (
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
			`,
		},
		{
			Version: 2,
			Name:    "create_check_results_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS check_results (
					id TEXT PRIMARY KEY,
					device_id TEXT NOT NULL,
					check_name TEXT NOT NULL,
					check_type TEXT NOT NULL,
					severity TEXT NOT NULL,
					status TEXT NOT NULL,
					message TEXT,
					evidence TEXT,
					checked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
				);
			`,
		},
		{
			Version: 3,
			Name:    "create_security_rules_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS security_rules (
					id TEXT PRIMARY KEY,
					name TEXT NOT NULL,
					description TEXT,
					vendor TEXT NOT NULL,
					command TEXT NOT NULL,
					expected_pattern TEXT,
					severity TEXT NOT NULL,
					enabled BOOLEAN DEFAULT TRUE,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);
			`,
		},
		{
			Version: 4,
			Name:    "create_app_settings_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS app_settings (
					key TEXT PRIMARY KEY,
					value TEXT NOT NULL,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);
			`,
		},
		{
			Version: 5,
			Name:    "create_schema_migrations_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS schema_migrations (
					version INTEGER PRIMARY KEY,
					name TEXT NOT NULL,
					applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);
			`,
		},
	}
}

// RunMigrations executes all pending migrations
func RunMigrations(db *sql.DB) error {
	// First, ensure the migrations table exists
	migrations := GetMigrations()

	// Create the migrations table first
	for _, migration := range migrations {
		if migration.Name == "create_schema_migrations_table" {
			if _, err := db.Exec(migration.SQL); err != nil {
				return fmt.Errorf("failed to create migrations table: %w", err)
			}
			break
		}
	}

	// Get applied migrations
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Run pending migrations
	for _, migration := range migrations {
		if migration.Name == "create_schema_migrations_table" {
			continue // Already applied above
		}

		if !contains(appliedMigrations, migration.Version) {
			if err := runMigration(db, migration); err != nil {
				return fmt.Errorf("failed to run migration %s: %w", migration.Name, err)
			}
		}
	}

	return nil
}

// getAppliedMigrations returns a list of applied migration versions
func getAppliedMigrations(db *sql.DB) ([]int, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []int
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, nil
}

// runMigration executes a single migration
func runMigration(db *sql.DB, migration Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute the migration SQL
	if _, err := tx.Exec(migration.SQL); err != nil {
		return err
	}

	// Record the migration as applied
	if _, err := tx.Exec("INSERT INTO schema_migrations (version, name) VALUES (?, ?)",
		migration.Version, migration.Name); err != nil {
		return err
	}

	return tx.Commit()
}

// contains checks if a slice contains a value
func contains(slice []int, value int) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
