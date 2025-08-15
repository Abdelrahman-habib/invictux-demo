package checker

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RuleManager handles security rule operations
type RuleManager struct {
	db *sql.DB
}

// NewRuleManager creates a new rule manager
func NewRuleManager(db *sql.DB) *RuleManager {
	return &RuleManager{db: db}
}

// LoadPredefinedRules loads predefined security rules for all vendors
func (rm *RuleManager) LoadPredefinedRules() error {
	rules := GetPredefinedRules()

	for _, rule := range rules {
		// Check if rule already exists
		exists, err := rm.ruleExists(rule.Name, rule.Vendor)
		if err != nil {
			return fmt.Errorf("failed to check if rule exists: %w", err)
		}

		if !exists {
			if err := rm.CreateRule(rule); err != nil {
				return fmt.Errorf("failed to create rule %s: %w", rule.Name, err)
			}
		}
	}

	return nil
}

// CreateRule creates a new security rule
func (rm *RuleManager) CreateRule(rule SecurityRule) error {
	if rule.ID == "" {
		rule.ID = uuid.New().String()
	}

	if rule.CreatedAt.IsZero() {
		rule.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO security_rules (id, name, description, vendor, command, expected_pattern, severity, enabled, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := rm.db.Exec(query, rule.ID, rule.Name, rule.Description, rule.Vendor,
		rule.Command, rule.ExpectedPattern, rule.Severity, rule.Enabled, rule.CreatedAt)

	return err
}

// GetAllRules retrieves all security rules
func (rm *RuleManager) GetAllRules() ([]SecurityRule, error) {
	query := `
		SELECT id, name, description, vendor, command, expected_pattern, severity, enabled, created_at
		FROM security_rules
		ORDER BY vendor, name
	`

	rows, err := rm.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []SecurityRule
	for rows.Next() {
		var rule SecurityRule
		err := rows.Scan(&rule.ID, &rule.Name, &rule.Description, &rule.Vendor,
			&rule.Command, &rule.ExpectedPattern, &rule.Severity, &rule.Enabled, &rule.CreatedAt)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// GetRulesByVendor retrieves security rules for a specific vendor
func (rm *RuleManager) GetRulesByVendor(vendor string) ([]SecurityRule, error) {
	query := `
		SELECT id, name, description, vendor, command, expected_pattern, severity, enabled, created_at
		FROM security_rules
		WHERE vendor = ? OR vendor = 'generic'
		ORDER BY name
	`

	rows, err := rm.db.Query(query, vendor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []SecurityRule
	for rows.Next() {
		var rule SecurityRule
		err := rows.Scan(&rule.ID, &rule.Name, &rule.Description, &rule.Vendor,
			&rule.Command, &rule.ExpectedPattern, &rule.Severity, &rule.Enabled, &rule.CreatedAt)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// UpdateRule updates an existing security rule
func (rm *RuleManager) UpdateRule(rule SecurityRule) error {
	query := `
		UPDATE security_rules 
		SET name = ?, description = ?, vendor = ?, command = ?, expected_pattern = ?, severity = ?, enabled = ?
		WHERE id = ?
	`

	result, err := rm.db.Exec(query, rule.Name, rule.Description, rule.Vendor,
		rule.Command, rule.ExpectedPattern, rule.Severity, rule.Enabled, rule.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("rule with ID %s not found", rule.ID)
	}

	return nil
}

// DeleteRule deletes a security rule
func (rm *RuleManager) DeleteRule(id string) error {
	query := "DELETE FROM security_rules WHERE id = ?"

	result, err := rm.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("rule with ID %s not found", id)
	}

	return nil
}

// EnableRule enables a security rule
func (rm *RuleManager) EnableRule(id string) error {
	query := "UPDATE security_rules SET enabled = TRUE WHERE id = ?"

	result, err := rm.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("rule with ID %s not found", id)
	}

	return nil
}

// DisableRule disables a security rule
func (rm *RuleManager) DisableRule(id string) error {
	query := "UPDATE security_rules SET enabled = FALSE WHERE id = ?"

	result, err := rm.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("rule with ID %s not found", id)
	}

	return nil
}

// ruleExists checks if a rule with the given name and vendor already exists
func (rm *RuleManager) ruleExists(name, vendor string) (bool, error) {
	query := "SELECT COUNT(*) FROM security_rules WHERE name = ? AND vendor = ?"

	var count int
	err := rm.db.QueryRow(query, name, vendor).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetPredefinedRules returns predefined security rules for various vendors
func GetPredefinedRules() []SecurityRule {
	var rules []SecurityRule

	// Add Cisco IOS specific rules
	rules = append(rules, getCiscoIOSRules()...)

	// Add generic rules that apply to all vendors
	rules = append(rules, getGenericRules()...)

	return rules
}

// getCiscoIOSRules returns Cisco IOS specific security rules
func getCiscoIOSRules() []SecurityRule {
	return []SecurityRule{
		{
			ID:              uuid.New().String(),
			Name:            "Check Default Enable Password",
			Description:     "Verify that the default enable password is not being used",
			Vendor:          "cisco",
			Command:         "show running-config | include enable password",
			ExpectedPattern: `^$|enable password \$1\$.*|enable secret \$.*`,
			Severity:        string(SeverityCritical),
			Enabled:         true,
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New().String(),
			Name:            "Check SSH vs Telnet Configuration",
			Description:     "Ensure SSH is enabled and Telnet is disabled for secure remote access",
			Vendor:          "cisco",
			Command:         "show ip ssh",
			ExpectedPattern: `SSH Enabled - version [12]\..*`,
			Severity:        string(SeverityHigh),
			Enabled:         true,
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New().String(),
			Name:            "Check Telnet VTY Lines",
			Description:     "Verify that Telnet access is disabled on VTY lines",
			Vendor:          "cisco",
			Command:         "show running-config | section line vty",
			ExpectedPattern: `transport input ssh|transport input none`,
			Severity:        string(SeverityHigh),
			Enabled:         true,
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New().String(),
			Name:            "Check Unused Interfaces",
			Description:     "Identify interfaces that are administratively up but not in use",
			Vendor:          "cisco",
			Command:         "show interfaces status | include notconnect",
			ExpectedPattern: `.*shutdown.*|^$`,
			Severity:        string(SeverityMedium),
			Enabled:         true,
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New().String(),
			Name:            "Check Console Password",
			Description:     "Verify that console access is password protected",
			Vendor:          "cisco",
			Command:         "show running-config | section line con",
			ExpectedPattern: `password .*|login local`,
			Severity:        string(SeverityHigh),
			Enabled:         true,
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New().String(),
			Name:            "Check SNMP Community Strings",
			Description:     "Verify that default SNMP community strings are not in use",
			Vendor:          "cisco",
			Command:         "show running-config | include snmp-server community",
			ExpectedPattern: `^$|snmp-server community [^p].*|snmp-server community p[^ru].*|snmp-server community pr[^i].*|snmp-server community pri[^v].*`,
			Severity:        string(SeverityCritical),
			Enabled:         true,
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New().String(),
			Name:            "Check Service Password Encryption",
			Description:     "Ensure password encryption service is enabled",
			Vendor:          "cisco",
			Command:         "show running-config | include service password-encryption",
			ExpectedPattern: `service password-encryption`,
			Severity:        string(SeverityMedium),
			Enabled:         true,
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New().String(),
			Name:            "Check Login Banner",
			Description:     "Verify that a login banner is configured for legal compliance",
			Vendor:          "cisco",
			Command:         "show running-config | include banner",
			ExpectedPattern: `banner (login|motd)`,
			Severity:        string(SeverityLow),
			Enabled:         true,
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New().String(),
			Name:            "Check HTTP/HTTPS Server Status",
			Description:     "Verify that HTTP server is disabled and HTTPS is used if web management is needed",
			Vendor:          "cisco",
			Command:         "show running-config | include ip http",
			ExpectedPattern: `no ip http server|ip http secure-server`,
			Severity:        string(SeverityHigh),
			Enabled:         true,
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New().String(),
			Name:            "Check CDP Configuration",
			Description:     "Verify CDP is disabled on interfaces facing untrusted networks",
			Vendor:          "cisco",
			Command:         "show cdp neighbors",
			ExpectedPattern: `.*`, // This rule requires manual review of CDP neighbors
			Severity:        string(SeverityMedium),
			Enabled:         true,
			CreatedAt:       time.Now(),
		},
	}
}

// getGenericRules returns generic security rules applicable to all vendors
func getGenericRules() []SecurityRule {
	return []SecurityRule{
		{
			ID:              uuid.New().String(),
			Name:            "Check System Uptime",
			Description:     "Monitor system uptime to identify devices that may need updates",
			Vendor:          "generic",
			Command:         "show version | include uptime",
			ExpectedPattern: `.*uptime.*`,
			Severity:        string(SeverityLow),
			Enabled:         true,
			CreatedAt:       time.Now(),
		},
		{
			ID:              uuid.New().String(),
			Name:            "Check Running Configuration",
			Description:     "Verify that running configuration can be accessed",
			Vendor:          "generic",
			Command:         "show running-config | head -5",
			ExpectedPattern: `.*version.*|.*hostname.*|.*!.*`,
			Severity:        string(SeverityLow),
			Enabled:         true,
			CreatedAt:       time.Now(),
		},
	}
}
