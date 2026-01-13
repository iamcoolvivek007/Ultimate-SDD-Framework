package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/learning"
)

var (
	interactionType string
	contextInfo     string
	actionTaken     string
	outcomeDesc     string
	successFlag     bool
	durationMs      int
)

func NewLearnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "learn",
		Short: "Learning and adaptation management",
		Long: `Manage the adaptive learning system:
- Record development interactions and outcomes
- View personalized suggestions based on learning
- Track coding patterns and preferences
- Monitor learning evolution over time
- Generate insights from development history`,
	}

	// Subcommands
	cmd.AddCommand(NewLearnRecordCmd())
	cmd.AddCommand(NewLearnSuggestCmd())
	cmd.AddCommand(NewLearnReportCmd())
	cmd.AddCommand(NewLearnEvolveCmd())

	return cmd
}

func NewLearnRecordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "Record a development interaction",
		Long: `Record a development interaction for learning:
- Track successful patterns and approaches
- Learn from failures and mistakes
- Build personalized coding preferences
- Improve future AI suggestions`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			if interactionType == "" || contextInfo == "" || actionTaken == "" {
				return fmt.Errorf("interaction type, context, and action are required")
			}

			if outcomeDesc == "" {
				outcomeDesc = "Completed successfully"
			}

			if durationMs == 0 {
				durationMs = 1000 // Default 1 second
			}

			fmt.Printf("🧠 Recording interaction: %s\n", interactionType)
			fmt.Printf("📝 Context: %s\n", contextInfo)
			fmt.Printf("🎯 Action: %s\n", actionTaken)
			fmt.Printf("📊 Outcome: %s (%v)\n", outcomeDesc, successFlag)

			// Create adaptive learner
			learner, err := learning.NewAdaptiveLearner(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize learner: %w", err)
			}

			// Record interaction
			if err := learner.LearnFromInteraction(interactionType, contextInfo, actionTaken, outcomeDesc, successFlag, durationMs); err != nil {
				return fmt.Errorf("failed to record interaction: %w", err)
			}

			fmt.Println("✅ Interaction recorded and learning updated!")

			// Show current learning summary briefly
			summary := learner.GetLearningSummary()
			lines := strings.Split(summary, "\n")
			if len(lines) > 10 {
				fmt.Println("📊 Current Learning Status:")
				for i := 0; i < 10 && i < len(lines); i++ {
					if strings.TrimSpace(lines[i]) != "" {
						fmt.Printf("  %s\n", lines[i])
					}
				}
				fmt.Println("  ... (use 'nexus learn report' for full summary)")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&interactionType, "type", "", "Interaction type (e.g., refactoring, debugging, testing)")
	cmd.Flags().StringVar(&contextInfo, "context", "", "Context of the interaction")
	cmd.Flags().StringVar(&actionTaken, "action", "", "Action taken")
	cmd.Flags().StringVar(&outcomeDesc, "outcome", "Completed successfully", "Outcome description")
	cmd.Flags().BoolVar(&successFlag, "success", true, "Whether the interaction was successful")
	cmd.Flags().IntVar(&durationMs, "duration", 1000, "Duration in milliseconds")

	return cmd
}

func NewLearnSuggestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "suggest [context] [task-type]",
		Short: "Get personalized suggestions",
		Long: `Receive personalized suggestions based on your learning history:
- Context-aware recommendations
- Pattern-based suggestions
- Preference-driven guidance
- Learning from past successes and failures`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			context := "general"
			taskType := "development"

			if len(args) >= 1 {
				context = args[0]
			}
			if len(args) >= 2 {
				taskType = args[1]
			}

			fmt.Printf("🎯 Getting personalized suggestions for: %s (%s)\n", context, taskType)

			// Create adaptive learner
			learner, err := learning.NewAdaptiveLearner(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize learner: %w", err)
			}

			// Get suggestions
			suggestions, err := learner.GetPersonalizedSuggestions(context, taskType)
			if err != nil {
				return fmt.Errorf("failed to get suggestions: %w", err)
			}

			if len(suggestions) == 0 {
				fmt.Println("🤔 No specific suggestions available yet.")
				fmt.Println("  Start using the framework to build your personalized profile!")
				fmt.Println("  Try: nexus learn record --type \"refactoring\" --context \"api\" --action \"extract method\" --success")
				return nil
			}

			fmt.Printf("💡 Found %d personalized suggestions:\n", len(suggestions))

			for i, suggestion := range suggestions {
				fmt.Printf("\n%d. **%s** (%.1f%% confidence)\n", i+1, suggestion.Title, suggestion.Confidence*100)
				fmt.Printf("   💭 %s\n", suggestion.Description)
				fmt.Printf("   📚 Reason: %s\n", suggestion.Reason)

				if len(suggestion.Examples) > 0 {
					fmt.Println("   📝 Examples:")
					for _, example := range suggestion.Examples {
						fmt.Printf("     • %s\n", example)
					}
				}
			}

			fmt.Println("\n🔄 Suggestions improve as you use the framework more!")

			return nil
		},
	}

	return cmd
}

func NewLearnReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "View learning summary and insights",
		Long:  "Display comprehensive learning report with patterns, preferences, and evolution insights.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			fmt.Println("🧠 Generating learning report...")

			// Create adaptive learner
			learner, err := learning.NewAdaptiveLearner(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize learner: %w", err)
			}

			// Get learning summary
			report := learner.GetLearningSummary()

			// Display report
			fmt.Println(report)

			// Save detailed report
			reportPath := ".sdd/learning_report.md"
			if err := os.WriteFile(reportPath, []byte(report), 0644); err != nil {
				fmt.Printf("Warning: Failed to save learning report: %v\n", err)
			} else {
				fmt.Printf("📄 Learning report saved to: %s\n", reportPath)
			}

			return nil
		},
	}

	return cmd
}

func NewLearnEvolveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "evolve",
		Short: "Evolve rules based on learning",
		Long: `Analyze learning data and suggest rule improvements:
- Identify patterns that should become rules
- Suggest modifications to existing rules
- Propose new best practices based on successes
- Recommend rule updates based on failure patterns`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			fmt.Println("🔄 Analyzing learning data for rule evolution...")

			// Create adaptive learner
			learner, err := learning.NewAdaptiveLearner(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize learner: %w", err)
			}

			// Get rule evolution suggestions
			suggestions, err := learner.EvolveRules()
			if err != nil {
				return fmt.Errorf("failed to analyze rule evolution: %w", err)
			}

			if len(suggestions) == 0 {
				fmt.Println("📊 No rule evolution suggestions available yet.")
				fmt.Println("  Continue using the framework to generate more learning data!")
				return nil
			}

			fmt.Printf("💡 Found %d rule evolution suggestions:\n", len(suggestions))

			for i, suggestion := range suggestions {
				fmt.Printf("\n%d. **Rule Evolution Opportunity** (%.1f%% confidence)\n", i+1, suggestion.Confidence*100)
				fmt.Printf("   📝 Reason: %s\n", suggestion.Reason)

				if suggestion.CurrentRule != "" {
					fmt.Printf("   📋 Current Rule: %s\n", suggestion.CurrentRule)
				}

				fmt.Printf("   ✨ Suggested Rule: %s\n", suggestion.SuggestedRule)

				if suggestion.Evidence != "" {
					fmt.Printf("   📊 Evidence: %s\n", suggestion.Evidence)
				}

				if suggestion.Mitigation != "" {
					fmt.Printf("   🛠️  Mitigation: %s\n", suggestion.Mitigation)
				}
			}

			fmt.Println("\n📝 To implement these suggestions:")
			fmt.Println("  1. Review each suggestion carefully")
			fmt.Println("  2. Update .sdd/rules/ files accordingly")
			fmt.Println("  3. Run: nexus evolve \"rule update description\"")
			fmt.Println("  4. The framework will learn from your rule updates!")

			return nil
		},
	}

	return cmd
}
