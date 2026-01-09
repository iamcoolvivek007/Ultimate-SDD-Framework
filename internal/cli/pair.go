package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/pair"
)

var (
	agentRole     string
	focusArea     string
	activeFile    string
	cursorLine     int
	contextCode   string
	requestType   string
	suggestionID  string
	userAction    string
)

func NewPairCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pair",
		Short: "Interactive pair programming with AI assistance",
		Long: `Start an interactive pair programming session with AI assistance:
- Real-time code suggestions and completions
- Refactoring recommendations
- Testing strategy guidance
- Best practice enforcement
- Learning from your coding patterns

Supports various interaction modes for different development needs.`,
	}

	// Subcommands
	cmd.AddCommand(NewPairStartCmd())
	cmd.AddCommand(NewPairSuggestCmd())
	cmd.AddCommand(NewPairActionCmd())
	cmd.AddCommand(NewPairEndCmd())
	cmd.AddCommand(NewPairReportCmd())

	return cmd
}

func NewPairStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [agent-role] [focus-area]",
		Short: "Start a new pair programming session",
		Long: `Begin an interactive pair programming session with AI assistance.

Available agent roles:
- developer: General development assistance
- architect: Architecture and design guidance
- qa: Testing and quality assurance
- system: Framework and tooling expertise

Focus areas help tailor suggestions to your current task.`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			// Parse arguments
			if len(args) >= 1 {
				agentRole = args[0]
			} else {
				agentRole = "developer" // Default
			}

			if len(args) >= 2 {
				focusArea = args[1]
			} else {
				focusArea = "general development"
			}

			fmt.Printf("🤝 Starting pair programming session with %s agent\n", agentRole)
			fmt.Printf("🎯 Focus: %s\n", focusArea)

			// Create pair programmer
			pairProgrammer, err := pair.NewPairProgrammer(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize pair programmer: %w", err)
			}

			// Start session
			session, err := pairProgrammer.StartSession(agentRole, focusArea)
			if err != nil {
				return fmt.Errorf("failed to start session: %w", err)
			}

			fmt.Printf("✅ Session started: %s\n", session.ID)
			fmt.Println("\n💡 Available commands:")
			fmt.Println("  nexus pair suggest --file <file> --line <line> --type <completion|refactor|test|explanation>")
			fmt.Println("  nexus pair action --id <suggestion-id> --action <accepted|rejected|modified>")
			fmt.Println("  nexus pair report  # View session summary")
			fmt.Println("  nexus pair end     # End session")

			return nil
		},
	}

	return cmd
}

func NewPairSuggestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "suggest",
		Short: "Get AI suggestions for current code context",
		Long: `Request AI assistance for your current coding context:
- completion: Code completion and continuation
- refactor: Refactoring suggestions and improvements
- test: Testing strategies and test code generation
- explanation: Code explanation and best practice guidance`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			if activeFile == "" {
				return fmt.Errorf("file must be specified with --file flag")
			}

			if requestType == "" {
				requestType = "completion" // Default
			}

			// Get code context (placeholder - would read from file)
			if contextCode == "" {
				// Read some context from the file
				contextCode = fmt.Sprintf("// Context from %s around line %d", activeFile, cursorLine)
			}

			fmt.Printf("🧠 Getting %s suggestion for %s:%d...\n", requestType, activeFile, cursorLine)

			// Create pair programmer
			pairProgrammer, err := pair.NewPairProgrammer(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize pair programmer: %w", err)
			}

			// Get suggestion
			suggestion, err := pairProgrammer.GetSuggestion(activeFile, cursorLine, contextCode, requestType)
			if err != nil {
				return fmt.Errorf("failed to get suggestion: %w", err)
			}

			// Display suggestion
			fmt.Printf("💡 Suggestion ID: %s\n", suggestion.ID)
			fmt.Printf("🎯 Type: %s (%.1f%% confidence)\n", suggestion.Type, suggestion.Confidence*100)
			fmt.Println("📝 Suggestion:")
			fmt.Println(suggestion.Content)

			if suggestion.Explanation != "" {
				fmt.Println("\n💭 Explanation:")
				fmt.Println(suggestion.Explanation)
			}

			if len(suggestion.Alternatives) > 0 {
				fmt.Println("\n🔄 Alternatives:")
				for i, alt := range suggestion.Alternatives {
					fmt.Printf("  %d. %s\n", i+1, alt)
				}
			}

			fmt.Printf("\n✅ Use: nexus pair action --id %s --action <accepted|rejected|modified>\n", suggestion.ID)

			return nil
		},
	}

	cmd.Flags().StringVar(&activeFile, "file", "", "File path for context")
	cmd.Flags().IntVar(&cursorLine, "line", 1, "Cursor line number")
	cmd.Flags().StringVar(&contextCode, "context", "", "Code context (auto-detected if not provided)")
	cmd.Flags().StringVar(&requestType, "type", "completion", "Suggestion type: completion, refactor, test, explanation")

	return cmd
}

func NewPairActionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "action",
		Short: "Record your response to a suggestion",
		Long: `Record how you responded to an AI suggestion:
- accepted: Used the suggestion as-is
- rejected: Ignored the suggestion
- modified: Adapted the suggestion

This helps improve future suggestions based on your preferences.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if suggestionID == "" {
				return fmt.Errorf("suggestion ID must be specified with --id flag")
			}

			if userAction == "" {
				return fmt.Errorf("action must be specified with --action flag")
			}

			validActions := []string{"accepted", "rejected", "modified"}
			actionValid := false
			for _, validAction := range validActions {
				if userAction == validAction {
					actionValid = true
					break
				}
			}

			if !actionValid {
				return fmt.Errorf("invalid action. Must be one of: %s", strings.Join(validActions, ", "))
			}

			projectRoot := "."

			// Create pair programmer
			pairProgrammer, err := pair.NewPairProgrammer(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize pair programmer: %w", err)
			}

			// Record action
			if err := pairProgrammer.RecordUserAction(userAction, suggestionID); err != nil {
				return fmt.Errorf("failed to record action: %w", err)
			}

			fmt.Printf("✅ Recorded: %s suggestion %s\n", userAction, suggestionID)
			fmt.Println("📊 This feedback improves future AI suggestions!")

			return nil
		},
	}

	cmd.Flags().StringVar(&suggestionID, "id", "", "Suggestion ID to respond to")
	cmd.Flags().StringVar(&userAction, "action", "", "Your action: accepted, rejected, modified")

	return cmd
}

func NewPairEndCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "end",
		Short: "End the current pair programming session",
		Long:  "Conclude the active pair programming session and generate a summary report.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			fmt.Println("🏁 Ending pair programming session...")

			// Create pair programmer
			pairProgrammer, err := pair.NewPairProgrammer(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize pair programmer: %w", err)
			}

			// End session
			session, err := pairProgrammer.EndSession()
			if err != nil {
				return fmt.Errorf("failed to end session: %w", err)
			}

			// Generate and display report
			report := pairProgrammer.GetSessionReport(session)
			fmt.Println(report)

			// Save report
			reportPath := ".sdd/pair_session_report.md"
			if err := os.WriteFile(reportPath, []byte(report), 0644); err != nil {
				fmt.Printf("Warning: Failed to save session report: %v\n", err)
			} else {
				fmt.Printf("📄 Session report saved to: %s\n", reportPath)
			}

			fmt.Printf("👋 Session ended after %v\n", session.Stats.TimeSpent.Round(0))

			return nil
		},
	}

	return cmd
}

func NewPairReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "View current session report",
		Long:  "Display statistics and insights from the current pair programming session.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			// Create pair programmer
			pairProgrammer, err := pair.NewPairProgrammer(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize pair programmer: %w", err)
			}

			session := pairProgrammer.GetActiveSession()
			if session == nil {
				fmt.Println("No active pair programming session.")
				fmt.Println("Start one with: nexus pair start [agent-role] [focus-area]")
				return nil
			}

			// Generate current report
			report := pairProgrammer.GetSessionReport(session)
			fmt.Println(report)

			return nil
		},
	}

	return cmd
}