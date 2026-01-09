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
	IsRaw      bool   `yaml:"-"`
}

// AgentManager handles loading and managing agents
type AgentManager struct {
	agentsDir    string
	autoMakerDir string
	agents       map[string]*Agent
}

// NewAgentManager creates a new agent manager
func NewAgentManager(projectRoot string) *AgentManager {
	return &AgentManager{
		agentsDir:    filepath.Join(projectRoot, ".agents"),
		autoMakerDir: filepath.Join(projectRoot, ".automaker", "system_prompts"),
		agents:       make(map[string]*Agent),
	}
}

// LoadAgents loads all agent definitions from directories
func (am *AgentManager) LoadAgents() error {
	// Strategy: Load legacy first, then overwrite with AutoMaker so AutoMaker takes precedence.

	// 1. Load from .agents (legacy/fallback)
	if _, err := os.Stat(am.agentsDir); err == nil {
		if err := am.loadFromDir(am.agentsDir, false); err != nil {
			return fmt.Errorf("failed to load legacy agents: %w", err)
		}
	}

	// 2. Load from .automaker (override/add)
	// Note: AutoMaker filenames might differ (engineer.md vs developer.md).
	if _, err := os.Stat(am.autoMakerDir); err == nil {
		if err := am.loadFromDir(am.autoMakerDir, true); err != nil {
			return fmt.Errorf("failed to load automaker agents: %w", err)
		}
	}

	return nil
}

func (am *AgentManager) loadFromDir(dir string, isRaw bool) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		agentName := strings.TrimSuffix(file.Name(), ".md")
		var agent *Agent
		var err error

		filePath := filepath.Join(dir, file.Name())
		if isRaw {
			agent, err = am.loadRawAgent(filePath)
		} else {
			agent, err = am.loadAgent(filePath)
		}

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
		// AutoMaker uses "pm", legacy might use "product_manager"
		if _, ok := am.agents["pm"]; ok {
			agentName = "pm"
		} else {
			agentName = "product_manager"
		}
	case "plan":
		agentName = "architect"
	case "task", "execute":
		// AutoMaker uses "engineer", legacy uses "developer"
		if _, ok := am.agents["engineer"]; ok {
			agentName = "engineer"
		} else {
			agentName = "developer"
		}
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

// loadRawAgent loads a single agent from a markdown file without frontmatter
func (am *AgentManager) loadRawAgent(filePath string) (*Agent, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return &Agent{
		Content: string(content),
		IsRaw:   true,
		// Role and other fields might be inferred or left empty for raw agents
		Role: "AutoMaker Agent",
	}, nil
}

// GetSystemPrompt generates a system prompt for the agent
func (a *Agent) GetSystemPrompt() string {
	if a.IsRaw {
		return a.Content
	}
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