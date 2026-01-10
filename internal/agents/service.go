package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ultimate-sdd-framework/internal/lsp"
	"ultimate-sdd-framework/internal/mcp"
)

// AgentService provides high-level agent operations with context awareness
type AgentService struct {
	agentMgr        *AgentManager
	mcpMgr          *mcp.MCPManager
	lspContext      *lsp.CodebaseContext
	brownfieldCtx   *lsp.BrownfieldContext
	projectRoot     string
	hasBrownfieldContext bool
}

// NewAgentService creates a new agent service
func NewAgentService(projectRoot string) *AgentService {
	return &AgentService{
		agentMgr:    NewAgentManager(projectRoot),
		mcpMgr:      mcp.NewMCPManager(projectRoot),
		projectRoot: projectRoot,
	}
}

// Initialize loads all components
func (as *AgentService) Initialize() error {
	// Load agents
	if err := as.agentMgr.LoadAgents(); err != nil {
		return fmt.Errorf("failed to load agents: %w", err)
	}

	// Load MCP configuration
	if err := as.mcpMgr.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load MCP config: %w", err)
	}

	// Check for brownfield context (CONTEXT.md)
	contextPath := filepath.Join(as.projectRoot, ".sdd", "CONTEXT.md")
	if _, err := os.Stat(contextPath); err == nil {
		// Brownfield context exists - use it
		as.brownfieldCtx = lsp.NewBrownfieldContext(as.projectRoot)
		if err := as.brownfieldCtx.AnalyzeBrownfield(); err != nil {
			return fmt.Errorf("failed to analyze brownfield context: %w", err)
		}
		as.hasBrownfieldContext = true

		// Still initialize regular LSP context for compatibility
		as.lspContext = &as.brownfieldCtx.CodebaseContext
	} else {
		// No brownfield context - use regular LSP analysis
		as.lspContext = lsp.NewCodebaseContext(as.projectRoot)
		if err := as.lspContext.AnalyzeProject(); err != nil {
			return fmt.Errorf("failed to analyze codebase: %w", err)
		}
		as.hasBrownfieldContext = false
	}

	return nil
}

// GetAgentResponse gets a response from an agent with full context
func (as *AgentService) GetAgentResponse(agentName, phase, userInput string) (string, error) {
	// Get the agent
	agent, err := as.agentMgr.GetAgent(agentName)
	if err != nil {
		return "", fmt.Errorf("agent not found: %w", err)
	}

	// Get context for this phase (brownfield-aware)
	contextInfo := ""
	if as.lspContext != nil {
		contextInfo = as.lspContext.GetContextForPhase(phase)

		// Add brownfield constraints if available
		if as.hasBrownfieldContext && as.brownfieldCtx != nil {
			contextInfo += as.getBrownfieldConstraintsForPhase(phase)
		}
	}

	// Add Conductor Context (persistent brain)
	contextInfo += as.getConductorContext()

	// Build the full prompt
	systemPrompt := agent.GetSystemPrompt()
	phasePrompt := agent.GetPhasePrompt(phase, contextInfo)

	// Combine with user input
	prompt := fmt.Sprintf("%s\n\n%s\n\nUser Request: %s", systemPrompt, phasePrompt, userInput)

	// Get MCP client
	client, err := as.mcpMgr.GetClient("")
	if err != nil {
		return "", fmt.Errorf("no MCP client available: %w", err)
	}

	// Send to AI model
	messages := []mcp.Message{
		{Role: "user", Content: prompt},
	}

	options := map[string]interface{}{
		"temperature": 0.7,
		"max_tokens":  4000,
	}

	response, err := client.Chat(messages, options)
	if err != nil {
		return "", fmt.Errorf("AI request failed: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from AI model")
	}

	return response.Choices[0].Message.Content, nil
}

// getConductorContext reads files from .sdd/context/ to inject persistent context
func (as *AgentService) getConductorContext() string {
	contextDir := filepath.Join(as.projectRoot, ".sdd", "context")
	if _, err := os.Stat(contextDir); os.IsNotExist(err) {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("\n\n## 🧠 PERSISTENT PROJECT CONTEXT (CONDUCTOR)\n")

	files, err := os.ReadDir(contextDir)
	if err != nil {
		return ""
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		content, err := os.ReadFile(filepath.Join(contextDir, file.Name()))
		if err != nil {
			continue
		}

		builder.WriteString(fmt.Sprintf("\n### %s\n%s\n", strings.ToUpper(strings.TrimSuffix(file.Name(), ".md")), string(content)))
	}

	return builder.String()
}

// getBrownfieldConstraintsForPhase provides brownfield-specific constraints for each phase
func (as *AgentService) getBrownfieldConstraintsForPhase(phase string) string {
	if !as.hasBrownfieldContext || as.brownfieldCtx == nil {
		return ""
	}

	var constraints strings.Builder
	constraints.WriteString("\n\n## 🔧 Brownfield Constraints\n\n")

	switch phase {
	case "specify":
		constraints.WriteString("### Legacy System Awareness\n")
		constraints.WriteString("- **Existing Architecture**: Must integrate with current system design\n")
		constraints.WriteString("- **Forbidden Patterns**: Avoid anti-patterns identified in CONTEXT.md\n")
		constraints.WriteString("- **Integration Points**: Consider existing API and database touchpoints\n")
		constraints.WriteString("- **Technical Debt**: Acknowledge known issues and limitations\n")

		// Add specific forbidden patterns
		if len(as.brownfieldCtx.ForbiddenPatterns) > 0 {
			constraints.WriteString("\n### Prohibited Approaches\n")
			for _, pattern := range as.brownfieldCtx.ForbiddenPatterns {
				constraints.WriteString(fmt.Sprintf("- **%s**: %s\n", pattern.Pattern, pattern.Recommended))
			}
		}

	case "plan":
		constraints.WriteString("### Legacy Integration Requirements\n")
		constraints.WriteString("- **Migration Strategy**: Plan how new features integrate with existing code\n")
		constraints.WriteString("- **Refactoring Steps**: Include necessary code restructuring\n")
		constraints.WriteString("- **Backwards Compatibility**: Ensure existing functionality remains intact\n")
		constraints.WriteString("- **Testing Strategy**: Include regression testing for legacy components\n")

		// Add integration points
		if len(as.brownfieldCtx.IntegrationPoints) > 0 {
			constraints.WriteString("\n### Key Integration Points\n")
			for _, point := range as.brownfieldCtx.IntegrationPoints {
				constraints.WriteString(fmt.Sprintf("- **%s**: %s\n", point.Name, point.Description))
			}
		}

	case "task":
		constraints.WriteString("### Legacy-Aware Task Creation\n")
		constraints.WriteString("- **File Path Requirements**: Every task must specify existing files to modify\n")
		constraints.WriteString("- **Regression Risk Assessment**: Identify potential side effects\n")
		constraints.WriteString("- **Legacy Pattern Compliance**: Follow established coding patterns\n")
		constraints.WriteString("- **Testing Obligations**: Include tests for modified legacy components\n")

	case "execute":
		constraints.WriteString("### Safe Modification Practices\n")
		constraints.WriteString("- **Pattern Compliance**: Use only established legacy patterns\n")
		constraints.WriteString("- **Integration Testing**: Verify changes don't break existing functionality\n")
		constraints.WriteString("- **Documentation Updates**: Update CONTEXT.md if patterns change\n")
		constraints.WriteString("- **Code Review Focus**: Pay special attention to legacy integration points\n")

		// Add relevant legacy patterns
		if len(as.brownfieldCtx.LegacyPatterns) > 0 {
			constraints.WriteString("\n### Required Pattern Usage\n")
			for _, pattern := range as.brownfieldCtx.LegacyPatterns {
				constraints.WriteString(fmt.Sprintf("- **%s**: %s\n", pattern.Pattern, pattern.Description))
			}
		}

	case "review":
		constraints.WriteString("### Legacy System Validation\n")
		constraints.WriteString("- **Regression Testing**: Run full existing test suite\n")
		constraints.WriteString("- **Integration Verification**: Test all identified integration points\n")
		constraints.WriteString("- **Pattern Compliance**: Ensure no forbidden patterns were introduced\n")
		constraints.WriteString("- **Technical Debt Assessment**: Evaluate impact on existing debt\n")

		// Add technical debt considerations
		if len(as.brownfieldCtx.TechnicalDebt) > 0 {
			constraints.WriteString("\n### Technical Debt Considerations\n")
			for _, debt := range as.brownfieldCtx.TechnicalDebt {
				if debt.Severity == "High" {
					constraints.WriteString(fmt.Sprintf("- **%s**: %s\n", debt.Issue, debt.Recommendation))
				}
			}
		}
	}

	// Add constitution compliance
	if len(as.brownfieldCtx.Constitution.ArchitecturalRules) > 0 {
		constraints.WriteString("\n### Constitution Compliance\n")
		for _, rule := range as.brownfieldCtx.Constitution.ArchitecturalRules {
			constraints.WriteString(fmt.Sprintf("- **Architectural Rule**: %s\n", rule))
		}
	}

	return constraints.String()
}

// HasBrownfieldContext returns whether brownfield context is available
func (as *AgentService) HasBrownfieldContext() bool {
	return as.hasBrownfieldContext
}

// GetBrownfieldSummary returns a summary of brownfield analysis
func (as *AgentService) GetBrownfieldSummary() string {
	if !as.hasBrownfieldContext || as.brownfieldCtx == nil {
		return "No brownfield context available. Run 'nexus discovery' to analyze the codebase."
	}

	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("## Brownfield Analysis Summary\n\n"))
	summary.WriteString(fmt.Sprintf("**System**: %s with %s\n", as.brownfieldCtx.Structure.MainLanguage, as.brownfieldCtx.Structure.Framework))
	summary.WriteString(fmt.Sprintf("**Files Analyzed**: %d\n", len(as.brownfieldCtx.Files)))
	summary.WriteString(fmt.Sprintf("**Legacy Patterns**: %d identified\n", len(as.brownfieldCtx.LegacyPatterns)))
	summary.WriteString(fmt.Sprintf("**Forbidden Patterns**: %d flagged\n", len(as.brownfieldCtx.ForbiddenPatterns)))
	summary.WriteString(fmt.Sprintf("**Integration Points**: %d mapped\n", len(as.brownfieldCtx.IntegrationPoints)))
	summary.WriteString(fmt.Sprintf("**Technical Debt Items**: %d identified\n", len(as.brownfieldCtx.TechnicalDebt)))

	return summary.String()
}

// GetAgentForPhase returns the appropriate agent for a phase
func (as *AgentService) GetAgentForPhase(phase string) (*Agent, error) {
	return as.agentMgr.GetAgentForPhase(phase)
}

// ListAgents returns available agents
func (as *AgentService) ListAgents() []string {
	return as.agentMgr.ListAgents()
}

// GetCodebaseSummary returns a summary of the codebase
func (as *AgentService) GetCodebaseSummary() string {
	if as.lspContext == nil {
		return "Codebase analysis not available"
	}

	var summary strings.Builder

	summary.WriteString(fmt.Sprintf("## Codebase Summary\n\n"))
	summary.WriteString(fmt.Sprintf("**Language:** %s\n", as.lspContext.Structure.MainLanguage))
	summary.WriteString(fmt.Sprintf("**Framework:** %s\n", as.lspContext.Structure.Framework))
	summary.WriteString(fmt.Sprintf("**Files:** %d\n", len(as.lspContext.Files)))

	features := []string{}
	if as.lspContext.Structure.HasAPI {
		features = append(features, "API")
	}
	if as.lspContext.Structure.HasDatabase {
		features = append(features, "Database")
	}
	if as.lspContext.Structure.HasFrontend {
		features = append(features, "Frontend")
	}
	if as.lspContext.Structure.HasTests {
		features = append(features, "Tests")
	}

	if len(features) > 0 {
		summary.WriteString(fmt.Sprintf("**Features:** %s\n", strings.Join(features, ", ")))
	}

	return summary.String()
}

// PrepareArchitectRequest builds a request for the Architect with architectural guidelines
func (as *AgentService) PrepareArchitectRequest(prd string) (*mcp.ChatRequest, error) {
	systemPrompt := as.loadPrompt("architect.md")

	// Ingest the whitepaper context
	archContext := as.loadContext("architecture.md")

	return &mcp.ChatRequest{
		System: systemPrompt,
		// The Architect gets the PRD + the Architectural Whitepaper
		Context: fmt.Sprintf("PRD: %s\nARCH_GUIDELINES: %s", prd, archContext),
		// Architect is forced to use the audit skill before finishing
		Instructions: "Equip [USE_SKILL: architecture-audit] to validate your plan before outputting.",
	}, nil
}

// loadPrompt loads a system prompt from the .agents/ directory
func (as *AgentService) loadPrompt(filename string) string {
	agentName := strings.TrimSuffix(filename, ".md")
	agent, err := as.agentMgr.GetAgent(agentName)
	if err != nil {
		// Fallback or log error? For now, return empty or a basic prompt
		return ""
	}
	return agent.GetSystemPrompt()
}

// loadContext loads a specific context file from .sdd/context/
func (as *AgentService) loadContext(filename string) string {
	contextPath := filepath.Join(as.projectRoot, ".sdd", "context", filename)
	content, err := os.ReadFile(contextPath)
	if err != nil {
		return ""
	}
	return string(content)
}

// ValidateSetup checks if all required components are configured
func (as *AgentService) ValidateSetup() []string {
	var issues []string

	// Check agents
	agents := as.agentMgr.ListAgents()
	requiredAgents := []string{"pm", "architect", "developer", "qa"}
	for _, required := range requiredAgents {
		found := false
		for _, agent := range agents {
			if agent == required {
				found = true
				break
			}
		}
		if !found {
			issues = append(issues, fmt.Sprintf("Required agent '%s' not found in .agents/ directory", required))
		}
	}

	// Check MCP configuration
	providers := as.mcpMgr.ListProviders()
	if len(providers) == 0 {
		issues = append(issues, "No AI providers configured. Run 'sdd mcp add <name> --provider <provider>'")
	} else {
		enabledProviders := 0
		for _, config := range providers {
			if config.Enabled {
				enabledProviders++
			}
		}
		if enabledProviders == 0 {
			issues = append(issues, "No AI providers are enabled")
		}
	}

	// Check LSP context
	if as.lspContext == nil {
		issues = append(issues, "Codebase analysis failed")
	}

	return issues
}