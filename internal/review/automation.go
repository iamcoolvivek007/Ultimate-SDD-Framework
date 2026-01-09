package review

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"ultimate-sdd-framework/internal/agents"
	"ultimate-sdd-framework/internal/analysis"
)

// CodeReview represents an automated code review
type CodeReview struct {
	Repository string                    `json:"repository"`
	Branch     string                    `json:"branch"`
	Files      []FileReview              `json:"files"`
	Summary    ReviewSummary             `json:"summary"`
	Agent      *agents.Agent            `json:"agent"`
}

// FileReview represents review of a single file
type FileReview struct {
	Path         string        `json:"path"`
	Status       string        `json:"status"`       // approved, changes_requested, commented
	Comments     []ReviewComment `json:"comments"`
	Suggestions  []string      `json:"suggestions"`
	Score        int           `json:"score"`        // 1-10 quality score
	Issues       []CodeIssue   `json:"issues"`
}

// ReviewComment represents a specific comment on code
type ReviewComment struct {
	Line     int    `json:"line"`
	Type    string `json:"type"`    // suggestion, issue, praise, question
	Message string `json:"message"`
	Severity string `json:"severity"` // info, warning, error
	RuleID  string `json:"rule_id"`
}

// CodeIssue represents a specific code issue
type CodeIssue struct {
	Type        string `json:"type"`        // security, performance, maintainability, style
	Severity    string `json:"severity"`    // low, medium, high, critical
	Message     string `json:"message"`
	Line         int    `json:"line"`
	Suggestion  string `json:"suggestion"`
	Category    string `json:"category"`
}

// ReviewSummary provides overall review assessment
type ReviewSummary struct {
	OverallScore     int               `json:"overall_score"`     // 1-10
	ApprovalStatus   string            `json:"approval_status"`   // approved, requested_changes, blocked
	RiskLevel        string            `json:"risk_level"`        // low, medium, high, critical
	IssuesByCategory map[string]int    `json:"issues_by_category"`
	KeyFindings      []string          `json:"key_findings"`
	Recommendations  []string          `json:"recommendations"`
}

// CodeReviewer performs automated code reviews
type CodeReviewer struct {
	agentSvc    *agents.AgentService
	analyzer    *analysis.CodeAnalyzer
	projectRoot string
}

// NewCodeReviewer creates a new code reviewer
func NewCodeReviewer(projectRoot string) (*CodeReviewer, error) {
	agentSvc := agents.NewAgentService(projectRoot)
	if err := agentSvc.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize agent service: %w", err)
	}

	analyzer := analysis.NewCodeAnalyzer(projectRoot)

	return &CodeReviewer{
		agentSvc:    agentSvc,
		analyzer:    analyzer,
		projectRoot: projectRoot,
	}, nil
}

// ReviewPullRequest performs automated review of a pull request
func (cr *CodeReviewer) ReviewPullRequest(prNumber int, changedFiles []string) (*CodeReview, error) {
	review := &CodeReview{
		Repository: "current",
		Branch:     "current",
		Files:      []FileReview{},
	}

	// Get QA agent for review
	qaAgent, err := cr.agentSvc.GetAgentForPhase("review")
	if err != nil {
		return nil, fmt.Errorf("failed to get QA agent: %w", err)
	}
	review.Agent = qaAgent

	// Analyze each changed file
	for _, filePath := range changedFiles {
		fileReview, err := cr.reviewFile(filePath)
		if err != nil {
			// Log error but continue with other files
			fmt.Printf("Warning: Failed to review %s: %v\n", filePath, err)
			continue
		}
		review.Files = append(review.Files, *fileReview)
	}

	// Generate overall summary
	review.Summary = cr.generateSummary(review.Files)

	return review, nil
}

// reviewFile performs detailed review of a single file
func (cr *CodeReviewer) reviewFile(filePath string) (*FileReview, error) {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	fileReview := &FileReview{
		Path:     filePath,
		Status:   "approved", // Default to approved
		Comments: []ReviewComment{},
		Suggestions: []string{},
		Score:    8, // Default good score
		Issues:   []CodeIssue{},
	}

	// Perform automated analysis
	issues := cr.analyzeFileIssues(filePath, string(content))
	fileReview.Issues = issues

	// Generate comments from issues
	comments := cr.generateCommentsFromIssues(issues)
	fileReview.Comments = comments

	// Generate suggestions
	suggestions := cr.generateSuggestions(filePath, string(content))
	fileReview.Suggestions = suggestions

	// Calculate file score
	fileReview.Score = cr.calculateFileScore(issues, comments)

	// Determine status based on issues
	fileReview.Status = cr.determineFileStatus(issues)

	return fileReview, nil
}

// analyzeFileIssues performs automated issue detection
func (cr *CodeReviewer) analyzeFileIssues(filePath, content string) []CodeIssue {
	issues := []CodeIssue{}

	lines := strings.Split(content, "\n")

	// Check for common issues based on file type
	if strings.HasSuffix(filePath, ".go") {
		issues = append(issues, cr.analyzeGoIssues(filePath, lines)...)
	} else if strings.HasSuffix(filePath, ".ts") || strings.HasSuffix(filePath, ".tsx") ||
	          strings.HasSuffix(filePath, ".js") || strings.HasSuffix(filePath, ".jsx") {
		issues = append(issues, cr.analyzeJSIssues(filePath, lines)...)
	}

	// Check for security issues
	issues = append(issues, cr.analyzeSecurityIssues(content)...)

	// Check for performance issues
	issues = append(issues, cr.analyzePerformanceIssues(content)...)

	return issues
}

// analyzeGoIssues checks for Go-specific issues
func (cr *CodeReviewer) analyzeGoIssues(filePath string, lines []string) []CodeIssue {
	issues := []CodeIssue{}

	for i, line := range lines {
		// Check for panic usage
		if strings.Contains(line, "panic(") {
			issues = append(issues, CodeIssue{
				Type:       "error-handling",
				Severity:   "medium",
				Message:    "Use of panic() detected - consider proper error handling",
				Line:        i + 1,
				Suggestion: "Return errors instead of panicking",
				Category:   "maintainability",
			})
		}

		// Check for TODO comments
		if strings.Contains(strings.ToLower(line), "todo") {
			issues = append(issues, CodeIssue{
				Type:       "documentation",
				Severity:   "low",
				Message:    "TODO comment found - consider addressing or creating issue",
				Line:        i + 1,
				Suggestion: "Resolve TODO or create tracking issue",
				Category:   "maintainability",
			})
		}

		// Check for long lines
		if len(line) > 120 {
			issues = append(issues, CodeIssue{
				Type:       "style",
				Severity:   "low",
				Message:    fmt.Sprintf("Line too long (%d characters)", len(line)),
				Line:        i + 1,
				Suggestion: "Break long lines for better readability",
				Category:   "style",
			})
		}
	}

	return issues
}

// analyzeJSIssues checks for JavaScript/TypeScript specific issues
func (cr *CodeReviewer) analyzeJSIssues(filePath string, lines []string) []CodeIssue {
	issues := []CodeIssue{}

	for i, line := range lines {
		// Check for console.log in production code
		if strings.Contains(line, "console.log") && !strings.Contains(filePath, "test") {
			issues = append(issues, CodeIssue{
				Type:       "logging",
				Severity:   "low",
				Message:    "console.log found in production code",
				Line:        i + 1,
				Suggestion: "Use proper logging library instead",
				Category:   "maintainability",
			})
		}

		// Check for any type usage in TypeScript
		if strings.HasSuffix(filePath, ".ts") && strings.Contains(line, ": any") {
			issues = append(issues, CodeIssue{
				Type:       "typescript",
				Severity:   "medium",
				Message:    "'any' type usage reduces type safety",
				Line:        i + 1,
				Suggestion: "Use specific types instead of 'any'",
				Category:   "maintainability",
			})
		}
	}

	return issues
}

// analyzeSecurityIssues checks for security vulnerabilities
func (cr *CodeReviewer) analyzeSecurityIssues(content string) []CodeIssue {
	issues := []CodeIssue{}

	// Check for hardcoded secrets
	secretPatterns := []string{
		`password\s*=\s*["'][^"']*["']`,
		`secret\s*=\s*["'][^"']*["']`,
		`token\s*=\s*["'][^"']*["']`,
		`key\s*=\s*["'][^"']*["']`,
	}

	for _, pattern := range secretPatterns {
		re := regexp.MustCompile("(?i)" + pattern)
		if re.MatchString(content) {
			issues = append(issues, CodeIssue{
				Type:       "security",
				Severity:   "high",
				Message:    "Potential hardcoded secret detected",
				Suggestion: "Use environment variables or secure credential storage",
				Category:   "security",
			})
		}
	}

	// Check for SQL injection vulnerabilities
	if strings.Contains(content, "sprintf") && strings.Contains(content, "query") {
		issues = append(issues, CodeIssue{
			Type:       "security",
			Severity:   "high",
			Message:    "Potential SQL injection vulnerability",
			Suggestion: "Use parameterized queries instead of string formatting",
			Category:   "security",
		})
	}

	return issues
}

// analyzePerformanceIssues checks for performance problems
func (cr *CodeReviewer) analyzePerformanceIssues(content string) []CodeIssue {
	issues := []CodeIssue{}

	// Check for N+1 query patterns
	if strings.Contains(content, "for") && strings.Contains(content, "query") {
		issues = append(issues, CodeIssue{
			Type:       "performance",
			Severity:   "medium",
			Message:    "Potential N+1 query pattern detected",
			Suggestion: "Consider batch queries or eager loading",
			Category:   "performance",
		})
	}

	// Check for memory leaks (simplified)
	if strings.Contains(content, "new ") && !strings.Contains(content, "delete") && strings.HasSuffix(content, ".cpp") {
		issues = append(issues, CodeIssue{
			Type:       "performance",
			Severity:   "medium",
			Message:    "Potential memory leak with 'new' allocation",
			Suggestion: "Ensure proper memory cleanup",
			Category:   "performance",
		})
	}

	return issues
}

// generateCommentsFromIssues creates review comments from issues
func (cr *CodeReviewer) generateCommentsFromIssues(issues []CodeIssue) []ReviewComment {
	comments := []ReviewComment{}

	for _, issue := range issues {
		comment := ReviewComment{
			Line:     issue.Line,
			Message: issue.Message,
			RuleID:  issue.Type,
		}

		// Map severity to comment type
		switch issue.Severity {
		case "critical", "high":
			comment.Type = "issue"
			comment.Severity = "error"
		case "medium":
			comment.Type = "issue"
			comment.Severity = "warning"
		case "low":
			comment.Type = "suggestion"
			comment.Severity = "info"
		default:
			comment.Type = "comment"
			comment.Severity = "info"
		}

		if issue.Suggestion != "" {
			comment.Message += fmt.Sprintf(" Suggestion: %s", issue.Suggestion)
		}

		comments = append(comments, comment)
	}

	return comments
}

// generateSuggestions creates improvement suggestions
func (cr *CodeReviewer) generateSuggestions(filePath, content string) []string {
	suggestions := []string{}

	// File-specific suggestions
	if strings.HasSuffix(filePath, "_test.go") {
		if !strings.Contains(content, "TestMain") {
			suggestions = append(suggestions, "Consider adding TestMain for test setup/cleanup")
		}
	}

	// General suggestions
	lineCount := len(strings.Split(content, "\n"))
	if lineCount > 300 {
		suggestions = append(suggestions, "Consider splitting this large file into smaller modules")
	}

	if strings.Count(content, "TODO") > 0 {
		suggestions = append(suggestions, "Address the TODO comments or create tracking issues")
	}

	return suggestions
}

// calculateFileScore computes a quality score for the file
func (cr *CodeReviewer) calculateFileScore(issues []CodeIssue, comments []ReviewComment) int {
	score := 10 // Start with perfect score

	// Deduct points for issues
	for _, issue := range issues {
		switch issue.Severity {
		case "critical":
			score -= 3
		case "high":
			score -= 2
		case "medium":
			score -= 1
		case "low":
			score -= 0
		}
	}

	// Ensure score stays within bounds
	if score < 1 {
		score = 1
	}

	return score
}

// determineFileStatus determines if the file should be approved or needs changes
func (cr *CodeReviewer) determineFileStatus(issues []CodeIssue) string {
	hasHighSeverity := false
	hasCritical := false

	for _, issue := range issues {
		if issue.Severity == "critical" {
			hasCritical = true
		} else if issue.Severity == "high" {
			hasHighSeverity = true
		}
	}

	if hasCritical {
		return "blocked"
	} else if hasHighSeverity {
		return "changes_requested"
	}

	return "approved"
}

// generateSummary creates overall review summary
func (cr *CodeReviewer) generateSummary(fileReviews []FileReview) ReviewSummary {
	summary := ReviewSummary{
		IssuesByCategory: make(map[string]int),
		KeyFindings:      []string{},
		Recommendations:  []string{},
	}

	totalScore := 0
	totalFiles := len(fileReviews)
	criticalIssues := 0
	highIssues := 0

	// Aggregate data from all files
	for _, file := range fileReviews {
		totalScore += file.Score

		for _, issue := range file.Issues {
			summary.IssuesByCategory[issue.Category]++

			if issue.Severity == "critical" {
				criticalIssues++
			} else if issue.Severity == "high" {
				highIssues++
			}
		}
	}

	// Calculate overall score
	if totalFiles > 0 {
		summary.OverallScore = totalScore / totalFiles
	} else {
		summary.OverallScore = 10
	}

	// Determine approval status and risk level
	if criticalIssues > 0 {
		summary.ApprovalStatus = "blocked"
		summary.RiskLevel = "critical"
	} else if highIssues > 0 || summary.OverallScore < 7 {
		summary.ApprovalStatus = "requested_changes"
		summary.RiskLevel = "high"
	} else if summary.OverallScore < 8 {
		summary.ApprovalStatus = "approved"
		summary.RiskLevel = "medium"
	} else {
		summary.ApprovalStatus = "approved"
		summary.RiskLevel = "low"
	}

	// Generate key findings
	if criticalIssues > 0 {
		summary.KeyFindings = append(summary.KeyFindings,
			fmt.Sprintf("🚨 %d critical issues require immediate attention", criticalIssues))
	}

	if highIssues > 0 {
		summary.KeyFindings = append(summary.KeyFindings,
			fmt.Sprintf("⚠️  %d high-severity issues need addressing", highIssues))
	}

	if summary.OverallScore >= 9 {
		summary.KeyFindings = append(summary.KeyFindings, "✅ High-quality code with excellent maintainability")
	}

	// Generate recommendations
	if summary.RiskLevel == "critical" || summary.RiskLevel == "high" {
		summary.Recommendations = append(summary.Recommendations,
			"Address high-severity issues before merging")
	}

	if summary.IssuesByCategory["security"] > 0 {
		summary.Recommendations = append(summary.Recommendations,
			"Review and fix security issues identified")
	}

	if summary.IssuesByCategory["performance"] > 0 {
		summary.Recommendations = append(summary.Recommendations,
			"Optimize performance issues for better user experience")
	}

	if summary.OverallScore < 7 {
		summary.Recommendations = append(summary.Recommendations,
			"Consider refactoring for better maintainability")
	}

	return summary
}

// GetReviewReport generates a formatted review report
func (cr *CodeReviewer) GetReviewReport(review *CodeReview) string {
	var report strings.Builder

	report.WriteString(fmt.Sprintf("# 🤖 Automated Code Review Report\n\n"))
	report.WriteString(fmt.Sprintf("**Repository:** %s\n", review.Repository))
	report.WriteString(fmt.Sprintf("**Branch:** %s\n", review.Branch))
	report.WriteString(fmt.Sprintf("**Agent:** %s\n", review.Agent.Role))
	report.WriteString(fmt.Sprintf("**Files Reviewed:** %d\n\n", len(review.Files)))

	// Summary section
	report.WriteString("## 📊 Review Summary\n\n")
	report.WriteString(fmt.Sprintf("**Overall Score:** %d/10\n", review.Summary.OverallScore))
	report.WriteString(fmt.Sprintf("**Status:** %s\n", review.Summary.ApprovalStatus))
	report.WriteString(fmt.Sprintf("**Risk Level:** %s\n\n", review.Summary.RiskLevel))

	// Issues by category
	if len(review.Summary.IssuesByCategory) > 0 {
		report.WriteString("### Issues by Category\n")
		for category, count := range review.Summary.IssuesByCategory {
			report.WriteString(fmt.Sprintf("- **%s:** %d issues\n", category, count))
		}
		report.WriteString("\n")
	}

	// Key findings
	if len(review.Summary.KeyFindings) > 0 {
		report.WriteString("### Key Findings\n")
		for _, finding := range review.Summary.KeyFindings {
			report.WriteString(fmt.Sprintf("- %s\n", finding))
		}
		report.WriteString("\n")
	}

	// Recommendations
	if len(review.Summary.Recommendations) > 0 {
		report.WriteString("### Recommendations\n")
		for _, rec := range review.Summary.Recommendations {
			report.WriteString(fmt.Sprintf("- %s\n", rec))
		}
		report.WriteString("\n")
	}

	// Detailed file reviews
	report.WriteString("## 📁 File Reviews\n\n")

	for i, file := range review.Files {
		report.WriteString(fmt.Sprintf("### %d. %s\n", i+1, file.Path))
		report.WriteString(fmt.Sprintf("**Status:** %s\n", file.Status))
		report.WriteString(fmt.Sprintf("**Score:** %d/10\n", file.Score))

		if len(file.Issues) > 0 {
			report.WriteString("\n**Issues:**\n")
			for _, issue := range file.Issues {
				report.WriteString(fmt.Sprintf("- **%s** (%s): %s\n",
					issue.Type, issue.Severity, issue.Message))
				if issue.Suggestion != "" {
					report.WriteString(fmt.Sprintf("  *Suggestion:* %s\n", issue.Suggestion))
				}
			}
		}

		if len(file.Suggestions) > 0 {
			report.WriteString("\n**Suggestions:**\n")
			for _, suggestion := range file.Suggestions {
				report.WriteString(fmt.Sprintf("- %s\n", suggestion))
			}
		}

		if len(file.Comments) > 0 {
			report.WriteString("\n**Comments:**\n")
			for _, comment := range file.Comments {
				report.WriteString(fmt.Sprintf("- Line %d: %s\n", comment.Line, comment.Message))
			}
		}

		report.WriteString("\n")
	}

	report.WriteString("---\n")
	report.WriteString("*Generated by Ultimate SDD Framework - Automated Code Review*\n")

	return report.String()
}