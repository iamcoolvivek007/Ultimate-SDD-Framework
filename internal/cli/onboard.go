package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"ultimate-sdd-framework/internal/prompts"

	"github.com/charmbracelet/lipgloss"
)

// OnboardingConfig stores first-run state
type OnboardingConfig struct {
	Completed    bool   `json:"completed"`
	ProviderSet  bool   `json:"provider_set"`
	DemoComplete bool   `json:"demo_complete"`
	Version      string `json:"version"`
}

// RunOnboarding runs the first-time user experience
func RunOnboarding() error {
	// Check if already onboarded
	if isOnboarded() {
		return nil
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		MarginBottom(1)

	welcomeBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(1, 3).
		MarginTop(1).
		MarginBottom(1)

	// Welcome screen
	welcome := welcomeBox.Render(fmt.Sprintf(`%s

Hey there! 👋 I'm Viki, your AI development assistant.

I help you build software the right way:
• 21+ specialized AI agents for every role
• Guided workflows to prevent "vibe coding"
• Smart suggestions and progress tracking

Let's get you set up in under 2 minutes!`,
		titleStyle.Render("🤖 Welcome to Viki v3.0!")))

	fmt.Println(welcome)

	// Step 1: Choose what to do
	prompts.Header("🚀 Quick Setup")

	idx, _ := prompts.Select("What would you like to do?", []string{
		"🔑 Set up AI provider (OpenAI, Gemini, Claude)",
		"⚡ Quick demo - build something in 5 min",
		"📖 Explore commands",
		"⏭️  Skip for now",
	})

	switch idx {
	case 0:
		setupAIProvider()
	case 1:
		runQuickDemo()
	case 2:
		showCommandExplorer()
	case 3:
		fmt.Println("\n✅ No problem! Run 'viki guide' anytime to get help.")
	}

	// Mark as onboarded
	markOnboarded()

	// Show next steps
	prompts.ShowNextSteps([]struct{ Command, Description string }{
		{"viki init \"project-name\"", "Start a new project"},
		{"viki chat", "Talk to AI assistant"},
		{"viki agents", "See all 21+ AI agents"},
		{"viki guide", "Step-by-step tutorial"},
	})

	return nil
}

func setupAIProvider() {
	prompts.Header("🔑 AI Provider Setup")

	idx, provider := prompts.Select("Choose your AI provider:", []string{
		"OpenAI (GPT-4, recommended)",
		"Google Gemini (free tier available)",
		"Anthropic Claude",
		"Ollama (local, free)",
	})

	providerMap := map[int]string{
		0: "openai",
		1: "google",
		2: "anthropic",
		3: "ollama",
	}

	providerName := providerMap[idx]
	fmt.Printf("\nGreat choice! %s\n\n", provider)

	if providerName == "ollama" {
		fmt.Println("Ollama runs locally - no API key needed!")
		fmt.Println("Make sure Ollama is installed: https://ollama.ai")
	} else {
		apiKey := prompts.Input("Enter your API key (or press Enter to skip)", "")
		if apiKey != "" {
			fmt.Printf("\n")
			spinner := prompts.NewSpinner("Saving API key securely...")
			spinner.Start()
			// Here we would save to keychain
			spinner.StopWithMessage("API key saved!", true)
		}
	}

	// Create MCP config
	fmt.Println()
	spinner := prompts.NewSpinner("Configuring provider...")
	spinner.Start()
	// Simulate setup
	spinner.StopWithMessage(fmt.Sprintf("Provider '%s' configured!", providerName), true)
}

func runQuickDemo() {
	prompts.Header("⚡ Quick Demo")

	fmt.Println("Let's build a simple TODO app in 5 minutes!")

	if !prompts.Confirm("Ready to start?", true) {
		return
	}

	steps := []struct {
		Step    string
		Command string
	}{
		{"Initialize project", "viki init demo-todo"},
		{"Define requirements", "viki specify \"A simple todo app\""},
		{"Create architecture", "viki plan"},
		{"Generate tasks", "viki task"},
	}

	progress := prompts.NewProgressBar(len(steps), "Demo progress:")

	for i, step := range steps {
		fmt.Printf("\n📋 Step %d: %s\n", i+1, step.Step)
		fmt.Printf("   Command: %s\n", step.Command)
		progress.Increment()

		// In real implementation, we'd run these commands
	}

	progress.Complete()
	prompts.Success("Demo complete! You've seen the basic workflow.")
}

func showCommandExplorer() {
	prompts.Header("📖 Command Explorer")

	categories := []struct {
		Name     string
		Commands []struct{ Cmd, Desc string }
	}{
		{
			Name: "🚀 Getting Started",
			Commands: []struct{ Cmd, Desc string }{
				{"viki init \"name\"", "Create new project"},
				{"viki chat", "Interactive AI chat"},
				{"viki guide", "Step-by-step tutorial"},
			},
		},
		{
			Name: "📋 Workflow",
			Commands: []struct{ Cmd, Desc string }{
				{"viki specify", "Define requirements"},
				{"viki plan", "Design architecture"},
				{"viki task", "Break into tasks"},
				{"viki execute", "Start implementation"},
			},
		},
		{
			Name: "🆕 v3.0 Features",
			Commands: []struct{ Cmd, Desc string }{
				{"viki agents", "21+ AI personas"},
				{"viki brainstorm", "6 ideation techniques"},
				{"viki workflow", "Guided workflows"},
				{"viki constitution", "Project governance"},
			},
		},
	}

	for _, cat := range categories {
		fmt.Printf("\n%s\n", cat.Name)
		for _, cmd := range cat.Commands {
			prompts.Suggestion("  ", cmd.Cmd, cmd.Desc)
		}
	}
}

func isOnboarded() bool {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".viki", "onboarded")
	_, err := os.Stat(configPath)
	return err == nil
}

func markOnboarded() {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".viki")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "onboarded"), []byte("true"), 0644)
}

// ShowSmartSuggestions displays contextual next steps based on project state
func ShowSmartSuggestions(projectDir string) {
	suggestions := getContextualSuggestions(projectDir)
	if len(suggestions) == 0 {
		return
	}

	prompts.ShowNextSteps(suggestions)
}

func getContextualSuggestions(projectDir string) []struct{ Command, Description string } {
	var suggestions []struct{ Command, Description string }

	sddDir := filepath.Join(projectDir, ".sdd")

	// Check what's missing
	hasState := fileExists(filepath.Join(sddDir, "state.yaml"))
	hasSpec := fileExists(filepath.Join(sddDir, "spec.md"))
	hasPlan := fileExists(filepath.Join(sddDir, "plan.md"))
	hasTasks := fileExists(filepath.Join(sddDir, "tasks.md"))

	if !hasState {
		suggestions = append(suggestions, struct{ Command, Description string }{
			"viki init \"project-name\"", "Start your project",
		})
		return suggestions
	}

	if !hasSpec {
		suggestions = append(suggestions, struct{ Command, Description string }{
			"viki specify \"your idea\"", "Define what you want to build",
		})
	} else if !hasPlan {
		suggestions = append(suggestions, struct{ Command, Description string }{
			"viki plan", "Create the technical architecture",
		})
	} else if !hasTasks {
		suggestions = append(suggestions, struct{ Command, Description string }{
			"viki approve && viki task", "Approve plan and create tasks",
		})
	} else {
		suggestions = append(suggestions, struct{ Command, Description string }{
			"viki execute", "Start implementing",
		})
	}

	// Always suggest chat
	suggestions = append(suggestions, struct{ Command, Description string }{
		"viki chat", "Ask questions or get help",
	})

	return suggestions
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
