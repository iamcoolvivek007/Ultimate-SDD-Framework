package agents

// ScaleLevel represents the complexity level of a project or task
type ScaleLevel int

const (
	// ScaleLevel0 - Bug fixes, typos, minimal changes (~5 min)
	ScaleLevel0 ScaleLevel = 0
	// ScaleLevel1 - Small features, single file changes (~15 min)
	ScaleLevel1 ScaleLevel = 1
	// ScaleLevel2 - Medium features, multi-file changes (~30 min)
	ScaleLevel2 ScaleLevel = 2
	// ScaleLevel3 - Product/Platform development (~1 hour)
	ScaleLevel3 ScaleLevel = 3
	// ScaleLevel4 - Enterprise systems, compliance-heavy (~2+ hours)
	ScaleLevel4 ScaleLevel = 4
)

// ScaleConfig defines behavior adjustments for each scale level
type ScaleConfig struct {
	Level                ScaleLevel `json:"level"`
	Name                 string     `json:"name"`
	Description          string     `json:"description"`
	PlanningDepth        string     `json:"planning_depth"` // minimal, light, standard, detailed, comprehensive
	RequiredApprovals    int        `json:"required_approvals"`
	DocumentationLevel   string     `json:"documentation_level"` // none, inline, summary, detailed, exhaustive
	TestingRequirement   string     `json:"testing_requirement"` // optional, recommended, required, extensive
	ReviewRequirement    bool       `json:"review_requirement"`
	AgentCount           int        `json:"agent_count"` // How many agents to involve
	EstimatedTimeMinutes int        `json:"estimated_time_minutes"`
}

// DefaultScaleConfigs returns the default configuration for each scale level
func DefaultScaleConfigs() map[ScaleLevel]*ScaleConfig {
	return map[ScaleLevel]*ScaleConfig{
		ScaleLevel0: {
			Level:                ScaleLevel0,
			Name:                 "Quick Fix",
			Description:          "Bug fixes, typos, minimal changes",
			PlanningDepth:        "minimal",
			RequiredApprovals:    0,
			DocumentationLevel:   "none",
			TestingRequirement:   "optional",
			ReviewRequirement:    false,
			AgentCount:           1,
			EstimatedTimeMinutes: 5,
		},
		ScaleLevel1: {
			Level:                ScaleLevel1,
			Name:                 "Small Feature",
			Description:          "Single file changes, small additions",
			PlanningDepth:        "light",
			RequiredApprovals:    0,
			DocumentationLevel:   "inline",
			TestingRequirement:   "recommended",
			ReviewRequirement:    false,
			AgentCount:           2,
			EstimatedTimeMinutes: 15,
		},
		ScaleLevel2: {
			Level:                ScaleLevel2,
			Name:                 "Medium Feature",
			Description:          "Multi-file changes, moderate complexity",
			PlanningDepth:        "standard",
			RequiredApprovals:    1,
			DocumentationLevel:   "summary",
			TestingRequirement:   "required",
			ReviewRequirement:    true,
			AgentCount:           3,
			EstimatedTimeMinutes: 30,
		},
		ScaleLevel3: {
			Level:                ScaleLevel3,
			Name:                 "Product Development",
			Description:          "Full feature implementation, architectural changes",
			PlanningDepth:        "detailed",
			RequiredApprovals:    2,
			DocumentationLevel:   "detailed",
			TestingRequirement:   "extensive",
			ReviewRequirement:    true,
			AgentCount:           5,
			EstimatedTimeMinutes: 60,
		},
		ScaleLevel4: {
			Level:                ScaleLevel4,
			Name:                 "Enterprise",
			Description:          "Large-scale implementation, compliance requirements",
			PlanningDepth:        "comprehensive",
			RequiredApprovals:    3,
			DocumentationLevel:   "exhaustive",
			TestingRequirement:   "extensive",
			ReviewRequirement:    true,
			AgentCount:           8,
			EstimatedTimeMinutes: 120,
		},
	}
}

// DetectScaleLevel analyzes a description and returns the appropriate scale level
func DetectScaleLevel(description string, fileCount int, hasBreakingChanges bool) ScaleLevel {
	// Simple heuristics for scale detection

	// Check for quick fix keywords
	quickFixKeywords := []string{"fix", "typo", "bug", "hotfix", "patch", "minor", "small"}
	for _, kw := range quickFixKeywords {
		if containsIgnoreCase(description, kw) && fileCount <= 1 {
			return ScaleLevel0
		}
	}

	// Check for enterprise keywords
	enterpriseKeywords := []string{"migration", "enterprise", "compliance", "security audit", "gdpr", "hipaa", "soc2"}
	for _, kw := range enterpriseKeywords {
		if containsIgnoreCase(description, kw) {
			return ScaleLevel4
		}
	}

	// Check for product-level keywords
	productKeywords := []string{"architecture", "redesign", "refactor", "new service", "microservice", "api"}
	for _, kw := range productKeywords {
		if containsIgnoreCase(description, kw) && fileCount >= 5 {
			return ScaleLevel3
		}
	}

	// Breaking changes typically need more scrutiny
	if hasBreakingChanges {
		if fileCount > 3 {
			return ScaleLevel3
		}
		return ScaleLevel2
	}

	// File count based detection
	switch {
	case fileCount <= 1:
		return ScaleLevel1
	case fileCount <= 3:
		return ScaleLevel2
	case fileCount <= 10:
		return ScaleLevel3
	default:
		return ScaleLevel4
	}
}

// containsIgnoreCase checks if s contains substr (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(stringContains(toLowerCase(s), toLowerCase(substr)))))
}

func toLowerCase(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// GetRecommendedAgents returns the recommended agents for a given scale level
func GetRecommendedAgents(level ScaleLevel) []string {
	switch level {
	case ScaleLevel0:
		return []string{"developer"}
	case ScaleLevel1:
		return []string{"developer", "qa"}
	case ScaleLevel2:
		return []string{"developer", "qa", "architect"}
	case ScaleLevel3:
		return []string{"pm", "architect", "developer", "qa", "devops"}
	case ScaleLevel4:
		return []string{"pm", "architect", "developer", "qa", "devops", "security", "tech_lead", "data_architect"}
	default:
		return []string{"developer"}
	}
}
