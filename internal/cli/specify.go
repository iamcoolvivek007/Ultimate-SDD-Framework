package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/agents"
	"ultimate-sdd-framework/internal/gates"
	"ultimate-sdd-framework/internal/tui"
)

func NewSpecifyCmd() *cobra.Command {
	var useTUI bool

	cmd := &cobra.Command{
		Use:   "specify [description]",
		Short: "Create feature specifications using the PM agent",
		Long: `Generate detailed technical specifications for a feature.

This command uses the Product Manager agent to analyze your request
and create comprehensive specifications including requirements,
constraints, and acceptance criteria.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			description := strings.Join(args, " ")

			// Check project state
			stateMgr := gates.NewStateManager(".")
			state, err := stateMgr.LoadState()
			if err != nil {
				return fmt.Errorf("project not initialized: %w", err)
			}

			if state.CurrentPhase != gates.PhaseInit && state.CurrentPhase != gates.PhaseSpecify {
				return fmt.Errorf("cannot specify: current phase is %s", state.CurrentPhase)
			}

			// Initialize agent service
			agentSvc := agents.NewAgentService(".")
			if err := agentSvc.Initialize(); err != nil {
				return fmt.Errorf("failed to initialize agent service: %w", err)
			}

			// Get PM agent
			pmAgent, err := agentSvc.GetAgentForPhase("specify")
			if err != nil {
				return fmt.Errorf("PM agent not available: %w", err)
			}

			// Transition to specify phase
			if err := stateMgr.TransitionPhase(gates.PhaseSpecify, "pm"); err != nil {
				return fmt.Errorf("failed to transition to specify phase: %w", err)
			}

			// Generate specifications
			if useTUI {
				return tui.RunSpecifyTUI(pmAgent, description, stateMgr)
			}

			// Generate specifications using AI
			specContent, err := agentSvc.GetAgentResponse("pm", "specify", description)
			if err != nil {
				return fmt.Errorf("failed to generate specification: %w", err)
			}

			// Save specification
			specPath := stateMgr.GetPhaseOutputPath(gates.PhaseSpecify)
			if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
				return fmt.Errorf("failed to save specification: %w", err)
			}

			// Complete phase
			if err := stateMgr.CompletePhase([]string{filepath.Base(specPath)}); err != nil {
				return fmt.Errorf("failed to complete specify phase: %w", err)
			}

			fmt.Printf("✅ Specification created: %s\n", specPath)
			fmt.Println("Next: Run 'sdd plan' to design the architecture")

			return nil
		},
	}

	cmd.Flags().BoolVarP(&useTUI, "tui", "t", false, "Use terminal UI for specification creation")

	return cmd
}
