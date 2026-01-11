package cli

import (
	"fmt"
	"strings"

	"ultimate-sdd-framework/internal/brainstorm"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// NewBrainstormCmd creates the brainstorm command
func NewBrainstormCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "brainstorm [topic]",
		Short: "💡 Interactive brainstorming session",
		Long: `Start an interactive brainstorming session with AI assistance.

Available Techniques:
- classic:     Traditional open brainstorming
- reverse:     Think of ways to fail, then reverse
- six_hats:    Explore from 6 perspectives
- scamper:     Systematic innovation (Substitute, Combine, etc.)
- starbursting: Question-based exploration (5W1H)
- party_mode:  Multi-agent collaborative discussion`,
		Example: `  viki brainstorm "How to improve API performance"
  viki brainstorm --technique reverse "Reduce technical debt"
  viki brainstorm --technique party_mode "Architecture decisions"
  viki brainstorm --list`,
		Run: runBrainstorm,
	}

	cmd.Flags().StringP("technique", "t", "", "Brainstorming technique to use")
	cmd.Flags().Bool("list", false, "List available techniques")
	cmd.Flags().Bool("random", false, "Use a random technique")

	return cmd
}

func runBrainstorm(cmd *cobra.Command, args []string) {
	listMode, _ := cmd.Flags().GetBool("list")
	randomMode, _ := cmd.Flags().GetBool("random")
	techniqueName, _ := cmd.Flags().GetString("technique")

	if listMode {
		listTechniques()
		return
	}

	topic := "General brainstorming"
	if len(args) > 0 {
		topic = strings.Join(args, " ")
	}

	// Select technique
	var technique *brainstorm.Technique
	if randomMode {
		technique = brainstorm.GetRandomTechnique()
	} else if techniqueName != "" {
		technique = brainstorm.GetTechniqueByID(techniqueName)
		if technique == nil {
			fmt.Printf("❌ Unknown technique: %s\n", techniqueName)
			listTechniques()
			return
		}
	} else {
		// Recommend based on topic
		technique = brainstorm.RecommendTechnique(topic)
	}

	startBrainstormSession(topic, technique)
}

func listTechniques() {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("220"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("249"))

	durationStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	fmt.Println()
	fmt.Println(titleStyle.Render("💡 Available Brainstorming Techniques"))
	fmt.Println(descStyle.Render("─────────────────────────────────────────────────"))

	for _, t := range brainstorm.GetTechniques() {
		fmt.Println()
		fmt.Printf("  %s %s\n", titleStyle.Render(t.Name), durationStyle.Render(fmt.Sprintf("(%s)", t.Duration)))
		fmt.Printf("  %s\n", descStyle.Render(t.Description))
		fmt.Printf("  ID: %s\n", t.ID)
		fmt.Printf("  Best for: %s\n", strings.Join(t.BestFor, ", "))
	}

	fmt.Println()
	fmt.Println(descStyle.Render("Usage: viki brainstorm --technique <id> \"Your topic\""))
	fmt.Println()
}

func startBrainstormSession(topic string, technique *brainstorm.Technique) {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("220"))

	techniqueStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("42"))

	stepStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("249"))

	fmt.Println()
	fmt.Println(titleStyle.Render("💡 Brainstorm Session Started"))
	fmt.Println("─────────────────────────────────────────────────")
	fmt.Printf("Topic: %s\n", topic)
	fmt.Printf("Technique: %s\n", techniqueStyle.Render(technique.Name))
	fmt.Printf("Duration: %s\n", technique.Duration)
	fmt.Println()

	// Display technique steps
	fmt.Println(titleStyle.Render("Steps:"))
	for i, step := range technique.Steps {
		fmt.Printf("  %d. %s\n", i+1, stepStyle.Render(step))
	}
	fmt.Println()

	// Generate prompt with topic
	prompt := strings.ReplaceAll(technique.Prompt, "$TOPIC", topic)

	fmt.Println(titleStyle.Render("Prompt to use:"))
	fmt.Println("┌─────────────────────────────────────────────────")
	for _, line := range strings.Split(prompt, "\n") {
		fmt.Printf("│ %s\n", line)
	}
	fmt.Println("└─────────────────────────────────────────────────")
	fmt.Println()
	fmt.Println("💡 Use this prompt in 'viki chat' to start brainstorming!")
	fmt.Println("   Or continue interactively with this session...")
}

// NewAgentSelectCmd creates the agent selection command
func NewAgentSelectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agents",
		Short: "👥 List and select AI agents",
		Long: `List available AI agents and their specializations.

Viki includes 21+ specialized agents covering:
- Core: PM, Architect, Developer, QA
- Product: UX Designer, Scrum Master, Business Analyst
- Engineering: DevOps, Security, Tech Lead, Data Architect, API Designer
- Quality: Test Automation, Performance
- Operations: SRE, Documentation
- Creative: Innovator, Reviewer, Debugger`,
		Run: runAgentList,
	}

	return cmd
}

func runAgentList(cmd *cobra.Command, args []string) {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99"))

	categoryStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("220"))

	agentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("249"))

	fmt.Println()
	fmt.Println(titleStyle.Render("👥 Available AI Agents"))
	fmt.Println(descStyle.Render("─────────────────────────────────────────────────"))

	categories := map[string][]struct {
		id   string
		name string
		role string
	}{
		"Core": {
			{"pm", "Product Manager", "Strategist"},
			{"architect", "System Architect", "Designer"},
			{"developer", "Software Developer", "Builder"},
			{"qa", "Quality Assurance", "Validator"},
		},
		"Product": {
			{"ux_designer", "UX Designer", "Experience Crafter"},
			{"scrum_master", "Scrum Master", "Facilitator"},
			{"business_analyst", "Business Analyst", "Translator"},
		},
		"Engineering": {
			{"devops", "DevOps Engineer", "Automator"},
			{"security", "Security Analyst", "Guardian"},
			{"data_architect", "Data Architect", "Data Designer"},
			{"tech_lead", "Tech Lead", "Technical Leader"},
			{"api_designer", "API Designer", "Interface Architect"},
			{"frontend", "Frontend Developer", "UI Builder"},
			{"backend", "Backend Developer", "Server Builder"},
		},
		"Quality": {
			{"test_automation", "Test Automation Engineer", "Automation Specialist"},
			{"performance", "Performance Engineer", "Optimizer"},
			{"reviewer", "Code Reviewer", "Quality Gatekeeper"},
		},
		"Operations": {
			{"sre", "Site Reliability Engineer", "Reliability Guardian"},
			{"documentation", "Technical Writer", "Documenter"},
		},
		"Creative": {
			{"innovator", "Innovation Catalyst", "Ideator"},
			{"debugger", "Debug Specialist", "Problem Solver"},
		},
	}

	order := []string{"Core", "Product", "Engineering", "Quality", "Operations", "Creative"}

	for _, cat := range order {
		agents := categories[cat]
		fmt.Println()
		fmt.Println(categoryStyle.Render(fmt.Sprintf("### %s", cat)))

		for _, agent := range agents {
			fmt.Printf("  %s (%s)\n", agentStyle.Render(agent.name), agent.id)
			fmt.Printf("    Role: %s\n", descStyle.Render(agent.role))
		}
	}

	fmt.Println()
	fmt.Println(descStyle.Render("Use agents in chat: /agent <id>"))
	fmt.Println()
}
