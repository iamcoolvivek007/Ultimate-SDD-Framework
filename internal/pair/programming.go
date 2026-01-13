package pair

import (
	"fmt"
	"strings"
	"time"

	"ultimate-sdd-framework/internal/agents"
	"ultimate-sdd-framework/internal/lsp"
)

// PairSession represents an active pair programming session
type PairSession struct {
	ID         string               `json:"id"`
	StartTime  time.Time            `json:"start_time"`
	ActiveFile string               `json:"active_file"`
	Context    *lsp.CodebaseContext `json:"context"`
	Agent      *agents.Agent        `json:"agent"`
	SessionLog []SessionEntry       `json:"session_log"`
	Stats      PairingStats         `json:"stats"`
	IsActive   bool                 `json:"is_active"`
}

// SessionEntry represents a single interaction in the session
type SessionEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	Type       string    `json:"type"` // suggestion, question, code_change, explanation
	Content    string    `json:"content"`
	File       string    `json:"file"`
	Line       int       `json:"line"`
	UserAction string    `json:"user_action"` // accepted, rejected, modified, ignored
	Duration   int       `json:"duration_ms"` // milliseconds spent on this interaction
}

// PairingStats tracks session statistics
type PairingStats struct {
	TotalInteractions     int           `json:"total_interactions"`
	AcceptedSuggestions   int           `json:"accepted_suggestions"`
	RejectedSuggestions   int           `json:"rejected_suggestions"`
	TimeSpent             time.Duration `json:"time_spent"`
	FilesTouched          []string      `json:"files_touched"`
	ProductivityScore     float64       `json:"productivity_score"`
	LearningOpportunities int           `json:"learning_opportunities"`
}

// PairProgrammer manages pair programming sessions
type PairProgrammer struct {
	projectRoot    string
	agentSvc       *agents.AgentService
	activeSession  *PairSession
	sessionHistory []PairSession
}

// NewPairProgrammer creates a new pair programming manager
func NewPairProgrammer(projectRoot string) (*PairProgrammer, error) {
	agentSvc := agents.NewAgentService(projectRoot)
	if err := agentSvc.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize agent service: %w", err)
	}

	return &PairProgrammer{
		projectRoot: projectRoot,
		agentSvc:    agentSvc,
	}, nil
}

// StartSession begins a new pair programming session
func (pp *PairProgrammer) StartSession(agentRole string, focusArea string) (*PairSession, error) {
	// Map role to phase for agent selection
	phase := "execute" // default
	switch agentRole {
	case "architect":
		phase = "plan"
	case "qa":
		phase = "review"
	}

	agent, err := pp.agentSvc.GetAgentForPhase(phase)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent for role '%s': %w", agentRole, err)
	}

	session := &PairSession{
		ID:         generateSessionID(),
		StartTime:  time.Now(),
		Context:    lsp.NewCodebaseContext(pp.projectRoot),
		Agent:      agent,
		SessionLog: []SessionEntry{},
		IsActive:   true,
	}

	// Initialize context
	if err := session.Context.AnalyzeProject(); err != nil {
		return nil, fmt.Errorf("failed to analyze codebase: %w", err)
	}

	pp.activeSession = session

	// Add initial session entry
	pp.logSessionEntry("session_start", fmt.Sprintf("Started pair programming session with %s agent focusing on %s", agentRole, focusArea), "", 0, "")

	return session, nil
}

// EndSession concludes the current pair programming session
func (pp *PairProgrammer) EndSession() (*PairSession, error) {
	if pp.activeSession == nil {
		return nil, fmt.Errorf("no active session to end")
	}

	session := pp.activeSession
	session.IsActive = false
	session.Stats.TimeSpent = time.Since(session.StartTime)

	// Calculate final statistics
	session.Stats = pp.calculateSessionStats(session)

	// Add final log entry
	pp.logSessionEntry("session_end", fmt.Sprintf("Ended session after %v", session.Stats.TimeSpent.Round(time.Second)), "", 0, "")

	// Move to history
	pp.sessionHistory = append(pp.sessionHistory, *session)
	pp.activeSession = nil

	return session, nil
}

// GetSuggestion requests AI assistance for current context
func (pp *PairProgrammer) GetSuggestion(filePath string, cursorLine int, context string, requestType string) (*PairSuggestion, error) {
	if pp.activeSession == nil {
		return nil, fmt.Errorf("no active pair programming session")
	}

	startTime := time.Now()

	// Update active file
	pp.activeSession.ActiveFile = filePath

	// Get context-aware suggestion
	suggestion, err := pp.generateSuggestion(filePath, cursorLine, context, requestType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate suggestion: %w", err)
	}

	duration := int(time.Since(startTime).Milliseconds())
	_ = duration // Reserve for future use in logging

	// Log the interaction
	pp.logSessionEntry("suggestion", suggestion.Content, filePath, cursorLine, "")

	return suggestion, nil
}

// RecordUserAction records how the user responded to a suggestion
func (pp *PairProgrammer) RecordUserAction(action string, suggestionID string) error {
	if pp.activeSession == nil {
		return fmt.Errorf("no active session")
	}

	// Find the last suggestion entry and update it
	for i := len(pp.activeSession.SessionLog) - 1; i >= 0; i-- {
		entry := &pp.activeSession.SessionLog[i]
		if entry.Type == "suggestion" {
			entry.UserAction = action

			// Update statistics
			switch action {
			case "accepted":
				pp.activeSession.Stats.AcceptedSuggestions++
			case "rejected":
				pp.activeSession.Stats.RejectedSuggestions++
			}
			break
		}
	}

	return nil
}

// PairSuggestion represents an AI-generated suggestion
type PairSuggestion struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"` // completion, refactor, explanation, test
	Content      string   `json:"content"`
	Explanation  string   `json:"explanation"`
	Confidence   float64  `json:"confidence"`
	Alternatives []string `json:"alternatives"`
	File         string   `json:"file"`
	Line         int      `json:"line"`
}

// generateSuggestion creates context-aware suggestions
func (pp *PairProgrammer) generateSuggestion(filePath string, cursorLine int, context string, requestType string) (*PairSuggestion, error) {
	suggestion := &PairSuggestion{
		ID:         generateSuggestionID(),
		File:       filePath,
		Line:       cursorLine,
		Type:       requestType,
		Confidence: 0.8,
	}

	// Build context-aware prompt
	prompt := pp.buildSuggestionPrompt(filePath, cursorLine, context, requestType)

	// Get AI response
	response, err := pp.agentSvc.GetAgentResponse(pp.activeSession.Agent.Role, "execute", prompt, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}

	// Parse and format the response
	suggestion.Content = pp.parseSuggestionContent(response, requestType)
	suggestion.Explanation = pp.extractExplanation(response)

	// Generate alternatives if applicable
	if requestType == "completion" {
		suggestion.Alternatives = pp.generateAlternatives(response)
	}

	return suggestion, nil
}

// buildSuggestionPrompt creates a detailed prompt for AI suggestions
func (pp *PairProgrammer) buildSuggestionPrompt(filePath string, cursorLine int, context string, requestType string) string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("You are pair programming with a developer. Current context:\n"))
	prompt.WriteString(fmt.Sprintf("- File: %s (line %d)\n", filePath, cursorLine))
	prompt.WriteString(fmt.Sprintf("- Request type: %s\n", requestType))
	prompt.WriteString(fmt.Sprintf("- Code context:\n```\n%s\n```\n\n", context))

	// Add project context
	if pp.activeSession.Context != nil {
		prompt.WriteString("Project context:\n")
		prompt.WriteString(fmt.Sprintf("- Language: %s\n", pp.activeSession.Context.Structure.MainLanguage))
		prompt.WriteString(fmt.Sprintf("- Framework: %s\n", pp.activeSession.Context.Structure.Framework))

		if pp.activeSession.Context.Structure.HasTests {
			prompt.WriteString("- Has test suite\n")
		}

		prompt.WriteString("\n")
	}

	// Add request-specific instructions
	switch requestType {
	case "completion":
		prompt.WriteString("Provide code completion suggestions. Focus on:\n")
		prompt.WriteString("- Following existing code patterns\n")
		prompt.WriteString("- Type safety and best practices\n")
		prompt.WriteString("- Performance considerations\n")
		prompt.WriteString("- Error handling\n\n")
		prompt.WriteString("Format: Show the completed code with clear indication of what was added.\n")

	case "refactor":
		prompt.WriteString("Suggest refactoring improvements. Consider:\n")
		prompt.WriteString("- Code readability and maintainability\n")
		prompt.WriteString("- Performance optimizations\n")
		prompt.WriteString("- Reducing complexity\n")
		prompt.WriteString("- Following SOLID principles\n\n")
		prompt.WriteString("Format: Show before/after code with explanations.\n")

	case "test":
		prompt.WriteString("Suggest test cases and testing approaches. Focus on:\n")
		prompt.WriteString("- Unit test coverage\n")
		prompt.WriteString("- Edge cases and error conditions\n")
		prompt.WriteString("- Integration testing\n")
		prompt.WriteString("- Test-driven development principles\n\n")
		prompt.WriteString("Format: Provide test code examples with explanations.\n")

	case "explanation":
		prompt.WriteString("Explain the code and suggest improvements. Cover:\n")
		prompt.WriteString("- What the code does\n")
		prompt.WriteString("- Potential issues or improvements\n")
		prompt.WriteString("- Best practices and alternatives\n")
		prompt.WriteString("- Performance and security considerations\n\n")

	default:
		prompt.WriteString("Provide helpful assistance for the developer's current task.\n")
	}

	prompt.WriteString("\nBe concise but thorough. Focus on practical, actionable suggestions.")

	return prompt.String()
}

// parseSuggestionContent extracts the main suggestion from AI response
func (pp *PairProgrammer) parseSuggestionContent(response, requestType string) string {
	// Simple parsing - in a real implementation, this would be more sophisticated
	lines := strings.Split(response, "\n")

	switch requestType {
	case "completion":
		// Look for code blocks
		for _, line := range lines {
			if strings.Contains(line, "```") {
				return strings.TrimSpace(response)
			}
		}
		return strings.TrimSpace(response)

	case "refactor":
		return strings.TrimSpace(response)

	default:
		return strings.TrimSpace(response)
	}
}

// extractExplanation pulls explanation text from response
func (pp *PairProgrammer) extractExplanation(response string) string {
	// Look for explanation sections
	if idx := strings.Index(strings.ToLower(response), "explanation"); idx != -1 {
		return strings.TrimSpace(response[idx:])
	}

	// Look for reasoning sections
	if idx := strings.Index(strings.ToLower(response), "reasoning"); idx != -1 {
		return strings.TrimSpace(response[idx:])
	}

	// Default: return first paragraph
	lines := strings.Split(response, "\n\n")
	if len(lines) > 1 {
		return strings.TrimSpace(lines[0])
	}

	return "AI-generated suggestion for your code."
}

// generateAlternatives creates alternative suggestions
func (pp *PairProgrammer) generateAlternatives(response string) []string {
	alternatives := []string{}

	// Simple approach: look for "alternative" or "option" mentions
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		line = strings.ToLower(line)
		if strings.Contains(line, "alternative") || strings.Contains(line, "option") ||
			strings.Contains(line, "or you could") {
			alternatives = append(alternatives, strings.TrimSpace(line))
		}
	}

	// If no explicit alternatives, generate some basic variations
	if len(alternatives) == 0 {
		alternatives = []string{
			"Consider using a more functional approach",
			"You could extract this into a separate function",
			"Think about using early returns for clarity",
		}
	}

	return alternatives
}

// calculateSessionStats computes session statistics
func (pp *PairProgrammer) calculateSessionStats(session *PairSession) PairingStats {
	stats := PairingStats{
		TotalInteractions: len(session.SessionLog),
		TimeSpent:         time.Since(session.StartTime),
		FilesTouched:      []string{},
	}

	// Count interactions and collect files
	fileSet := make(map[string]bool)

	for _, entry := range session.SessionLog {
		if entry.File != "" {
			fileSet[entry.File] = true
		}
	}

	// Convert set to slice
	for file := range fileSet {
		stats.FilesTouched = append(stats.FilesTouched, file)
	}

	// Calculate productivity score (simplified)
	totalSuggestions := stats.AcceptedSuggestions + stats.RejectedSuggestions
	if totalSuggestions > 0 {
		acceptanceRate := float64(stats.AcceptedSuggestions) / float64(totalSuggestions)
		stats.ProductivityScore = acceptanceRate * 10 // Scale to 0-10
	} else {
		stats.ProductivityScore = 7.5 // Default neutral score
	}

	// Count learning opportunities (questions, explanations)
	for _, entry := range session.SessionLog {
		if entry.Type == "question" || entry.Type == "explanation" {
			stats.LearningOpportunities++
		}
	}

	return stats
}

// logSessionEntry adds an entry to the session log
func (pp *PairProgrammer) logSessionEntry(entryType, content, file string, line int, userAction string) {
	if pp.activeSession == nil {
		return
	}

	entry := SessionEntry{
		Timestamp:  time.Now(),
		Type:       entryType,
		Content:    content,
		File:       file,
		Line:       line,
		UserAction: userAction,
	}

	pp.activeSession.SessionLog = append(pp.activeSession.SessionLog, entry)
}

// GetSessionReport generates a detailed session report
func (pp *PairProgrammer) GetSessionReport(session *PairSession) string {
	var report strings.Builder

	report.WriteString(fmt.Sprintf("# 👥 Pair Programming Session Report\n\n"))
	report.WriteString(fmt.Sprintf("**Session ID:** %s\n", session.ID))
	report.WriteString(fmt.Sprintf("**Agent:** %s\n", session.Agent.Role))
	report.WriteString(fmt.Sprintf("**Duration:** %v\n", session.Stats.TimeSpent.Round(time.Second)))
	report.WriteString(fmt.Sprintf("**Files Touched:** %d\n\n", len(session.Stats.FilesTouched)))

	// Statistics
	report.WriteString("## 📊 Session Statistics\n\n")
	report.WriteString(fmt.Sprintf("- **Total Interactions:** %d\n", session.Stats.TotalInteractions))
	report.WriteString(fmt.Sprintf("- **Suggestions Accepted:** %d\n", session.Stats.AcceptedSuggestions))
	report.WriteString(fmt.Sprintf("- **Suggestions Rejected:** %d\n", session.Stats.RejectedSuggestions))
	report.WriteString(fmt.Sprintf("- **Productivity Score:** %.1f/10\n", session.Stats.ProductivityScore))
	report.WriteString(fmt.Sprintf("- **Learning Opportunities:** %d\n\n", session.Stats.LearningOpportunities))

	// Acceptance rate
	totalSuggestions := session.Stats.AcceptedSuggestions + session.Stats.RejectedSuggestions
	if totalSuggestions > 0 {
		acceptanceRate := float64(session.Stats.AcceptedSuggestions) / float64(totalSuggestions) * 100
		report.WriteString(fmt.Sprintf("- **Suggestion Acceptance Rate:** %.1f%%\n\n", acceptanceRate))
	}

	// Files worked on
	if len(session.Stats.FilesTouched) > 0 {
		report.WriteString("## 📁 Files Modified\n\n")
		for _, file := range session.Stats.FilesTouched {
			report.WriteString(fmt.Sprintf("- %s\n", file))
		}
		report.WriteString("\n")
	}

	// Session timeline
	report.WriteString("## ⏰ Session Timeline\n\n")

	for i, entry := range session.SessionLog {
		report.WriteString(fmt.Sprintf("**%d.** %s - %s\n",
			i+1,
			entry.Timestamp.Format("15:04:05"),
			strings.Title(entry.Type)))

		if entry.File != "" {
			report.WriteString(fmt.Sprintf("   File: %s", entry.File))
			if entry.Line > 0 {
				report.WriteString(fmt.Sprintf(":%d", entry.Line))
			}
			report.WriteString("\n")
		}

		// Truncate content for readability
		content := entry.Content
		if len(content) > 100 {
			content = content[:97] + "..."
		}
		report.WriteString(fmt.Sprintf("   %s\n", content))

		if entry.UserAction != "" {
			report.WriteString(fmt.Sprintf("   *User action: %s*\n", entry.UserAction))
		}

		report.WriteString("\n")
	}

	report.WriteString("---\n")
	report.WriteString("*Generated by Ultimate SDD Framework - Pair Programming Mode*\n")

	return report.String()
}

// GetActiveSession returns the current active session
func (pp *PairProgrammer) GetActiveSession() *PairSession {
	return pp.activeSession
}

// GetSessionHistory returns all completed sessions
func (pp *PairProgrammer) GetSessionHistory() []PairSession {
	return pp.sessionHistory
}

// Utility functions
func generateSessionID() string {
	return fmt.Sprintf("pair_%d", time.Now().Unix())
}

func generateSuggestionID() string {
	return fmt.Sprintf("sugg_%d", time.Now().UnixNano())
}
