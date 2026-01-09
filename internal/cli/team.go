package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/collaboration"
)

var (
	teamName        string
	teamDescription string
	memberName      string
	memberEmail     string
	memberRole      string
	memberSkills    []string
	ruleCategory    string
	ruleTitle       string
	ruleDescription string
	ruleSeverity    string
	ruleExamples    []string
	knowledgeTitle  string
	knowledgeContent string
	knowledgeCategory string
	knowledgeTags   []string
	patternName     string
	patternDesc     string
	patternLang     string
	patternCode     string
	patternUseCase  string
	searchQuery     string
	searchCategory  string
)

func NewTeamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "team",
		Short: "Team collaboration and knowledge management",
		Long: `Manage team collaboration features:
- Team member management
- Shared coding standards and rules
- Knowledge base management
- Code pattern sharing
- Decision logging
- Collaborative development workflows`,
	}

	// Subcommands
	cmd.AddCommand(NewTeamInitCmd())
	cmd.AddCommand(NewTeamMemberCmd())
	cmd.AddCommand(NewTeamRuleCmd())
	cmd.AddCommand(NewTeamKnowledgeCmd())
	cmd.AddCommand(NewTeamPatternCmd())
	cmd.AddCommand(NewTeamDecisionCmd())
	cmd.AddCommand(NewTeamSearchCmd())
	cmd.AddCommand(NewTeamReportCmd())

	return cmd
}

func NewTeamInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize team collaboration",
		Long:  "Create a new team or initialize team features for the project.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			if teamName == "" {
				teamName = "Development Team"
			}

			if teamDescription == "" {
				teamDescription = "Collaborative development team"
			}

			fmt.Printf("👥 Initializing team: %s\n", teamName)

			// Create team collaboration
			teamCollab, err := collaboration.NewTeamCollaboration(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize team collaboration: %w", err)
			}

			// Create team
			team, err := teamCollab.CreateTeam(teamName, teamDescription)
			if err != nil {
				return fmt.Errorf("failed to create team: %w", err)
			}

			fmt.Printf("✅ Team created: %s (%s)\n", team.Name, team.ID)
			fmt.Println("\n🚀 Next steps:")
			fmt.Println("  nexus team member add --name \"Your Name\" --email \"your.email@company.com\" --role \"developer\"")
			fmt.Println("  nexus team rule add --category \"coding_standards\" --title \"Use meaningful variable names\"")
			fmt.Println("  nexus team knowledge add --title \"API Design Patterns\" --category \"best_practices\"")

			return nil
		},
	}

	cmd.Flags().StringVar(&teamName, "name", "", "Team name")
	cmd.Flags().StringVar(&teamDescription, "description", "", "Team description")

	return cmd
}

func NewTeamMemberCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member",
		Short: "Manage team members",
	}

	cmd.AddCommand(NewTeamMemberAddCmd())
	cmd.AddCommand(NewTeamMemberListCmd())

	return cmd
}

func NewTeamMemberAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a team member",
		Long:  "Add a new member to the development team.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			if memberName == "" || memberEmail == "" {
				return fmt.Errorf("name and email are required")
			}

			if memberRole == "" {
				memberRole = "developer"
			}

			fmt.Printf("👤 Adding team member: %s (%s)\n", memberName, memberRole)

			// Create team collaboration
			teamCollab, err := collaboration.NewTeamCollaboration(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize team collaboration: %w", err)
			}

			// Add member
			member, err := teamCollab.AddTeamMember(memberName, memberEmail, memberRole, memberSkills)
			if err != nil {
				return fmt.Errorf("failed to add team member: %w", err)
			}

			fmt.Printf("✅ Added: %s (%s) - %s\n", member.Name, member.Email, member.Role)
			if len(member.Skills) > 0 {
				fmt.Printf("🛠️  Skills: %s\n", strings.Join(member.Skills, ", "))
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&memberName, "name", "", "Member name")
	cmd.Flags().StringVar(&memberEmail, "email", "", "Member email")
	cmd.Flags().StringVar(&memberRole, "role", "developer", "Member role (lead, senior, junior, qa, etc.)")
	cmd.Flags().StringSliceVar(&memberSkills, "skills", []string{}, "Member skills (comma-separated)")

	return cmd
}

func NewTeamMemberListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List team members",
		Long:  "Display all members of the development team.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			// Create team collaboration
			_, err := collaboration.NewTeamCollaboration(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize team collaboration: %w", err)
			}

			// This would need to be added to the TeamCollaboration struct
			// For now, just show that team features are available
			fmt.Println("👥 Team Members:")
			fmt.Println("  (Team member listing would be implemented here)")

			return nil
		},
	}

	return cmd
}

func NewTeamRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "Manage team coding rules",
	}

	cmd.AddCommand(NewTeamRuleAddCmd())
	cmd.AddCommand(NewTeamRuleListCmd())

	return cmd
}

func NewTeamRuleAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a team coding rule",
		Long:  "Add a new coding standard or rule for the team.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			if ruleCategory == "" || ruleTitle == "" {
				return fmt.Errorf("category and title are required")
			}

			if ruleSeverity == "" {
				ruleSeverity = "recommended"
			}

			fmt.Printf("📋 Adding team rule: %s\n", ruleTitle)

			// Create team collaboration
			teamCollab, err := collaboration.NewTeamCollaboration(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize team collaboration: %w", err)
			}

			// Add rule
			rule, err := teamCollab.AddTeamRule(ruleCategory, ruleTitle, ruleDescription, ruleSeverity, "current_user", ruleExamples)
			if err != nil {
				return fmt.Errorf("failed to add team rule: %w", err)
			}

			fmt.Printf("✅ Added rule: %s (%s - %s)\n", rule.Title, rule.Category, rule.Severity)

			return nil
		},
	}

	cmd.Flags().StringVar(&ruleCategory, "category", "", "Rule category (coding_standards, code_review, testing, security, performance, documentation)")
	cmd.Flags().StringVar(&ruleTitle, "title", "", "Rule title")
	cmd.Flags().StringVar(&ruleDescription, "description", "", "Rule description")
	cmd.Flags().StringVar(&ruleSeverity, "severity", "recommended", "Rule severity (mandatory, recommended, optional)")
	cmd.Flags().StringSliceVar(&ruleExamples, "examples", []string{}, "Rule examples")

	return cmd
}

func NewTeamRuleListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List team rules",
		Long:  "Display all team coding standards and rules.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			// Create team collaboration
			teamCollab, err := collaboration.NewTeamCollaboration(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize team collaboration: %w", err)
			}

			rules := teamCollab.GetTeamRules()

			fmt.Println("📋 Team Rules:")

			ruleCategories := map[string][]collaboration.RuleDefinition{
				"Coding Standards": rules.CodingStandards,
				"Code Review":      rules.CodeReviewRules,
				"Testing":          rules.TestingStandards,
				"Security":         rules.SecurityPolicies,
				"Performance":      rules.PerformanceRules,
				"Documentation":    rules.DocumentationRules,
			}

			totalRules := 0
			for category, ruleList := range ruleCategories {
				if len(ruleList) > 0 {
					fmt.Printf("\n### %s (%d rules)\n", category, len(ruleList))
					for _, rule := range ruleList {
						fmt.Printf("  • **%s** (%s): %s\n", rule.Title, rule.Severity, rule.Description)
					}
					totalRules += len(ruleList)
				}
			}

			if totalRules == 0 {
				fmt.Println("  No team rules defined yet.")
				fmt.Println("  Add some with: nexus team rule add --category \"coding_standards\" --title \"Your Rule\"")
			}

			return nil
		},
	}

	return cmd
}

func NewTeamKnowledgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "knowledge",
		Short: "Manage team knowledge base",
	}

	cmd.AddCommand(NewTeamKnowledgeAddCmd())
	cmd.AddCommand(NewTeamKnowledgeListCmd())

	return cmd
}

func NewTeamKnowledgeAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add knowledge to team base",
		Long:  "Add a new knowledge item to the team knowledge base.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			if knowledgeTitle == "" || knowledgeCategory == "" {
				return fmt.Errorf("title and category are required")
			}

			if knowledgeContent == "" {
				// Try to read from stdin or file
				fmt.Println("Enter knowledge content (Ctrl+D to finish):")
				content, err := readFromStdin()
				if err != nil {
					return fmt.Errorf("failed to read content: %w", err)
				}
				knowledgeContent = content
			}

			fmt.Printf("🧠 Adding knowledge: %s\n", knowledgeTitle)

			// Create team collaboration
			teamCollab, err := collaboration.NewTeamCollaboration(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize team collaboration: %w", err)
			}

			// Add knowledge
			item, err := teamCollab.AddKnowledgeItem(knowledgeTitle, knowledgeContent, knowledgeCategory, "current_user", knowledgeTags)
			if err != nil {
				return fmt.Errorf("failed to add knowledge: %w", err)
			}

			fmt.Printf("✅ Added: %s (%s)\n", item.Title, item.Category)
			if len(item.Tags) > 0 {
				fmt.Printf("🏷️  Tags: %s\n", strings.Join(item.Tags, ", "))
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&knowledgeTitle, "title", "", "Knowledge title")
	cmd.Flags().StringVar(&knowledgeContent, "content", "", "Knowledge content")
	cmd.Flags().StringVar(&knowledgeCategory, "category", "", "Category (best_practices, common_issues, architecture)")
	cmd.Flags().StringSliceVar(&knowledgeTags, "tags", []string{}, "Tags (comma-separated)")

	return cmd
}

func NewTeamKnowledgeListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List team knowledge",
		Long:  "Display items from the team knowledge base.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			// Create team collaboration
			teamCollab, err := collaboration.NewTeamCollaboration(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize team collaboration: %w", err)
			}

			knowledge := teamCollab.GetTeamKnowledge()

			fmt.Println("🧠 Team Knowledge Base:")

			knowledgeCategories := map[string][]collaboration.KnowledgeItem{
				"Best Practices":   knowledge.BestPractices,
				"Common Issues":    knowledge.CommonIssues,
				"Architecture":     knowledge.ArchitectureDocs,
			}

			totalItems := 0
			for category, items := range knowledgeCategories {
				if len(items) > 0 {
					fmt.Printf("\n### %s (%d items)\n", category, len(items))
					for _, item := range items {
						fmt.Printf("  • **%s**: %s\n", item.Title, truncateString(item.Content, 100))
					}
					totalItems += len(items)
				}
			}

			if len(knowledge.CodePatterns) > 0 {
				fmt.Printf("\n### Code Patterns (%d patterns)\n", len(knowledge.CodePatterns))
				for _, pattern := range knowledge.CodePatterns {
					fmt.Printf("  • **%s** (%s): Used %d times\n", pattern.Name, pattern.Language, pattern.UsageCount)
				}
				totalItems += len(knowledge.CodePatterns)
			}

			if totalItems == 0 {
				fmt.Println("  No knowledge items yet.")
				fmt.Println("  Add some with: nexus team knowledge add --title \"Your Knowledge\" --category \"best_practices\"")
			}

			return nil
		},
	}

	return cmd
}

func NewTeamPatternCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pattern",
		Short: "Manage code patterns",
	}

	cmd.AddCommand(NewTeamPatternAddCmd())
	cmd.AddCommand(NewTeamPatternListCmd())

	return cmd
}

func NewTeamPatternAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a code pattern",
		Long:  "Add a reusable code pattern to the team library.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			if patternName == "" || patternLang == "" {
				return fmt.Errorf("name and language are required")
			}

			if patternCode == "" {
				return fmt.Errorf("code content is required")
			}

			fmt.Printf("🔧 Adding code pattern: %s (%s)\n", patternName, patternLang)

			// Create team collaboration
			teamCollab, err := collaboration.NewTeamCollaboration(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize team collaboration: %w", err)
			}

			// Add pattern
			pattern, err := teamCollab.AddCodePattern(patternName, patternDesc, patternLang, patternCode, patternUseCase, "current_user")
			if err != nil {
				return fmt.Errorf("failed to add pattern: %w", err)
			}

			fmt.Printf("✅ Added pattern: %s\n", pattern.Name)
			fmt.Printf("📊 Usage: %d times\n", pattern.UsageCount)

			return nil
		},
	}

	cmd.Flags().StringVar(&patternName, "name", "", "Pattern name")
	cmd.Flags().StringVar(&patternDesc, "description", "", "Pattern description")
	cmd.Flags().StringVar(&patternLang, "language", "", "Programming language")
	cmd.Flags().StringVar(&patternCode, "code", "", "Pattern code")
	cmd.Flags().StringVar(&patternUseCase, "usecase", "general", "Use case for the pattern")

	return cmd
}

func NewTeamPatternListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List code patterns",
		Long:  "Display available code patterns from the team library.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			// Create team collaboration
			teamCollab, err := collaboration.NewTeamCollaboration(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize team collaboration: %w", err)
			}

			patterns := teamCollab.GetCodePatterns("", "")

			if len(patterns) == 0 {
				fmt.Println("🔧 No code patterns in the team library yet.")
				fmt.Println("  Add some with: nexus team pattern add --name \"Your Pattern\" --language \"go\" --code \"...\"")
				return nil
			}

			fmt.Printf("🔧 Team Code Patterns (%d total):\n", len(patterns))

			for i, pattern := range patterns {
				fmt.Printf("\n%d. **%s** (%s)\n", i+1, pattern.Name, pattern.Language)
				fmt.Printf("   Usage: %d times\n", pattern.UsageCount)
				if pattern.Description != "" {
					fmt.Printf("   Description: %s\n", pattern.Description)
				}
				if pattern.UseCase != "" {
					fmt.Printf("   Use Case: %s\n", pattern.UseCase)
				}
			}

			return nil
		},
	}

	return cmd
}

func NewTeamDecisionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decision",
		Short: "Record team decisions",
		Long:  "Document important architectural or design decisions made by the team.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// This would be implemented with proper argument parsing
			fmt.Println("📋 Team decision recording would be implemented here")
			fmt.Println("Usage: nexus team decision --title \"Decision Title\" --context \"Background\" --decision \"What was decided\"")
			return nil
		},
	}

	return cmd
}

func NewTeamSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search team knowledge base",
		Long:  "Search through team knowledge, patterns, and documentation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			if len(args) == 0 && searchQuery == "" {
				return fmt.Errorf("search query is required")
			}

			query := searchQuery
			if len(args) > 0 {
				query = args[0]
			}

			fmt.Printf("🔍 Searching team knowledge for: %s\n", query)

			// Create team collaboration
			teamCollab, err := collaboration.NewTeamCollaboration(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize team collaboration: %w", err)
			}

			// Search knowledge
			results := teamCollab.SearchKnowledge(query, searchCategory)

			if len(results) == 0 {
				fmt.Println("❌ No matching knowledge found.")
				return nil
			}

			fmt.Printf("📚 Found %d matching items:\n", len(results))

			for i, item := range results {
				fmt.Printf("\n%d. **%s** (%s)\n", i+1, item.Title, item.Category)
				fmt.Printf("   Content: %s\n", truncateString(item.Content, 150))
				if len(item.Tags) > 0 {
					fmt.Printf("   Tags: %s\n", strings.Join(item.Tags, ", "))
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&searchQuery, "query", "", "Search query")
	cmd.Flags().StringVar(&searchCategory, "category", "", "Search category filter")

	return cmd
}

func NewTeamReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate team collaboration report",
		Long:  "Create a comprehensive report of team activities, knowledge, and collaboration metrics.",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			fmt.Println("📊 Generating team collaboration report...")

			// Create team collaboration
			teamCollab, err := collaboration.NewTeamCollaboration(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to initialize team collaboration: %w", err)
			}

			// Generate report
			report := teamCollab.GenerateTeamReport()

			// Display report
			fmt.Println(report)

			// Save report
			reportPath := ".sdd/team_report.md"
			if err := os.WriteFile(reportPath, []byte(report), 0644); err != nil {
				fmt.Printf("Warning: Failed to save team report: %v\n", err)
			} else {
				fmt.Printf("📄 Team report saved to: %s\n", reportPath)
			}

			return nil
		},
	}

	return cmd
}

// Helper functions

func readFromStdin() (string, error) {
	// Simple implementation - would need proper stdin reading
	return "", fmt.Errorf("stdin reading not implemented yet")
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}