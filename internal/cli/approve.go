package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/gates"
)

func NewApproveCmd() *cobra.Command {
	var comments string

	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Approve the current phase to proceed",
		Long: `Approve the current phase for transition to the next phase.

Some phases require explicit approval before proceeding:
- Plan phase must be approved before creating tasks
- Review phase must be approved to complete the feature`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check project state
			stateMgr := gates.NewStateManager(".")
			state, err := stateMgr.LoadState()
			if err != nil {
				return fmt.Errorf("project not initialized: %w", err)
			}

			currentPhase := state.CurrentPhase
			phaseState := state.Phases[currentPhase]

			// Check if approval is needed
			if phaseState.Status.IsComplete() {
				fmt.Printf("Phase %s is already approved.\n", currentPhase)
				return nil
			}

			// Get approval comments if not provided
			if comments == "" {
				fmt.Printf("Approving phase: %s\n", currentPhase)
				fmt.Print("Comments (optional): ")

				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					comments = strings.TrimSpace(scanner.Text())
				}
			}

			// Get approver name (in real implementation, this would be from git config or auth)
			approver := "developer" // Default for demo
			if user := os.Getenv("USER"); user != "" {
				approver = user
			}

			// Approve the phase
			if err := stateMgr.ApprovePhase(approver, comments); err != nil {
				return fmt.Errorf("failed to approve phase: %w", err)
			}

			fmt.Printf("✅ Phase %s approved by %s\n", currentPhase, approver)
			if comments != "" {
				fmt.Printf("Comments: %s\n", comments)
			}

			// Show next steps
			nextPhase := currentPhase.NextPhase()
			if nextPhase != currentPhase {
				fmt.Printf("\nNext: Run 'sdd %s' to proceed\n", nextPhase)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&comments, "comments", "c", "", "Approval comments")

	return cmd
}
