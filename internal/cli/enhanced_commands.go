package cli

import (
	"fmt"
	"strings"

	"ultimate-sdd-framework/internal/agents"
	"ultimate-sdd-framework/internal/prompts"

	"github.com/charmbracelet/lipgloss"
)

// EnhancedSlashCommands adds more interactive slash commands to chat
type EnhancedSlashCommands struct {
	session      *ChatSession
	currentAgent *agents.ExtendedAgent
}

// NewEnhancedSlashCommands creates enhanced command handler
func NewEnhancedSlashCommands(session *ChatSession) *EnhancedSlashCommands {
	return &EnhancedSlashCommands{
		session: session,
	}
}

// HandleCommand processes enhanced slash commands
func (e *EnhancedSlashCommands) HandleCommand(input string) (handled bool, continueChat bool) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false, true
	}

	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "/agent":
		e.handleAgentSwitch(parts)
		return true, true

	case "/agents":
		e.showAgentList()
		return true, true

	case "/suggest":
		e.showSuggestions()
		return true, true

	case "/status":
		e.showStatus()
		return true, true

	case "/undo":
		e.handleUndo()
		return true, true

	case "/diff":
		e.showRecentChanges()
		return true, true

	case "/brainstorm":
		e.startBrainstorm(parts)
		return true, true

	case "/workflow":
		e.showWorkflowProgress()
		return true, true

	default:
		return false, true
	}
}

func (e *EnhancedSlashCommands) handleAgentSwitch(parts []string) {
	if len(parts) < 2 {
		// Show selection menu
		prompts.Header("👥 Switch Agent")

		agentOptions := []string{}
		agentList := agents.AllExtendedAgents()

		for _, agent := range agentList {
			agentOptions = append(agentOptions,
				fmt.Sprintf("%s (%s) - %s", agent.Name, agent.ID, agent.Role))
		}

		idx, _ := prompts.Select("Choose an agent:", agentOptions[:8]) // Show first 8 for brevity

		e.currentAgent = agentList[idx]
		prompts.Success(fmt.Sprintf("Switched to %s agent", e.currentAgent.Name))
		fmt.Printf("💬 %s is ready to help!\n", e.currentAgent.Name)
		return
	}

	// Direct agent ID provided
	agentID := parts[1]
	agent := agents.GetAgentByID(agentID)
	if agent == nil {
		prompts.Error(fmt.Sprintf("Agent '%s' not found. Use /agents to see available agents.", agentID))
		return
	}

	e.currentAgent = agent
	prompts.Success(fmt.Sprintf("Switched to %s agent", agent.Name))
	fmt.Printf("💬 %s: %s\n", agent.Name, getAgentGreeting(agent))
}

func getAgentGreeting(agent *agents.ExtendedAgent) string {
	greetings := map[string]string{
		"pm":          "Ready to help with requirements and user stories!",
		"architect":   "Let's design a solid architecture together.",
		"developer":   "Ready to write some clean, efficient code.",
		"qa":          "Time to ensure quality and catch any issues.",
		"devops":      "Let's get your deployment pipeline sorted.",
		"security":    "I'll help identify and mitigate security risks.",
		"ux_designer": "Let's create an amazing user experience!",
		"tech_lead":   "I'll guide the technical decisions.",
	}

	if greeting, ok := greetings[agent.ID]; ok {
		return greeting
	}
	return fmt.Sprintf("Ready to help with %s!", strings.Join(agent.Expertise, ", "))
}

func (e *EnhancedSlashCommands) showAgentList() {
	prompts.Header("👥 Available Agents")

	categories := map[string][]*agents.ExtendedAgent{
		"Core":        agents.GetAgentsByCategory("core"),
		"Product":     agents.GetAgentsByCategory("product"),
		"Engineering": agents.GetAgentsByCategory("engineering"),
		"Quality":     agents.GetAgentsByCategory("quality"),
		"Operations":  agents.GetAgentsByCategory("operations"),
		"Creative":    agents.GetAgentsByCategory("creative"),
	}

	order := []string{"Core", "Product", "Engineering", "Quality", "Operations", "Creative"}

	for _, cat := range order {
		agentList := categories[cat]
		if len(agentList) == 0 {
			continue
		}

		fmt.Printf("\n%s:\n", cat)
		for _, agent := range agentList {
			icon := "💼"
			switch agent.Category {
			case "core":
				icon = "⭐"
			case "engineering":
				icon = "⚙️"
			case "quality":
				icon = "✅"
			case "creative":
				icon = "💡"
			}
			fmt.Printf("  %s %s (%s)\n", icon, agent.Name, agent.ID)
		}
	}

	fmt.Println()
	prompts.Info("Use /agent <id> to switch, e.g., /agent architect")
}

func (e *EnhancedSlashCommands) showSuggestions() {
	prompts.Header("💡 Smart Suggestions")

	currentContext := "Based on your current project state:"
	fmt.Println(currentContext)
	fmt.Println()

	suggestions := []struct {
		Icon   string
		Action string
		Reason string
	}{
		{"📝", "Define acceptance criteria", "Your spec lacks specific success metrics"},
		{"🔍", "Consider edge cases", "No error handling scenarios mentioned"},
		{"🧪", "Add test strategy", "Testing approach not documented"},
		{"📊", "Review performance needs", "No performance requirements specified"},
	}

	for _, s := range suggestions {
		fmt.Printf("  %s %s\n", s.Icon, s.Action)
		dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		fmt.Printf("     %s\n", dimStyle.Render(s.Reason))
	}

	fmt.Println()
	prompts.Info("Ask me about any of these to get help!")
}

func (e *EnhancedSlashCommands) showStatus() {
	prompts.Header("📊 Project Status")

	// Simulate checking project state
	steps := []struct {
		Name string
		Done bool
	}{
		{"Initialize", true},
		{"Specify", true},
		{"Plan", true},
		{"Tasks", true},
		{"Execute", false},
		{"Review", false},
	}

	completed := 0
	for _, step := range steps {
		icon := "⬜"
		if step.Done {
			icon = "✅"
			completed++
		}
		fmt.Printf("  %s %s\n", icon, step.Name)
	}

	fmt.Println()

	progress := prompts.NewProgressBar(len(steps), "Progress:")
	progress.Update(completed)
	fmt.Println()
}

func (e *EnhancedSlashCommands) handleUndo() {
	prompts.Header("⏪ Undo Changes")

	// Show recent changes (simulated)
	changes := []string{
		"Modified: internal/cli/chat.go (5 min ago)",
		"Created: internal/prompts/prompts.go (10 min ago)",
		"Modified: cmd/sdd/main.go (15 min ago)",
	}

	for i, change := range changes {
		fmt.Printf("  %d. %s\n", i+1, change)
	}

	fmt.Println()
	if prompts.Confirm("Undo the most recent change?", false) {
		spinner := prompts.NewSpinner("Reverting changes...")
		spinner.Start()
		spinner.StopWithMessage("Changes reverted!", true)
	}
}

func (e *EnhancedSlashCommands) showRecentChanges() {
	prompts.Header("📝 Recent Changes")

	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	addStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	delStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	// Simulated diff output
	fmt.Println("internal/cli/chat.go:")
	fmt.Println(dimStyle.Render("@@ -15,6 +15,8 @@"))
	fmt.Println(addStyle.Render("+   // Enhanced slash commands"))
	fmt.Println(addStyle.Render("+   e.HandleCommand(input)"))
	fmt.Println(delStyle.Render("-   // Basic handling"))
	fmt.Println()
}

func (e *EnhancedSlashCommands) startBrainstorm(parts []string) {
	technique := "classic"
	if len(parts) > 1 {
		technique = parts[1]
	}

	prompts.Header("💡 Brainstorm Mode")
	fmt.Printf("Technique: %s\n\n", technique)

	topic := prompts.Input("What would you like to brainstorm about?", "")
	if topic == "" {
		return
	}

	fmt.Println()
	spinner := prompts.NewSpinner("Generating ideas...")
	spinner.Start()
	spinner.StopWithMessage("Ideas generated!", true)

	// Show sample ideas
	ideas := []string{
		"Consider user onboarding flow",
		"Add progress visualization",
		"Implement smart defaults",
		"Create contextual help",
	}

	fmt.Println()
	for i, idea := range ideas {
		fmt.Printf("  %d. 💡 %s\n", i+1, idea)
	}
}

func (e *EnhancedSlashCommands) showWorkflowProgress() {
	prompts.Header("🔄 Workflow Progress")

	stages := []struct {
		Name     string
		Status   string
		Duration string
	}{
		{"Quick Flow", "Available", "~5 min"},
		{"Standard Method", "In Progress", "~15 min"},
		{"Enterprise", "Available", "~30 min"},
	}

	for _, stage := range stages {
		icon := "⬜"
		if stage.Status == "In Progress" {
			icon = "🔄"
		}
		fmt.Printf("  %s %s (%s) - %s\n", icon, stage.Name, stage.Duration, stage.Status)
	}
}
