package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/cli"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "sdd",
		Short: "Ultimate SDD Framework - Spec-Driven Development with AI Agents",
		Long: `The Ultimate SDD Framework combines structured gating, expert AI personas,
and terminal-native execution to enforce rigorous development practices.

Complete the sequence: Specify → Plan → Task → Execute with AI assistance.`,
	}

	// Add subcommands
	rootCmd.AddCommand(cli.NewInitCmd())
	rootCmd.AddCommand(cli.NewSpecifyCmd())
	rootCmd.AddCommand(cli.NewPlanCmd())
	rootCmd.AddCommand(cli.NewTaskCmd())
	rootCmd.AddCommand(cli.NewExecuteCmd())
	rootCmd.AddCommand(cli.NewReviewCmd())
	rootCmd.AddCommand(cli.NewStatusCmd())
	rootCmd.AddCommand(cli.NewApproveCmd())
	rootCmd.AddCommand(cli.NewMCPCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}