package checker

import (
	"testing"
	"time"

	"invictux-demo/internal/device"

	"github.com/stretchr/testify/assert"
)

// TestEngine_Integration tests the engine with real components
func TestEngine_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	engine := NewEngine()

	// Test engine configuration
	engine.SetWorkerCount(3)
	engine.SetTimeout(10 * time.Second)

	// Load test rules
	rules := []SecurityRule{
		{
			ID:              "test-rule-1",
			Name:            "Version Check",
			Description:     "Check device version information",
			Vendor:          "cisco",
			Command:         "show version",
			ExpectedPattern: ".*", // Accept any output for integration test
			Severity:        string(SeverityMedium),
			Enabled:         true,
			CreatedAt:       time.Now(),
		},
		{
			ID:              "test-rule-2",
			Name:            "Generic Config Check",
			Description:     "Generic configuration check",
			Vendor:          "generic",
			Command:         "echo 'test'",
			ExpectedPattern: "test",
			Severity:        string(SeverityLow),
			Enabled:         true,
			CreatedAt:       time.Now(),
		},
		{
			ID:              "test-rule-3",
			Name:            "Disabled Rule",
			Description:     "This rule should be skipped",
			Vendor:          "cisco",
			Command:         "show interfaces",
			ExpectedPattern: ".*",
			Severity:        string(SeverityHigh),
			Enabled:         false, // Disabled rule
			CreatedAt:       time.Now(),
		},
	}

	engine.LoadRules(rules)

	// Test rule filtering
	t.Run("Rule Filtering", func(t *testing.T) {
		ciscoRules := engine.GetSecurityRules("cisco")
		assert.Len(t, ciscoRules, 3) // 2 cisco rules + 1 generic rule

		juniperRules := engine.GetSecurityRules("juniper")
		assert.Len(t, juniperRules, 1) // Only generic rule

		unknownRules := engine.GetSecurityRules("unknown")
		assert.Len(t, unknownRules, 1) // Only generic rule
	})

	// Test progress tracking with mock device
	t.Run("Progress Tracking", func(t *testing.T) {
		testDevice := &device.Device{
			ID:        "integration-test-device",
			Name:      "Integration Test Device",
			IPAddress: "127.0.0.1", // Localhost for testing
			Vendor:    "cisco",
			Username:  "test",
			SSHPort:   22,
		}

		var progressUpdates []*CheckProgress
		progressCallback := func(progress *CheckProgress) {
			// Make a copy to avoid race conditions
			progressCopy := *progress
			progressUpdates = append(progressUpdates, &progressCopy)
		}

		// This will attempt to connect and should track progress
		results, err := engine.RunChecksWithProgress(testDevice, progressCallback)

		// Should have received progress updates
		assert.NotEmpty(t, progressUpdates)

		// Check progress structure
		firstUpdate := progressUpdates[0]
		assert.Equal(t, testDevice.ID, firstUpdate.DeviceID)
		assert.Equal(t, testDevice.Name, firstUpdate.DeviceName)
		assert.Equal(t, "running", firstUpdate.Status)
		assert.Equal(t, 0, firstUpdate.Progress)
		assert.Greater(t, firstUpdate.Total, 0)

		// Should have a final update
		lastUpdate := progressUpdates[len(progressUpdates)-1]
		assert.Contains(t, []string{"completed", "error"}, lastUpdate.Status)

		// Results should be returned even if there are connection errors
		assert.NotNil(t, results)

		// Connection errors are handled at the individual check level
		// The function should complete successfully even with SSH failures
		if err != nil {
			t.Logf("Expected connection error occurred: %v", err)
		}
	})

	// Test bulk operations with multiple devices
	t.Run("Bulk Operations", func(t *testing.T) {
		devices := []device.Device{
			{
				ID:        "bulk-test-1",
				Name:      "Bulk Test Device 1",
				IPAddress: "192.168.1.1",
				Vendor:    "cisco",
				Username:  "admin",
				SSHPort:   22,
			},
			{
				ID:        "bulk-test-2",
				Name:      "Bulk Test Device 2",
				IPAddress: "192.168.1.2",
				Vendor:    "juniper",
				Username:  "admin",
				SSHPort:   22,
			},
			{
				ID:        "bulk-test-3",
				Name:      "Bulk Test Device 3",
				IPAddress: "192.168.1.3",
				Vendor:    "unknown",
				Username:  "admin",
				SSHPort:   22,
			},
		}

		var allProgressUpdates []*CheckProgress
		progressCallback := func(progress *CheckProgress) {
			progressCopy := *progress
			allProgressUpdates = append(allProgressUpdates, &progressCopy)
		}

		results, err := engine.RunBulkChecksWithProgress(devices, progressCallback)

		// Should not error on the bulk operation itself
		assert.NoError(t, err)
		assert.NotNil(t, results)

		// Should have received progress updates for each device
		assert.NotEmpty(t, allProgressUpdates)

		// Check that we have progress for each device
		deviceProgressMap := make(map[string][]*CheckProgress)
		for _, progress := range allProgressUpdates {
			deviceProgressMap[progress.DeviceID] = append(deviceProgressMap[progress.DeviceID], progress)
		}

		// Should have progress for each device
		for _, dev := range devices {
			assert.Contains(t, deviceProgressMap, dev.ID)
			deviceProgress := deviceProgressMap[dev.ID]
			assert.NotEmpty(t, deviceProgress)

			// First update should be queued or running
			firstUpdate := deviceProgress[0]
			assert.Contains(t, []string{"queued", "running"}, firstUpdate.Status)

			// Last update should be completed, error, or cancelled
			lastUpdate := deviceProgress[len(deviceProgress)-1]
			assert.Contains(t, []string{"completed", "error", "cancelled"}, lastUpdate.Status)
		}
	})

	// Test worker pool configuration
	t.Run("Worker Pool Configuration", func(t *testing.T) {
		// Test different worker counts
		testCounts := []int{1, 3, 5, 10}

		for _, count := range testCounts {
			engine.SetWorkerCount(count)
			assert.Equal(t, count, engine.workerCount)
		}

		// Test invalid worker counts
		engine.SetWorkerCount(0)
		assert.Equal(t, 10, engine.workerCount) // Should remain unchanged

		engine.SetWorkerCount(-1)
		assert.Equal(t, 10, engine.workerCount) // Should remain unchanged
	})

	// Test timeout configuration
	t.Run("Timeout Configuration", func(t *testing.T) {
		originalTimeout := engine.timeout

		newTimeout := 60 * time.Second
		engine.SetTimeout(newTimeout)
		assert.Equal(t, newTimeout, engine.timeout)

		// Restore original timeout
		engine.SetTimeout(originalTimeout)
	})
}

// TestEngine_RuleEvaluation tests rule evaluation with various patterns
func TestEngine_RuleEvaluation(t *testing.T) {
	engine := NewEngine()

	testCases := []struct {
		name           string
		output         string
		rule           SecurityRule
		expectedStatus CheckStatus
		description    string
	}{
		{
			name:   "Cisco IOS Version Check",
			output: "Cisco IOS Software, C2960X-STACK Software (C2960X-UNIVERSALK9-M), Version 15.2(4)E10",
			rule: SecurityRule{
				Name:            "IOS Version Check",
				ExpectedPattern: "Cisco IOS.*Version \\d+\\.\\d+",
			},
			expectedStatus: StatusPass,
			description:    "Should pass when Cisco IOS version is found",
		},
		{
			name:   "SSH Configuration Check",
			output: "ip ssh version 2\nip ssh time-out 60",
			rule: SecurityRule{
				Name:            "SSH Version Check",
				ExpectedPattern: "ip ssh version 2",
			},
			expectedStatus: StatusPass,
			description:    "Should pass when SSH version 2 is configured",
		},
		{
			name:   "Password Policy Check",
			output: "enable secret 5 $1$abcd$efghijklmnopqrstuvwxyz",
			rule: SecurityRule{
				Name:            "Enable Secret Check",
				ExpectedPattern: "enable secret [5-9]",
			},
			expectedStatus: StatusPass,
			description:    "Should pass when encrypted enable secret is configured",
		},
		{
			name:   "Weak Password Check",
			output: "enable password cisco123",
			rule: SecurityRule{
				Name:            "Weak Password Check",
				ExpectedPattern: "enable secret",
			},
			expectedStatus: StatusFail,
			description:    "Should fail when plain text password is used instead of secret",
		},
		{
			name:   "Interface Security Check",
			output: "interface GigabitEthernet0/1\n shutdown\n switchport mode access",
			rule: SecurityRule{
				Name:            "Unused Interface Check",
				ExpectedPattern: "shutdown",
			},
			expectedStatus: StatusPass,
			description:    "Should pass when unused interfaces are shutdown",
		},
		{
			name:   "SNMP Community Check",
			output: "snmp-server community public RO\nsnmp-server community private RW",
			rule: SecurityRule{
				Name:            "Default SNMP Community Check",
				ExpectedPattern: "snmp-server community (public|private)",
			},
			expectedStatus: StatusPass,
			description:    "Should detect default SNMP communities (this would typically be a FAIL in real scenarios)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, message := engine.evaluateRuleResult(tc.output, tc.rule)
			assert.Equal(t, tc.expectedStatus, status, tc.description)
			assert.NotEmpty(t, message)
		})
	}
}

// TestEngine_ConcurrentAccess tests concurrent access to the engine
func TestEngine_ConcurrentAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent access test in short mode")
	}

	engine := NewEngine()

	// Load rules
	rules := []SecurityRule{
		{
			ID:              "concurrent-rule-1",
			Name:            "Concurrent Test Rule 1",
			Vendor:          "generic",
			Command:         "echo test1",
			ExpectedPattern: "test1",
			Severity:        string(SeverityLow),
			Enabled:         true,
		},
		{
			ID:              "concurrent-rule-2",
			Name:            "Concurrent Test Rule 2",
			Vendor:          "generic",
			Command:         "echo test2",
			ExpectedPattern: "test2",
			Severity:        string(SeverityLow),
			Enabled:         true,
		},
	}

	engine.LoadRules(rules)

	// Test concurrent rule access
	t.Run("Concurrent Rule Access", func(t *testing.T) {
		done := make(chan bool, 10)

		// Start multiple goroutines accessing rules
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()

				// Access rules multiple times
				for j := 0; j < 100; j++ {
					rules := engine.GetSecurityRules("generic")
					assert.Len(t, rules, 2)
				}
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	// Test concurrent rule evaluation
	t.Run("Concurrent Rule Evaluation", func(t *testing.T) {
		done := make(chan bool, 5)

		testOutputs := []string{
			"test1 output",
			"test2 output",
			"different output",
			"another test",
			"final test",
		}

		// Start multiple goroutines evaluating rules
		for i, output := range testOutputs {
			go func(output string, index int) {
				defer func() { done <- true }()

				rule := SecurityRule{
					Name:            "Concurrent Eval Rule",
					ExpectedPattern: "test\\d+",
				}

				status, message := engine.evaluateRuleResult(output, rule)
				assert.NotEmpty(t, message)
				assert.Contains(t, []CheckStatus{StatusPass, StatusFail}, status)
			}(output, i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < len(testOutputs); i++ {
			<-done
		}
	})
}
