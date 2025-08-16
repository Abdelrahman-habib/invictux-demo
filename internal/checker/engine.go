package checker

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	"invictux-demo/internal/device"
	"invictux-demo/internal/ssh"

	"github.com/google/uuid"
)

// Engine handles security check execution
type Engine struct {
	sshClient   ssh.SSHClientInterface
	ruleManager *RuleManager
	workerCount int
	timeout     time.Duration
}

// CheckJob represents a security check job for a device
type CheckJob struct {
	Device *device.Device
	Rules  []SecurityRule
}

// CheckProgress represents the progress of security checks
type CheckProgress struct {
	DeviceID    string    `json:"deviceId"`
	DeviceName  string    `json:"deviceName"`
	Status      string    `json:"status"`
	Progress    int       `json:"progress"`
	Total       int       `json:"total"`
	CurrentRule string    `json:"currentRule"`
	Error       string    `json:"error,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// BulkCheckResult represents the result of bulk security checks
type BulkCheckResult struct {
	DeviceResults map[string][]CheckResult  `json:"deviceResults"`
	Progress      map[string]*CheckProgress `json:"progress"`
	Errors        map[string]error          `json:"errors"`
}

// ProgressCallback is called to report progress updates
type ProgressCallback func(progress *CheckProgress)

// NewEngine creates a new security check engine
func NewEngine(ruleManager *RuleManager) *Engine {
	return &Engine{
		sshClient:   ssh.NewSSHClient(nil), // Use default config
		ruleManager: ruleManager,
		workerCount: 5, // Default worker pool size
		timeout:     30 * time.Second,
	}
}

// NewEngineWithSSHClient creates a new engine with a custom SSH client
func NewEngineWithSSHClient(ruleManager *RuleManager, sshClient ssh.SSHClientInterface) *Engine {
	return &Engine{
		sshClient:   sshClient,
		ruleManager: ruleManager,
		workerCount: 5,
		timeout:     30 * time.Second,
	}
}

// SetWorkerCount sets the number of workers for parallel processing
func (e *Engine) SetWorkerCount(count int) {
	if count > 0 {
		e.workerCount = count
	}
}

// SetTimeout sets the timeout for security checks
func (e *Engine) SetTimeout(timeout time.Duration) {
	e.timeout = timeout
}

// RunChecks executes security checks on a device
func (e *Engine) RunChecks(device *device.Device) ([]CheckResult, error) {
	return e.RunChecksWithProgress(device, nil)
}

// RunChecksWithProgress executes security checks on a device with progress reporting
func (e *Engine) RunChecksWithProgress(device *device.Device, progressCallback ProgressCallback) ([]CheckResult, error) {
	var results []CheckResult

	// Get applicable rules for this device
	applicableRules := e.GetSecurityRules(device.Vendor)

	// Initialize progress tracking
	progress := &CheckProgress{
		DeviceID:   device.ID,
		DeviceName: device.Name,
		Status:     "running",
		Progress:   0,
		Total:      len(applicableRules),
		UpdatedAt:  time.Now(),
	}

	if progressCallback != nil {
		progressCallback(progress)
	}

	if len(applicableRules) == 0 {
		// Update progress to show completion even with no rules
		progress.Status = "completed"
		progress.UpdatedAt = time.Now()
		if progressCallback != nil {
			progressCallback(progress)
		}
		return results, fmt.Errorf("no security rules found for vendor: %s", device.Vendor)
	}

	// Execute each rule
	for i, rule := range applicableRules {
		if !rule.Enabled {
			continue
		}

		progress.CurrentRule = rule.Name
		progress.Progress = i
		progress.UpdatedAt = time.Now()

		if progressCallback != nil {
			progressCallback(progress)
		}

		result, err := e.executeRule(device, rule)
		if err != nil {
			// Create error result
			result = CheckResult{
				ID:        uuid.New().String(),
				DeviceID:  device.ID,
				CheckName: rule.Name,
				CheckType: "configuration",
				Severity:  rule.Severity,
				Status:    string(StatusError),
				Message:   fmt.Sprintf("Check execution failed: %s", err.Error()),
				Evidence:  "",
				CheckedAt: time.Now(),
			}
		}

		results = append(results, result)
	}

	// Update final progress
	progress.Status = "completed"
	progress.Progress = len(applicableRules)
	progress.CurrentRule = ""
	progress.UpdatedAt = time.Now()

	if progressCallback != nil {
		progressCallback(progress)
	}

	return results, nil
}

// executeRule executes a single security rule against a device
func (e *Engine) executeRule(device *device.Device, rule SecurityRule) (CheckResult, error) {
	result := CheckResult{
		ID:        uuid.New().String(),
		DeviceID:  device.ID,
		CheckName: rule.Name,
		CheckType: "configuration",
		Severity:  rule.Severity,
		Status:    string(StatusError),
		Message:   "",
		Evidence:  "",
		CheckedAt: time.Now(),
	}

	// Create connection info for the advanced SSH client
	connInfo := &ssh.ConnectionInfo{
		Host:       device.IPAddress,
		Port:       device.SSHPort,
		Username:   device.Username,
		Password:   "placeholder", // TODO: Decrypt device.PasswordEncrypted
		AuthMethod: ssh.AuthPassword,
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	// Connect to device via SSH
	conn, err := e.sshClient.Connect(ctx, connInfo)
	if err != nil {
		result.Message = fmt.Sprintf("SSH connection failed: %s", err.Error())
		return result, nil // Return result with error status, don't fail the entire check
	}
	defer e.sshClient.Disconnect(conn)

	// Execute the command
	cmdResult, err := e.sshClient.ExecuteCommand(ctx, conn, rule.Command)
	if err != nil {
		result.Message = fmt.Sprintf("Command execution failed: %s", err.Error())
		return result, nil
	}

	result.Evidence = cmdResult.Output

	// Evaluate the result against expected pattern
	status, message := e.evaluateRuleResult(cmdResult.Output, rule)
	result.Status = string(status)
	result.Message = message

	return result, nil
}

// evaluateRuleResult evaluates command output against rule expectations
func (e *Engine) evaluateRuleResult(output string, rule SecurityRule) (CheckStatus, string) {
	if rule.ExpectedPattern == "" {
		return StatusWarning, "No expected pattern defined for rule"
	}

	// Compile regex pattern
	regex, err := regexp.Compile(rule.ExpectedPattern)
	if err != nil {
		return StatusError, fmt.Sprintf("Invalid regex pattern: %s", err.Error())
	}

	// Check if pattern matches
	if regex.MatchString(output) {
		return StatusPass, "Configuration check passed"
	}

	// Pattern doesn't match - this could be a security issue
	return StatusFail, fmt.Sprintf("Configuration does not match expected pattern: %s", rule.ExpectedPattern)
}

// RunBulkChecks executes checks on multiple devices with parallel processing
func (e *Engine) RunBulkChecks(devices []device.Device) (map[string][]CheckResult, error) {
	return e.RunBulkChecksWithProgress(devices, nil)
}

// RunBulkChecksWithProgress executes checks on multiple devices with progress reporting
func (e *Engine) RunBulkChecksWithProgress(devices []device.Device, progressCallback ProgressCallback) (map[string][]CheckResult, error) {
	if len(devices) == 0 {
		return make(map[string][]CheckResult), nil
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout*time.Duration(len(devices)))
	defer cancel()

	// Initialize result structures
	results := make(map[string][]CheckResult)
	progress := make(map[string]*CheckProgress)
	errors := make(map[string]error)

	// Mutex for thread-safe access to shared data
	var mu sync.Mutex

	// Create job channel
	jobs := make(chan CheckJob, len(devices))

	// Create worker pool
	var wg sync.WaitGroup
	for i := 0; i < e.workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.worker(ctx, jobs, &mu, results, progress, errors, progressCallback)
		}()
	}

	// Send jobs to workers
	for _, dev := range devices {
		deviceCopy := dev // Create copy to avoid race conditions
		applicableRules := e.GetSecurityRules(deviceCopy.Vendor)

		// Initialize progress for this device
		mu.Lock()
		progress[deviceCopy.ID] = &CheckProgress{
			DeviceID:   deviceCopy.ID,
			DeviceName: deviceCopy.Name,
			Status:     "queued",
			Progress:   0,
			Total:      len(applicableRules),
			UpdatedAt:  time.Now(),
		}
		mu.Unlock()

		if progressCallback != nil {
			progressCallback(progress[deviceCopy.ID])
		}

		jobs <- CheckJob{
			Device: &deviceCopy,
			Rules:  applicableRules,
		}
	}
	close(jobs)

	// Wait for all workers to complete
	wg.Wait()

	return results, nil
}

// worker processes security check jobs from the job channel
func (e *Engine) worker(ctx context.Context, jobs <-chan CheckJob, mu *sync.Mutex,
	results map[string][]CheckResult, progress map[string]*CheckProgress,
	errors map[string]error, progressCallback ProgressCallback) {

	for job := range jobs {
		select {
		case <-ctx.Done():
			// Context cancelled, stop processing
			mu.Lock()
			if prog, exists := progress[job.Device.ID]; exists {
				prog.Status = "cancelled"
				prog.Error = "Operation cancelled due to timeout"
				prog.UpdatedAt = time.Now()
			}
			mu.Unlock()
			return
		default:
			// Process the job
			deviceResults, err := e.runChecksForJob(job, mu, progress, progressCallback)

			mu.Lock()
			if err != nil {
				errors[job.Device.ID] = err
				if prog, exists := progress[job.Device.ID]; exists {
					prog.Status = "error"
					prog.Error = err.Error()
					prog.UpdatedAt = time.Now()
				}
			} else {
				results[job.Device.ID] = deviceResults
				if prog, exists := progress[job.Device.ID]; exists {
					prog.Status = "completed"
					prog.Progress = prog.Total
					prog.CurrentRule = ""
					prog.UpdatedAt = time.Now()
				}
			}
			mu.Unlock()

			// Report final progress
			if progressCallback != nil {
				mu.Lock()
				if prog, exists := progress[job.Device.ID]; exists {
					progressCallback(prog)
				}
				mu.Unlock()
			}
		}
	}
}

// runChecksForJob executes security checks for a specific job
func (e *Engine) runChecksForJob(job CheckJob, mu *sync.Mutex,
	progress map[string]*CheckProgress, progressCallback ProgressCallback) ([]CheckResult, error) {

	var results []CheckResult

	// Update progress to running
	mu.Lock()
	if prog, exists := progress[job.Device.ID]; exists {
		prog.Status = "running"
		prog.UpdatedAt = time.Now()
	}
	mu.Unlock()

	if progressCallback != nil {
		mu.Lock()
		if prog, exists := progress[job.Device.ID]; exists {
			progressCallback(prog)
		}
		mu.Unlock()
	}

	// Execute each rule
	for i, rule := range job.Rules {
		if !rule.Enabled {
			continue
		}

		// Update progress
		mu.Lock()
		if prog, exists := progress[job.Device.ID]; exists {
			prog.CurrentRule = rule.Name
			prog.Progress = i
			prog.UpdatedAt = time.Now()
		}
		mu.Unlock()

		if progressCallback != nil {
			mu.Lock()
			if prog, exists := progress[job.Device.ID]; exists {
				progressCallback(prog)
			}
			mu.Unlock()
		}

		result, err := e.executeRule(job.Device, rule)
		if err != nil {
			// Create error result but continue with other rules
			result = CheckResult{
				ID:        uuid.New().String(),
				DeviceID:  job.Device.ID,
				CheckName: rule.Name,
				CheckType: "configuration",
				Severity:  rule.Severity,
				Status:    string(StatusError),
				Message:   fmt.Sprintf("Check execution failed: %s", err.Error()),
				Evidence:  "",
				CheckedAt: time.Now(),
			}
		}

		results = append(results, result)
	}

	return results, nil
}

// GetSecurityRules returns security rules for a specific vendor
func (e *Engine) GetSecurityRules(vendorType string) []SecurityRule {
	if e.ruleManager == nil {
		return []SecurityRule{}
	}

	rules, err := e.ruleManager.GetRulesByVendor(vendorType)
	if err != nil {
		// Log error and return empty slice
		return []SecurityRule{}
	}

	// Filter only enabled rules
	var enabledRules []SecurityRule
	for _, rule := range rules {
		if rule.Enabled {
			enabledRules = append(enabledRules, rule)
		}
	}

	return enabledRules
}

// LoadCustomRules loads custom security rules into the database
func (e *Engine) LoadCustomRules(rules []SecurityRule) error {
	if e.ruleManager == nil {
		return fmt.Errorf("rule manager not initialized")
	}

	for _, rule := range rules {
		if err := e.ruleManager.CreateRule(rule); err != nil {
			return fmt.Errorf("failed to create rule %s: %w", rule.Name, err)
		}
	}

	return nil
}

// GetProgress returns the current progress for all devices
func (e *Engine) GetProgress() map[string]*CheckProgress {
	// This would typically be stored in a shared state manager
	// For now, return empty map as progress is handled in the callback
	return make(map[string]*CheckProgress)
}
