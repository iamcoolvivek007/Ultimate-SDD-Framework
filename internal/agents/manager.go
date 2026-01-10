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
	agents       map[string]*Agent
}

// NewAgentManager creates a new agent manager
func NewAgentManager(projectRoot string) *AgentManager {
	return &AgentManager{
		// Moved from .agents to .sdd/role
		agentsDir:    filepath.Join(projectRoot, ".sdd", "role"),
		agents:       make(map[string]*Agent),
	}
}

// LoadAgents loads all agent definitions from directories
func (am *AgentManager) LoadAgents() error {
	if _, err := os.Stat(am.agentsDir); err == nil {
		if err := am.loadFromDir(am.agentsDir, false); err != nil {
			return fmt.Errorf("failed to load agents from .sdd/role: %w", err)
		}
	} else {
		return fmt.Errorf("role directory .sdd/role not found")
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
			// .sdd/role/* files are primarily markdown with frontmatter?
			// The provided prompts (guardian.md) had simple markdown headers, not YAML frontmatter.
			// However, legacy implementation expected YAML frontmatter.
			// I should support both or raw markdown.
			// Given the user provided prompts look like raw markdown without YAML frontmatter:
			// I will treat them as raw agents but extract Role from filename or content.
			agent, err = am.loadRawAgent(filePath)
		}

		if err != nil {
			return fmt.Errorf("failed to load agent %s: %w", agentName, err)
		}

		// Set the ID/Name
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
	case "discover":
		agentName = "scout"
	case "specify":
		agentName = "strategist"
	case "design", "plan":
		agentName = "designer"
	case "audit":
		agentName = "guardian"
	case "execute":
		agentName = "builder"
	case "validate", "review":
		agentName = "inspector"
	case "evolve":
		agentName = "librarian"
	case "deploy":
		agentName = "sre"
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

	name := strings.TrimSuffix(filepath.Base(filePath), ".md")

	return &Agent{
		Content: string(content),
		IsRaw:   true,
		Role:    name, // Default to filename
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
	// With the new role-based prompts, the instructions are often baked into the system prompt or skill.
	// We can keep this generic wrapper.
	return fmt.Sprintf(`Current Phase: %s
Context: %s
`, strings.Title(phase), context)
}
