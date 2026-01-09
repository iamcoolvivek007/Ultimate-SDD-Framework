package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/agents"
	"ultimate-sdd-framework/internal/gates"
)

func NewEvolveCmd() *cobra.Command {
	var (
		bugDescription string
		ruleCategory   string
		autoUpdate     bool
	)

	cmd := &cobra.Command{
		Use:   "evolve [bug-description]",
		Short: "Analyze bug and evolve system rules to prevent recurrence",
		Long: `Analyze a bug report and update the system rules to prevent similar issues.

This command implements the "System Evolution" philosophy - every bug becomes
a learning opportunity that improves the development system permanently.

Examples:
  nexus evolve "User registration fails with null pointer exception"
  nexus evolve --category frontend "Login form doesn't validate email format"
  nexus evolve --auto-update "Database connection timeout causes app crash"`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get bug description from args or flag
			if len(args) > 0 {
				bugDescription = args[0]
			}

			if bugDescription == "" {
				return fmt.Errorf("bug description is required (use --help for examples)")
			}

			// Check project state
			stateMgr := gates.NewStateManager(".")
			_, err := stateMgr.LoadState()
			if err != nil {
				return fmt.Errorf("project not initialized: %w", err)
			}

			// Initialize agent service
			agentSvc := agents.NewAgentService(".")
			if err := agentSvc.Initialize(); err != nil {
				return fmt.Errorf("failed to initialize agent service: %w", err)
			}

			fmt.Printf("🔄 Analyzing bug: %s\n", bugDescription)

			// Perform bug analysis and rule evolution
			evolution := analyzeBugAndEvolve(bugDescription, ruleCategory, agentSvc)

			// Display analysis results
			displayEvolutionResults(evolution)

			// Auto-update rules if requested
			if autoUpdate {
				if err := applyRuleEvolution(evolution); err != nil {
					return fmt.Errorf("failed to apply rule evolution: %w", err)
				}
				fmt.Println("✅ Rules automatically updated!")
			} else {
				fmt.Println("\n💡 To apply these rule updates, run with --auto-update")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&ruleCategory, "category", "c", "", "Rule category to update (global, frontend, backend, api)")
	cmd.Flags().BoolVar(&autoUpdate, "auto-update", false, "Automatically apply rule updates")
	cmd.Flags().StringVarP(&bugDescription, "bug", "b", "", "Bug description (alternative to positional argument)")

	return cmd
}

type RuleEvolution struct {
	BugAnalysis    BugAnalysis
	RuleUpdates    []RuleUpdate
	PreventionRules []PreventionRule
	TestCases      []string
}

type BugAnalysis struct {
	Description   string
	RootCause     string
	Category      string
	Severity      string
	Reproducibility string
	Impact        string
}

type RuleUpdate struct {
	RuleFile    string
	Section     string
	NewRule     string
	Explanation string
}

type PreventionRule struct {
	Pattern     string
	Prevention  string
	CodeExample string
}

func analyzeBugAndEvolve(bugDescription, category string, agentSvc *agents.AgentService) *RuleEvolution {
	// This would normally use the system evolution agent
	// For now, we'll create a structured analysis

	evolution := &RuleEvolution{
		BugAnalysis: BugAnalysis{
			Description:     bugDescription,
			RootCause:       analyzeRootCause(bugDescription),
			Category:        determineCategory(bugDescription, category),
			Severity:        assessSeverity(bugDescription),
			Reproducibility: assessReproducibility(bugDescription),
			Impact:          assessImpact(bugDescription),
		},
	}

	// Generate rule updates based on analysis
	evolution.RuleUpdates = generateRuleUpdates(evolution.BugAnalysis)
	evolution.PreventionRules = generatePreventionRules(evolution.BugAnalysis)
	evolution.TestCases = generateTestCases(evolution.BugAnalysis)

	return evolution
}

func analyzeRootCause(bugDescription string) string {
	desc := strings.ToLower(bugDescription)

	// Analyze common root causes
	if strings.Contains(desc, "null") || strings.Contains(desc, "undefined") {
		return "Input validation gap - missing null/undefined checks"
	}
	if strings.Contains(desc, "concurrent") || strings.Contains(desc, "race") {
		return "Concurrency issue - shared state not properly synchronized"
	}
	if strings.Contains(desc, "timeout") || strings.Contains(desc, "connection") {
		return "Resource management issue - improper error handling for external dependencies"
	}
	if strings.Contains(desc, "validation") || strings.Contains(desc, "format") {
		return "Input validation rule gap - insufficient client/server-side validation"
	}
	if strings.Contains(desc, "memory") || strings.Contains(desc, "leak") {
		return "Resource leak - not properly closing/disposing resources"
	}

	return "General error handling - insufficient error boundaries or recovery mechanisms"
}

func determineCategory(bugDescription, suggestedCategory string) string {
	if suggestedCategory != "" {
		return suggestedCategory
	}

	desc := strings.ToLower(bugDescription)

	if strings.Contains(desc, "ui") || strings.Contains(desc, "component") || strings.Contains(desc, "render") || strings.Contains(desc, "form") {
		return "frontend"
	}
	if strings.Contains(desc, "api") || strings.Contains(desc, "endpoint") || strings.Contains(desc, "request") || strings.Contains(desc, "response") {
		return "api"
	}
	if strings.Contains(desc, "database") || strings.Contains(desc, "query") || strings.Contains(desc, "connection") || strings.Contains(desc, "server") {
		return "backend"
	}

	return "global"
}

func assessSeverity(bugDescription string) string {
	desc := strings.ToLower(bugDescription)

	if strings.Contains(desc, "crash") || strings.Contains(desc, "break") || strings.Contains(desc, "fail") {
		return "High"
	}
	if strings.Contains(desc, "error") || strings.Contains(desc, "exception") {
		return "Medium"
	}

	return "Low"
}

func assessReproducibility(bugDescription string) string {
	desc := strings.ToLower(bugDescription)

	if strings.Contains(desc, "always") || strings.Contains(desc, "every time") {
		return "Always"
	}
	if strings.Contains(desc, "sometimes") || strings.Contains(desc, "intermittent") {
		return "Intermittent"
	}

	return "Rare"
}

func assessImpact(bugDescription string) string {
	desc := strings.ToLower(bugDescription)

	if strings.Contains(desc, "all users") || strings.Contains(desc, "system") {
		return "System-wide"
	}
	if strings.Contains(desc, "user") || strings.Contains(desc, "function") {
		return "Feature-specific"
	}

	return "Minor"
}

func generateRuleUpdates(analysis BugAnalysis) []RuleUpdate {
	var updates []RuleUpdate

	ruleFile := fmt.Sprintf(".sdd/rules/%s.md", analysis.Category)

	switch analysis.RootCause {
	case "Input validation gap - missing null/undefined checks":
		updates = append(updates, RuleUpdate{
			RuleFile:    ruleFile,
			Section:     "Input Validation",
			NewRule:     "### Null Safety\n**Bug Pattern**: Missing null/undefined checks cause runtime failures\n**Prevention**: Always validate input parameters for null/undefined values",
			Explanation: "Adds null safety checks to prevent runtime errors",
		})

	case "Concurrency issue - shared state not properly synchronized":
		updates = append(updates, RuleUpdate{
			RuleFile:    ruleFile,
			Section:     "Concurrency",
			NewRule:     "### Race Condition Prevention\n**Bug Pattern**: Concurrent access to shared state causes data corruption\n**Prevention**: Use proper synchronization primitives for shared resources",
			Explanation: "Adds concurrency safety rules to prevent race conditions",
		})

	case "Input validation rule gap - insufficient client/server-side validation":
		updates = append(updates, RuleUpdate{
			RuleFile:    ruleFile,
			Section:     "Validation",
			NewRule:     "### Input Validation\n**Bug Pattern**: Malformed input causes system failures\n**Prevention**: Validate all inputs at both client and server boundaries",
			Explanation: "Strengthens input validation requirements",
		})
	}

	return updates
}

func generatePreventionRules(analysis BugAnalysis) []PreventionRule {
	var rules []PreventionRule

	switch analysis.RootCause {
	case "Input validation gap - missing null/undefined checks":
		rules = append(rules, PreventionRule{
			Pattern:     "Null pointer exceptions",
			Prevention:  "Add null checks for all input parameters",
			CodeExample: "// JavaScript/TypeScript\nif (!input || input.trim() === '') {\n  throw new Error('Input is required');\n}\n\n// Go\nif input == nil {\n  return errors.New(\"input cannot be nil\")\n}",
		})

	case "Concurrency issue - shared state not properly synchronized":
		rules = append(rules, PreventionRule{
			Pattern:     "Race conditions",
			Prevention:  "Use atomic operations or mutexes for shared state",
			CodeExample: "// Go\natomic.AddInt64(&counter, 1)\n\n// JavaScript\n// Use immutable state updates\nconst [count, setCount] = useState(0);\nsetCount(prev => prev + 1);",
		})
	}

	return rules
}

func generateTestCases(analysis BugAnalysis) []string {
	var tests []string

	switch analysis.RootCause {
	case "Input validation gap - missing null/undefined checks":
		tests = append(tests,
			"Test with null input parameters",
			"Test with undefined/empty inputs",
			"Test with malformed data structures",
		)

	case "Concurrency issue - shared state not properly synchronized":
		tests = append(tests,
			"Test concurrent access to shared resources",
			"Test race conditions with multiple goroutines/threads",
			"Test atomic operations under load",
		)
	}

	return tests
}

func displayEvolutionResults(evolution *RuleEvolution) {
	fmt.Println("\n🔍 Bug Analysis Results")
	fmt.Println("=======================")

	analysis := evolution.BugAnalysis
	fmt.Printf("📝 Description: %s\n", analysis.Description)
	fmt.Printf("🔍 Root Cause: %s\n", analysis.RootCause)
	fmt.Printf("📂 Category: %s\n", analysis.Category)
	fmt.Printf("⚠️  Severity: %s\n", analysis.Severity)
	fmt.Printf("🔄 Reproducibility: %s\n", analysis.Reproducibility)
	fmt.Printf("💥 Impact: %s\n", analysis.Impact)

	if len(evolution.RuleUpdates) > 0 {
		fmt.Println("\n📋 Rule Updates Required")
		fmt.Println("========================")
		for i, update := range evolution.RuleUpdates {
			fmt.Printf("%d. %s (%s)\n", i+1, update.Section, update.RuleFile)
			fmt.Printf("   %s\n", update.Explanation)
		}
	}

	if len(evolution.PreventionRules) > 0 {
		fmt.Println("\n🛡️  Prevention Rules")
		fmt.Println("===================")
		for i, rule := range evolution.PreventionRules {
			fmt.Printf("%d. %s\n", i+1, rule.Pattern)
			fmt.Printf("   Prevention: %s\n", rule.Prevention)
			if rule.CodeExample != "" {
				fmt.Printf("   Example:\n%s\n", rule.CodeExample)
			}
		}
	}

	if len(evolution.TestCases) > 0 {
		fmt.Println("\n🧪 Required Test Cases")
		fmt.Println("=====================")
		for i, test := range evolution.TestCases {
			fmt.Printf("%d. %s\n", i+1, test)
		}
	}
}

func applyRuleEvolution(evolution *RuleEvolution) error {
	for _, update := range evolution.RuleUpdates {
		// Read existing rule file
		content, err := os.ReadFile(update.RuleFile)
		if err != nil {
			return fmt.Errorf("failed to read rule file %s: %w", update.RuleFile, err)
		}

		// Append new rule
		newContent := string(content) + "\n\n" + update.NewRule + "\n"

		// Write back to file
		if err := os.WriteFile(update.RuleFile, []byte(newContent), 0644); err != nil {
			return fmt.Errorf("failed to update rule file %s: %w", update.RuleFile, err)
		}

		fmt.Printf("✅ Updated %s with new %s rule\n", update.RuleFile, update.Section)
	}

	return nil
}