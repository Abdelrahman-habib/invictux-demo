package checker

import (
	"fmt"
	"qwin/internal/device"
)

// Engine handles security check execution
type Engine struct {
	sshClient *SSHClient
	rules     []SecurityRule
}

// NewEngine creates a new security check engine
func NewEngine() *Engine {
	return &Engine{
		sshClient: NewSSHClient(),
		rules:     []SecurityRule{}, // Will be loaded from database
	}
}

// RunChecks executes security checks on a device
func (e *Engine) RunChecks(device *device.Device) ([]CheckResult, error) {
	var results []CheckResult

	// TODO: Implement actual security checks
	// For now, return mock results
	mockResult := CheckResult{
		DeviceID:  device.ID,
		CheckName: "Mock Security Check",
		CheckType: "configuration",
		Severity:  string(SeverityMedium),
		Status:    string(StatusPass),
		Message:   "Mock check passed",
		Evidence:  "Mock evidence",
	}

	results = append(results, mockResult)
	return results, nil
}

// RunBulkChecks executes checks on multiple devices
func (e *Engine) RunBulkChecks(devices []device.Device) (map[string][]CheckResult, error) {
	results := make(map[string][]CheckResult)

	for _, dev := range devices {
		deviceResults, err := e.RunChecks(&dev)
		if err != nil {
			return nil, fmt.Errorf("failed to check device %s: %w", dev.Name, err)
		}
		results[dev.ID] = deviceResults
	}

	return results, nil
}

// GetSecurityRules returns security rules for a specific vendor
func (e *Engine) GetSecurityRules(vendorType string) []SecurityRule {
	var filteredRules []SecurityRule
	for _, rule := range e.rules {
		if rule.Vendor == vendorType || rule.Vendor == "generic" {
			filteredRules = append(filteredRules, rule)
		}
	}
	return filteredRules
}

// LoadRules loads security rules from database
func (e *Engine) LoadRules(rules []SecurityRule) {
	e.rules = rules
}
