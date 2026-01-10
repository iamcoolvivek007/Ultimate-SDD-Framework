package learning

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

// LearningData represents accumulated learning from development sessions
type LearningData struct {
	UserPreferences   UserPreferences     `json:"user_preferences"`
	CodePatterns      []CodePattern       `json:"code_patterns"`
	SuccessMetrics    []SuccessMetric     `json:"success_metrics"`
	FailurePatterns   []FailurePattern    `json:"failure_patterns"`
	RuleEvolutions    []RuleEvolution     `json:"rule_evolutions"`
	LastUpdated       time.Time           `json:"last_updated"`
}

// UserPreferences captures developer preferences and habits
type UserPreferences struct {
	PreferredLanguages   []string          `json:"preferred_languages"`
	CodingStyle          map[string]string `json:"coding_style"`
	TestingPreferences   []string          `json:"testing_preferences"`
	ErrorHandlingStyle   string            `json:"error_handling_style"`
	NamingConventions    map[string]string `json:"naming_conventions"`
	FavoriteLibraries    []string          `json:"favorite_libraries"`
	AvoidedPatterns      []string          `json:"avoided_patterns"`
}

// CodePattern represents learned code patterns and best practices
type CodePattern struct {
	Pattern     string    `json:"pattern"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	Category    string    `json:"category"` // error_handling, performance, readability, etc.
	Confidence  float64   `json:"confidence"`
	Examples    []string  `json:"examples"`
	LastUsed    time.Time `json:"last_used"`
	SuccessRate float64   `json:"success_rate"`
}

// SuccessMetric tracks successful patterns and approaches
type SuccessMetric struct {
	Action       string    `json:"action"`
	Context      string    `json:"context"`
	Outcome      string    `json:"outcome"`
	Score        float64   `json:"score"`
	Timestamp    time.Time `json:"timestamp"`
	Duration     int       `json:"duration_ms"`
}

// FailurePattern tracks patterns that lead to issues
type FailurePattern struct {
	Pattern      string    `json:"pattern"`
	Description  string    `json:"description"`
	Consequence  string    `json:"consequence"`
	Frequency    int       `json:"frequency"`
	LastOccurred time.Time `json:"last_occurred"`
	Mitigation   string    `json:"mitigation"`
}

// RuleEvolution tracks how project rules have evolved
type RuleEvolution struct {
	OriginalRule    string    `json:"original_rule"`
	EvolvedRule     string    `json:"evolved_rule"`
	Reason          string    `json:"reason"`
	DateEvolved     time.Time `json:"date_evolved"`
	Improvement     string    `json:"improvement"`
	ValidationScore float64   `json:"validation_score"`
}

// AdaptiveLearner manages the learning and adaptation system
type AdaptiveLearner struct {
	projectRoot  string
	learningData LearningData
	dataPath     string
	agentSvc     *agents.AgentService
}

// NewAdaptiveLearner creates a new adaptive learning system
func NewAdaptiveLearner(projectRoot string) (*AdaptiveLearner, error) {
	dataPath := filepath.Join(projectRoot, ".sdd", "learning.json")

	agentSvc := agents.NewAgentService(projectRoot)
	if err := agentSvc.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize agent service: %w", err)
	}

	learner := &AdaptiveLearner{
		projectRoot: projectRoot,
		dataPath:    dataPath,
		agentSvc:    agentSvc,
	}

	// Load existing learning data
	if err := learner.loadLearningData(); err != nil {
		// If file doesn't exist, start with empty data
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load learning data: %w", err)
		}
		learner.learningData = LearningData{
			UserPreferences: UserPreferences{
				CodingStyle:       make(map[string]string),
				NamingConventions: make(map[string]string),
			},
		}
	}

	return learner, nil
}

// LearnFromInteraction records learning from a development interaction
func (al *AdaptiveLearner) LearnFromInteraction(interactionType, context, action, outcome string, success bool, duration int) error {
	timestamp := time.Now()

	// Record success/failure metrics
	metric := SuccessMetric{
		Action:    action,
		Context:   context,
		Outcome:   outcome,
		Timestamp: timestamp,
		Duration:  duration,
	}

	if success {
		metric.Score = 1.0
		al.learningData.SuccessMetrics = append(al.learningData.SuccessMetrics, metric)

		// Learn successful patterns
		al.learnSuccessfulPattern(action, context, outcome)
	} else {
		metric.Score = 0.0
		al.learningData.SuccessMetrics = append(al.learningData.SuccessMetrics, metric)

		// Learn from failures
		al.learnFailurePattern(action, context, outcome)
	}

	// Update user preferences based on patterns
	al.updateUserPreferences(action, context, success)

	// Clean up old data (keep last 1000 interactions)
	if len(al.learningData.SuccessMetrics) > 1000 {
		al.learningData.SuccessMetrics = al.learningData.SuccessMetrics[len(al.learningData.SuccessMetrics)-1000:]
	}

	al.learningData.LastUpdated = timestamp
	return al.saveLearningData()
}

// LearnFromCodeReview incorporates insights from code reviews
func (al *AdaptiveLearner) LearnFromCodeReview(reviewResults map[string]interface{}) error {
	// Extract patterns from review comments and suggestions
	if comments, ok := reviewResults["comments"].([]string); ok {
		for _, comment := range comments {
			al.learnFromReviewComment(comment)
		}
	}

	// Learn from approved/rejected suggestions
	if suggestions, ok := reviewResults["suggestions"].(map[string]bool); ok {
		for suggestion, accepted := range suggestions {
			if accepted {
				al.learnSuccessfulPattern("code_review_suggestion", "review", suggestion)
			} else {
				al.learnFailurePattern("code_review_suggestion", "review", suggestion)
			}
		}
	}

	return al.saveLearningData()
}

// LearnFromPairProgramming records insights from pair programming sessions
func (al *AdaptiveLearner) LearnFromPairProgramming(sessionData map[string]interface{}) error {
	if interactions, ok := sessionData["interactions"].([]map[string]interface{}); ok {
		for _, interaction := range interactions {
			action, _ := interaction["action"].(string)
			context, _ := interaction["context"].(string)
			outcome, _ := interaction["outcome"].(string)
			success, _ := interaction["success"].(bool)
			duration, _ := interaction["duration"].(int)

			al.LearnFromInteraction("pair_programming", context, action, outcome, success, duration)
		}
	}

	// Learn user preferences from session
	if preferences, ok := sessionData["preferences"].(map[string]interface{}); ok {
		al.updatePreferencesFromSession(preferences)
	}

	return al.saveLearningData()
}

// GetPersonalizedSuggestions provides context-aware suggestions based on learning
func (al *AdaptiveLearner) GetPersonalizedSuggestions(context, taskType string) ([]PersonalizedSuggestion, error) {
	suggestions := []PersonalizedSuggestion{}

	// Get patterns relevant to the context
	relevantPatterns := al.getRelevantPatterns(context, taskType)

	for _, pattern := range relevantPatterns {
		if pattern.Confidence > 0.7 && pattern.SuccessRate > 0.8 {
			suggestion := PersonalizedSuggestion{
				Type:        "pattern",
				Title:       pattern.Pattern,
				Description: pattern.Description,
				Confidence:  pattern.Confidence,
				Examples:    pattern.Examples,
				Reason:      fmt.Sprintf("Based on your successful use of this pattern (%.1f%% success rate)", pattern.SuccessRate*100),
			}
			suggestions = append(suggestions, suggestion)
		}
	}

	// Add preference-based suggestions
	preferenceSuggestions := al.getPreferenceBasedSuggestions(context, taskType)
	suggestions = append(suggestions, preferenceSuggestions...)

	// Add avoidance suggestions (patterns to avoid)
	avoidanceSuggestions := al.getAvoidanceSuggestions(context)
	suggestions = append(suggestions, avoidanceSuggestions...)

	// Sort by confidence
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Confidence > suggestions[j].Confidence
	})

	// Limit to top 5 suggestions
	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	return suggestions, nil
}

// EvolveRules analyzes patterns and suggests rule improvements
func (al *AdaptiveLearner) EvolveRules() ([]RuleEvolutionSuggestion, error) {
	suggestions := []RuleEvolutionSuggestion{}

	// Analyze failure patterns for rule evolution opportunities
	for _, failure := range al.learningData.FailurePatterns {
		if failure.Frequency >= 3 { // Pattern occurs frequently
			suggestion := RuleEvolutionSuggestion{
				CurrentRule:  fmt.Sprintf("Avoid: %s", failure.Pattern),
				SuggestedRule: fmt.Sprintf("Proactively prevent: %s", failure.Pattern),
				Reason:       fmt.Sprintf("This pattern has caused issues %d times", failure.Frequency),
				Evidence:     failure.Consequence,
				Mitigation:   failure.Mitigation,
				Confidence:   0.8,
			}
			suggestions = append(suggestions, suggestion)
		}
	}

	// Analyze successful patterns for promotion to rules
	successPatterns := al.getHighSuccessPatterns()
	for _, pattern := range successPatterns {
		if pattern.SuccessRate > 0.9 && len(pattern.Examples) >= 3 {
			suggestion := RuleEvolutionSuggestion{
				CurrentRule:  "No specific rule",
				SuggestedRule: fmt.Sprintf("Best Practice: %s", pattern.Pattern),
				Reason:       fmt.Sprintf("Highly successful pattern with %.1f%% success rate", pattern.SuccessRate*100),
				Evidence:     fmt.Sprintf("Successfully applied %d times", len(pattern.Examples)),
				Confidence:   pattern.Confidence,
			}
			suggestions = append(suggestions, suggestion)
		}
	}

	return suggestions, nil
}

// PersonalizedSuggestion represents a personalized recommendation
type PersonalizedSuggestion struct {
	Type        string   `json:"type"`        // pattern, preference, avoidance
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Confidence  float64  `json:"confidence"`
	Examples    []string `json:"examples"`
	Reason      string   `json:"reason"`
}

// RuleEvolutionSuggestion represents a suggested rule improvement
type RuleEvolutionSuggestion struct {
	CurrentRule  string  `json:"current_rule"`
	SuggestedRule string  `json:"suggested_rule"`
	Reason       string  `json:"reason"`
	Evidence     string  `json:"evidence"`
	Mitigation   string  `json:"mitigation"`
	Confidence   float64 `json:"confidence"`
}

// Private methods

func (al *AdaptiveLearner) learnSuccessfulPattern(action, context, outcome string) {
	// Find existing pattern or create new one
	patternKey := fmt.Sprintf("%s:%s", action, context)
	found := false

	for i, pattern := range al.learningData.CodePatterns {
		if pattern.Pattern == patternKey {
			// Update existing pattern
			al.learningData.CodePatterns[i].Examples = append(al.learningData.CodePatterns[i].Examples, outcome)
			al.learningData.CodePatterns[i].LastUsed = time.Now()
			al.learningData.CodePatterns[i].SuccessRate = al.calculateSuccessRate(patternKey)
			al.learningData.CodePatterns[i].Confidence += 0.1 // Increase confidence
			if al.learningData.CodePatterns[i].Confidence > 1.0 {
				al.learningData.CodePatterns[i].Confidence = 1.0
			}
			found = true
			break
		}
	}

	if !found {
		// Create new pattern
		pattern := CodePattern{
			Pattern:     patternKey,
			Description: fmt.Sprintf("Successfully applied %s in %s context", action, context),
			Language:    al.detectLanguageFromContext(context),
			Category:    al.categorizePattern(action),
			Confidence:  0.6, // Start with moderate confidence
			Examples:    []string{outcome},
			LastUsed:    time.Now(),
			SuccessRate: 1.0, // First success
		}
		al.learningData.CodePatterns = append(al.learningData.CodePatterns, pattern)
	}
}

func (al *AdaptiveLearner) learnFailurePattern(action, context, outcome string) {
	patternKey := fmt.Sprintf("%s:%s", action, context)

	// Update failure patterns
	found := false
	for i, failure := range al.learningData.FailurePatterns {
		if failure.Pattern == patternKey {
			al.learningData.FailurePatterns[i].Frequency++
			al.learningData.FailurePatterns[i].LastOccurred = time.Now()
			found = true
			break
		}
	}

	if !found {
		failure := FailurePattern{
			Pattern:      patternKey,
			Description:  fmt.Sprintf("Failed attempt: %s", outcome),
			Consequence:  outcome,
			Frequency:    1,
			LastOccurred: time.Now(),
			Mitigation:   "Needs analysis", // Will be filled by AI analysis
		}
		al.learningData.FailurePatterns = append(al.learningData.FailurePatterns, failure)

		// Request AI analysis for mitigation
		al.analyzeFailureForMitigation(&failure)
	}
}

func (al *AdaptiveLearner) updateUserPreferences(action string, context string, success bool) {
	// Analyze successful actions to learn preferences
	if success {
		// Learn language preferences
		if strings.Contains(action, "javascript") || strings.Contains(action, "typescript") {
			al.addPreferredLanguage("JavaScript/TypeScript")
		} else if strings.Contains(action, "go") {
			al.addPreferredLanguage("Go")
		}

		// Learn testing preferences
		if strings.Contains(action, "test") {
			al.learningData.UserPreferences.TestingPreferences = append(
				al.learningData.UserPreferences.TestingPreferences, action)
		}
	}
}

func (al *AdaptiveLearner) learnFromReviewComment(comment string) {
	// Extract patterns from review comments
	comment = strings.ToLower(comment)

	if strings.Contains(comment, "naming") || strings.Contains(comment, "name") {
		al.learningData.UserPreferences.NamingConventions["review_feedback"] = comment
	}

	if strings.Contains(comment, "error") || strings.Contains(comment, "handle") {
		if al.learningData.UserPreferences.ErrorHandlingStyle == "" {
			al.learningData.UserPreferences.ErrorHandlingStyle = "review_suggested"
		}
	}
}

func (al *AdaptiveLearner) updatePreferencesFromSession(preferences map[string]interface{}) {
	// Update preferences based on pair programming session data
	for key, value := range preferences {
		switch key {
		case "preferred_language":
			if lang, ok := value.(string); ok {
				al.addPreferredLanguage(lang)
			}
		case "coding_style":
			if style, ok := value.(string); ok {
				al.learningData.UserPreferences.CodingStyle["session_learned"] = style
			}
		}
	}
}

func (al *AdaptiveLearner) getRelevantPatterns(context, taskType string) []CodePattern {
	relevant := []CodePattern{}

	for _, pattern := range al.learningData.CodePatterns {
		// Match context and task type
		if strings.Contains(strings.ToLower(context), strings.ToLower(pattern.Language)) ||
		   strings.Contains(strings.ToLower(taskType), strings.ToLower(pattern.Category)) {
			relevant = append(relevant, pattern)
		}
	}

	return relevant
}

func (al *AdaptiveLearner) getPreferenceBasedSuggestions(context, taskType string) []PersonalizedSuggestion {
	suggestions := []PersonalizedSuggestion{}

	// Language preferences
	if len(al.learningData.UserPreferences.PreferredLanguages) > 0 {
		suggestion := PersonalizedSuggestion{
			Type:        "preference",
			Title:       "Language Preference",
			Description: fmt.Sprintf("You prefer working with %s", strings.Join(al.learningData.UserPreferences.PreferredLanguages, ", ")),
			Confidence:  0.9,
			Reason:      "Based on your successful past projects",
		}
		suggestions = append(suggestions, suggestion)
	}

	// Testing preferences
	if len(al.learningData.UserPreferences.TestingPreferences) > 0 && strings.Contains(taskType, "test") {
		suggestion := PersonalizedSuggestion{
			Type:        "preference",
			Title:       "Testing Approach",
			Description: "Based on your testing preferences",
			Confidence:  0.8,
			Reason:      "Matches your preferred testing patterns",
		}
		suggestions = append(suggestions, suggestion)
	}

	return suggestions
}

func (al *AdaptiveLearner) getAvoidanceSuggestions(context string) []PersonalizedSuggestion {
	suggestions := []PersonalizedSuggestion{}

	for _, failure := range al.learningData.FailurePatterns {
		if failure.Frequency >= 2 {
			suggestion := PersonalizedSuggestion{
				Type:        "avoidance",
				Title:       "Avoid This Pattern",
				Description: fmt.Sprintf("Avoid: %s", failure.Pattern),
				Confidence:  0.85,
				Reason:      fmt.Sprintf("This pattern has caused issues %d times", failure.Frequency),
			}
			suggestions = append(suggestions, suggestion)
		}
	}

	return suggestions
}

func (al *AdaptiveLearner) getHighSuccessPatterns() []CodePattern {
	highSuccess := []CodePattern{}

	for _, pattern := range al.learningData.CodePatterns {
		if pattern.SuccessRate > 0.9 && len(pattern.Examples) >= 3 {
			highSuccess = append(highSuccess, pattern)
		}
	}

	return highSuccess
}

// Helper methods

func (al *AdaptiveLearner) addPreferredLanguage(language string) {
	found := false
	for _, lang := range al.learningData.UserPreferences.PreferredLanguages {
		if lang == language {
			found = true
			break
		}
	}
	if !found {
		al.learningData.UserPreferences.PreferredLanguages = append(
			al.learningData.UserPreferences.PreferredLanguages, language)
	}
}

func (al *AdaptiveLearner) calculateSuccessRate(patternKey string) float64 {
	successCount := 0
	totalCount := 0

	for _, metric := range al.learningData.SuccessMetrics {
		if strings.Contains(metric.Action+":"+metric.Context, patternKey) {
			totalCount++
			if metric.Score > 0.8 {
				successCount++
			}
		}
	}

	if totalCount == 0 {
		return 0.0
	}
	return float64(successCount) / float64(totalCount)
}

func (al *AdaptiveLearner) detectLanguageFromContext(context string) string {
	context = strings.ToLower(context)

	if strings.Contains(context, "go") || strings.Contains(context, ".go") {
		return "go"
	}
	if strings.Contains(context, "javascript") || strings.Contains(context, "typescript") ||
	   strings.Contains(context, ".js") || strings.Contains(context, ".ts") {
		return "javascript"
	}
	if strings.Contains(context, "python") || strings.Contains(context, ".py") {
		return "python"
	}

	return "unknown"
}

func (al *AdaptiveLearner) categorizePattern(action string) string {
	action = strings.ToLower(action)

	if strings.Contains(action, "error") || strings.Contains(action, "handle") {
		return "error_handling"
	}
	if strings.Contains(action, "test") {
		return "testing"
	}
	if strings.Contains(action, "performance") || strings.Contains(action, "optimize") {
		return "performance"
	}

	return "general"
}

func (al *AdaptiveLearner) analyzeFailureForMitigation(failure *FailurePattern) {
	// Use AI to analyze failure and suggest mitigation
	prompt := fmt.Sprintf("Analyze this development failure and suggest mitigation:\nPattern: %s\nConsequence: %s\n\nProvide a specific mitigation strategy.",
		failure.Pattern, failure.Consequence)

	if response, err := al.agentSvc.GetAgentResponse("system", "analyze", prompt, "", ""); err == nil {
		failure.Mitigation = strings.TrimSpace(response)
	} else {
		failure.Mitigation = "AI analysis failed - manual review needed"
	}
}

func (al *AdaptiveLearner) loadLearningData() error {
	data, err := os.ReadFile(al.dataPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &al.learningData)
}

func (al *AdaptiveLearner) saveLearningData() error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(al.dataPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(al.learningData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(al.dataPath, data, 0644)
}

// GetLearningSummary provides a summary of learned patterns and preferences
func (al *AdaptiveLearner) GetLearningSummary() string {
	var summary strings.Builder

	summary.WriteString("# 🧠 Learning & Adaptation Summary\n\n")

	// User preferences
	if len(al.learningData.UserPreferences.PreferredLanguages) > 0 {
		summary.WriteString("## 🌍 Preferred Languages\n")
		for _, lang := range al.learningData.UserPreferences.PreferredLanguages {
			summary.WriteString(fmt.Sprintf("- %s\n", lang))
		}
		summary.WriteString("\n")
	}

	// Successful patterns
	if len(al.learningData.CodePatterns) > 0 {
		summary.WriteString("## ✅ Successful Patterns\n")
		for _, pattern := range al.learningData.CodePatterns {
			if pattern.SuccessRate > 0.8 {
				summary.WriteString(fmt.Sprintf("- **%s**: %.1f%% success rate (%d examples)\n",
					pattern.Pattern, pattern.SuccessRate*100, len(pattern.Examples)))
			}
		}
		summary.WriteString("\n")
	}

	// Failure patterns
	if len(al.learningData.FailurePatterns) > 0 {
		summary.WriteString("## ❌ Patterns to Avoid\n")
		for _, failure := range al.learningData.FailurePatterns {
			if failure.Frequency >= 2 {
				summary.WriteString(fmt.Sprintf("- **%s**: Occurred %d times\n  *%s*\n",
					failure.Pattern, failure.Frequency, failure.Consequence))
			}
		}
		summary.WriteString("\n")
	}

	// Success metrics
	totalInteractions := len(al.learningData.SuccessMetrics)
	if totalInteractions > 0 {
		successfulInteractions := 0
		for _, metric := range al.learningData.SuccessMetrics {
			if metric.Score > 0.8 {
				successfulInteractions++
			}
		}

		successRate := float64(successfulInteractions) / float64(totalInteractions) * 100
		summary.WriteString("## 📊 Success Metrics\n")
		summary.WriteString(fmt.Sprintf("- **Total Interactions**: %d\n", totalInteractions))
		summary.WriteString(fmt.Sprintf("- **Success Rate**: %.1f%%\n", successRate))
		summary.WriteString(fmt.Sprintf("- **Learning Patterns**: %d\n", len(al.learningData.CodePatterns)))
		summary.WriteString("\n")
	}

	if al.learningData.LastUpdated.IsZero() {
		summary.WriteString("*No learning data available yet. Start using the framework to build your personalized profile.*\n")
	} else {
		summary.WriteString(fmt.Sprintf("*Last updated: %s*\n", al.learningData.LastUpdated.Format("2006-01-02 15:04:05")))
	}

	return summary.String()
}