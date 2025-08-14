package checker

import "time"

// CheckResult represents the result of a security check
type CheckResult struct {
	ID        string    `json:"id" db:"id"`
	DeviceID  string    `json:"deviceId" db:"device_id"`
	CheckName string    `json:"checkName" db:"check_name"`
	CheckType string    `json:"checkType" db:"check_type"`
	Severity  string    `json:"severity" db:"severity"`
	Status    string    `json:"status" db:"status"`
	Message   string    `json:"message" db:"message"`
	Evidence  string    `json:"evidence" db:"evidence"`
	CheckedAt time.Time `json:"checkedAt" db:"checked_at"`
}

// SecurityRule represents a security check rule
type SecurityRule struct {
	ID              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Description     string    `json:"description" db:"description"`
	Vendor          string    `json:"vendor" db:"vendor"`
	Command         string    `json:"command" db:"command"`
	ExpectedPattern string    `json:"expectedPattern" db:"expected_pattern"`
	Severity        string    `json:"severity" db:"severity"`
	Enabled         bool      `json:"enabled" db:"enabled"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}

// CheckStatus represents the status of a security check
type CheckStatus string

const (
	StatusPass    CheckStatus = "PASS"
	StatusFail    CheckStatus = "FAIL"
	StatusWarning CheckStatus = "WARNING"
	StatusError   CheckStatus = "ERROR"
)

// Severity levels for security checks
type Severity string

const (
	SeverityCritical Severity = "Critical"
	SeverityHigh     Severity = "High"
	SeverityMedium   Severity = "Medium"
	SeverityLow      Severity = "Low"
)
