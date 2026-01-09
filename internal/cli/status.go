package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/gates"
)

var (
	statusStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	phaseStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("15"))

	approvedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("46"))

	pendingStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("11"))

	blockedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))
)

func NewStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current project status",
		Long:  "Display the current phase, progress, and status of all project phases.",
		RunE: func(cmd *cobra.Command, args []string) error {
			stateMgr := gates.NewStateManager(".")
			state, err := stateMgr.LoadState()
			if err != nil {
				return fmt.Errorf("failed to load project state: %w", err)
			}

			fmt.Println(statusStyle.Render("🚀 Ultimate SDD Framework Status"))
			fmt.Println(strings.Repeat("=", 50))
			fmt.Printf("Project: %s\n", state.ProjectName)
			fmt.Printf("Current Phase: %s\n", state.CurrentPhase)
			fmt.Printf("Last Updated: %s\n", state.UpdatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println()

			// Show phase status
			fmt.Println("Phase Status:")
			for _, phase := range []gates.Phase{
				gates.PhaseInit,
				gates.PhaseSpecify,
				gates.PhasePlan,
				gates.PhaseTask,
				gates.PhaseExecute,
				gates.PhaseReview,
				gates.PhaseComplete,
			} {
				phaseState := state.Phases[phase]
				fmt.Printf("  %s: ", phaseStyle.Render(string(phase)))

				switch phaseState.Status {
				case gates.StatusApproved:
					fmt.Printf("%s", approvedStyle.Render("✓ APPROVED"))
					if phaseState.CompletedAt != nil {
						fmt.Printf(" (%s)", phaseState.CompletedAt.Format("01-02 15:04"))
					}
				case gates.StatusInProgress:
					fmt.Printf("%s", pendingStyle.Render("⟳ IN PROGRESS"))
					if phaseState.StartedAt != nil {
						fmt.Printf(" (started %s)", phaseState.StartedAt.Format("01-02 15:04"))
					}
				case gates.StatusRejected:
					fmt.Printf("%s", blockedStyle.Render("✗ REJECTED"))
				case gates.StatusBlocked:
					fmt.Printf("%s", blockedStyle.Render("🚫 BLOCKED"))
				default:
					fmt.Printf("%s", pendingStyle.Render("○ PENDING"))
				}

				if phaseState.AgentUsed != "" {
					fmt.Printf(" [%s]", phaseState.AgentUsed)
				}

				if len(phaseState.Approvals) > 0 {
					fmt.Printf(" (approved by %s)", phaseState.Approvals[len(phaseState.Approvals)-1].ApprovedBy)
				}

				fmt.Println()

				// Show output files
				if len(phaseState.OutputFiles) > 0 {
					for _, file := range phaseState.OutputFiles {
						fmt.Printf("    📄 %s\n", file)
					}
				}
			}

			fmt.Println()
			fmt.Println("Next Steps:")
			switch state.CurrentPhase {
			case gates.PhaseInit:
				fmt.Println("  Run: sdd specify \"your feature description\"")
			case gates.PhaseSpecify:
				fmt.Println("  Run: sdd plan")
			case gates.PhasePlan:
				fmt.Println("  Run: sdd approve  # then sdd task")
			case gates.PhaseTask:
				fmt.Println("  Run: sdd execute")
			case gates.PhaseExecute:
				fmt.Println("  Run: sdd review")
			case gates.PhaseReview:
				fmt.Println("  Run: sdd approve  # to complete")
			case gates.PhaseComplete:
				fmt.Println("  🎉 Project complete! Start a new feature with sdd specify")
			}

			return nil
		},
	}

	return cmd
}