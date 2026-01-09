package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

// Agent represents a BMAD-style persona agent
type Agent struct {
	Role       string `yaml:"role"`
	Expertise  string `yaml:"expertise"`
	Personality string `yaml:"personality"`
	Tone       string `yaml:"tone"`
	Content    string `yaml:"-"`
}

// AgentManager handles loading and managing agents
type AgentManager struct {
	agentsDir string
	agents    map[string]*Agent
}

// NewAgentManager creates a new agent manager
func NewAgentManager(projectRoot string) *AgentManager {
	return &AgentManager{
		agentsDir: filepath.Join(projectRoot, ".agents"),
		agents:    make(map[string]*Agent),
	}
}

// LoadAgents loads all agent definitions from the .agents directory
func (am *AgentManager) LoadAgents() error {
	if _, err := os.Stat(am.agentsDir); os.IsNotExist(err) {
		return fmt.Errorf("agents directory not found: %s", am.agentsDir)
	}

	files, err := os.ReadDir(am.agentsDir)
	if err != nil {
		return fmt.Errorf("failed to read agents directory: %w", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		agentName := strings.TrimSuffix(file.Name(), ".md")
		agent, err := am.loadAgent(filepath.Join(am.agentsDir, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to load agent %s: %w", agentName, err)
		}

		am.agents[agentName] = agent
	}

	return nil
}

// GetAgent returns an agent by name
func (am *AgentManager) GetAgent(name string) (*Agent, error) {
	agent, exists := am.agents[name]
	if !exists {
		return nil, fmt.Errorf("agent not found: %s", name)
	}
	return agent, nil
}

// ListAgents returns a list of available agent names
func (am *AgentManager) ListAgents() []string {
	names := make([]string, 0, len(am.agents))
	for name := range am.agents {
		names = append(names, name)
	}
	return names
}

// GetAgentForPhase returns the appropriate agent for a given phase
func (am *AgentManager) GetAgentForPhase(phase string) (*Agent, error) {
	var agentName string
	switch phase {
	case "specify":
		agentName = "pm"
	case "plan":
		agentName = "architect"
	case "task", "execute":
		agentName = "developer"
	case "review":
		agentName = "qa"
	default:
		return nil, fmt.Errorf("no agent defined for phase: %s", phase)
	}

	return am.GetAgent(agentName)
}

// loadAgent loads a single agent from a markdown file with frontmatter
func (am *AgentManager) loadAgent(filePath string) (*Agent, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse frontmatter and content
	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid agent file format: missing frontmatter")
	}

	frontmatter := parts[1]
	markdownContent := parts[2]

	var agent Agent
	if err := yaml.Unmarshal([]byte(frontmatter), &agent); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	agent.Content = strings.TrimSpace(markdownContent)
	return &agent, nil
}

// GetSystemPrompt generates a system prompt for the agent
func (a *Agent) GetSystemPrompt() string {
	return fmt.Sprintf(`You are a %s agent with the following characteristics:

ROLE: %s
EXPERTISE: %s
PERSONALITY: %s
TONE: %s

%s

Always respond in character and maintain your role throughout the interaction.`, a.Role, a.Role, a.Expertise, a.Personality, a.Tone, a.Content)
}

// GetPhasePrompt returns a phase-specific prompt for the agent
func (a *Agent) GetPhasePrompt(phase, context string) string {
	var phaseInstruction string

	switch phase {
	case "specify":
		phaseInstruction = `Create detailed technical specifications based on the user's request.
		Include requirements, constraints, acceptance criteria, and edge cases.`
	case "plan":
		phaseInstruction = `Design the system architecture based on the specifications.
		Define components, technologies, data flow, and implementation approach.`
	case "task":
		phaseInstruction = `Break down the plan into specific, actionable tasks.
		Create a detailed checklist with clear deliverables and acceptance criteria.`
	case "execute":
		phaseInstruction = `Implement the feature according to the plan and specifications.
		Write clean, testable code following best practices.`
	case "review":
		phaseInstruction = `Review the implementation for quality, completeness, and adherence to requirements.
		Identify issues, suggest improvements, and validate against acceptance criteria.`
	}

	return fmt.Sprintf(`%s

Current Phase: %s
Context: %s

%s`, a.GetSystemPrompt(), strings.Title(phase), context, phaseInstruction)
}