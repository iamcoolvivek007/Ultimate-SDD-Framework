package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/cli"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "nexus",
		Short: "Ultimate SDD Framework - System over Snippets Development",
		Long: `The Ultimate SDD Framework implements "System over Snippets" philosophy:
structured development with modular rules, context reset planning, and system evolution.

Complete the sequence: Specify → Plan → Task → Execute → Evolve with AI assistance.`,
	}

	// Add subcommands
	rootCmd.AddCommand(cli.NewInitCmd())
	rootCmd.AddCommand(cli.NewSpecifyCmd())
	rootCmd.AddCommand(cli.NewPlanCmd())
	rootCmd.AddCommand(cli.NewTaskCmd())
	rootCmd.AddCommand(cli.NewExecuteCmd())
	rootCmd.AddCommand(cli.NewReviewCmd())
	rootCmd.AddCommand(cli.NewEvolveCmd())
	rootCmd.AddCommand(cli.NewStatusCmd())
	rootCmd.AddCommand(cli.NewApproveCmd())
	rootCmd.AddCommand(cli.NewMCPCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}