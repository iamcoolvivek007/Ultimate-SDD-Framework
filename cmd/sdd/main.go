package main

import (
	"fmt"
	"os"

	"ultimate-sdd-framework/internal/cli"

	"github.com/spf13/cobra"
)

const version = "2.0.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "viki",
		Short: "🤖 Viki - Your AI Development Assistant",
		Long: `✨ Welcome to Viki - the friendly AI that helps you build software!

🎯 What Viki does:
• Takes your ideas and turns them into working code
• Guides you through the development process step-by-step
• Uses AI to help with planning, coding, and testing
• Works with your existing projects or helps start new ones

🚀 Quick Start for New Users:
1. viki init "my-awesome-app"    # Start a new project
2. viki specify "what you want to build"  # Tell Viki your idea
3. Follow the guided workflow!   # Viki will help with the rest

💡 Pro Tips:
• Viki works best when you describe what you want, not how to do it
• You can ask Viki to explain anything you don't understand
• Viki remembers your project and helps you continue where you left off

Ready to build something amazing? Let's get started! 🚀`,
	}

	// Check if this is first run and show welcome message
	if len(os.Args) == 1 {
		fmt.Println(`🤖 Welcome to Viki - Your AI Development Assistant!

✨ Viki helps you build software using AI. Whether you're new to coding or a seasoned developer,
Viki guides you through the development process with friendly AI assistants.

🚀 Quick Start:
1. viki init "your-project-name"     # Start a new project
2. viki mcp add my-ai --provider openai --model gpt-4  # Add AI provider
3. viki specify "what you want to build"               # Describe your idea
4. Follow Viki's guidance!                             # Let AI help you code

💡 Need help? Run 'viki --help' for all commands, or visit our docs!

Available commands:`)
		fmt.Println()
	}

	// Add subcommands
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

	// New commands (v2.0)
	rootCmd.AddCommand(cli.NewChatCmd())      // Interactive chat mode
	rootCmd.AddCommand(cli.NewUndoCmd())      // Undo file changes
	rootCmd.AddCommand(cli.NewSecretsCmd())   // Secrets management
	rootCmd.AddCommand(cli.NewNewCmd())       // Project templates
	rootCmd.AddCommand(cli.NewDashboardCmd()) // Web dashboard
	rootCmd.AddCommand(cli.NewConfigCmd())    // Global config
	rootCmd.AddCommand(cli.NewPluginCmd())    // Plugin management
	rootCmd.AddCommand(cli.NewIndexCmd())     // Codebase indexing

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
