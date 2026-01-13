package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ultimate-sdd-framework/internal/lsp"
	"ultimate-sdd-framework/internal/mcp"

	"github.com/goccy/go-yaml"
)

// AgentService provides high-level agent operations with context awareness
type AgentService struct {
	agentMgr             *AgentManager
	mcpMgr               *mcp.MCPManager
	lspContext           *lsp.CodebaseContext
	brownfieldCtx        *lsp.BrownfieldContext
	projectRoot          string
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

	// Check for brownfield context (CONTEXT.md) in .sdd/context/
	// (migrated from .sdd/CONTEXT.md based on new structure, but keeping logic resilient)
	contextPath := filepath.Join(as.projectRoot, ".sdd", "context", "current_state.md")
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

// Orchestrate handles the 7-Gate SDD Workflow
func (as *AgentService) Orchestrate(phase string, trackID string, userInput string) (string, error) {
	// 1. Identify Role and Artifacts based on Phase
	roleName, prevArtifact, currentArtifact, skill := as.getPhaseConfig(phase)

	// 2. Gatekeeper Check: Ensure previous phase artifact is APPROVED
	if prevArtifact != "" {
		approved, err := as.checkGateApproval(trackID, prevArtifact)
		if err != nil {
			return "", fmt.Errorf("gate check failed: %w", err)
		}
		if !approved {
			return "", fmt.Errorf("403 FORBIDDEN: Previous gate artifact '%s' is missing or not APPROVED", prevArtifact)
		}
	}

	// 3. Prepare Context
	contextInfo, err := as.prepareContext(phase, trackID, prevArtifact)
	if err != nil {
		return "", fmt.Errorf("failed to prepare context: %w", err)
	}

	// 4. Special Handling for Security Gate (Guardian)
	if phase == "audit" {
		return as.runSecurityGate(trackID, contextInfo)
	}

	// 5. Get Agent Response
	response, err := as.GetAgentResponse(roleName, phase, userInput, contextInfo, skill)
	if err != nil {
		return "", err
	}

	// 6. Save Artifact (Draft)
	if err := as.SaveArtifact(trackID, currentArtifact, response, "PENDING"); err != nil {
		return "", fmt.Errorf("failed to save artifact: %w", err)
	}

	return response, nil
}

func (as *AgentService) getPhaseConfig(phase string) (role, prev, curr, skill string) {
	switch phase {
	case "discover":
		return "scout", "", "0_discovery.md", "research-codebase"
	case "specify":
		return "strategist", "0_discovery.md", "1_prd.md", "create-prd"
	case "design":
		return "designer", "1_prd.md", "2_architecture.md", "plan-feature"
	case "audit":
		return "guardian", "2_architecture.md", "3_security_report.md", "architecture-audit"
	case "task":
		return "taskmaster", "2_architecture.md", "gsd.json", "plan-feature"
	case "execute":
		return "builder", "gsd.json", "source_code", "gsd-execute" // Builder follows GSD checklist
	case "validate":
		return "inspector", "source_code", "5_validation_report.md", "piv-validate"
	case "evolve":
		return "librarian", "5_validation_report.md", "context_update", "system-evolution"
	default:
		return "", "", "", ""
	}
}

func (as *AgentService) checkGateApproval(trackID, artifactName string) (bool, error) {
	// For "source_code", we assume implicit approval if validation is running,
	// or we might check git status. For now, skip file check for source_code.
	if artifactName == "source_code" {
		return true, nil
	}

	path := filepath.Join(as.projectRoot, ".sdd", "tracks", trackID, artifactName)
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// Parse frontmatter
	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) < 3 {
		return false, nil // No frontmatter
	}

	var metadata map[string]interface{}
	if err := yaml.Unmarshal([]byte(parts[1]), &metadata); err != nil {
		return false, err
	}

	status, ok := metadata["status"].(string)
	if !ok {
		return false, nil
	}

	return strings.ToUpper(status) == "APPROVED", nil
}

func (as *AgentService) prepareContext(phase, trackID, prevArtifact string) (string, error) {
	var contextBuilder strings.Builder

	// 1. Ingest previous artifact if exists
	if prevArtifact != "" && prevArtifact != "source_code" {
		path := filepath.Join(as.projectRoot, ".sdd", "tracks", trackID, prevArtifact)
		content, err := os.ReadFile(path)
		if err == nil {
			contextBuilder.WriteString(fmt.Sprintf("\n\n## INPUT ARTIFACT (%s)\n%s\n", prevArtifact, string(content)))
		}
	}

	// 2. Add Scout's Landscape for Strategist
	if phase == "specify" {
		// Already handled by prevArtifact="0_discovery.md"
	}

	// 3. Add Builder Constraints (Blind to PRD, sees GSD + Arch Spec + Security Report)
	if phase == "execute" {
		// GSD is in prevArtifact.
		// Need to inject Arch Spec and Security Report as well.
		archPath := filepath.Join(as.projectRoot, ".sdd", "tracks", trackID, "2_architecture.md")
		archContent, err := os.ReadFile(archPath)
		if err == nil {
			contextBuilder.WriteString(fmt.Sprintf("\n\n## ARCHITECTURE SPECIFICATION\n%s\n", string(archContent)))
		}

		secPath := filepath.Join(as.projectRoot, ".sdd", "tracks", trackID, "3_security_report.md")
		secContent, err := os.ReadFile(secPath)
		if err == nil {
			contextBuilder.WriteString(fmt.Sprintf("\n\n## SECURITY CONSTRAINTS (MANDATORY)\n%s\n", string(secContent)))
		}
	}

	// 4. Inject Brownfield/Legacy Context (Moved logic here or rely on Scout's finding)
	// If Scout has run, the Brownfield info is in 0_discovery.md.
	// However, we can also inject the raw brownfield constraints for the Scout to *create* that discovery.
	if phase == "discover" {
		contextBuilder.WriteString(as.getBrownfieldConstraintsForPhase("discover"))
	}

	// 5. Inject Conductor Context
	contextBuilder.WriteString(as.getConductorContext())

	return contextBuilder.String(), nil
}

// runSecurityGate is the specialized logic for the Guardian
func (as *AgentService) runSecurityGate(trackID, contextInfo string) (string, error) {
	fmt.Println("🛡️  Gate 3: Security Guardian is auditing the design...")

	// The contextInfo already contains the ARCH_SPEC (prevArtifact)

	agentName := "guardian"
	skill := "architecture-audit"

	// Get Agent
	agent, err := as.agentMgr.GetAgent(agentName)
	if err != nil {
		return "", fmt.Errorf("agent not found: %w", err)
	}

	systemPrompt := agent.GetSystemPrompt()
	// Add skill instruction
	systemPrompt += fmt.Sprintf("\n\n[SYSTEM]: You have equipped the skill '%s'. Use it to perform your task.", skill)

	prompt := fmt.Sprintf("%s\n\nCONTEXT:\n%s\n\nINSTRUCTIONS: Perform a deep security audit. Find at least one risk. Issue a PASS/FAIL verdict.", systemPrompt, contextInfo)

	// Call AI
	client, err := as.mcpMgr.GetClient("")
	if err != nil {
		return "", err
	}

	messages := []mcp.Message{
		{Role: "user", Content: prompt},
	}

	resp, err := client.Chat(messages, map[string]interface{}{"temperature": 0.0}) // Low temp for audit
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response")
	}

	report := resp.Choices[0].Message.Content

	// Check for FAIL (Gate Blocking)
	if strings.Contains(report, "[STATUS: FAIL]") {
		fmt.Println("❌ SECURITY GATE BLOCKED: Implementation cannot proceed.")
		// We still save the report so the Architect sees it
		as.SaveArtifact(trackID, "3_security_report.md", report, "REJECTED")
		return report, nil // Return report but user needs to revise Arch Spec
	}

	fmt.Println("✅ SECURITY GATE PASSED.")
	as.SaveArtifact(trackID, "3_security_report.md", report, "APPROVED") // Auto-approve if passed? Or wait for human?
	// Prompt says: "Human approves the security hardening."
	// But Guardian output says "Status: PASS".
	// Let's set it to PENDING so human can confirm.
	as.SaveArtifact(trackID, "3_security_report.md", report, "PENDING")

	return report, nil
}

// GetAgentResponse gets a response from an agent with full context
func (as *AgentService) GetAgentResponse(agentName, phase, userInput, contextInfo, skill string) (string, error) {
	// Get the agent
	agent, err := as.agentMgr.GetAgent(agentName)
	if err != nil {
		return "", fmt.Errorf("agent not found: %w", err)
	}

	// Build the full prompt
	systemPrompt := agent.GetSystemPrompt()

	// Inject Skill
	if skill != "" {
		systemPrompt += fmt.Sprintf("\n\n[SYSTEM]: You have equipped the skill '%s'. Use it to perform your task.", skill)
		// Optionally load skill instructions from .sdd/skill/<skill>/SKILL.md and append
		skillPath := filepath.Join(as.projectRoot, ".sdd", "skill", skill, "SKILL.md")
		skillContent, err := os.ReadFile(skillPath)
		if err == nil {
			systemPrompt += fmt.Sprintf("\n\nSKILL INSTRUCTIONS:\n%s", string(skillContent))
		}
	}

	// Phase Prompt (Agent's internal logic)
	phasePrompt := agent.GetPhasePrompt(phase, contextInfo)

	// Combine with user input
	prompt := fmt.Sprintf("%s\n\n%s\n\nUser Input: %s", systemPrompt, phasePrompt, userInput)

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

// SaveArtifact writes content to the track folder with frontmatter
func (as *AgentService) SaveArtifact(trackID, filename, content, status string) error {
	dir := filepath.Join(as.projectRoot, ".sdd", "tracks", trackID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	fullContent := fmt.Sprintf("---\nstatus: %s\nphase: %s\n---\n\n%s", status, strings.TrimSuffix(filename, ".md"), content)

	return os.WriteFile(filepath.Join(dir, filename), []byte(fullContent), 0644)
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
// Adapted to be used primarily by the Scout or inserted into Discover phase
func (as *AgentService) getBrownfieldConstraintsForPhase(phase string) string {
	if !as.hasBrownfieldContext || as.brownfieldCtx == nil {
		return ""
	}

	var constraints strings.Builder
	constraints.WriteString("\n\n## 🔧 Brownfield Constraints (System Detected)\n\n")

	// General constraints applicable to discovery/scouting
	constraints.WriteString("- **Existing Architecture**: Must integrate with current system design\n")
	constraints.WriteString("- **Forbidden Patterns**: Avoid anti-patterns identified in CONTEXT.md\n")
	constraints.WriteString("- **Integration Points**: Consider existing API and database touchpoints\n")
	constraints.WriteString("- **Technical Debt**: Acknowledge known issues and limitations\n")

	if len(as.brownfieldCtx.ForbiddenPatterns) > 0 {
		constraints.WriteString("\n### Prohibited Approaches\n")
		for _, pattern := range as.brownfieldCtx.ForbiddenPatterns {
			constraints.WriteString(fmt.Sprintf("- **%s**: %s\n", pattern.Pattern, pattern.Recommended))
		}
	}

	if len(as.brownfieldCtx.IntegrationPoints) > 0 {
		constraints.WriteString("\n### Key Integration Points\n")
		for _, point := range as.brownfieldCtx.IntegrationPoints {
			constraints.WriteString(fmt.Sprintf("- **%s**: %s\n", point.Name, point.Description))
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

// loadPrompt loads a system prompt from the .agents/ directory
// DEPRECATED: Use AgentManager
func (as *AgentService) loadPrompt(filename string) string {
	agentName := strings.TrimSuffix(filename, ".md")
	agent, err := as.agentMgr.GetAgent(agentName)
	if err != nil {
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
	requiredAgents := []string{"strategist", "designer", "builder", "inspector", "scout", "guardian", "librarian"}
	for _, required := range requiredAgents {
		found := false
		for _, agent := range agents {
			if agent == required {
				found = true
				break
			}
		}
		if !found {
			issues = append(issues, fmt.Sprintf("Required agent '%s' not found in .sdd/role/ directory", required))
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
