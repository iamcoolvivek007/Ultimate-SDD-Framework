package agents

import (
	"fmt"
	"strings"

	"ultimate-sdd-framework/internal/lsp"
	"ultimate-sdd-framework/internal/mcp"
)

// AgentService provides high-level agent operations with context awareness
type AgentService struct {
	agentMgr    *AgentManager
	mcpMgr      *mcp.MCPManager
	lspContext  *lsp.CodebaseContext
	projectRoot string
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

	// Initialize LSP context
	as.lspContext = lsp.NewCodebaseContext(as.projectRoot)
	if err := as.lspContext.AnalyzeProject(); err != nil {
		return fmt.Errorf("failed to analyze codebase: %w", err)
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

	// Get LSP context for this phase
	contextInfo := ""
	if as.lspContext != nil {
		contextInfo = as.lspContext.GetContextForPhase(phase)
	}

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