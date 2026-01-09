package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/cli"
)

const version = "1.0.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "nexus",
		Short: "Ultimate SDD Framework - System over Snippets Development",
		Long: `The Ultimate SDD Framework implements "System over Snippets" philosophy:
structured development with modular rules, context reset planning, and system evolution.

Supports both greenfield (new projects) and brownfield (existing codebases) development:
- Greenfield: Standard PRD → Plan → Task → Execute → Evolve workflow
- Brownfield: Discovery → Legacy-Aware Specification → PIV Planning → Safeguard Execution

Complete the sequence: [Discovery] → Specify → Plan → Task → Execute → Evolve with AI assistance.`,
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
	rootCmd.AddCommand(cli.NewEvolveCmd())
	rootCmd.AddCommand(cli.NewStatusCmd())
	rootCmd.AddCommand(cli.NewApproveCmd())
	rootCmd.AddCommand(cli.NewMCPCommand())
	rootCmd.AddCommand(newVersionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Display the current version of the Ultimate SDD Framework",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Ultimate SDD Framework v%s\n", version)
			fmt.Println("The most advanced AI-powered development platform")
			fmt.Println("Built with ❤️ using Go and Charm")
		},
	}
}