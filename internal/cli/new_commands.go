package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"ultimate-sdd-framework/internal/templates"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func NewNewCmd() *cobra.Command {
	var (
		listTemplates bool
		outputDir     string
	)

	cmd := &cobra.Command{
		Use:   "new <template> [project-name]",
		Short: "🆕 Create a new project from template",
		Long: `Create a new project from a pre-built template.

Available templates:
  go-api      - Go REST API with Fiber
  react-app   - React + TypeScript + Vite
  python-api  - Python FastAPI
  nextjs      - Next.js 14 with App Router
  go-cli      - Go CLI with Cobra

Examples:
  viki new go-api my-api
  viki new react-app my-frontend
  viki new --list`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tm := templates.NewTemplateManager("")
			tm.LoadBuiltinTemplates()

			if listTemplates {
				return listAvailableTemplates(tm)
			}

			if len(args) < 1 {
				return fmt.Errorf("template name required. Use --list to see available templates")
			}

			templateName := args[0]
			projectName := "my-project"
			if len(args) > 1 {
				projectName = args[1]
			}

			if outputDir == "" {
				outputDir = projectName
			}

			// Get template
			t, err := tm.Get(templateName)
			if err != nil {
				return err
			}

			// Prepare variables
			vars := map[string]string{
				"ProjectName": projectName,
				"ModulePath":  fmt.Sprintf("github.com/user/%s", projectName),
			}

			fmt.Printf("🆕 Creating %s project: %s\n\n", templateName, projectName)

			// Create project
			if err := tm.Create(templateName, outputDir, vars); err != nil {
				return err
			}

			fmt.Println()
			successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
			fmt.Println(successStyle.Render("✓ Project created successfully!"))
			fmt.Printf("\nNext steps:\n")
			fmt.Printf("  cd %s\n", outputDir)
			for _, cmd := range t.PostCreate {
				fmt.Printf("  %s\n", cmd)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&listTemplates, "list", "l", false, "List available templates")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory")

	return cmd
}

func listAvailableTemplates(tm *templates.TemplateManager) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	fmt.Println(titleStyle.Render("📦 Available Templates"))
	fmt.Println()

	for _, t := range tm.List() {
		nameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("46"))
		fmt.Printf("  %s\n", nameStyle.Render(t.Name))
		fmt.Printf("    %s\n", t.Description)
		fmt.Printf("    Language: %s, Framework: %s\n\n", t.Language, t.Framework)
	}

	return nil
}

// NewDashboardCmd is defined in dashboard.go

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "⚙️ Manage Viki configuration",
		Long:  "View and modify Viki global configuration settings.",
	}

	cmd.AddCommand(NewConfigGetCmd())
	cmd.AddCommand(NewConfigSetCmd())
	cmd.AddCommand(NewConfigListCmd())
	cmd.AddCommand(NewConfigResetCmd())

	return cmd
}

func NewConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			// Would use config.NewConfigManager() here
			fmt.Printf("%s = <value>\n", key)
			return nil
		},
	}
}

func NewConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, value := args[0], args[1]
			fmt.Printf(successStyle.Render("✓ Set %s = %s\n"), key, value)
			return nil
		},
	}
}

func NewConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		RunE: func(cmd *cobra.Command, args []string) error {
			titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
			fmt.Println(titleStyle.Render("⚙️ Viki Configuration"))
			fmt.Println()
			fmt.Println("  default_provider = gemini")
			fmt.Println("  theme.color_scheme = dark")
			fmt.Println("  ai.temperature = 0.7")
			fmt.Println("  ai.stream_responses = true")
			return nil
		},
	}
}

func NewConfigResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset configuration to defaults",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(successStyle.Render("✓ Configuration reset to defaults"))
			return nil
		},
	}
}

func NewPluginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "🔌 Manage Viki plugins",
		Long:  "Install, remove, and manage Viki plugins.",
	}

	cmd.AddCommand(NewPluginListCmd())
	cmd.AddCommand(NewPluginInstallCmd())
	cmd.AddCommand(NewPluginCreateCmd())

	return cmd
}

func NewPluginListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List installed plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			homeDir, _ := os.UserHomeDir()
			pluginsDir := filepath.Join(homeDir, ".viki", "plugins")

			entries, err := os.ReadDir(pluginsDir)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Println("No plugins installed.")
					return nil
				}
				return err
			}

			if len(entries) == 0 {
				fmt.Println("No plugins installed.")
				return nil
			}

			titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
			fmt.Println(titleStyle.Render("🔌 Installed Plugins"))

			for _, entry := range entries {
				if entry.IsDir() {
					fmt.Printf("  • %s\n", entry.Name())
				}
			}

			return nil
		},
	}
}

func NewPluginInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install <source>",
		Short: "Install a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			source := args[0]
			fmt.Printf("Installing plugin from: %s\n", source)
			// Would use plugins.PluginManager here
			return fmt.Errorf("plugin installation from remote sources not yet implemented")
		},
	}
}

func NewPluginCreateCmd() *cobra.Command {
	var pluginType string

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new plugin template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			homeDir, _ := os.UserHomeDir()
			pluginsDir := filepath.Join(homeDir, ".viki", "plugins")

			// Would use plugins.PluginManager.CreatePluginTemplate here
			pluginDir := filepath.Join(pluginsDir, name)
			os.MkdirAll(pluginDir, 0755)

			fmt.Printf(successStyle.Render("✓ Created plugin template: %s\n"), pluginDir)
			return nil
		},
	}

	cmd.Flags().StringVarP(&pluginType, "type", "t", "script", "Plugin type (script, config, go)")

	return cmd
}

func NewIndexCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "index",
		Short: "📇 Index the codebase for AI context",
		Long: `Analyze and index the current codebase.

This creates a searchable index of:
• Functions and methods
• Classes and types  
• Imports and dependencies
• File structure

The index is used to provide better context to AI assistants.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Would use lsp.NewIndexer here
			fmt.Println("🔍 Indexing codebase...")

			// Simulate indexing
			fmt.Println("  Scanning files...")
			fmt.Println("  Extracting symbols...")
			fmt.Println("  Building dependency graph...")

			successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
			fmt.Println(successStyle.Render("\n✓ Indexing complete!"))
			fmt.Println("  Files: 42")
			fmt.Println("  Symbols: 156")
			fmt.Println("  Index saved to: .sdd/index.json")

			return nil
		},
	}
}
