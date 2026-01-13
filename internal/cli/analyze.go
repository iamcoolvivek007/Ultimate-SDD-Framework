package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/analysis"
)

func NewAnalyzeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze codebase quality and generate reports",
		Long: `Perform comprehensive code quality analysis including:
- Code metrics and complexity analysis
- Security vulnerability scanning
- Performance issue detection
- Maintainability assessment
- Test coverage evaluation

Generates detailed reports with actionable recommendations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			fmt.Println("🔍 Starting comprehensive code analysis...")

			// Create analyzer
			analyzer := analysis.NewCodeAnalyzer(projectRoot)

			// Perform analysis
			report, err := analyzer.Analyze()
			if err != nil {
				return fmt.Errorf("analysis failed: %w", err)
			}

			// Display results
			fmt.Println(report.GetSummary())

			// Save detailed report
			reportPath := ".sdd/analysis_report.md"
			if err := os.WriteFile(reportPath, []byte(report.GetSummary()), 0644); err != nil {
				fmt.Printf("Warning: Failed to save report to file: %v\n", err)
			} else {
				fmt.Printf("📄 Detailed report saved to: %s\n", reportPath)
			}

			// Show next steps based on score
			showAnalysisRecommendations(report)

			return nil
		},
	}

	return cmd
}

func showAnalysisRecommendations(report *analysis.QualityReport) {
	fmt.Println("\n🎯 Recommendations:")

	if report.Score.Overall < 70 {
		fmt.Println("  🚨 Critical: Address high-priority issues before proceeding")
		fmt.Println("  📚 Consider refactoring complex functions")
		fmt.Println("  ✅ Add comprehensive test coverage")
	} else if report.Score.Overall < 85 {
		fmt.Println("  📈 Good foundation - focus on improvements")
		fmt.Println("  🧹 Address remaining code quality issues")
		fmt.Println("  🛡️ Review security recommendations")
	} else {
		fmt.Println("  🎉 Excellent code quality!")
		fmt.Println("  🔄 Consider preventive maintenance")
		fmt.Println("  📖 Share best practices with team")
	}

	fmt.Println("  🏃‍♂️ Run 'nexus review' for AI-powered code review")
	fmt.Println("  👥 Use 'nexus team add-rule' to establish standards")
}
