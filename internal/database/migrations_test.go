package database

import (
	"database/sql"
	"os"
	"testing"
)

func TestGetMigrations(t *testing.T) {
	migrations := GetMigrations()

	if len(migrations) == 0 {
		t.Fatal("Expected at least one migration")
	}

	// Check that migrations are properly structured
	for i, migration := range migrations {
		if migration.Version <= 0 {
			t.Errorf("Migration %d has invalid version: %d", i, migration.Version)
		}

		if migration.Name == "" {
			t.Errorf("Migration %d has empty name", i)
		}

		if migration.SQL == "" {
			t.Errorf("Migration %d has empty SQL", i)
		}
	}

	// Check that versions are unique
	versions := make(map[int]bool)
	for _, migration := range migrations {
		if versions[migration.Version] {
			t.Errorf("Duplicate migration version: %d", migration.Version)
		}
		versions[migration.Version] = true
	}

	// Check that versions are sequential (starting from 1)
	for i := 1; i <= len(migrations); i++ {
		if !versions[i] {
			t.Errorf("Missing migration version: %d", i)
		}
	}
}

func TestRunMigrations(t *testing.T) {
	// Create temporary database
	tempDir, err := os.MkdirTemp("", "test_migrations_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db, err := NewSQLiteDB(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Run migrations
	err = RunMigrations(db.DB)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Verify that all expected tables exist
	expectedTables := []string{
		"devices",
		"check_results",
		"security_rules",
		"app_settings",
		"schema_migrations",
	}

	for _, tableName := range expectedTables {
		var count int
		query := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?"
		err = db.QueryRow(query, tableName).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to check for table %s: %v", tableName, err)
		}

		if count != 1 {
			t.Errorf("Expected table %s to exist", tableName)
		}
	}

	// Verify that all migrations are recorded
	migrations := GetMigrations()
	for _, migration := range migrations {
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", migration.Version).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to check migration record for version %d: %v", migration.Version, err)
		}

		if count != 1 {
			t.Errorf("Expected migration version %d to be recorded", migration.Version)
		}
	}
}

func TestRunMigrationsIdempotent(t *testing.T) {
	// Create temporary database
	tempDir, err := os.MkdirTemp("", "test_migrations_idempotent_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db, err := NewSQLiteDB(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Run migrations first time
	err = RunMigrations(db.DB)
	if err != nil {
		t.Fatalf("Failed to run migrations first time: %v", err)
	}

	// Count migrations applied
	var firstCount int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&firstCount)
	if err != nil {
		t.Fatalf("Failed to count migrations: %v", err)
	}

	// Run migrations second time (should be idempotent)
	err = RunMigrations(db.DB)
	if err != nil {
		t.Fatalf("Failed to run migrations second time: %v", err)
	}

	// Count migrations applied again
	var secondCount int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&secondCount)
	if err != nil {
		t.Fatalf("Failed to count migrations second time: %v", err)
	}

	if firstCount != secondCount {
		t.Errorf("Expected same number of migrations after second run. First: %d, Second: %d", firstCount, secondCount)
	}
}

func TestGetAppliedMigrations(t *testing.T) {
	// Create temporary database
	tempDir, err := os.MkdirTemp("", "test_applied_migrations_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db, err := NewSQLiteDB(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create migrations table manually
	_, err = db.Exec(`CREATE TABLE schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		t.Fatalf("Failed to create migrations table: %v", err)
	}

	// Insert some test migration records
	testVersions := []int{1, 3, 5}
	for _, version := range testVersions {
		_, err = db.Exec("INSERT INTO schema_migrations (version, name) VALUES (?, ?)",
			version, "test_migration_"+string(rune(version)))
		if err != nil {
			t.Fatalf("Failed to insert test migration: %v", err)
		}
	}

	// Test getAppliedMigrations
	appliedVersions, err := getAppliedMigrations(db.DB)
	if err != nil {
		t.Fatalf("Failed to get applied migrations: %v", err)
	}

	if len(appliedVersions) != len(testVersions) {
		t.Errorf("Expected %d applied migrations, got %d", len(testVersions), len(appliedVersions))
	}

	// Check that all test versions are present
	for _, expectedVersion := range testVersions {
		found := false
		for _, actualVersion := range appliedVersions {
			if actualVersion == expectedVersion {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected version %d to be in applied migrations", expectedVersion)
		}
	}
}

func TestRunMigration(t *testing.T) {
	// Create temporary database
	tempDir, err := os.MkdirTemp("", "test_run_migration_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db, err := NewSQLiteDB(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create migrations table
	_, err = db.Exec(`CREATE TABLE schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		t.Fatalf("Failed to create migrations table: %v", err)
	}

	// Test migration
	testMigration := Migration{
		Version: 999,
		Name:    "test_migration",
		SQL:     "CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT)",
	}

	// Run the migration
	err = runMigration(db.DB, testMigration)
	if err != nil {
		t.Fatalf("Failed to run test migration: %v", err)
	}

	// Verify table was created
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='test_table'").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to check for test table: %v", err)
	}

	if count != 1 {
		t.Error("Expected test table to be created")
	}

	// Verify migration was recorded
	var recordedName string
	err = db.QueryRow("SELECT name FROM schema_migrations WHERE version = ?", testMigration.Version).Scan(&recordedName)
	if err != nil {
		t.Fatalf("Failed to check migration record: %v", err)
	}

	if recordedName != testMigration.Name {
		t.Errorf("Expected migration name %s, got %s", testMigration.Name, recordedName)
	}
}

func TestRunMigrationRollback(t *testing.T) {
	// Create temporary database
	tempDir, err := os.MkdirTemp("", "test_migration_rollback_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db, err := NewSQLiteDB(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create migrations table
	_, err = db.Exec(`CREATE TABLE schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		t.Fatalf("Failed to create migrations table: %v", err)
	}

	// Test migration with invalid SQL (should rollback)
	invalidMigration := Migration{
		Version: 998,
		Name:    "invalid_migration",
		SQL:     "CREATE TABLE invalid_table (invalid_syntax",
	}

	// Run the invalid migration (should fail)
	err = runMigration(db.DB, invalidMigration)
	if err == nil {
		t.Fatal("Expected invalid migration to fail")
	}

	// Verify migration was not recorded (due to rollback)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", invalidMigration.Version).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to check migration record: %v", err)
	}

	if count != 0 {
		t.Error("Expected invalid migration to not be recorded due to rollback")
	}
}

func TestContains(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []int
		value    int
		expected bool
	}{
		{"empty slice", []int{}, 1, false},
		{"value present", []int{1, 2, 3}, 2, true},
		{"value not present", []int{1, 2, 3}, 4, false},
		{"single element present", []int{5}, 5, true},
		{"single element not present", []int{5}, 6, false},
		{"duplicate values", []int{1, 2, 2, 3}, 2, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := contains(tc.slice, tc.value)
			if result != tc.expected {
				t.Errorf("Expected %v for contains(%v, %d), got %v", tc.expected, tc.slice, tc.value, result)
			}
		})
	}
}

func TestMigrationTableStructure(t *testing.T) {
	// Create temporary database
	tempDir, err := os.MkdirTemp("", "test_migration_structure_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	db, err := NewSQLiteDB(tempDir)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Run migrations
	err = RunMigrations(db.DB)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Test devices table structure
	rows, err := db.Query("PRAGMA table_info(devices)")
	if err != nil {
		t.Fatalf("Failed to get devices table info: %v", err)
	}
	defer rows.Close()

	deviceColumns := make(map[string]bool)
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, pk int
		var defaultValue sql.NullString

		err = rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
		if err != nil {
			t.Fatalf("Failed to scan column info: %v", err)
		}
		deviceColumns[name] = true
	}

	expectedDeviceColumns := []string{
		"id", "name", "ip_address", "device_type", "vendor",
		"username", "password_encrypted", "ssh_port", "snmp_community",
		"tags", "created_at", "updated_at",
	}

	for _, expectedCol := range expectedDeviceColumns {
		if !deviceColumns[expectedCol] {
			t.Errorf("Expected column %s in devices table", expectedCol)
		}
	}
}

func BenchmarkRunMigrations(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tempDir, err := os.MkdirTemp("", "bench_migrations_*")
		if err != nil {
			b.Fatalf("Failed to create temp dir: %v", err)
		}

		db, err := NewSQLiteDB(tempDir)
		if err != nil {
			b.Fatalf("Failed to create database: %v", err)
		}

		err = RunMigrations(db.DB)
		if err != nil {
			b.Fatalf("Failed to run migrations: %v", err)
		}

		db.Close()
		os.RemoveAll(tempDir)
	}
}
