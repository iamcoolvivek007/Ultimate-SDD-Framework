package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/review"
)

var (
	prNumber   int
	reviewDeep bool
)

func NewReviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "review [pr-number]",
		Short: "AI-powered automated code review",
		Long: `Perform comprehensive automated code review using AI analysis:
- Code quality assessment
- Security vulnerability detection
- Performance issue identification
- Best practice compliance
- Maintainability evaluation

Supports both PR review and general codebase analysis.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			// Get changed files (simplified - would integrate with Git in real implementation)
			changedFiles := []string{}
			if len(args) > 0 {
				// If PR number provided, get changed files from PR
				if pr, err := strconv.Atoi(args[0]); err == nil {
					prNumber = pr
					// In real implementation, would fetch from GitHub/GitLab API
					changedFiles = []string{"main.go", "internal/cli/review.go"} // Placeholder
				}
			} else {
				// General review - analyze recent changes
				changedFiles = []string{"internal/analysis/metrics.go"} // Placeholder
			}

			if len(changedFiles) == 0 {
				fmt.Println("No files to review. Specify a PR number or ensure there are changes.")
				return nil
			}

			fmt.Printf("🤖 Starting AI-powered code review of %d files...\n", len(changedFiles))

			// Create reviewer
			reviewer, err := review.NewCodeReviewer(projectRoot)
			if err != nil {
				return fmt.Errorf("failed to create reviewer: %w", err)
			}

			// Perform review
			codeReview, err := reviewer.ReviewPullRequest(prNumber, changedFiles)
			if err != nil {
				return fmt.Errorf("review failed: %w", err)
			}

			// Display results
			fmt.Println(reviewer.GetReviewReport(codeReview))

			// Save detailed report
			reportPath := ".sdd/review_report.md"
			if err := os.WriteFile(reportPath, []byte(reviewer.GetReviewReport(codeReview)), 0644); err != nil {
				fmt.Printf("Warning: Failed to save review report: %v\n", err)
			} else {
				fmt.Printf("📄 Review report saved to: %s\n", reportPath)
			}

			// Show approval status
			showReviewStatus(codeReview)

			return nil
		},
	}

	cmd.Flags().BoolVar(&reviewDeep, "deep", false, "Perform deep analysis with AI reasoning")

	return cmd
}

func showReviewStatus(review *review.CodeReview) {
	fmt.Println("\n📊 Review Status:")

	switch review.Summary.ApprovalStatus {
	case "approved":
		fmt.Printf("  ✅ **APPROVED** - Ready to merge\n")
	case "requested_changes":
		fmt.Printf("  🔄 **CHANGES REQUESTED** - Address issues before merging\n")
	case "blocked":
		fmt.Printf("  🚫 **BLOCKED** - Critical issues must be resolved\n")
	}

	fmt.Printf("  📈 Overall Score: %d/10\n", review.Summary.OverallScore)
	fmt.Printf("  🎯 Risk Level: %s\n", review.Summary.RiskLevel)

	if len(review.Summary.KeyFindings) > 0 {
		fmt.Println("  📋 Key Findings:")
		for _, finding := range review.Summary.KeyFindings {
			fmt.Printf("    • %s\n", finding)
		}
	}

	if len(review.Summary.Recommendations) > 0 {
		fmt.Println("  💡 Recommendations:")
		for _, rec := range review.Summary.Recommendations {
			fmt.Printf("    • %s\n", rec)
		}
	}
}