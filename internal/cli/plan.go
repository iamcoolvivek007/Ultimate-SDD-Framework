package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/agents"
	"ultimate-sdd-framework/internal/gates"
)

func NewPlanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Create architecture plan using the Designer agent",
		Long: `Generate a detailed system architecture plan based on specifications.

This command uses the Architect agent to design the system components,
technology choices, data flow, and implementation strategy.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check project state
			stateMgr := gates.NewStateManager(".")
			state, err := stateMgr.LoadState()
			if err != nil {
				return fmt.Errorf("project not initialized: %w", err)
			}

			if state.CurrentPhase != gates.PhaseSpecify {
				return fmt.Errorf("cannot plan: current phase is %s (need %s)", state.CurrentPhase, gates.PhaseSpecify)
			}

			// Check if specification exists
			specPath := stateMgr.GetPhaseOutputPath(gates.PhaseSpecify)
			if _, err := os.Stat(specPath); os.IsNotExist(err) {
				return fmt.Errorf("specification not found: %s", specPath)
			}

			// Load specification
			specContent, err := os.ReadFile(specPath)
			if err != nil {
				return fmt.Errorf("failed to read specification: %w", err)
			}

			// Initialize agent service
			agentSvc := agents.NewAgentService(".")
			if err := agentSvc.Initialize(); err != nil {
				return fmt.Errorf("failed to initialize agent service: %w", err)
			}

			// Validate designer agent is available
			_, err = agentSvc.GetAgentForPhase("plan")
			if err != nil {
				return fmt.Errorf("designer agent not available: %w", err)
			}

			// Transition to plan phase
			if err := stateMgr.TransitionPhase(gates.PhasePlan, "designer"); err != nil {
				return fmt.Errorf("failed to transition to plan phase: %w", err)
			}

			// Generate architecture plan
			planContent, err := agentSvc.GetAgentResponse("designer", "plan", string(specContent), "", "")
			if err != nil {
				return fmt.Errorf("failed to generate architecture plan: %w", err)
			}

			// Save plan
			planPath := stateMgr.GetPhaseOutputPath(gates.PhasePlan)
			if err := os.WriteFile(planPath, []byte(planContent), 0644); err != nil {
				return fmt.Errorf("failed to save plan: %w", err)
			}

			// Complete phase
			if err := stateMgr.CompletePhase([]string{filepath.Base(planPath)}); err != nil {
				return fmt.Errorf("failed to complete plan phase: %w", err)
			}

			fmt.Printf("✅ Architecture plan created: %s\n", planPath)
			fmt.Println("Next: Run 'sdd approve' to approve the plan, then 'sdd task' to break it down")

			return nil
		},
	}

	return cmd
}
