package checker

import (
	"database/sql"
	"regexp"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestEngine_EvaluateRuleResult(t *testing.T) {
	rm := setupTestRuleManager(t)
	engine := NewEngine(rm)

	tests := []struct {
		name           string
		output         string
		rule           SecurityRule
		expectedStatus CheckStatus
		expectedMsg    string
	}{
		{
			name:   "Cisco Enable Password - Secure Configuration",
			output: "enable secret $1$abcd$xyz123",
			rule: SecurityRule{
				Name:            "Check Default Enable Password",
				ExpectedPattern: `^$|enable password \$1\$.*|enable secret \$.*`,
			},
			expectedStatus: StatusPass,
			expectedMsg:    "Configuration check passed",
		},
		{
			name:   "Cisco Enable Password - Insecure Configuration",
			output: "enable password cisco",
			rule: SecurityRule{
				Name:            "Check Default Enable Password",
				ExpectedPattern: `^$|enable password \$1\$.*|enable secret \$.*`,
			},
			expectedStatus: StatusFail,
			expectedMsg:    "Configuration does not match expected pattern: ^$|enable password \\$1\\$.*|enable secret \\$.*",
		},
		{
			name:   "SSH Configuration - Secure",
			output: "SSH Enabled - version 2.0",
			rule: SecurityRule{
				Name:            "Check SSH vs Telnet Configuration",
				ExpectedPattern: `SSH Enabled - version [12]\..*`,
			},
			expectedStatus: StatusPass,
			expectedMsg:    "Configuration check passed",
		},
		{
			name:   "SSH Configuration - Insecure",
			output: "SSH Disabled",
			rule: SecurityRule{
				Name:            "Check SSH vs Telnet Configuration",
				ExpectedPattern: `SSH Enabled - version [12]\..*`,
			},
			expectedStatus: StatusFail,
			expectedMsg:    "Configuration does not match expected pattern: SSH Enabled - version [12]\\..*",
		},
		{
			name:   "VTY Lines - Secure Configuration",
			output: "line vty 0 4\n transport input ssh",
			rule: SecurityRule{
				Name:            "Check Telnet VTY Lines",
				ExpectedPattern: `transport input ssh|transport input none`,
			},
			expectedStatus: StatusPass,
			expectedMsg:    "Configuration check passed",
		},
		{
			name:   "VTY Lines - Insecure Configuration",
			output: "line vty 0 4\n transport input telnet",
			rule: SecurityRule{
				Name:            "Check Telnet VTY Lines",
				ExpectedPattern: `transport input ssh|transport input none`,
			},
			expectedStatus: StatusFail,
			expectedMsg:    "Configuration does not match expected pattern: transport input ssh|transport input none",
		},
		{
			name:   "SNMP Community - Secure Configuration",
			output: "snmp-server community MySecureString RO",
			rule: SecurityRule{
				Name:            "Check SNMP Community Strings",
				ExpectedPattern: `snmp-server community [^p].*|snmp-server community p[^ru].*|snmp-server community pr[^i].*|snmp-server community pri[^v].*`,
			},
			expectedStatus: StatusPass,
			expectedMsg:    "Configuration check passed",
		},
		{
			name:   "SNMP Community - Default Public String",
			output: "snmp-server community public RO",
			rule: SecurityRule{
				Name:            "Check SNMP Community Strings",
				ExpectedPattern: `snmp-server community [^p].*|snmp-server community p[^ru].*|snmp-server community pr[^i].*|snmp-server community pri[^v].*`,
			},
			expectedStatus: StatusFail,
			expectedMsg:    "Configuration does not match expected pattern: snmp-server community [^p].*|snmp-server community p[^ru].*|snmp-server community pr[^i].*|snmp-server community pri[^v].*",
		},
		{
			name:   "SNMP Community - Default Private String",
			output: "snmp-server community private RW",
			rule: SecurityRule{
				Name:            "Check SNMP Community Strings",
				ExpectedPattern: `snmp-server community [^p].*|snmp-server community p[^ru].*|snmp-server community pr[^i].*|snmp-server community pri[^v].*`,
			},
			expectedStatus: StatusFail,
			expectedMsg:    "Configuration does not match expected pattern: snmp-server community [^p].*|snmp-server community p[^ru].*|snmp-server community pr[^i].*|snmp-server community pri[^v].*",
		},
		{
			name:   "Password Encryption - Enabled",
			output: "service password-encryption",
			rule: SecurityRule{
				Name:            "Check Service Password Encryption",
				ExpectedPattern: `^service password-encryption$`,
			},
			expectedStatus: StatusPass,
			expectedMsg:    "Configuration check passed",
		},
		{
			name:   "Password Encryption - Disabled",
			output: "no service password-encryption",
			rule: SecurityRule{
				Name:            "Check Service Password Encryption",
				ExpectedPattern: `^service password-encryption$`,
			},
			expectedStatus: StatusFail,
			expectedMsg:    "Configuration does not match expected pattern: ^service password-encryption$",
		},
		{
			name:   "HTTP Server - Secure Configuration",
			output: "no ip http server\nip http secure-server",
			rule: SecurityRule{
				Name:            "Check HTTP/HTTPS Server Status",
				ExpectedPattern: `no ip http server|ip http secure-server`,
			},
			expectedStatus: StatusPass,
			expectedMsg:    "Configuration check passed",
		},
		{
			name:   "HTTP Server - Insecure Configuration",
			output: "ip http server",
			rule: SecurityRule{
				Name:            "Check HTTP/HTTPS Server Status",
				ExpectedPattern: `no ip http server|ip http secure-server`,
			},
			expectedStatus: StatusFail,
			expectedMsg:    "Configuration does not match expected pattern: no ip http server|ip http secure-server",
		},
		{
			name:   "Console Password - Secure Configuration",
			output: "line con 0\n password 7 $1$abcd$xyz\n login",
			rule: SecurityRule{
				Name:            "Check Console Password",
				ExpectedPattern: `password .*|login local`,
			},
			expectedStatus: StatusPass,
			expectedMsg:    "Configuration check passed",
		},
		{
			name:   "Console Password - No Password",
			output: "line con 0\n no login",
			rule: SecurityRule{
				Name:            "Check Console Password",
				ExpectedPattern: `password .*|login local`,
			},
			expectedStatus: StatusFail,
			expectedMsg:    "Configuration does not match expected pattern: password .*|login local",
		},
		{
			name:   "Login Banner - Configured",
			output: "banner login ^C\nAuthorized access only\n^C",
			rule: SecurityRule{
				Name:            "Check Login Banner",
				ExpectedPattern: `^banner (login|motd)`,
			},
			expectedStatus: StatusPass,
			expectedMsg:    "Configuration check passed",
		},
		{
			name:   "Login Banner - Not Configured",
			output: "no banner login",
			rule: SecurityRule{
				Name:            "Check Login Banner",
				ExpectedPattern: `^banner (login|motd)`,
			},
			expectedStatus: StatusFail,
			expectedMsg:    "Configuration does not match expected pattern: ^banner (login|motd)",
		},
		{
			name:   "Empty Pattern - Warning",
			output: "some output",
			rule: SecurityRule{
				Name:            "Test Rule",
				ExpectedPattern: "",
			},
			expectedStatus: StatusWarning,
			expectedMsg:    "No expected pattern defined for rule",
		},
		{
			name:   "Invalid Regex Pattern - Error",
			output: "some output",
			rule: SecurityRule{
				Name:            "Test Rule",
				ExpectedPattern: "[invalid regex",
			},
			expectedStatus: StatusError,
			expectedMsg:    "Invalid regex pattern: error parsing regexp: missing closing ]: `[invalid regex`",
		},
		{
			name:   "Empty Output - Pass for Empty Pattern",
			output: "",
			rule: SecurityRule{
				Name:            "Check Unused Interfaces",
				ExpectedPattern: `.*shutdown.*|^$`,
			},
			expectedStatus: StatusPass,
			expectedMsg:    "Configuration check passed",
		},
		{
			name:   "Unused Interfaces - Shutdown Configured",
			output: "FastEthernet0/1 is administratively down, line protocol is down\n shutdown",
			rule: SecurityRule{
				Name:            "Check Unused Interfaces",
				ExpectedPattern: `.*shutdown.*|^$`,
			},
			expectedStatus: StatusPass,
			expectedMsg:    "Configuration check passed",
		},
		{
			name:   "Unused Interfaces - Not Shutdown",
			output: "FastEthernet0/1 is up, line protocol is down",
			rule: SecurityRule{
				Name:            "Check Unused Interfaces",
				ExpectedPattern: `.*shutdown.*|^$`,
			},
			expectedStatus: StatusFail,
			expectedMsg:    "Configuration does not match expected pattern: .*shutdown.*|^$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, msg := engine.evaluateRuleResult(tt.output, tt.rule)

			if status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, status)
			}

			if msg != tt.expectedMsg {
				t.Errorf("Expected message %q, got %q", tt.expectedMsg, msg)
			}
		})
	}
}

func TestRegexPatterns(t *testing.T) {
	// Test individual regex patterns used in Cisco rules
	tests := []struct {
		name    string
		pattern string
		input   string
		match   bool
	}{
		{
			name:    "Enable Password - Encrypted",
			pattern: `^$|enable password \$1\$.*|enable secret \$.*`,
			input:   "enable secret $1$abcd$xyz123",
			match:   true,
		},
		{
			name:    "Enable Password - Plain Text",
			pattern: `^$|enable password \$1\$.*|enable secret \$.*`,
			input:   "enable password cisco",
			match:   false,
		},
		{
			name:    "SSH Version Check",
			pattern: `SSH Enabled - version [12]\..*`,
			input:   "SSH Enabled - version 2.0",
			match:   true,
		},
		{
			name:    "SSH Version Check - Version 1",
			pattern: `SSH Enabled - version [12]\..*`,
			input:   "SSH Enabled - version 1.99",
			match:   true,
		},
		{
			name:    "SSH Disabled",
			pattern: `SSH Enabled - version [12]\..*`,
			input:   "SSH Disabled",
			match:   false,
		},
		{
			name:    "Transport Input SSH",
			pattern: `transport input ssh|transport input none`,
			input:   "transport input ssh",
			match:   true,
		},
		{
			name:    "Transport Input None",
			pattern: `transport input ssh|transport input none`,
			input:   "transport input none",
			match:   true,
		},
		{
			name:    "Transport Input Telnet",
			pattern: `transport input ssh|transport input none`,
			input:   "transport input telnet",
			match:   false,
		},
		{
			name:    "SNMP Community - Not Default",
			pattern: `snmp-server community [^p].*|snmp-server community p[^ru].*|snmp-server community pr[^i].*|snmp-server community pri[^v].*`,
			input:   "snmp-server community MySecureString RO",
			match:   true,
		},
		{
			name:    "SNMP Community - Public Default",
			pattern: `snmp-server community [^p].*|snmp-server community p[^ru].*|snmp-server community pr[^i].*|snmp-server community pri[^v].*`,
			input:   "snmp-server community public RO",
			match:   false,
		},
		{
			name:    "SNMP Community - Private Default",
			pattern: `snmp-server community [^p].*|snmp-server community p[^ru].*|snmp-server community pr[^i].*|snmp-server community pri[^v].*`,
			input:   "snmp-server community private RW",
			match:   false,
		},
		{
			name:    "HTTP Server - Disabled",
			pattern: `no ip http server|ip http secure-server`,
			input:   "no ip http server",
			match:   true,
		},
		{
			name:    "HTTPS Server - Enabled",
			pattern: `no ip http server|ip http secure-server`,
			input:   "ip http secure-server",
			match:   true,
		},
		{
			name:    "HTTP Server - Enabled (Insecure)",
			pattern: `no ip http server|ip http secure-server`,
			input:   "ip http server",
			match:   false,
		},
		{
			name:    "Banner Login",
			pattern: `^banner (login|motd)`,
			input:   "banner login ^C\nAuthorized access only\n^C",
			match:   true,
		},
		{
			name:    "Banner MOTD",
			pattern: `^banner (login|motd)`,
			input:   "banner motd ^C\nWelcome\n^C",
			match:   true,
		},
		{
			name:    "No Banner",
			pattern: `^banner (login|motd)`,
			input:   "no banner login",
			match:   false,
		},
		{
			name:    "Console Password",
			pattern: `password .*|login local`,
			input:   "password 7 $1$abcd$xyz",
			match:   true,
		},
		{
			name:    "Console Login Local",
			pattern: `password .*|login local`,
			input:   "login local",
			match:   true,
		},
		{
			name:    "Console No Login",
			pattern: `password .*|login local`,
			input:   "no login",
			match:   false,
		},
		{
			name:    "Shutdown Interface",
			pattern: `.*shutdown.*|^$`,
			input:   "FastEthernet0/1 is administratively down\n shutdown",
			match:   true,
		},
		{
			name:    "Empty Output",
			pattern: `.*shutdown.*|^$`,
			input:   "",
			match:   true,
		},
		{
			name:    "Interface Up",
			pattern: `.*shutdown.*|^$`,
			input:   "FastEthernet0/1 is up, line protocol is up",
			match:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regex, err := regexp.Compile(tt.pattern)
			if err != nil {
				t.Fatalf("Failed to compile regex pattern %s: %v", tt.pattern, err)
			}

			match := regex.MatchString(tt.input)
			if match != tt.match {
				t.Errorf("Pattern %s with input %q: expected match=%t, got match=%t",
					tt.pattern, tt.input, tt.match, match)
			}
		})
	}
}

func TestRuleEvaluationEdgeCases(t *testing.T) {
	rm := setupTestRuleManager(t)
	engine := NewEngine(rm)

	tests := []struct {
		name           string
		output         string
		rule           SecurityRule
		expectedStatus CheckStatus
	}{
		{
			name:   "Multiline Output - SSH Check",
			output: "SSH Enabled - version 2.0\nAuthentication timeout: 120 secs; Authentication retries: 3",
			rule: SecurityRule{
				Name:            "Check SSH Configuration",
				ExpectedPattern: `SSH Enabled - version [12]\..*`,
			},
			expectedStatus: StatusPass,
		},
		{
			name:   "Case Sensitive Pattern",
			output: "ssh enabled - version 2.0",
			rule: SecurityRule{
				Name:            "Check SSH Configuration",
				ExpectedPattern: `SSH Enabled - version [12]\..*`,
			},
			expectedStatus: StatusFail,
		},
		{
			name:   "Special Characters in Output",
			output: "banner login ^C\nUNAUTHORIZED ACCESS TO THIS DEVICE IS PROHIBITED!\n^C",
			rule: SecurityRule{
				Name:            "Check Login Banner",
				ExpectedPattern: `banner (login|motd)`,
			},
			expectedStatus: StatusPass,
		},
		{
			name:   "Very Long Output",
			output: "show running-config\n" + generateLongString(1000) + "\nservice password-encryption",
			rule: SecurityRule{
				Name:            "Check Service Password Encryption",
				ExpectedPattern: `service password-encryption`,
			},
			expectedStatus: StatusPass,
		},
		{
			name:   "Empty Output with Non-Empty Pattern",
			output: "",
			rule: SecurityRule{
				Name:            "Check SSH Configuration",
				ExpectedPattern: `SSH Enabled - version [12]\..*`,
			},
			expectedStatus: StatusFail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, _ := engine.evaluateRuleResult(tt.output, tt.rule)

			if status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, status)
			}
		})
	}
}

// generateLongString creates a string of specified length for testing
func generateLongString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = 'a'
	}
	return string(result)
}

func BenchmarkRuleEvaluation(b *testing.B) {
	// Create test database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		b.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	// Create table
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
		b.Fatalf("Failed to create test table: %v", err)
	}

	rm := NewRuleManager(db)
	engine := NewEngine(rm)

	rule := SecurityRule{
		Name:            "Check SSH Configuration",
		ExpectedPattern: `SSH Enabled - version [12]\..*`,
	}

	output := "SSH Enabled - version 2.0\nAuthentication timeout: 120 secs"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.evaluateRuleResult(output, rule)
	}
}

func BenchmarkComplexRegexEvaluation(b *testing.B) {
	// Create test database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		b.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	// Create table
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
		b.Fatalf("Failed to create test table: %v", err)
	}

	rm := NewRuleManager(db)
	engine := NewEngine(rm)

	rule := SecurityRule{
		Name:            "Check SNMP Community Strings",
		ExpectedPattern: `snmp-server community [^p].*|snmp-server community p[^ru].*|snmp-server community pr[^i].*|snmp-server community pri[^v].*`,
	}

	output := "snmp-server community MyVeryLongAndComplexCommunityStringThatShouldNotMatchDefaults RO"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.evaluateRuleResult(output, rule)
	}
}
