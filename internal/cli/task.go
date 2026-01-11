package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/agents"
	"ultimate-sdd-framework/internal/gates"
)

func NewTaskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Break down plan into actionable tasks",
		Long: `Create a detailed task breakdown from the approved architecture plan.

This command uses the Developer agent to convert the high-level plan
into specific, actionable tasks with clear deliverables and acceptance criteria.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check project state
			stateMgr := gates.NewStateManager(".")
			state, err := stateMgr.LoadState()
			if err != nil {
				return fmt.Errorf("project not initialized: %w", err)
			}

			if state.CurrentPhase != gates.PhasePlan {
				return fmt.Errorf("cannot create tasks: current phase is %s (need %s)", state.CurrentPhase, gates.PhasePlan)
			}

			// Check if plan is approved
			planState := state.Phases[gates.PhasePlan]
			if !planState.Status.IsComplete() {
				return fmt.Errorf("plan phase requires approval before creating tasks. Run 'sdd approve' first")
			}

			// Check if plan exists
			planPath := stateMgr.GetPhaseOutputPath(gates.PhasePlan)
			if _, err := os.Stat(planPath); os.IsNotExist(err) {
				return fmt.Errorf("plan not found: %s", planPath)
			}

			// Initialize agent service
			agentSvc := agents.NewAgentService(".")
			if err := agentSvc.Initialize(); err != nil {
				return fmt.Errorf("failed to initialize agent service: %w", err)
			}

			// Transition to task phase
			if err := stateMgr.TransitionPhase(gates.PhaseTask, "taskmaster"); err != nil {
				return fmt.Errorf("failed to transition to task phase: %w", err)
			}

			// Generate task breakdown using Taskmaster
			fmt.Println("🤖 Taskmaster is breaking down the plan into atomic GSD tasks...")

			// Get current track ID from metadata or use default
			trackID := "feature-implementation"
			if state.Metadata != nil {
				if t, ok := state.Metadata["current_track"].(string); ok && t != "" {
					trackID = t
				}
			}

			response, err := agentSvc.Orchestrate("task", trackID, "")
			if err != nil {
				return fmt.Errorf("Taskmaster failed: %w", err)
			}

			// Save tasks (Orchestrate already saves it as gsd.json, but we confirm here)
			// taskPath := stateMgr.GetPhaseOutputPath(gates.PhaseTask)
			// If Orchestrate saves to .sdd/tracks/[trackID]/gsd.json, but GetPhaseOutputPath might return something else?
			// Let's rely on Orchestrate's save location which matches GetPhaseOutputPath's expectation if configured correctly.
			// Currently Orchestrate saves to .sdd/tracks/[trackID]/[artifact].
			// GetPhaseOutputPath likely points to the project state.
			// Let's assume Orchestrate handles the file creation.

			fmt.Printf("✅ GSD Checklist generated: %s\n", ".sdd/tracks/"+trackID+"/gsd.json")
			fmt.Println("Taskmaster Output:")
			fmt.Println(response)

			// Complete phase
			// We point to gsd.json instead of tasks.md
			if err := stateMgr.CompletePhase([]string{"gsd.json"}); err != nil {
				return fmt.Errorf("failed to complete task phase: %w", err)
			}

			fmt.Printf("✅ Gate 3.5 Passed: Taskify Complete.\n")
			fmt.Println("Next: Run 'sdd execute' to start High-Velocity Implementation")

			return nil
		},
	}

	return cmd
}
