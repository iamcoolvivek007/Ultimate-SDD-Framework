package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// Constitution represents a project constitution
type Constitution struct {
	Version          string          `json:"version"`
	ProjectName      string          `json:"project_name"`
	Description      string          `json:"description"`
	RatificationDate string          `json:"ratification_date"`
	LastAmended      string          `json:"last_amended"`
	Principles       []Principle     `json:"principles"`
	CodingStandards  []string        `json:"coding_standards"`
	QualityRules     []string        `json:"quality_rules"`
	Governance       GovernanceRules `json:"governance"`
}

// Principle represents a core principle
type Principle struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Rationale   string `json:"rationale"`
}

// GovernanceRules defines how the constitution is maintained
type GovernanceRules struct {
	AmendmentProcess string   `json:"amendment_process"`
	VersioningPolicy string   `json:"versioning_policy"`
	ReviewSchedule   string   `json:"review_schedule"`
	Maintainers      []string `json:"maintainers"`
}

// NewConstitutionCmd creates the constitution command
func NewConstitutionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "constitution [description]",
		Short: "📜 Create or update project constitution",
		Long: `Create or update the project constitution - the governing document that defines
project principles, coding standards, and quality requirements.

The constitution serves as the source of truth for all development decisions
and is referenced by AI agents during planning and implementation.`,
		Example: `  viki constitution "Create principles for code quality and testing"
  viki constitution --view
  viki constitution --amend "Add new security principle"`,
		Run: runConstitution,
	}

	cmd.Flags().Bool("view", false, "View current constitution")
	cmd.Flags().Bool("amend", false, "Amend existing constitution")
	cmd.Flags().Bool("interactive", false, "Interactive mode")

	return cmd
}

func runConstitution(cmd *cobra.Command, args []string) {
	viewMode, _ := cmd.Flags().GetBool("view")
	amendMode, _ := cmd.Flags().GetBool("amend")

	constitutionPath := filepath.Join(".viki", "constitution.md")

	// Ensure .viki directory exists
	os.MkdirAll(".viki", 0755)

	if viewMode {
		viewConstitution(constitutionPath)
		return
	}

	description := ""
	if len(args) > 0 {
		description = strings.Join(args, " ")
	}

	if amendMode {
		amendConstitution(constitutionPath, description)
		return
	}

	createConstitution(constitutionPath, description)
}

func viewConstitution(path string) {
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		fmt.Println("❌ No constitution found. Create one with: viki constitution \"your principles\"")
		return
	}
	if err != nil {
		fmt.Printf("❌ Error reading constitution: %v\n", err)
		return
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99"))

	fmt.Println(titleStyle.Render("\n📜 Project Constitution"))
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println(string(content))
}

func createConstitution(path string, description string) {
	// Check if constitution already exists
	if _, err := os.Stat(path); err == nil {
		fmt.Println("⚠️  Constitution already exists. Use --amend to update.")
		return
	}

	projectName := filepath.Base(mustGetwd())
	today := time.Now().Format("2006-01-02")

	constitution := generateConstitutionTemplate(projectName, today, description)

	if err := os.WriteFile(path, []byte(constitution), 0644); err != nil {
		fmt.Printf("❌ Error creating constitution: %v\n", err)
		return
	}

	successStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("42"))

	fmt.Println(successStyle.Render("\n✅ Constitution created!"))
	fmt.Printf("📄 Location: %s\n", path)
	fmt.Println("\n💡 Next steps:")
	fmt.Println("   1. Review and customize the constitution")
	fmt.Println("   2. Add project-specific principles")
	fmt.Println("   3. Run 'viki specify' to start development")
}

func amendConstitution(path string, amendment string) {
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		fmt.Println("❌ No constitution to amend. Create one first.")
		return
	}
	if err != nil {
		fmt.Printf("❌ Error reading constitution: %v\n", err)
		return
	}

	// Add amendment section
	today := time.Now().Format("2006-01-02")
	amendmentSection := fmt.Sprintf(`

---

## Amendment (%s)

%s
`, today, amendment)

	newContent := string(content) + amendmentSection

	if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
		fmt.Printf("❌ Error updating constitution: %v\n", err)
		return
	}

	fmt.Println("✅ Constitution amended successfully!")
	fmt.Printf("📄 Amendment added: %s\n", amendment)
}

func generateConstitutionTemplate(projectName, date, description string) string {
	if description == "" {
		description = "A well-crafted software project built with quality and maintainability in mind."
	}

	return fmt.Sprintf(`---
title: Project Constitution
version: 1.0.0
ratification_date: %s
last_amended: %s
---

# 📜 %s Constitution

## Project Description

%s

---

## 🎯 Core Principles

### Principle 1: Code Quality First

**Statement**: All code must be clean, readable, and maintainable.

**Rationale**: Quality code reduces technical debt, improves developer experience, and ensures long-term project health.

**Rules**:
- MUST use meaningful variable and function names
- MUST keep functions focused and under 50 lines
- MUST add comments for non-obvious logic
- SHOULD follow established design patterns

---

### Principle 2: Test-Driven Development

**Statement**: Tests are written before or alongside implementation.

**Rationale**: TDD ensures code correctness, provides documentation, and enables safe refactoring.

**Rules**:
- MUST have unit tests for all business logic
- MUST maintain minimum 80%% test coverage
- SHOULD write integration tests for critical paths
- MUST run tests before committing

---

### Principle 3: Security by Design

**Statement**: Security is a first-class concern, not an afterthought.

**Rationale**: Security vulnerabilities can have severe consequences; prevention is better than remediation.

**Rules**:
- MUST validate all user input
- MUST never store secrets in code
- MUST use secure authentication methods
- SHOULD conduct regular security reviews

---

### Principle 4: Documentation as Code

**Statement**: Documentation is maintained alongside code and stays current.

**Rationale**: Good documentation improves onboarding, reduces cognitive load, and preserves knowledge.

**Rules**:
- MUST document public APIs
- MUST keep README up to date
- SHOULD use inline documentation
- MUST document architectural decisions

---

## 💻 Coding Standards

1. Use consistent formatting (enforced by linters)
2. Follow language-specific conventions and idioms
3. Handle errors explicitly - no silent failures
4. Use dependency injection for testability
5. Prefer composition over inheritance
6. Keep dependencies minimal and updated

---

## ✅ Quality Requirements

1. All PRs require code review
2. CI/CD must pass before merge
3. No known security vulnerabilities
4. Performance benchmarks must pass
5. Accessibility standards must be met (if applicable)

---

## 🏛️ Governance

### Amendment Process

1. Propose amendment via pull request
2. Discuss in team meeting
3. Require majority approval
4. Update version number appropriately

### Versioning Policy

- **MAJOR**: Breaking changes to principles or governance
- **MINOR**: New principles or significant expansions
- **PATCH**: Clarifications and minor updates

### Review Schedule

- Quarterly review of all principles
- Annual comprehensive audit

---

*This constitution guides all development decisions. When in doubt, refer to these principles.*
`, date, date, projectName, description)
}

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "project"
	}
	return wd
}

// NewClarifyCmd creates the clarify command
func NewClarifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clarify",
		Short: "❓ Clarify specifications with structured Q&A",
		Long: `Analyze current specifications and identify areas that need clarification.

Generates structured questions to fill gaps in:
- Requirements completeness
- Edge cases
- Technical constraints
- User flows
- Error handling`,
		Run: runClarify,
	}
}

func runClarify(cmd *cobra.Command, args []string) {
	fmt.Println("🔍 Analyzing specifications for gaps...")

	specPath := filepath.Join(".sdd", "spec.md")
	planPath := filepath.Join(".sdd", "plan.md")

	// Check if spec exists
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		fmt.Println("❌ No specification found. Run 'viki specify' first.")
		return
	}

	specContent, err := os.ReadFile(specPath)
	if err != nil {
		fmt.Printf("❌ Error reading spec: %v\n", err)
		return
	}

	// Generate clarification questions
	questions := generateClarificationQuestions(string(specContent))

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("220"))

	fmt.Println(titleStyle.Render("\n❓ Clarification Questions"))
	fmt.Println(strings.Repeat("─", 50))

	for i, q := range questions {
		fmt.Printf("\n%d. %s\n", i+1, q.Question)
		fmt.Printf("   📝 Why: %s\n", q.Reason)
		if q.Suggestion != "" {
			fmt.Printf("   💡 Suggestion: %s\n", q.Suggestion)
		}
	}

	// Save clarification report
	reportPath := filepath.Join(".sdd", "clarifications.md")
	if planPath != "" {
		// Check plan status
		_, planErr := os.Stat(planPath)
		if planErr == nil {
			fmt.Println("\n⚠️  Plan already exists. Address these before proceeding.")
		}
	}

	report := generateClarificationReport(questions)
	os.WriteFile(reportPath, []byte(report), 0644)
	fmt.Printf("\n📄 Report saved to: %s\n", reportPath)
}

type ClarificationQuestion struct {
	Question   string
	Reason     string
	Suggestion string
	Category   string
}

func generateClarificationQuestions(specContent string) []ClarificationQuestion {
	questions := []ClarificationQuestion{}

	// Check for missing sections
	if !strings.Contains(specContent, "User Stor") {
		questions = append(questions, ClarificationQuestion{
			Question:   "Are there specific user stories that define the core workflows?",
			Reason:     "No user stories section detected in specification",
			Suggestion: "Add user stories in 'As a [user], I want to [action], so that [benefit]' format",
			Category:   "requirements",
		})
	}

	if !strings.Contains(strings.ToLower(specContent), "error") {
		questions = append(questions, ClarificationQuestion{
			Question:   "How should the system handle error conditions?",
			Reason:     "No error handling strategy mentioned",
			Suggestion: "Define error types, messages, and recovery procedures",
			Category:   "error_handling",
		})
	}

	if !strings.Contains(strings.ToLower(specContent), "edge case") {
		questions = append(questions, ClarificationQuestion{
			Question:   "What edge cases need to be handled?",
			Reason:     "No edge cases documented",
			Suggestion: "Consider empty states, maximum limits, concurrent access, etc.",
			Category:   "edge_cases",
		})
	}

	if !strings.Contains(strings.ToLower(specContent), "performance") {
		questions = append(questions, ClarificationQuestion{
			Question:   "What are the performance requirements?",
			Reason:     "No performance criteria specified",
			Suggestion: "Define response times, throughput, and scalability needs",
			Category:   "requirements",
		})
	}

	if !strings.Contains(strings.ToLower(specContent), "security") {
		questions = append(questions, ClarificationQuestion{
			Question:   "What security requirements apply?",
			Reason:     "No security considerations documented",
			Suggestion: "Define authentication, authorization, and data protection needs",
			Category:   "security",
		})
	}

	// Add default questions if none found
	if len(questions) == 0 {
		questions = append(questions, ClarificationQuestion{
			Question: "Is the current specification complete enough to proceed to planning?",
			Reason:   "Basic sections appear present",
			Category: "general",
		})
	}

	return questions
}

func generateClarificationReport(questions []ClarificationQuestion) string {
	var sb strings.Builder

	sb.WriteString("# Clarification Report\n\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n\n", time.Now().Format("2006-01-02 15:04")))
	sb.WriteString("---\n\n")

	// Group by category
	categories := map[string][]ClarificationQuestion{}
	for _, q := range questions {
		categories[q.Category] = append(categories[q.Category], q)
	}

	for category, qs := range categories {
		sb.WriteString(fmt.Sprintf("## %s\n\n", strings.Title(strings.ReplaceAll(category, "_", " "))))
		for _, q := range qs {
			sb.WriteString(fmt.Sprintf("### %s\n\n", q.Question))
			sb.WriteString(fmt.Sprintf("**Reason**: %s\n\n", q.Reason))
			if q.Suggestion != "" {
				sb.WriteString(fmt.Sprintf("**Suggestion**: %s\n\n", q.Suggestion))
			}
			sb.WriteString("**Answer**: _[To be filled]_\n\n")
			sb.WriteString("---\n\n")
		}
	}

	return sb.String()
}

// NewChecklistCmd creates the checklist command
func NewChecklistCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "checklist",
		Short: "✅ Generate quality checklist",
		Long: `Generate a quality checklist for the current specification and plan.

Creates checklists for:
- Requirements completeness
- Technical readiness
- Security considerations
- Testing strategy
- Documentation needs`,
		Run: runChecklist,
	}
}

func runChecklist(cmd *cobra.Command, args []string) {
	fmt.Println("📋 Generating quality checklist...")

	os.MkdirAll(filepath.Join(".sdd", "checklists"), 0755)

	// Generate various checklists
	checklists := map[string]string{
		"requirements.md": generateRequirementsChecklist(),
		"technical.md":    generateTechnicalChecklist(),
		"security.md":     generateSecurityChecklist(),
		"testing.md":      generateTestingChecklist(),
	}

	for filename, content := range checklists {
		path := filepath.Join(".sdd", "checklists", filename)
		os.WriteFile(path, []byte(content), 0644)
		fmt.Printf("  ✅ Created: %s\n", path)
	}

	fmt.Println("\n✅ Checklists generated!")
	fmt.Println("💡 Review each checklist and mark items as complete [X] before proceeding.")
}

func generateRequirementsChecklist() string {
	return `# Requirements Checklist

## Completeness
- [ ] All user stories documented
- [ ] Acceptance criteria defined for each story
- [ ] Scope boundaries clearly defined
- [ ] Out-of-scope items listed
- [ ] Dependencies identified

## Clarity
- [ ] Requirements are unambiguous
- [ ] Business terminology defined
- [ ] User roles and permissions clear
- [ ] Data requirements specified

## Validation
- [ ] Stakeholder approval obtained
- [ ] Feasibility assessed
- [ ] Conflicts resolved
`
}

func generateTechnicalChecklist() string {
	return `# Technical Readiness Checklist

## Architecture
- [ ] System design documented
- [ ] Component responsibilities defined
- [ ] Data flow documented
- [ ] Integration points identified
- [ ] Technology choices justified

## Infrastructure
- [ ] Development environment ready
- [ ] CI/CD pipeline configured
- [ ] Deployment strategy defined
- [ ] Monitoring approach planned

## Code Quality
- [ ] Coding standards documented
- [ ] Linting/formatting configured
- [ ] Code review process defined
`
}

func generateSecurityChecklist() string {
	return `# Security Checklist

## Authentication & Authorization
- [ ] Authentication method defined
- [ ] Authorization model documented
- [ ] Session management planned
- [ ] Password requirements specified

## Data Protection
- [ ] Sensitive data identified
- [ ] Encryption requirements defined
- [ ] Data retention policy set
- [ ] PII handling documented

## API Security
- [ ] Input validation planned
- [ ] Rate limiting considered
- [ ] CORS configuration defined
- [ ] API authentication method chosen
`
}

func generateTestingChecklist() string {
	return `# Testing Checklist

## Test Strategy
- [ ] Test approach documented
- [ ] Test environments defined
- [ ] Coverage targets set
- [ ] Test data strategy planned

## Test Types
- [ ] Unit tests required for business logic
- [ ] Integration tests for critical paths
- [ ] E2E tests for user flows
- [ ] Performance tests for SLAs

## Quality Gates
- [ ] Minimum coverage threshold defined
- [ ] All tests must pass before merge
- [ ] Security scans automated
`
}
