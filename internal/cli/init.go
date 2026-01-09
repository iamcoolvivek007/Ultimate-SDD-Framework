package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/agents"
	"ultimate-sdd-framework/internal/gates"
)

func NewInitCmd() *cobra.Command {
	var projectName string

	cmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a new SDD project",
		Long: `Initialize a new Spec-Driven Development project.

This creates the .sdd/ directory to track project state and validates
that all required agent personas are available.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				projectName = args[0]
			}

			if projectName == "" {
				return fmt.Errorf("project name is required")
			}

			// Check if agents are available
			agentMgr := agents.NewAgentManager(".")
			if err := agentMgr.LoadAgents(); err != nil {
				return fmt.Errorf("failed to load agents: %w", err)
			}

			requiredAgents := []string{"pm", "architect", "developer", "qa"}
			availableAgents := agentMgr.ListAgents()

			for _, required := range requiredAgents {
				found := false
				for _, available := range availableAgents {
					if available == required {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("required agent not found: %s", required)
				}
			}

			// Initialize project state
			stateMgr := gates.NewStateManager(".")
			if err := stateMgr.InitializeProject(projectName); err != nil {
				return fmt.Errorf("failed to initialize project: %w", err)
			}

			fmt.Printf("✅ Successfully initialized SDD project: %s\n", projectName)
			fmt.Println("Available agents:", availableAgents)
			fmt.Println("\nNext steps:")
			fmt.Println("  sdd specify \"your feature description\"")
			fmt.Println("  sdd status  # to check project status")

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name")

	return cmd
}