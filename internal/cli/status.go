package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/gates"
	"ultimate-sdd-framework/internal/ui"
)

func NewStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Launch Nexus UI Dashboard",
		Long:  "Launch the interactive Nexus UI Dashboard to manage project status and workflow.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize state manager
			stateMgr := gates.NewStateManager(".")

			// Initialize UI model
			model, err := ui.NewSDDModel(stateMgr)
			if err != nil {
				return fmt.Errorf("failed to initialize UI: %w", err)
			}

			// Load available skills
			model.LoadSkills()

			// Run Bubble Tea program
			p := tea.NewProgram(model, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				return fmt.Errorf("failed to run UI: %w", err)
			}

			return nil
		},
	}

	return cmd
}
