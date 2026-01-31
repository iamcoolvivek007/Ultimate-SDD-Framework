package main

import (
	"fmt"
	"os"

	"ultimate-sdd-framework/internal/cli"

	"github.com/spf13/cobra"
)

const version = "3.0.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "viki",
		Short: "🤖 Viki - Your AI Development Assistant",
		Long: `✨ Welcome to Viki - the Ultimate AI Development Framework!

🎯 What Viki does:
• Takes your ideas and turns them into working code
• Guides you through the development process step-by-step  
• Uses 21+ specialized AI agents for every role
• Features SQLite-persistent sessions and advanced tooling
• Supports scale-adaptive workflows (Quick → Enterprise)

🚀 Quick Start:
1. viki init "my-awesome-app"     # Start a new project
2. viki specify "your idea"       # Describe what you want
3. viki plan                      # Let AI design it
4. viki task                      # Break into tasks
5. viki execute                   # Generate the code!

💡 New in v3.0:
• viki session - Manage persistent chat sessions
• viki workflow - Structured development workflows
• viki brainstorm - Interactive ideation
• viki constitution - Project governance
• viki agents - 21+ specialized AI personas

Ready to build something amazing? Let's get started! 🚀`,
	}

	// Check if this is first run and show welcome message
	if len(os.Args) == 1 {
		fmt.Print(`🤖 Welcome to Viki v3.0 - The Ultimate AI Development Framework!

✨ Viki helps you build software using AI with 21+ specialized agents.

🚀 Quick Start:
1. viki init "project-name"       # Initialize project
2. viki mcp add ai --provider openai --model gpt-4  # Add AI
3. viki specify "your idea"       # Describe your vision
4. viki workflow init             # Get guided workflow

💡 New Commands:
• viki session - Manage AI chat sessions
• viki workflow - Development workflows (Quick/Standard/Enterprise)
• viki brainstorm - Ideation with 6 techniques
• viki constitution - Project principles
• viki agents - View 21+ specialized agents

Run 'viki --help' for all commands!

Available commands:
`)
	}

	// Core SDD commands
	rootCmd.AddCommand(cli.NewInitCmd())
	rootCmd.AddCommand(cli.NewDiscoveryCmd())
	rootCmd.AddCommand(cli.NewSpecifyCmd())
	rootCmd.AddCommand(cli.NewPlanCmd())
	rootCmd.AddCommand(cli.NewTaskCmd())
	rootCmd.AddCommand(cli.NewExecuteCmd())
	rootCmd.AddCommand(cli.NewAnalyzeCmd())
	rootCmd.AddCommand(cli.NewReviewCmd())
	rootCmd.AddCommand(cli.NewPairCmd())
	rootCmd.AddCommand(cli.NewTeamCmd())
	rootCmd.AddCommand(cli.NewLearnCmd())
	rootCmd.AddCommand(cli.NewVisionCmd())
	rootCmd.AddCommand(cli.NewPerformanceCmd())
	rootCmd.AddCommand(cli.NewEvolveCmd())
	rootCmd.AddCommand(cli.NewStatusCmd())
	rootCmd.AddCommand(cli.NewApproveCmd())
	rootCmd.AddCommand(cli.NewMCPCommand())
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newGuideCmd())

	// v2.0 commands
	rootCmd.AddCommand(cli.NewChatCmd())      // Interactive chat mode
	rootCmd.AddCommand(cli.NewUndoCmd())      // Undo file changes
	rootCmd.AddCommand(cli.NewSecretsCmd())   // Secrets management
	rootCmd.AddCommand(cli.NewNewCmd())       // Project templates
	rootCmd.AddCommand(cli.NewDashboardCmd()) // Web dashboard
	rootCmd.AddCommand(cli.NewConfigCmd())    // Global config
	rootCmd.AddCommand(cli.NewPluginCmd())    // Plugin management
	rootCmd.AddCommand(cli.NewIndexCmd())     // Codebase indexing

	// v3.0 commands - Enhanced with competitor features
	rootCmd.AddCommand(cli.NewSessionCmd())      // Session management (from OpenCode)
	rootCmd.AddCommand(cli.NewWorkflowCmd())     // Workflow engine (from BMAD)
	rootCmd.AddCommand(cli.NewBrainstormCmd())   // Brainstorming (from BMAD)
	rootCmd.AddCommand(cli.NewAgentSelectCmd())  // Agent selection (from BMAD)
	rootCmd.AddCommand(cli.NewConstitutionCmd()) // Constitution (from Spec-Kit)
	rootCmd.AddCommand(cli.NewClarifyCmd())      // Clarify specs (from Spec-Kit)
	rootCmd.AddCommand(cli.NewChecklistCmd())    // Quality checklists (from Spec-Kit)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Display the current version of Viki - The Ultimate SDD Framework",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Viki v%s - Ultimate SDD Framework\n", version)
			fmt.Println("The most advanced AI-powered development platform")
			fmt.Println("Built with ❤️ using Go and Charm")
		},
	}
}

func newGuideCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "guide",
		Short: "📚 Step-by-step guide for new users",
		Long:  "Get a friendly, step-by-step guide to start using Viki",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(`🎉 Welcome to Viki - Your AI Development Guide!

This guide will help you build your first app with Viki. Let's get started!

┌─────────────────────────────────────────────────────────────────────┐
│                           STEP 1: SETUP                             │
└─────────────────────────────────────────────────────────────────────┘

1️⃣  First, create a new project:
    viki init "my-awesome-app"

2️⃣  Add an AI assistant (choose one):
    # For OpenAI (recommended)
    viki mcp add my-openai --provider openai --model gpt-4

    # For Google Gemini (free)
    viki mcp add my-gemini --provider google --model gemini-1.5-pro

    # For Anthropic Claude
    viki mcp add my-claude --provider anthropic --model claude-3-sonnet-20240229

┌─────────────────────────────────────────────────────────────────────┐
│                        STEP 2: DESCRIBE YOUR IDEA                   │
└─────────────────────────────────────────────────────────────────────┘

3️⃣  Tell Viki what you want to build:
    viki specify "Create a todo list app where users can add, edit, delete, and mark tasks as complete"

    💡 Tip: Be specific about what you want, but don't worry about technical details!

┌─────────────────────────────────────────────────────────────────────┐
│                        STEP 3: LET VIKI WORK                        │
└─────────────────────────────────────────────────────────────────────┘

4️⃣  Viki will guide you through the rest:
    viki plan    # Create a technical plan
    viki task    # Break it into steps
    viki execute # Generate code
    viki review  # Check quality

┌─────────────────────────────────────────────────────────────────────┐
│                          💡 PRO TIPS                               │
└─────────────────────────────────────────────────────────────────────┘

• Start simple: "Build a basic blog" or "Create a weather app"
• Be descriptive: Include who will use it and key features
• Ask questions: Viki loves to explain things!
• Take breaks: You can always continue where you left off
• Experiment: Try different ideas and see what Viki creates

┌─────────────────────────────────────────────────────────────────────┐
│                        🚨 NEED HELP?                               │
└─────────────────────────────────────────────────────────────────────┘

• Run 'viki --help' to see all commands
• Run 'viki status' to see your project progress
• Visit our docs for detailed guides
• Join our community for support

Ready to build something amazing? Let's go! 🚀

Start with: viki init "your-first-app"
`)
		},
	}
}
