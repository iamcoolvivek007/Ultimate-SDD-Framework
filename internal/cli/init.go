package cli

import (
	"fmt"
	"os"
	"path/filepath"

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

			// Generate default roles if missing
			if err := generateDefaultRoles("."); err != nil {
				return fmt.Errorf("failed to generate default roles: %w", err)
			}

			// Check if agents are available
			agentMgr := agents.NewAgentManager(".")
			if err := agentMgr.LoadAgents(); err != nil {
				return fmt.Errorf("failed to load agents: %w", err)
			}

			requiredAgents := []string{"scout", "strategist", "designer", "guardian", "taskmaster", "builder", "inspector", "librarian"}
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

			// Initialize Conductor Context
			if err := initializeConductorContext("."); err != nil {
				fmt.Printf("⚠️ Warning: Failed to initialize Conductor context: %v\n", err)
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

func initializeConductorContext(root string) error {
	contextDir := filepath.Join(root, ".sdd", "context")
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		return err
	}

	// Default files content
	defaults := map[string]string{
		"product.md":   "# Product Goals & Personas\n\n## Vision\n[Define your product vision here]\n\n## User Personas\n1. [Persona 1]\n2. [Persona 2]",
		"techstack.md": "# Tech Stack & Architecture\n\n## Core Stack\n- Backend: [Language/Framework]\n- Frontend: [Library/Framework]\n- Database: [Database]",
		"workflow.md":  "# Team Workflow & Rules\n\n## Quality Gates\n- Coverage: [e.g. 80%]\n\n## Intent Gating\n- No code can be written without an approved `track_spec.md`.",
		"CONSTITUTION.md": `# SYSTEM CONSTITUTION

## I. THE PRIME DIRECTIVE
- Protect the integrity of the codebase.

## VII. THE GSD MANDATE (Action over Talk)
- Execution is handled via the ` + "`gsd.json`" + ` protocol.
- AI is forbidden from explaining its code during the Implementation phase.
- Success is defined by a green checkmark in the GSD panel, not a conversational response.
- **GET SHIT DONE.**
`,
	}

	for filename, content := range defaults {
		path := filepath.Join(contextDir, filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

func generateDefaultRoles(root string) error {
	roleDir := filepath.Join(root, ".sdd", "role")
	if err := os.MkdirAll(roleDir, 0755); err != nil {
		return err
	}

	// Write roles
	for filename, content := range agents.DefaultRoles {
		path := filepath.Join(roleDir, filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				return err
			}
		}
	}

	// Write GSD Skill
	skillDir := filepath.Join(root, ".sdd", "skill", "gsd-execute")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return err
	}
	skillPath := filepath.Join(skillDir, "SKILL.md")
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		if err := os.WriteFile(skillPath, []byte(agents.GSDSkill), 0644); err != nil {
			return err
		}
	}

	return nil
}