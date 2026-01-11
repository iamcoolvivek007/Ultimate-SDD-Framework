package workflows

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Workflow represents a multi-step development workflow
type Workflow struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Track       string            `json:"track"` // quick, standard, enterprise
	Steps       []*WorkflowStep   `json:"steps"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// WorkflowStep represents a single step in a workflow
type WorkflowStep struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Command      string     `json:"command,omitempty"`      // CLI command to run
	Agent        string     `json:"agent,omitempty"`        // Agent to use
	Dependencies []string   `json:"dependencies,omitempty"` // Step IDs that must complete first
	Parallel     bool       `json:"parallel,omitempty"`     // Can run in parallel with others
	Required     bool       `json:"required,omitempty"`     // Must complete to proceed
	Condition    string     `json:"condition,omitempty"`    // Conditional execution
	Outputs      []string   `json:"outputs,omitempty"`      // Expected output files
	Status       StepStatus `json:"status"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}

// StepStatus represents the status of a workflow step
type StepStatus string

const (
	StepStatusPending    StepStatus = "pending"
	StepStatusInProgress StepStatus = "in_progress"
	StepStatusCompleted  StepStatus = "completed"
	StepStatusFailed     StepStatus = "failed"
	StepStatusSkipped    StepStatus = "skipped"
)

// Track represents a workflow track (complexity level)
type Track struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TimeToFirst string `json:"time_to_first"` // Time to first story
}

// GetTracks returns available workflow tracks
func GetTracks() []*Track {
	return []*Track{
		{
			ID:          "quick",
			Name:        "Quick Flow",
			Description: "Bug fixes, small features",
			TimeToFirst: "~5 minutes",
		},
		{
			ID:          "standard",
			Name:        "Standard Method",
			Description: "Products and platforms",
			TimeToFirst: "~15 minutes",
		},
		{
			ID:          "enterprise",
			Name:        "Enterprise",
			Description: "Compliance-heavy systems",
			TimeToFirst: "~30 minutes",
		},
	}
}

// GetBuiltinWorkflows returns built-in workflow definitions
func GetBuiltinWorkflows() []*Workflow {
	return []*Workflow{
		{
			ID:          "quick-fix",
			Name:        "Quick Fix",
			Description: "Fast path for bug fixes and small changes",
			Track:       "quick",
			Steps: []*WorkflowStep{
				{ID: "identify", Name: "Identify Issue", Description: "Locate and understand the problem", Agent: "developer", Required: true},
				{ID: "fix", Name: "Implement Fix", Description: "Make the necessary changes", Agent: "developer", Dependencies: []string{"identify"}, Required: true},
				{ID: "test", Name: "Verify Fix", Description: "Test that the fix works", Agent: "qa", Dependencies: []string{"fix"}, Required: true},
			},
		},
		{
			ID:          "standard",
			Name:        "Standard Development",
			Description: "Full development workflow for features",
			Track:       "standard",
			Steps: []*WorkflowStep{
				{ID: "init", Name: "Initialize", Description: "Set up project structure", Command: "viki init", Required: true},
				{ID: "discover", Name: "Discover", Description: "Analyze existing codebase", Command: "viki discovery", Condition: "brownfield"},
				{ID: "specify", Name: "Specify", Description: "Define requirements", Command: "viki specify", Agent: "pm", Dependencies: []string{"init"}, Required: true},
				{ID: "plan", Name: "Plan", Description: "Design architecture", Command: "viki plan", Agent: "architect", Dependencies: []string{"specify"}, Required: true},
				{ID: "approve-plan", Name: "Approve Plan", Description: "Review and approve plan", Command: "viki approve", Dependencies: []string{"plan"}, Required: true},
				{ID: "task", Name: "Break Down Tasks", Description: "Create task list", Command: "viki task", Agent: "developer", Dependencies: []string{"approve-plan"}, Required: true},
				{ID: "execute", Name: "Execute", Description: "Implement solution", Command: "viki execute", Agent: "developer", Dependencies: []string{"task"}, Required: true},
				{ID: "review", Name: "Review", Description: "Quality review", Command: "viki review", Agent: "qa", Dependencies: []string{"execute"}, Required: true},
				{ID: "approve-final", Name: "Final Approval", Description: "Final sign-off", Command: "viki approve", Dependencies: []string{"review"}, Required: true},
			},
		},
		{
			ID:          "enterprise",
			Name:        "Enterprise Development",
			Description: "Comprehensive workflow with compliance",
			Track:       "enterprise",
			Steps: []*WorkflowStep{
				{ID: "init", Name: "Initialize", Description: "Set up project structure", Command: "viki init", Required: true},
				{ID: "constitution", Name: "Constitution", Description: "Define governance", Command: "viki constitution", Required: true},
				{ID: "discover", Name: "Discover", Description: "Deep codebase analysis", Command: "viki discovery --deep", Dependencies: []string{"init"}},
				{ID: "specify", Name: "Specify", Description: "Define requirements", Command: "viki specify", Agent: "pm", Dependencies: []string{"constitution"}, Required: true},
				{ID: "clarify", Name: "Clarify", Description: "Clarify requirements", Command: "viki clarify", Agent: "business_analyst", Dependencies: []string{"specify"}, Required: true},
				{ID: "checklist-req", Name: "Requirements Checklist", Description: "Validate requirements", Command: "viki checklist", Dependencies: []string{"clarify"}, Required: true},
				{ID: "plan", Name: "Plan", Description: "Design architecture", Command: "viki plan", Agent: "architect", Dependencies: []string{"checklist-req"}, Required: true},
				{ID: "security-review", Name: "Security Review", Description: "Security assessment", Agent: "security", Dependencies: []string{"plan"}, Required: true},
				{ID: "approve-plan", Name: "Approve Plan", Description: "Multi-stakeholder approval", Command: "viki approve", Dependencies: []string{"security-review"}, Required: true},
				{ID: "task", Name: "Break Down Tasks", Description: "Create task list", Command: "viki task", Agent: "developer", Dependencies: []string{"approve-plan"}, Required: true},
				{ID: "checklist-tech", Name: "Technical Checklist", Description: "Technical validation", Command: "viki checklist", Dependencies: []string{"task"}, Required: true},
				{ID: "execute", Name: "Execute", Description: "Implement solution", Command: "viki execute", Agent: "developer", Dependencies: []string{"checklist-tech"}, Required: true},
				{ID: "test", Name: "Test", Description: "Comprehensive testing", Agent: "test_automation", Dependencies: []string{"execute"}, Required: true},
				{ID: "performance", Name: "Performance Test", Description: "Performance validation", Agent: "performance", Dependencies: []string{"execute"}, Parallel: true},
				{ID: "review", Name: "Review", Description: "Quality review", Command: "viki review", Agent: "qa", Dependencies: []string{"test", "performance"}, Required: true},
				{ID: "documentation", Name: "Documentation", Description: "Update documentation", Agent: "documentation", Dependencies: []string{"execute"}, Parallel: true},
				{ID: "approve-final", Name: "Final Approval", Description: "Final sign-off", Command: "viki approve", Dependencies: []string{"review", "documentation"}, Required: true},
			},
		},
	}
}

// WorkflowEngine manages workflow execution
type WorkflowEngine struct {
	workflow *Workflow
	stateDir string
}

// NewWorkflowEngine creates a new workflow engine
func NewWorkflowEngine(workflow *Workflow, stateDir string) *WorkflowEngine {
	return &WorkflowEngine{
		workflow: workflow,
		stateDir: stateDir,
	}
}

// GetReadySteps returns steps that are ready to execute
func (e *WorkflowEngine) GetReadySteps() []*WorkflowStep {
	var ready []*WorkflowStep

	for _, step := range e.workflow.Steps {
		if step.Status != StepStatusPending {
			continue
		}

		// Check dependencies
		depsComplete := true
		for _, depID := range step.Dependencies {
			depStep := e.GetStep(depID)
			if depStep == nil || depStep.Status != StepStatusCompleted {
				depsComplete = false
				break
			}
		}

		if depsComplete {
			ready = append(ready, step)
		}
	}

	return ready
}

// GetStep retrieves a step by ID
func (e *WorkflowEngine) GetStep(id string) *WorkflowStep {
	for _, step := range e.workflow.Steps {
		if step.ID == id {
			return step
		}
	}
	return nil
}

// StartStep marks a step as in progress
func (e *WorkflowEngine) StartStep(id string) error {
	step := e.GetStep(id)
	if step == nil {
		return fmt.Errorf("step not found: %s", id)
	}

	now := time.Now()
	step.Status = StepStatusInProgress
	step.StartedAt = &now

	return e.SaveState()
}

// CompleteStep marks a step as completed
func (e *WorkflowEngine) CompleteStep(id string) error {
	step := e.GetStep(id)
	if step == nil {
		return fmt.Errorf("step not found: %s", id)
	}

	now := time.Now()
	step.Status = StepStatusCompleted
	step.CompletedAt = &now

	return e.SaveState()
}

// FailStep marks a step as failed
func (e *WorkflowEngine) FailStep(id string) error {
	step := e.GetStep(id)
	if step == nil {
		return fmt.Errorf("step not found: %s", id)
	}

	step.Status = StepStatusFailed

	return e.SaveState()
}

// GetProgress returns workflow progress
func (e *WorkflowEngine) GetProgress() (completed, total int) {
	for _, step := range e.workflow.Steps {
		total++
		if step.Status == StepStatusCompleted {
			completed++
		}
	}
	return
}

// IsComplete returns true if workflow is complete
func (e *WorkflowEngine) IsComplete() bool {
	for _, step := range e.workflow.Steps {
		if step.Required && step.Status != StepStatusCompleted {
			return false
		}
	}
	return true
}

// SaveState persists workflow state
func (e *WorkflowEngine) SaveState() error {
	statePath := filepath.Join(e.stateDir, "workflow_state.json")
	// In a full implementation, this would serialize the workflow state
	_ = statePath
	return nil
}

// LoadState loads workflow state
func (e *WorkflowEngine) LoadState() error {
	statePath := filepath.Join(e.stateDir, "workflow_state.json")
	// In a full implementation, this would deserialize the workflow state
	_ = statePath
	return nil
}

// FormatProgress formats workflow progress for display
func (e *WorkflowEngine) FormatProgress() string {
	var sb strings.Builder

	completed, total := e.GetProgress()
	percentage := 0
	if total > 0 {
		percentage = (completed * 100) / total
	}

	sb.WriteString(fmt.Sprintf("\n📊 Workflow: %s\n", e.workflow.Name))
	sb.WriteString(fmt.Sprintf("   Progress: %d/%d (%d%%)\n\n", completed, total, percentage))

	for _, step := range e.workflow.Steps {
		var icon string
		switch step.Status {
		case StepStatusCompleted:
			icon = "✅"
		case StepStatusInProgress:
			icon = "🔄"
		case StepStatusFailed:
			icon = "❌"
		case StepStatusSkipped:
			icon = "⏭️"
		default:
			icon = "⬜"
		}

		sb.WriteString(fmt.Sprintf("   %s %s\n", icon, step.Name))
	}

	return sb.String()
}

// DetectTrack analyzes a project and recommends a workflow track
func DetectTrack(projectDir string) *Track {
	tracks := GetTracks()

	// Simple heuristics for track detection
	_, err := os.Stat(filepath.Join(projectDir, ".viki", "constitution.md"))
	hasConstitution := err == nil

	// Count source files
	fileCount := 0
	filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext == ".go" || ext == ".py" || ext == ".js" || ext == ".ts" {
			fileCount++
		}
		return nil
	})

	// Recommend track based on heuristics
	if hasConstitution || fileCount > 50 {
		return tracks[2] // Enterprise
	}
	if fileCount > 10 {
		return tracks[1] // Standard
	}
	return tracks[0] // Quick
}
