package gates

import (
	"time"
)

// Phase represents the current phase in the SDD workflow
type Phase string

const (
	PhaseInit     Phase = "init"     // Project initialized
	PhaseSpecify  Phase = "specify"  // Requirements gathered
	PhasePlan     Phase = "plan"     // Architecture designed
	PhaseTask     Phase = "task"     // Tasks broken down
	PhaseExecute  Phase = "execute"  // Implementation in progress
	PhaseReview   Phase = "review"   // QA review
	PhaseComplete Phase = "complete" // Feature complete
)

// Status represents the approval status of a phase
type Status string

const (
	StatusPending   Status = "pending"   // Not yet started
	StatusInProgress Status = "in_progress" // Currently working
	StatusApproved  Status = "approved"  // Approved to proceed
	StatusRejected  Status = "rejected"  // Needs revision
	StatusBlocked   Status = "blocked"   // Cannot proceed
)

// Approval represents an approval record
type Approval struct {
	ApprovedBy string    `yaml:"approved_by"`
	ApprovedAt time.Time `yaml:"approved_at"`
	Comments   string    `yaml:"comments,omitempty"`
}

// PhaseState represents the state of a single phase
type PhaseState struct {
	Phase       Phase     `yaml:"phase"`
	Status      Status    `yaml:"status"`
	StartedAt   *time.Time `yaml:"started_at,omitempty"`
	CompletedAt *time.Time `yaml:"completed_at,omitempty"`
	Approvals   []Approval `yaml:"approvals,omitempty"`
	AgentUsed   string     `yaml:"agent_used,omitempty"`
	OutputFiles []string   `yaml:"output_files,omitempty"`
}

// ProjectState represents the overall state of an SDD project
type ProjectState struct {
	ProjectID   string                 `yaml:"project_id"`
	ProjectName string                 `yaml:"project_name"`
	CreatedAt   time.Time              `yaml:"created_at"`
	UpdatedAt   time.Time              `yaml:"updated_at"`
	CurrentPhase Phase                 `yaml:"current_phase"`
	Phases       map[Phase]PhaseState  `yaml:"phases"`
	Metadata     map[string]interface{} `yaml:"metadata,omitempty"`
}

// PhaseTransition represents a valid phase transition
type PhaseTransition struct {
	From Phase
	To   Phase
	RequiresApproval bool
}

// ValidTransitions defines the allowed phase transitions
var ValidTransitions = []PhaseTransition{
	{From: PhaseInit, To: PhaseSpecify, RequiresApproval: false},
	{From: PhaseSpecify, To: PhasePlan, RequiresApproval: false},
	{From: PhasePlan, To: PhaseTask, RequiresApproval: true}, // Requires approval before task breakdown
	{From: PhaseTask, To: PhaseExecute, RequiresApproval: false},
	{From: PhaseExecute, To: PhaseReview, RequiresApproval: false},
	{From: PhaseReview, To: PhaseComplete, RequiresApproval: false},
	// Allow revisions
	{From: PhasePlan, To: PhaseSpecify, RequiresApproval: false},
	{From: PhaseTask, To: PhasePlan, RequiresApproval: false},
	{From: PhaseExecute, To: PhaseTask, RequiresApproval: false},
	{From: PhaseReview, To: PhaseExecute, RequiresApproval: false},
}

// CanTransition checks if a transition from one phase to another is valid
func CanTransition(from, to Phase) bool {
	for _, transition := range ValidTransitions {
		if transition.From == from && transition.To == to {
			return true
		}
	}
	return false
}

// RequiresApproval checks if a transition requires approval
func RequiresApproval(from, to Phase) bool {
	for _, transition := range ValidTransitions {
		if transition.From == from && transition.To == to {
			return transition.RequiresApproval
		}
	}
	return false
}

// NextPhase returns the next logical phase after the current one
func (p Phase) NextPhase() Phase {
	switch p {
	case PhaseInit:
		return PhaseSpecify
	case PhaseSpecify:
		return PhasePlan
	case PhasePlan:
		return PhaseTask
	case PhaseTask:
		return PhaseExecute
	case PhaseExecute:
		return PhaseReview
	case PhaseReview:
		return PhaseComplete
	default:
		return p
	}
}

// IsComplete returns true if the phase is in a completed state
func (s Status) IsComplete() bool {
	return s == StatusApproved
}

// IsBlocked returns true if the phase cannot proceed
func (s Status) IsBlocked() bool {
	return s == StatusRejected || s == StatusBlocked
}