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
		Short: "📝 Tell Viki what you want to build",
		Long: `💭 Time to describe your idea!

Tell Viki what you want to create. Be as clear as you can about:
• What the app/website should do
• Who will use it
• Any specific features you want

Viki will think through your idea and create a clear plan that both humans and computers can understand.

Examples:
• "Build a todo list app where users can add, edit, and delete tasks"
• "Create a simple blog with posts and comments"
• "Make a weather app that shows the forecast for my city"

Don't worry about technical details - just describe what you want! ✨`,
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
				return fmt.Errorf("🤖 Oops! Viki's AI assistants aren't ready. Try running 'viki init' first: %w", err)
			}

			// Get PM agent
			pmAgent, err := agentSvc.GetAgentForPhase("specify")
			if err != nil {
				return fmt.Errorf("PM agent not available: %w", err)
			}

			// Check current phase and only transition if not already in specify phase
			state, err = stateMgr.LoadState()
			if err != nil {
				return fmt.Errorf("📁 Can't find your project info. Did you run 'viki init' first?: %w", err)
			}

			if state.CurrentPhase != gates.PhaseSpecify {
				// Only transition if not already in specify phase
				if err := stateMgr.TransitionPhase(gates.PhaseSpecify, "strategist"); err != nil {
					return fmt.Errorf("failed to transition to specify phase: %w", err)
				}
			}

			// Generate specifications
			if useTUI {
				return tui.RunSpecifyTUI(pmAgent, description, stateMgr)
			}

			// Generate specifications using AI
			specContent, err := agentSvc.GetAgentResponse("strategist", "specify", description, "", "")
			if err != nil {
				return fmt.Errorf("🤔 Viki had trouble understanding your request. Try rephrasing it or check your AI provider setup: %w", err)
			}

			// Add Status to specification
			specContentWithStatus := fmt.Sprintf("---\nstatus: pending\n---\n\n%s", specContent)

			// Save specification
			specPath := stateMgr.GetPhaseOutputPath(gates.PhaseSpecify)
			if err := os.WriteFile(specPath, []byte(specContentWithStatus), 0644); err != nil {
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
