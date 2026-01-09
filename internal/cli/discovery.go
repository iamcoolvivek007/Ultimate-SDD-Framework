package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/lsp"
)

func NewDiscoveryCmd() *cobra.Command {
	var deepAnalysis bool

	cmd := &cobra.Command{
		Use:   "discovery",
		Short: "Analyze existing codebase and generate system context",
		Long: `Perform comprehensive brownfield analysis of the existing codebase.

This command creates a CONTEXT.md file that serves as the source of truth
for the current system state, helping AI agents understand legacy patterns,
forbidden practices, integration points, and technical debt.

Use --deep flag for thorough analysis including code patterns and dependencies.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			fmt.Println("🔍 Starting brownfield discovery analysis...")

			// Create brownfield context analyzer
			bfc := lsp.NewBrownfieldContext(projectRoot)

			// Perform analysis
			if err := bfc.AnalyzeBrownfield(); err != nil {
				return fmt.Errorf("failed to analyze codebase: %w", err)
			}

			fmt.Printf("✅ Analyzed %d files\n", len(bfc.Files))

			// Generate CONTEXT.md
			contextContent := bfc.GenerateCONTEXTFile()

			// Ensure .sdd directory exists
			sddDir := filepath.Join(projectRoot, ".sdd")
			if err := os.MkdirAll(sddDir, 0755); err != nil {
				return fmt.Errorf("failed to create .sdd directory: %w", err)
			}

			// Save CONTEXT.md
			contextPath := filepath.Join(sddDir, "CONTEXT.md")
			if err := os.WriteFile(contextPath, []byte(contextContent), 0644); err != nil {
				return fmt.Errorf("failed to save context file: %w", err)
			}

			fmt.Printf("📄 Generated system context: %s\n", contextPath)

			// Show summary
			showDiscoverySummary(bfc)

			fmt.Println("\n🎯 Next steps:")
			fmt.Println("  1. Review CONTEXT.md to understand system constraints")
			fmt.Println("  2. Run: nexus specify \"your feature description\"")
			fmt.Println("  3. The system will now validate requests against legacy patterns")

			return nil
		},
	}

	cmd.Flags().BoolVar(&deepAnalysis, "deep", false, "Perform deep analysis including code patterns and dependencies")

	return cmd
}

func showDiscoverySummary(bfc *lsp.BrownfieldContext) {
	fmt.Println("\n📊 Discovery Summary")
	fmt.Println("===================")

	// System overview
	fmt.Printf("🏗️  Architecture: %s with %s\n", bfc.Structure.MainLanguage, bfc.Structure.Framework)
	fmt.Printf("📁 Files analyzed: %d\n", len(bfc.Files))

	// Features detected
	features := []string{}
	if bfc.Structure.HasAPI {
		features = append(features, "API")
	}
	if bfc.Structure.HasDatabase {
		features = append(features, "Database")
	}
	if bfc.Structure.HasFrontend {
		features = append(features, "Frontend")
	}
	if bfc.Structure.HasTests {
		features = append(features, "Tests")
	}

	if len(features) > 0 {
		fmt.Printf("✨ Features: %s\n", fmt.Sprintf("%v", features))
	}

	// Legacy patterns found
	if len(bfc.LegacyPatterns) > 0 {
		fmt.Printf("\n📚 Legacy Patterns: %d identified\n", len(bfc.LegacyPatterns))
		for i, pattern := range bfc.LegacyPatterns {
			if i < 3 { // Show first 3
				fmt.Printf("  • %s (%d files)\n", pattern.Pattern, len(pattern.Files))
			}
		}
		if len(bfc.LegacyPatterns) > 3 {
			fmt.Printf("  ... and %d more\n", len(bfc.LegacyPatterns)-3)
		}
	}

	// Forbidden patterns found
	if len(bfc.ForbiddenPatterns) > 0 {
		fmt.Printf("\n🚫 Forbidden Patterns: %d identified\n", len(bfc.ForbiddenPatterns))
		severityCount := make(map[string]int)
		for _, pattern := range bfc.ForbiddenPatterns {
			severityCount[pattern.Severity]++
		}

		for severity, count := range severityCount {
			fmt.Printf("  • %s: %d\n", severity, count)
		}
	}

	// Integration points
	if len(bfc.IntegrationPoints) > 0 {
		fmt.Printf("\n🔗 Integration Points: %d mapped\n", len(bfc.IntegrationPoints))
		pointTypes := make(map[string]int)
		for _, point := range bfc.IntegrationPoints {
			pointTypes[point.Type]++
		}

		for pointType, count := range pointTypes {
			fmt.Printf("  • %s: %d\n", pointType, count)
		}
	}

	// Technical debt
	if len(bfc.TechnicalDebt) > 0 {
		fmt.Printf("\n💸 Technical Debt: %d items identified\n", len(bfc.TechnicalDebt))
		severityCount := make(map[string]int)
		for _, debt := range bfc.TechnicalDebt {
			severityCount[debt.Severity]++
		}

		for severity, count := range severityCount {
			fmt.Printf("  • %s: %d\n", severity, count)
		}
	}

	// Constitution status
	if len(bfc.Constitution.TechStack) > 0 {
		fmt.Printf("\n🏛️  Constitution: %d tech stack items defined\n", len(bfc.Constitution.TechStack))
		fmt.Printf("📋 Rules: %d architectural, %d coding, %d integration\n",
			len(bfc.Constitution.ArchitecturalRules),
			len(bfc.Constitution.CodingStandards),
			len(bfc.Constitution.IntegrationRules))
	} else {
		fmt.Println("\n🏛️  Constitution: Default rules applied (create CONSTITUTION.md for custom rules)")
	}
}