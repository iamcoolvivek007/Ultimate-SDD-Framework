package collaboration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"ultimate-sdd-framework/internal/agents"
)

// Team represents a collaborative development team
type Team struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Members     []TeamMember      `json:"members"`
	Rules       TeamRules         `json:"rules"`
	Projects    []TeamProject     `json:"projects"`
	Knowledge   TeamKnowledge     `json:"knowledge"`
	Created     time.Time         `json:"created"`
	LastUpdated time.Time         `json:"last_updated"`
}

// TeamMember represents a team member
type TeamMember struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`     // lead, senior, junior, qa, etc.
	Skills   []string  `json:"skills"`
	Joined   time.Time `json:"joined"`
	LastActive time.Time `json:"last_active"`
}

// TeamRules represents shared team coding standards and rules
type TeamRules struct {
	CodingStandards    []RuleDefinition `json:"coding_standards"`
	CodeReviewRules    []RuleDefinition `json:"code_review_rules"`
	TestingStandards   []RuleDefinition `json:"testing_standards"`
	SecurityPolicies   []RuleDefinition `json:"security_policies"`
	PerformanceRules   []RuleDefinition `json:"performance_rules"`
	DocumentationRules []RuleDefinition `json:"documentation_rules"`
}

// RuleDefinition represents a team rule
type RuleDefinition struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Severity    string    `json:"severity"` // mandatory, recommended, optional
	Examples    []string  `json:"examples"`
	Exceptions  []string  `json:"exceptions"`
	CreatedBy   string    `json:"created_by"`
	Created     time.Time `json:"created"`
	Votes       int       `json:"votes"`
}

// TeamProject represents a project within the team
type TeamProject struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Path        string            `json:"path"`
	TechStack   []string          `json:"tech_stack"`
	Status      string            `json:"status"` // active, archived, planning
	Members     []string          `json:"members"` // member IDs
	Created     time.Time         `json:"created"`
	LastCommit  time.Time         `json:"last_commit"`
}

// TeamKnowledge represents shared knowledge base
type TeamKnowledge struct {
	BestPractices    []KnowledgeItem `json:"best_practices"`
	CommonIssues     []KnowledgeItem `json:"common_issues"`
	ArchitectureDocs []KnowledgeItem `json:"architecture_docs"`
	CodePatterns     []CodePattern   `json:"code_patterns"`
	DecisionLog      []Decision      `json:"decision_log"`
}

// KnowledgeItem represents a piece of team knowledge
type KnowledgeItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	Author      string    `json:"author"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Views       int       `json:"views"`
	Helpful     int       `json:"helpful"`
}

// CodePattern represents a reusable code pattern
type CodePattern struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	Code        string    `json:"code"`
	UseCase     string    `json:"use_case"`
	Author      string    `json:"author"`
	Created     time.Time `json:"created"`
	UsageCount  int       `json:"usage_count"`
}

// Decision represents an architectural or important decision
type Decision struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Context     string    `json:"context"`
	Decision    string    `json:"decision"`
	Alternatives []string `json:"alternatives"`
	Consequences string   `json:"consequences"`
	MadeBy      string    `json:"made_by"`
	Date        time.Time `json:"date"`
	Status      string    `json:"status"` // implemented, pending, rejected
}

// TeamCollaboration manages team-based development
type TeamCollaboration struct {
	teamData    Team
	dataPath    string
	projectRoot string
	agentSvc    *agents.AgentService
}

// NewTeamCollaboration creates a new team collaboration system
func NewTeamCollaboration(projectRoot string) (*TeamCollaboration, error) {
	dataPath := filepath.Join(projectRoot, ".sdd", "team.json")

	agentSvc := agents.NewAgentService(projectRoot)
	if err := agentSvc.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize agent service: %w", err)
	}

	tc := &TeamCollaboration{
		dataPath:    dataPath,
		projectRoot: projectRoot,
		agentSvc:    agentSvc,
	}

	// Load existing team data
	if err := tc.loadTeamData(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load team data: %w", err)
		}
		// Initialize default team
		tc.teamData = Team{
			ID:          "default",
			Name:        "Development Team",
			Description: "Default development team",
			Created:     time.Now(),
			LastUpdated: time.Now(),
		}
	}

	return tc, nil
}

// CreateTeam initializes a new team
func (tc *TeamCollaboration) CreateTeam(name, description string) (*Team, error) {
	tc.teamData = Team{
		ID:          generateTeamID(),
		Name:        name,
		Description: description,
		Members:     []TeamMember{},
		Rules:       TeamRules{},
		Projects:    []TeamProject{},
		Knowledge:   TeamKnowledge{},
		Created:     time.Now(),
		LastUpdated: time.Now(),
	}

	return &tc.teamData, tc.saveTeamData()
}

// AddTeamMember adds a new member to the team
func (tc *TeamCollaboration) AddTeamMember(name, email, role string, skills []string) (*TeamMember, error) {
	member := TeamMember{
		ID:         generateMemberID(),
		Name:       name,
		Email:      email,
		Role:       role,
		Skills:     skills,
		Joined:     time.Now(),
		LastActive: time.Now(),
	}

	tc.teamData.Members = append(tc.teamData.Members, member)
	tc.teamData.LastUpdated = time.Now()

	return &member, tc.saveTeamData()
}

// AddTeamRule adds a new team rule
func (tc *TeamCollaboration) AddTeamRule(category, title, description, severity, createdBy string, examples []string) (*RuleDefinition, error) {
	rule := RuleDefinition{
		ID:          generateRuleID(),
		Title:       title,
		Description: description,
		Category:    category,
		Severity:    severity,
		Examples:    examples,
		CreatedBy:   createdBy,
		Created:     time.Now(),
		Votes:       1, // Creator automatically votes
	}

	// Add to appropriate category
	switch category {
	case "coding_standards":
		tc.teamData.Rules.CodingStandards = append(tc.teamData.Rules.CodingStandards, rule)
	case "code_review":
		tc.teamData.Rules.CodeReviewRules = append(tc.teamData.Rules.CodeReviewRules, rule)
	case "testing":
		tc.teamData.Rules.TestingStandards = append(tc.teamData.Rules.TestingStandards, rule)
	case "security":
		tc.teamData.Rules.SecurityPolicies = append(tc.teamData.Rules.SecurityPolicies, rule)
	case "performance":
		tc.teamData.Rules.PerformanceRules = append(tc.teamData.Rules.PerformanceRules, rule)
	case "documentation":
		tc.teamData.Rules.DocumentationRules = append(tc.teamData.Rules.DocumentationRules, rule)
	}

	tc.teamData.LastUpdated = time.Now()
	return &rule, tc.saveTeamData()
}

// AddKnowledgeItem adds a new knowledge item to the team knowledge base
func (tc *TeamCollaboration) AddKnowledgeItem(title, content, category, author string, tags []string) (*KnowledgeItem, error) {
	item := KnowledgeItem{
		ID:       generateKnowledgeID(),
		Title:    title,
		Content:  content,
		Category: category,
		Tags:     tags,
		Author:   author,
		Created:  time.Now(),
		Updated:  time.Now(),
		Views:    0,
		Helpful:  0,
	}

	// Add to appropriate category
	switch category {
	case "best_practices":
		tc.teamData.Knowledge.BestPractices = append(tc.teamData.Knowledge.BestPractices, item)
	case "common_issues":
		tc.teamData.Knowledge.CommonIssues = append(tc.teamData.Knowledge.CommonIssues, item)
	case "architecture":
		tc.teamData.Knowledge.ArchitectureDocs = append(tc.teamData.Knowledge.ArchitectureDocs, item)
	}

	tc.teamData.LastUpdated = time.Now()
	return &item, tc.saveTeamData()
}

// AddCodePattern adds a reusable code pattern
func (tc *TeamCollaboration) AddCodePattern(name, description, language, code, useCase, author string) (*CodePattern, error) {
	pattern := CodePattern{
		ID:          generatePatternID(),
		Name:        name,
		Description: description,
		Language:    language,
		Code:        code,
		UseCase:     useCase,
		Author:      author,
		Created:     time.Now(),
		UsageCount:  0,
	}

	tc.teamData.Knowledge.CodePatterns = append(tc.teamData.Knowledge.CodePatterns, pattern)
	tc.teamData.LastUpdated = time.Now()

	return &pattern, tc.saveTeamData()
}

// RecordDecision records an important team decision
func (tc *TeamCollaboration) RecordDecision(title, context, decision string, alternatives []string, consequences, madeBy string) (*Decision, error) {
	teamDecision := Decision{
		ID:           generateDecisionID(),
		Title:        title,
		Context:      context,
		Decision:     decision,
		Alternatives: alternatives,
		Consequences: consequences,
		MadeBy:       madeBy,
		Date:         time.Now(),
		Status:       "implemented",
	}

	tc.teamData.Knowledge.DecisionLog = append(tc.teamData.Knowledge.DecisionLog, teamDecision)
	tc.teamData.LastUpdated = time.Now()

	return &teamDecision, tc.saveTeamData()
}

// GetTeamRules returns all team rules organized by category
func (tc *TeamCollaboration) GetTeamRules() TeamRules {
	return tc.teamData.Rules
}

// GetTeamKnowledge returns the team knowledge base
func (tc *TeamCollaboration) GetTeamKnowledge() TeamKnowledge {
	return tc.teamData.Knowledge
}

// SearchKnowledge searches the team knowledge base
func (tc *TeamCollaboration) SearchKnowledge(query string, category string) []KnowledgeItem {
	results := []KnowledgeItem{}
	query = strings.ToLower(query)

	searchInItems := func(items []KnowledgeItem) {
		for _, item := range items {
			if (category == "" || item.Category == category) &&
			   (strings.Contains(strings.ToLower(item.Title), query) ||
			    strings.Contains(strings.ToLower(item.Content), query)) {
				results = append(results, item)
			}
		}
	}

	searchInItems(tc.teamData.Knowledge.BestPractices)
	searchInItems(tc.teamData.Knowledge.CommonIssues)
	searchInItems(tc.teamData.Knowledge.ArchitectureDocs)

	// Sort by relevance (simplified - just by title match)
	sort.Slice(results, func(i, j int) bool {
		iTitle := strings.Contains(strings.ToLower(results[i].Title), query)
		jTitle := strings.Contains(strings.ToLower(results[j].Title), query)
		return iTitle && !jTitle // Title matches first
	})

	return results
}

// GetCodePatterns returns code patterns filtered by criteria
func (tc *TeamCollaboration) GetCodePatterns(language, useCase string) []CodePattern {
	patterns := []CodePattern{}

	for _, pattern := range tc.teamData.Knowledge.CodePatterns {
		if (language == "" || pattern.Language == language) &&
		   (useCase == "" || strings.Contains(strings.ToLower(pattern.UseCase), strings.ToLower(useCase))) {
			patterns = append(patterns, pattern)
		}
	}

	// Sort by usage count (most used first)
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].UsageCount > patterns[j].UsageCount
	})

	return patterns
}

// GenerateTeamReport creates a comprehensive team report
func (tc *TeamCollaboration) GenerateTeamReport() string {
	var report strings.Builder

	report.WriteString(fmt.Sprintf("# 👥 Team Collaboration Report - %s\n\n", tc.teamData.Name))
	report.WriteString(fmt.Sprintf("**Description:** %s\n\n", tc.teamData.Description))
	report.WriteString(fmt.Sprintf("**Members:** %d\n", len(tc.teamData.Members)))
	report.WriteString(fmt.Sprintf("**Projects:** %d\n", len(tc.teamData.Projects)))
	report.WriteString(fmt.Sprintf("**Created:** %s\n\n", tc.teamData.Created.Format("2006-01-02")))

	// Team Members
	if len(tc.teamData.Members) > 0 {
		report.WriteString("## 👤 Team Members\n\n")
		for _, member := range tc.teamData.Members {
			report.WriteString(fmt.Sprintf("### %s (%s)\n", member.Name, member.Role))
			report.WriteString(fmt.Sprintf("**Email:** %s\n", member.Email))
			report.WriteString(fmt.Sprintf("**Joined:** %s\n", member.Joined.Format("2006-01-02")))
			if len(member.Skills) > 0 {
				report.WriteString(fmt.Sprintf("**Skills:** %s\n", strings.Join(member.Skills, ", ")))
			}
			report.WriteString("\n")
		}
	}

	// Team Rules Summary
	totalRules := len(tc.teamData.Rules.CodingStandards) + len(tc.teamData.Rules.CodeReviewRules) +
	              len(tc.teamData.Rules.TestingStandards) + len(tc.teamData.Rules.SecurityPolicies) +
	              len(tc.teamData.Rules.PerformanceRules) + len(tc.teamData.Rules.DocumentationRules)

	if totalRules > 0 {
		report.WriteString("## 📋 Team Rules\n\n")
		report.WriteString(fmt.Sprintf("**Total Rules:** %d\n\n", totalRules))

		ruleCategories := map[string][]RuleDefinition{
			"Coding Standards":    tc.teamData.Rules.CodingStandards,
			"Code Review":         tc.teamData.Rules.CodeReviewRules,
			"Testing Standards":   tc.teamData.Rules.TestingStandards,
			"Security Policies":   tc.teamData.Rules.SecurityPolicies,
			"Performance Rules":   tc.teamData.Rules.PerformanceRules,
			"Documentation Rules": tc.teamData.Rules.DocumentationRules,
		}

		for category, rules := range ruleCategories {
			if len(rules) > 0 {
				report.WriteString(fmt.Sprintf("### %s (%d rules)\n", category, len(rules)))
				for _, rule := range rules {
					report.WriteString(fmt.Sprintf("- **%s**: %s\n", rule.Title, rule.Description))
				}
				report.WriteString("\n")
			}
		}
	}

	// Knowledge Base Summary
	totalKnowledge := len(tc.teamData.Knowledge.BestPractices) + len(tc.teamData.Knowledge.CommonIssues) +
	                  len(tc.teamData.Knowledge.ArchitectureDocs) + len(tc.teamData.Knowledge.CodePatterns) +
	                  len(tc.teamData.Knowledge.DecisionLog)

	if totalKnowledge > 0 {
		report.WriteString("## 🧠 Knowledge Base\n\n")
		report.WriteString(fmt.Sprintf("**Total Items:** %d\n\n", totalKnowledge))

		report.WriteString(fmt.Sprintf("- **Best Practices:** %d\n", len(tc.teamData.Knowledge.BestPractices)))
		report.WriteString(fmt.Sprintf("- **Common Issues:** %d\n", len(tc.teamData.Knowledge.CommonIssues)))
		report.WriteString(fmt.Sprintf("- **Architecture Docs:** %d\n", len(tc.teamData.Knowledge.ArchitectureDocs)))
		report.WriteString(fmt.Sprintf("- **Code Patterns:** %d\n", len(tc.teamData.Knowledge.CodePatterns)))
		report.WriteString(fmt.Sprintf("- **Decision Log:** %d\n", len(tc.teamData.Knowledge.DecisionLog)))
		report.WriteString("\n")
	}

	// Popular Code Patterns
	if len(tc.teamData.Knowledge.CodePatterns) > 0 {
		report.WriteString("### 🔥 Popular Code Patterns\n\n")
		patterns := tc.GetCodePatterns("", "")
		for i, pattern := range patterns {
			if i >= 5 { // Top 5
				break
			}
			report.WriteString(fmt.Sprintf("- **%s** (%s): Used %d times\n",
				pattern.Name, pattern.Language, pattern.UsageCount))
		}
		report.WriteString("\n")
	}

	// Recent Decisions
	if len(tc.teamData.Knowledge.DecisionLog) > 0 {
		report.WriteString("### 🏛️ Recent Decisions\n\n")
		// Get last 3 decisions
		decisions := tc.teamData.Knowledge.DecisionLog
		start := len(decisions) - 3
		if start < 0 {
			start = 0
		}

		for i := len(decisions) - 1; i >= start; i-- {
			decision := decisions[i]
			report.WriteString(fmt.Sprintf("- **%s** (%s): %s\n",
				decision.Title, decision.Date.Format("2006-01-02"), decision.Decision))
		}
		report.WriteString("\n")
	}

	report.WriteString(fmt.Sprintf("*Report generated: %s*\n", time.Now().Format("2006-01-02 15:04:05")))

	return report.String()
}

// Private methods

func (tc *TeamCollaboration) loadTeamData() error {
	data, err := os.ReadFile(tc.dataPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &tc.teamData)
}

func (tc *TeamCollaboration) saveTeamData() error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(tc.dataPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(tc.teamData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tc.dataPath, data, 0644)
}

// ID generation functions
func generateTeamID() string      { return fmt.Sprintf("team_%d", time.Now().Unix()) }
func generateMemberID() string    { return fmt.Sprintf("member_%d", time.Now().UnixNano()) }
func generateRuleID() string      { return fmt.Sprintf("rule_%d", time.Now().UnixNano()) }
func generateKnowledgeID() string { return fmt.Sprintf("knowledge_%d", time.Now().UnixNano()) }
func generatePatternID() string   { return fmt.Sprintf("pattern_%d", time.Now().UnixNano()) }
func generateDecisionID() string  { return fmt.Sprintf("decision_%d", time.Now().UnixNano()) }