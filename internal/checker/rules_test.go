package checker

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create security_rules table
	createTableSQL := `
		CREATE TABLE security_rules (
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
	`

	if _, err := db.Exec(createTableSQL); err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	return db
}

func TestRuleManager_CreateRule(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := NewRuleManager(db)

	rule := SecurityRule{
		Name:            "Test Rule",
		Description:     "Test Description",
		Vendor:          "cisco",
		Command:         "show version",
		ExpectedPattern: ".*IOS.*",
		Severity:        string(SeverityHigh),
		Enabled:         true,
	}

	err := rm.CreateRule(rule)
	if err != nil {
		t.Fatalf("Failed to create rule: %v", err)
	}

	// Verify rule was created
	rules, err := rm.GetAllRules()
	if err != nil {
		t.Fatalf("Failed to get rules: %v", err)
	}

	if len(rules) != 1 {
		t.Fatalf("Expected 1 rule, got %d", len(rules))
	}

	createdRule := rules[0]
	if createdRule.Name != rule.Name {
		t.Errorf("Expected name %s, got %s", rule.Name, createdRule.Name)
	}
	if createdRule.Vendor != rule.Vendor {
		t.Errorf("Expected vendor %s, got %s", rule.Vendor, createdRule.Vendor)
	}
	if createdRule.ID == "" {
		t.Error("Expected rule ID to be generated")
	}
}

func TestRuleManager_GetRulesByVendor(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := NewRuleManager(db)

	// Create rules for different vendors
	ciscoRule := SecurityRule{
		ID:              uuid.New().String(),
		Name:            "Cisco Rule",
		Vendor:          "cisco",
		Command:         "show version",
		ExpectedPattern: ".*IOS.*",
		Severity:        string(SeverityHigh),
		Enabled:         true,
		CreatedAt:       time.Now(),
	}

	genericRule := SecurityRule{
		ID:              uuid.New().String(),
		Name:            "Generic Rule",
		Vendor:          "generic",
		Command:         "show system",
		ExpectedPattern: ".*",
		Severity:        string(SeverityMedium),
		Enabled:         true,
		CreatedAt:       time.Now(),
	}

	juniperRule := SecurityRule{
		ID:              uuid.New().String(),
		Name:            "Juniper Rule",
		Vendor:          "juniper",
		Command:         "show version",
		ExpectedPattern: ".*JUNOS.*",
		Severity:        string(SeverityHigh),
		Enabled:         true,
		CreatedAt:       time.Now(),
	}

	// Create all rules
	if err := rm.CreateRule(ciscoRule); err != nil {
		t.Fatalf("Failed to create Cisco rule: %v", err)
	}
	if err := rm.CreateRule(genericRule); err != nil {
		t.Fatalf("Failed to create generic rule: %v", err)
	}
	if err := rm.CreateRule(juniperRule); err != nil {
		t.Fatalf("Failed to create Juniper rule: %v", err)
	}

	// Test getting Cisco rules (should include generic rules too)
	ciscoRules, err := rm.GetRulesByVendor("cisco")
	if err != nil {
		t.Fatalf("Failed to get Cisco rules: %v", err)
	}

	if len(ciscoRules) != 2 {
		t.Fatalf("Expected 2 rules for Cisco (including generic), got %d", len(ciscoRules))
	}

	// Verify we got the right rules
	foundCisco := false
	foundGeneric := false
	for _, rule := range ciscoRules {
		if rule.Name == "Cisco Rule" {
			foundCisco = true
		}
		if rule.Name == "Generic Rule" {
			foundGeneric = true
		}
	}

	if !foundCisco {
		t.Error("Expected to find Cisco rule")
	}
	if !foundGeneric {
		t.Error("Expected to find generic rule")
	}

	// Test getting Juniper rules
	juniperRules, err := rm.GetRulesByVendor("juniper")
	if err != nil {
		t.Fatalf("Failed to get Juniper rules: %v", err)
	}

	if len(juniperRules) != 2 {
		t.Fatalf("Expected 2 rules for Juniper (including generic), got %d", len(juniperRules))
	}
}

func TestRuleManager_UpdateRule(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := NewRuleManager(db)

	// Create initial rule
	rule := SecurityRule{
		ID:              uuid.New().String(),
		Name:            "Original Rule",
		Description:     "Original Description",
		Vendor:          "cisco",
		Command:         "show version",
		ExpectedPattern: ".*IOS.*",
		Severity:        string(SeverityHigh),
		Enabled:         true,
		CreatedAt:       time.Now(),
	}

	if err := rm.CreateRule(rule); err != nil {
		t.Fatalf("Failed to create rule: %v", err)
	}

	// Update the rule
	rule.Name = "Updated Rule"
	rule.Description = "Updated Description"
	rule.Severity = string(SeverityCritical)

	if err := rm.UpdateRule(rule); err != nil {
		t.Fatalf("Failed to update rule: %v", err)
	}

	// Verify update
	rules, err := rm.GetAllRules()
	if err != nil {
		t.Fatalf("Failed to get rules: %v", err)
	}

	if len(rules) != 1 {
		t.Fatalf("Expected 1 rule, got %d", len(rules))
	}

	updatedRule := rules[0]
	if updatedRule.Name != "Updated Rule" {
		t.Errorf("Expected name 'Updated Rule', got %s", updatedRule.Name)
	}
	if updatedRule.Description != "Updated Description" {
		t.Errorf("Expected description 'Updated Description', got %s", updatedRule.Description)
	}
	if updatedRule.Severity != string(SeverityCritical) {
		t.Errorf("Expected severity %s, got %s", string(SeverityCritical), updatedRule.Severity)
	}
}

func TestRuleManager_DeleteRule(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := NewRuleManager(db)

	// Create rule
	rule := SecurityRule{
		ID:              uuid.New().String(),
		Name:            "Test Rule",
		Vendor:          "cisco",
		Command:         "show version",
		ExpectedPattern: ".*IOS.*",
		Severity:        string(SeverityHigh),
		Enabled:         true,
		CreatedAt:       time.Now(),
	}

	if err := rm.CreateRule(rule); err != nil {
		t.Fatalf("Failed to create rule: %v", err)
	}

	// Delete the rule
	if err := rm.DeleteRule(rule.ID); err != nil {
		t.Fatalf("Failed to delete rule: %v", err)
	}

	// Verify deletion
	rules, err := rm.GetAllRules()
	if err != nil {
		t.Fatalf("Failed to get rules: %v", err)
	}

	if len(rules) != 0 {
		t.Fatalf("Expected 0 rules after deletion, got %d", len(rules))
	}
}

func TestRuleManager_EnableDisableRule(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := NewRuleManager(db)

	// Create rule
	rule := SecurityRule{
		ID:              uuid.New().String(),
		Name:            "Test Rule",
		Vendor:          "cisco",
		Command:         "show version",
		ExpectedPattern: ".*IOS.*",
		Severity:        string(SeverityHigh),
		Enabled:         true,
		CreatedAt:       time.Now(),
	}

	if err := rm.CreateRule(rule); err != nil {
		t.Fatalf("Failed to create rule: %v", err)
	}

	// Disable the rule
	if err := rm.DisableRule(rule.ID); err != nil {
		t.Fatalf("Failed to disable rule: %v", err)
	}

	// Verify rule is disabled
	rules, err := rm.GetAllRules()
	if err != nil {
		t.Fatalf("Failed to get rules: %v", err)
	}

	if len(rules) != 1 {
		t.Fatalf("Expected 1 rule, got %d", len(rules))
	}

	if rules[0].Enabled {
		t.Error("Expected rule to be disabled")
	}

	// Enable the rule
	if err := rm.EnableRule(rule.ID); err != nil {
		t.Fatalf("Failed to enable rule: %v", err)
	}

	// Verify rule is enabled
	rules, err = rm.GetAllRules()
	if err != nil {
		t.Fatalf("Failed to get rules: %v", err)
	}

	if !rules[0].Enabled {
		t.Error("Expected rule to be enabled")
	}
}

func TestRuleManager_LoadPredefinedRules(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	rm := NewRuleManager(db)

	// Load predefined rules
	if err := rm.LoadPredefinedRules(); err != nil {
		t.Fatalf("Failed to load predefined rules: %v", err)
	}

	// Verify rules were loaded
	rules, err := rm.GetAllRules()
	if err != nil {
		t.Fatalf("Failed to get rules: %v", err)
	}

	if len(rules) == 0 {
		t.Fatal("Expected predefined rules to be loaded")
	}

	// Verify we have Cisco rules
	ciscoRules, err := rm.GetRulesByVendor("cisco")
	if err != nil {
		t.Fatalf("Failed to get Cisco rules: %v", err)
	}

	if len(ciscoRules) == 0 {
		t.Fatal("Expected Cisco rules to be loaded")
	}

	// Verify we have generic rules
	genericRules, err := rm.GetRulesByVendor("generic")
	if err != nil {
		t.Fatalf("Failed to get generic rules: %v", err)
	}

	if len(genericRules) == 0 {
		t.Fatal("Expected generic rules to be loaded")
	}

	// Test loading again (should not create duplicates)
	if err := rm.LoadPredefinedRules(); err != nil {
		t.Fatalf("Failed to load predefined rules second time: %v", err)
	}

	// Verify no duplicates were created
	rulesAfterSecondLoad, err := rm.GetAllRules()
	if err != nil {
		t.Fatalf("Failed to get rules after second load: %v", err)
	}

	if len(rulesAfterSecondLoad) != len(rules) {
		t.Errorf("Expected same number of rules after second load, got %d vs %d",
			len(rulesAfterSecondLoad), len(rules))
	}
}

func TestGetPredefinedRules(t *testing.T) {
	rules := GetPredefinedRules()

	if len(rules) == 0 {
		t.Fatal("Expected predefined rules to be returned")
	}

	// Verify we have Cisco rules
	foundCisco := false
	foundGeneric := false

	for _, rule := range rules {
		if rule.Vendor == "cisco" {
			foundCisco = true
		}
		if rule.Vendor == "generic" {
			foundGeneric = true
		}

		// Verify required fields are set
		if rule.Name == "" {
			t.Error("Rule name should not be empty")
		}
		if rule.Command == "" {
			t.Error("Rule command should not be empty")
		}
		if rule.Severity == "" {
			t.Error("Rule severity should not be empty")
		}
		if rule.ID == "" {
			t.Error("Rule ID should not be empty")
		}
	}

	if !foundCisco {
		t.Error("Expected to find Cisco rules")
	}
	if !foundGeneric {
		t.Error("Expected to find generic rules")
	}
}

func TestGetCiscoIOSRules(t *testing.T) {
	rules := getCiscoIOSRules()

	if len(rules) == 0 {
		t.Fatal("Expected Cisco IOS rules to be returned")
	}

	// Verify specific Cisco rules exist
	expectedRules := map[string]bool{
		"Check Default Enable Password":     false,
		"Check SSH vs Telnet Configuration": false,
		"Check Telnet VTY Lines":            false,
		"Check Unused Interfaces":           false,
		"Check Console Password":            false,
		"Check SNMP Community Strings":      false,
		"Check Service Password Encryption": false,
		"Check Login Banner":                false,
		"Check HTTP/HTTPS Server Status":    false,
		"Check CDP Configuration":           false,
	}

	for _, rule := range rules {
		if rule.Vendor != "cisco" {
			t.Errorf("Expected vendor 'cisco', got %s", rule.Vendor)
		}

		if _, exists := expectedRules[rule.Name]; exists {
			expectedRules[rule.Name] = true
		}

		// Verify rule has required fields
		if rule.Command == "" {
			t.Errorf("Rule %s should have a command", rule.Name)
		}
		if rule.ExpectedPattern == "" {
			t.Errorf("Rule %s should have an expected pattern", rule.Name)
		}
		if rule.Severity == "" {
			t.Errorf("Rule %s should have a severity", rule.Name)
		}
	}

	// Verify all expected rules were found
	for ruleName, found := range expectedRules {
		if !found {
			t.Errorf("Expected rule %s not found", ruleName)
		}
	}
}

func TestGetGenericRules(t *testing.T) {
	rules := getGenericRules()

	if len(rules) == 0 {
		t.Fatal("Expected generic rules to be returned")
	}

	for _, rule := range rules {
		if rule.Vendor != "generic" {
			t.Errorf("Expected vendor 'generic', got %s", rule.Vendor)
		}

		// Verify rule has required fields
		if rule.Name == "" {
			t.Error("Rule name should not be empty")
		}
		if rule.Command == "" {
			t.Error("Rule command should not be empty")
		}
		if rule.ExpectedPattern == "" {
			t.Error("Rule expected pattern should not be empty")
		}
		if rule.Severity == "" {
			t.Error("Rule severity should not be empty")
		}
	}
}
