package gates

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/goccy/go-yaml"
)

// StateManager handles project state persistence and transitions
type StateManager struct {
	projectRoot string
	sddDir      string
}

// NewStateManager creates a new state manager for the given project root
func NewStateManager(projectRoot string) *StateManager {
	return &StateManager{
		projectRoot: projectRoot,
		sddDir:      filepath.Join(projectRoot, ".sdd"),
	}
}

// GetProjectRoot returns the project root directory
func (sm *StateManager) GetProjectRoot() string {
	return sm.projectRoot
}

// InitializeProject creates the initial .sdd directory and state file
func (sm *StateManager) InitializeProject(projectName string) error {
	// Create .sdd directory
	if err := os.MkdirAll(sm.sddDir, 0755); err != nil {
		return fmt.Errorf("failed to create .sdd directory: %w", err)
	}

	// Create initial state
	now := time.Now()
	state := &ProjectState{
		ProjectID:    generateProjectID(),
		ProjectName:  projectName,
		CreatedAt:    now,
		UpdatedAt:    now,
		CurrentPhase: PhaseInit,
		Phases: map[Phase]PhaseState{
			PhaseInit: {
				Phase:       PhaseInit,
				Status:      StatusApproved, // Init is always approved
				CompletedAt: &now,
			},
			PhaseSpecify: {
				Phase:  PhaseSpecify,
				Status: StatusPending,
			},
			PhasePlan: {
				Phase:  PhasePlan,
				Status: StatusPending,
			},
			PhaseTask: {
				Phase:  PhaseTask,
				Status: StatusPending,
			},
			PhaseExecute: {
				Phase:  PhaseExecute,
				Status: StatusPending,
			},
			PhaseReview: {
				Phase:  PhaseReview,
				Status: StatusPending,
			},
			PhaseComplete: {
				Phase:  PhaseComplete,
				Status: StatusPending,
			},
		},
	}

	// Save state
	return sm.saveState(state)
}

// LoadState loads the current project state
func (sm *StateManager) LoadState() (*ProjectState, error) {
	stateFile := filepath.Join(sm.sddDir, "state.yaml")

	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("project not initialized, run 'sdd init' first")
	}

	data, err := os.ReadFile(stateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state ProjectState
	if err := yaml.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	return &state, nil
}

// TransitionPhase attempts to transition to a new phase
func (sm *StateManager) TransitionPhase(targetPhase Phase, agentUsed string) error {
	state, err := sm.LoadState()
	if err != nil {
		return err
	}

	currentPhase := state.CurrentPhase

	// Check if transition is valid
	if !CanTransition(currentPhase, targetPhase) {
		return fmt.Errorf("invalid transition from %s to %s", currentPhase, targetPhase)
	}

	// Check if approval is required
	if RequiresApproval(currentPhase, targetPhase) {
		currentPhaseState := state.Phases[currentPhase]
		if !currentPhaseState.Status.IsComplete() {
			return fmt.Errorf("cannot transition to %s: %s phase requires approval", targetPhase, currentPhase)
		}
	}

	// Update the target phase
	now := time.Now()
	targetPhaseState := state.Phases[targetPhase]
	targetPhaseState.Status = StatusInProgress
	targetPhaseState.StartedAt = &now
	targetPhaseState.AgentUsed = agentUsed
	state.Phases[targetPhase] = targetPhaseState

	// Update current phase if moving forward
	if targetPhase != currentPhase {
		state.CurrentPhase = targetPhase
	}

	state.UpdatedAt = now

	return sm.saveState(state)
}

// ApprovePhase approves the current phase for transition
func (sm *StateManager) ApprovePhase(approvedBy, comments string) error {
	state, err := sm.LoadState()
	if err != nil {
		return err
	}

	currentPhase := state.CurrentPhase
	phaseState := state.Phases[currentPhase]

	// Add approval
	approval := Approval{
		ApprovedBy: approvedBy,
		ApprovedAt: time.Now(),
		Comments:   comments,
	}

	phaseState.Approvals = append(phaseState.Approvals, approval)
	phaseState.Status = StatusApproved
	phaseState.CompletedAt = &approval.ApprovedAt

	state.Phases[currentPhase] = phaseState
	state.UpdatedAt = approval.ApprovedAt

	return sm.saveState(state)
}

// CompletePhase marks the current phase as completed
func (sm *StateManager) CompletePhase(outputFiles []string) error {
	state, err := sm.LoadState()
	if err != nil {
		return err
	}

	currentPhase := state.CurrentPhase
	phaseState := state.Phases[currentPhase]

	now := time.Now()
	phaseState.Status = StatusApproved
	phaseState.CompletedAt = &now
	phaseState.OutputFiles = append(phaseState.OutputFiles, outputFiles...)

	state.Phases[currentPhase] = phaseState
	state.UpdatedAt = now

	return sm.saveState(state)
}

// GetPhaseOutputPath returns the expected output path for a phase
func (sm *StateManager) GetPhaseOutputPath(phase Phase) string {
	var filename string
	switch phase {
	case PhaseSpecify:
		filename = "spec.md"
	case PhasePlan:
		filename = "plan.md"
	case PhaseTask:
		filename = "tasks.md"
	case PhaseExecute:
		filename = "implementation.md"
	case PhaseReview:
		filename = "review.md"
	default:
		filename = fmt.Sprintf("%s.md", phase)
	}
	return filepath.Join(sm.sddDir, filename)
}

// saveState saves the project state to disk
func (sm *StateManager) saveState(state *ProjectState) error {
	stateFile := filepath.Join(sm.sddDir, "state.yaml")

	data, err := yaml.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(stateFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// generateProjectID generates a unique project ID
func generateProjectID() string {
	return fmt.Sprintf("sdd_%d", time.Now().Unix())
}
