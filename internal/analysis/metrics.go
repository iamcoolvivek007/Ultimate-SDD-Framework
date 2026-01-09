package analysis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// CodeMetrics represents comprehensive code quality metrics
type CodeMetrics struct {
	Files           int     `json:"files"`
	TotalLines      int     `json:"total_lines"`
	CodeLines       int     `json:"code_lines"`
	CommentLines    int     `json:"comment_lines"`
	BlankLines      int     `json:"blank_lines"`
	Complexity      float64 `json:"complexity"`
	TestCoverage    float64 `json:"test_coverage"`
	Maintainability float64 `json:"maintainability"`
	Duplication     float64 `json:"duplication"`
}

// QualityReport represents a comprehensive code quality analysis
type QualityReport struct {
	Metrics       CodeMetrics            `json:"metrics"`
	Issues        []QualityIssue         `json:"issues"`
	Recommendations []QualityRecommendation `json:"recommendations"`
	Score         QualityScore           `json:"score"`
}

// QualityIssue represents a specific code quality issue
type QualityIssue struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Type     string `json:"type"`     // warning, error, info
	Category string `json:"category"` // complexity, style, security, performance
	Message  string `json:"message"`
	Severity string `json:"severity"` // low, medium, high, critical
	RuleID   string `json:"rule_id"`
}

// QualityRecommendation represents improvement suggestions
type QualityRecommendation struct {
	Category    string   `json:"category"`
	Priority    string   `json:"priority"`
	Description string   `json:"description"`
	Files       []string `json:"files"`
	Suggestions []string `json:"suggestions"`
}

// QualityScore represents overall quality assessment
type QualityScore struct {
	Overall     float64 `json:"overall"`      // 0-100
	CodeQuality float64 `json:"code_quality"` // 0-100
	Maintainability float64 `json:"maintainability"` // 0-100
	TestQuality float64 `json:"test_quality"` // 0-100
	Security    float64 `json:"security"`     // 0-100
	Performance float64 `json:"performance"` // 0-100
	Grade       string  `json:"grade"`        // A+, A, B, C, D, F
}

// CodeAnalyzer performs comprehensive code analysis
type CodeAnalyzer struct {
	RootPath string
	Metrics  CodeMetrics
	Issues   []QualityIssue
}

// NewCodeAnalyzer creates a new code analyzer
func NewCodeAnalyzer(rootPath string) *CodeAnalyzer {
	return &CodeAnalyzer{
		RootPath: rootPath,
		Metrics:  CodeMetrics{},
		Issues:   []QualityIssue{},
	}
}

// Analyze performs comprehensive code quality analysis
func (ca *CodeAnalyzer) Analyze() (*QualityReport, error) {
	// Walk through all source files
	err := filepath.Walk(ca.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-source files
		if info.IsDir() || !ca.isSourceFile(path) {
			return nil
		}

		return ca.analyzeFile(path)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to analyze codebase: %w", err)
	}

	// Calculate derived metrics
	ca.calculateDerivedMetrics()

	// Generate recommendations
	recommendations := ca.generateRecommendations()

	// Calculate quality score
	score := ca.calculateQualityScore()

	report := &QualityReport{
		Metrics:        ca.Metrics,
		Issues:         ca.Issues,
		Recommendations: recommendations,
		Score:          score,
	}

	return report, nil
}

// analyzeFile analyzes a single source file
func (ca *CodeAnalyzer) analyzeFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	ca.Metrics.Files++

	// Basic line counting
	for _, line := range lines {
		ca.Metrics.TotalLines++
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			ca.Metrics.BlankLines++
		} else if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
			ca.Metrics.CommentLines++
		} else {
			ca.Metrics.CodeLines++
		}
	}

	// Language-specific analysis
	if strings.HasSuffix(filePath, ".go") {
		return ca.analyzeGoFile(filePath, content)
	} else if strings.HasSuffix(filePath, ".ts") || strings.HasSuffix(filePath, ".tsx") ||
	          strings.HasSuffix(filePath, ".js") || strings.HasSuffix(filePath, ".jsx") {
		return ca.analyzeTypeScriptFile(filePath, content)
	}

	return nil
}

// analyzeGoFile performs Go-specific analysis
func (ca *CodeAnalyzer) analyzeGoFile(filePath string, content []byte) error {
	// Parse Go AST
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		// Not a parse error we can analyze, skip
		return nil
	}

	// Analyze function complexity
	ast.Inspect(file, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			complexity := ca.calculateFunctionComplexity(fn)
			if complexity > 10 {
				ca.Issues = append(ca.Issues, QualityIssue{
					File:     filePath,
					Line:     fset.Position(fn.Pos()).Line,
					Type:     "warning",
					Category: "complexity",
					Message:  fmt.Sprintf("Function '%s' has high complexity (%d)", fn.Name.Name, complexity),
					Severity: "medium",
					RuleID:   "complexity-high",
				})
			}
		}
		return true
	})

	// Check for common Go issues
	ca.checkGoIssues(filePath, content)

	return nil
}

// analyzeTypeScriptFile performs TypeScript/JavaScript specific analysis
func (ca *CodeAnalyzer) analyzeTypeScriptFile(filePath string, content []byte) error {
	contentStr := string(content)

	// Check for console.log statements in production code
	if strings.Contains(contentStr, "console.log") && !strings.Contains(filePath, "test") {
		ca.Issues = append(ca.Issues, QualityIssue{
			File:     filePath,
			Type:     "warning",
			Category: "logging",
			Message:  "console.log statements found in production code",
			Severity: "low",
			RuleID:   "console-log-production",
		})
	}

	// Check for any usage
	if strings.Contains(contentStr, "any") && strings.Contains(filePath, ".ts") {
		ca.Issues = append(ca.Issues, QualityIssue{
			File:     filePath,
			Type:     "info",
			Category: "typescript",
			Message:  "'any' type usage detected - consider using specific types",
			Severity: "low",
			RuleID:   "typescript-any-usage",
		})
	}

	return nil
}

// calculateFunctionComplexity calculates cyclomatic complexity
func (ca *CodeAnalyzer) calculateFunctionComplexity(fn *ast.FuncDecl) int {
	complexity := 1 // base complexity

	ast.Inspect(fn, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause:
			complexity++
		case *ast.BinaryExpr:
			// Logical operators increase complexity
			if be, ok := n.(*ast.BinaryExpr); ok {
				if be.Op == token.LAND || be.Op == token.LOR {
					complexity++
				}
			}
		}
		return true
	})

	return complexity
}

// checkGoIssues performs Go-specific quality checks
func (ca *CodeAnalyzer) checkGoIssues(filePath string, content []byte) {
	contentStr := string(content)

	// Check for panic usage
	if strings.Contains(contentStr, "panic(") {
		ca.Issues = append(ca.Issues, QualityIssue{
			File:     filePath,
			Type:     "warning",
			Category: "error-handling",
			Message:  "panic() usage detected - consider proper error handling",
			Severity: "medium",
			RuleID:   "go-panic-usage",
		})
	}

	// Check for global variables
	lines := strings.Split(contentStr, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "var ") && !strings.Contains(line, "=") {
			// Global variable declaration
			ca.Issues = append(ca.Issues, QualityIssue{
				File:     filePath,
				Line:     i + 1,
				Type:     "info",
				Category: "architecture",
				Message:  "Global variable detected - consider encapsulation",
				Severity: "low",
				RuleID:   "go-global-variables",
			})
		}
	}
}

// calculateDerivedMetrics calculates metrics derived from raw data
func (ca *CodeAnalyzer) calculateDerivedMetrics() {
	if ca.Metrics.TotalLines > 0 {
		// Comment ratio
		commentRatio := float64(ca.Metrics.CommentLines) / float64(ca.Metrics.TotalLines)

		// Simple maintainability index (placeholder - could be more sophisticated)
		if commentRatio > 0.2 {
			ca.Metrics.Maintainability = 85.0
		} else if commentRatio > 0.1 {
			ca.Metrics.Maintainability = 70.0
		} else {
			ca.Metrics.Maintainability = 50.0
		}
	}

	// Estimate test coverage (placeholder - would need actual test analysis)
	ca.Metrics.TestCoverage = 75.0 // This would be calculated from actual test execution

	// Estimate code duplication (placeholder - would need clone detection)
	ca.Metrics.Duplication = 5.0 // Percentage of duplicated code
}

// generateRecommendations creates improvement suggestions
func (ca *CodeAnalyzer) generateRecommendations() []QualityRecommendation {
	recommendations := []QualityRecommendation{}

	// Code quality recommendations
	if ca.Metrics.Maintainability < 70 {
		recommendations = append(recommendations, QualityRecommendation{
			Category:    "Code Quality",
			Priority:    "High",
			Description: "Improve code maintainability through better documentation and structure",
			Suggestions: []string{
				"Add comprehensive code comments",
				"Break down complex functions",
				"Implement consistent naming conventions",
				"Add API documentation",
			},
		})
	}

	// Testing recommendations
	if ca.Metrics.TestCoverage < 80 {
		recommendations = append(recommendations, QualityRecommendation{
			Category:    "Testing",
			Priority:    "High",
			Description: "Increase test coverage to ensure code reliability",
			Suggestions: []string{
				"Add unit tests for all public functions",
				"Implement integration tests",
				"Add end-to-end test scenarios",
				"Set up automated testing pipeline",
			},
		})
	}

	// Complexity recommendations
	complexityIssues := 0
	for _, issue := range ca.Issues {
		if issue.Category == "complexity" {
			complexityIssues++
		}
	}

	if complexityIssues > 0 {
		recommendations = append(recommendations, QualityRecommendation{
			Category:    "Code Complexity",
			Priority:    "Medium",
			Description: fmt.Sprintf("Reduce complexity in %d functions", complexityIssues),
			Suggestions: []string{
				"Extract methods from complex functions",
				"Implement early returns",
				"Use design patterns to reduce complexity",
				"Add comprehensive unit tests",
			},
		})
	}

	return recommendations
}

// calculateQualityScore computes overall quality metrics
func (ca *CodeAnalyzer) calculateQualityScore() QualityScore {
	score := QualityScore{}

	// Code quality based on issues
	issuePenalty := 0
	for _, issue := range ca.Issues {
		switch issue.Severity {
		case "critical":
			issuePenalty += 20
		case "high":
			issuePenalty += 10
		case "medium":
			issuePenalty += 5
		case "low":
			issuePenalty += 2
		}
	}

	score.CodeQuality = 100 - float64(issuePenalty)
	if score.CodeQuality < 0 {
		score.CodeQuality = 0
	}

	// Maintainability based on metrics
	score.Maintainability = ca.Metrics.Maintainability

	// Test quality based on coverage
	score.TestQuality = ca.Metrics.TestCoverage

	// Security and performance (placeholder - would need specific analysis)
	score.Security = 85.0
	score.Performance = 80.0

	// Overall score (weighted average)
	score.Overall = (score.CodeQuality*0.3 + score.Maintainability*0.2 +
		score.TestQuality*0.2 + score.Security*0.15 + score.Performance*0.15)

	// Assign grade
	switch {
	case score.Overall >= 95:
		score.Grade = "A+"
	case score.Overall >= 90:
		score.Grade = "A"
	case score.Overall >= 80:
		score.Grade = "B"
	case score.Overall >= 70:
		score.Grade = "C"
	case score.Overall >= 60:
		score.Grade = "D"
	default:
		score.Grade = "F"
	}

	return score
}

// isSourceFile checks if a file is a source code file we should analyze
func (ca *CodeAnalyzer) isSourceFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))

	// Supported file extensions
	sourceExts := []string{".go", ".ts", ".tsx", ".js", ".jsx", ".py", ".rs", ".java", ".cpp", ".c"}

	for _, sourceExt := range sourceExts {
		if ext == sourceExt {
			return true
		}
	}

	return false
}

// GetSummary returns a human-readable summary of the analysis
func (report *QualityReport) GetSummary() string {
	var summary strings.Builder

	summary.WriteString(fmt.Sprintf("## Code Quality Analysis Summary\n\n"))
	summary.WriteString(fmt.Sprintf("**Overall Grade:** %s (%.1f/100)\n\n", report.Score.Grade, report.Score.Overall))

	summary.WriteString("### Metrics\n")
	summary.WriteString(fmt.Sprintf("- Files Analyzed: %d\n", report.Metrics.Files))
	summary.WriteString(fmt.Sprintf("- Total Lines: %d\n", report.Metrics.TotalLines))
	summary.WriteString(fmt.Sprintf("- Code Lines: %d\n", report.Metrics.CodeLines))
	summary.WriteString(fmt.Sprintf("- Comment Lines: %d (%.1f%%)\n", report.Metrics.CommentLines,
		float64(report.Metrics.CommentLines)/float64(report.Metrics.TotalLines)*100))
	summary.WriteString(fmt.Sprintf("- Test Coverage: %.1f%%\n", report.Metrics.TestCoverage))
	summary.WriteString(fmt.Sprintf("- Maintainability: %.1f/100\n", report.Metrics.Maintainability))

	summary.WriteString("\n### Issues by Severity\n")
	severityCount := make(map[string]int)
	for _, issue := range report.Issues {
		severityCount[issue.Severity]++
	}

	for _, severity := range []string{"critical", "high", "medium", "low"} {
		if count, exists := severityCount[severity]; exists {
			summary.WriteString(fmt.Sprintf("- %s: %d\n", strings.Title(severity), count))
		}
	}

	if len(report.Recommendations) > 0 {
		summary.WriteString("\n### Key Recommendations\n")
		for i, rec := range report.Recommendations {
			if i >= 3 { // Show only top 3
				break
			}
			summary.WriteString(fmt.Sprintf("**%s** (%s): %s\n", rec.Category, rec.Priority, rec.Description))
		}
	}

	return summary.String()
}