package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/ultimate-sdd-framework/internal/mcp"
)

var (
	mcpStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("46"))

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))

	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("33"))
)

func NewMCPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Manage AI model providers and connections",
		Long: `Configure and manage AI model providers for the Ultimate SDD Framework.

Supports multiple providers: OpenAI, Anthropic, Google Gemini, Ollama, and Azure OpenAI.`,
	}

	cmd.AddCommand(NewMCPAddCmd())
	cmd.AddCommand(NewMCPRemoveCmd())
	cmd.AddCommand(NewMCPListCmd())
	cmd.AddCommand(NewMCPDefaultCmd())
	cmd.AddCommand(NewMCPTestCmd())
	cmd.AddCommand(NewMCPChatCmd())

	return cmd
}

func NewMCPAddCmd() *cobra.Command {
	var (
		provider string
		model    string
		baseURL  string
		setDefault bool
	)

	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Add a new AI provider configuration",
		Long: `Add a new AI provider configuration with API key.

Supported providers: openai, anthropic, google, ollama, azure

Example:
  sdd mcp add my-openai --provider openai --model gpt-4`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			// Validate provider
			modelProvider := mcp.ModelProvider(provider)
			validProviders := mcp.GetAvailableProviders()
			valid := false
			for _, p := range validProviders {
				if p == modelProvider {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid provider '%s'. Valid providers: %v", provider, validProviders)
			}

			// Get API key from environment or prompt
			apiKey := os.Getenv("SDD_API_KEY")
			if apiKey == "" {
				fmt.Printf("Enter API key for %s: ", mcp.GetProviderDisplayName(modelProvider))
				var err error
				apiKey, err = readPassword()
				if err != nil {
					return fmt.Errorf("failed to read API key: %w", err)
				}
				apiKey = strings.TrimSpace(apiKey)
			}

			if apiKey == "" {
				return fmt.Errorf("API key is required")
			}

			// Set default model if not provided
			if model == "" {
				model = mcp.GetDefaultModelForProvider(modelProvider)
			}

			// Initialize MCP manager
			mcpMgr := mcp.NewMCPManager(".")
			if err := mcpMgr.LoadConfig(); err != nil {
				return fmt.Errorf("failed to load MCP config: %w", err)
			}

			// Add provider
			options := make(map[string]interface{})
			if baseURL != "" {
				options["base_url"] = baseURL
			}

			if err := mcpMgr.AddProvider(name, modelProvider, apiKey, model, options); err != nil {
				return fmt.Errorf("failed to add provider: %w", err)
			}

			// Set as default if requested
			if setDefault {
				if err := mcpMgr.SetDefaultProvider(name); err != nil {
					return fmt.Errorf("failed to set default provider: %w", err)
				}
			}

			fmt.Printf(successStyle.Render("✅ Successfully added provider '%s'\n"), name)
			fmt.Printf("Provider: %s\n", mcp.GetProviderDisplayName(modelProvider))
			fmt.Printf("Model: %s\n", model)
			if setDefault {
				fmt.Println("Set as default provider")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&provider, "provider", "p", "", "AI provider (openai, anthropic, google, ollama, azure)")
	cmd.Flags().StringVarP(&model, "model", "m", "", "Model name (provider-specific)")
	cmd.Flags().StringVar(&baseURL, "base-url", "", "Custom base URL for the provider")
	cmd.Flags().BoolVar(&setDefault, "default", false, "Set this provider as the default")

	cmd.MarkFlagRequired("provider")

	return cmd
}

func NewMCPRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove an AI provider configuration",
		Long:  "Remove an AI provider configuration and its API key.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			mcpMgr := mcp.NewMCPManager(".")
			if err := mcpMgr.LoadConfig(); err != nil {
				return fmt.Errorf("failed to load MCP config: %w", err)
			}

			if err := mcpMgr.RemoveProvider(name); err != nil {
				return fmt.Errorf("failed to remove provider: %w", err)
			}

			fmt.Printf(successStyle.Render("✅ Successfully removed provider '%s'\n"), name)
			return nil
		},
	}

	return cmd
}

func NewMCPListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List configured AI providers",
		Long:  "Display all configured AI providers and their status.",
		RunE: func(cmd *cobra.Command, args []string) error {
			mcpMgr := mcp.NewMCPManager(".")
			if err := mcpMgr.LoadConfig(); err != nil {
				return fmt.Errorf("failed to load MCP config: %w", err)
			}

			providers := mcpMgr.ListProviders()

			if len(providers) == 0 {
				fmt.Println(infoStyle.Render("No AI providers configured."))
				fmt.Println("Use 'sdd mcp add <name> --provider <provider>' to add one.")
				return nil
			}

			fmt.Println(mcpStyle.Render("🤖 Configured AI Providers"))
			fmt.Println(strings.Repeat("=", 50))

			for name, config := range providers {
				status := "❌ Disabled"
				if config.Enabled {
					status = "✅ Enabled"
				}

				defaultMark := ""
				if name == mcpMgr.GetDefaultProvider() {
					defaultMark = " (default)"
				}

				fmt.Printf("%s %s\n", successStyle.Render(name+defaultMark), status)
				fmt.Printf("  Provider: %s\n", mcp.GetProviderDisplayName(config.Provider))
				fmt.Printf("  Model: %s\n", config.Model)
				if config.BaseURL != "" {
					fmt.Printf("  Base URL: %s\n", config.BaseURL)
				}
				fmt.Println()
			}

			return nil
		},
	}

	return cmd
}

func NewMCPDefaultCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default <name>",
		Short: "Set the default AI provider",
		Long:  "Set which AI provider to use by default for all operations.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			mcpMgr := mcp.NewMCPManager(".")
			if err := mcpMgr.LoadConfig(); err != nil {
				return fmt.Errorf("failed to load MCP config: %w", err)
			}

			if err := mcpMgr.SetDefaultProvider(name); err != nil {
				return fmt.Errorf("failed to set default provider: %w", err)
			}

			fmt.Printf(successStyle.Render("✅ Set '%s' as the default provider\n"), name)
			return nil
		},
	}

	return cmd
}

func NewMCPTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test [name]",
		Short: "Test AI provider connection",
		Long: `Test the connection to an AI provider.

If no provider name is specified, tests the default provider.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			providerName := ""
			if len(args) > 0 {
				providerName = args[0]
			}

			mcpMgr := mcp.NewMCPManager(".")
			if err := mcpMgr.LoadConfig(); err != nil {
				return fmt.Errorf("failed to load MCP config: %w", err)
			}

			client, err := mcpMgr.GetClient(providerName)
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}

			fmt.Printf("Testing connection to %s...\n", mcp.GetProviderDisplayName(client.Provider))

			if err := client.ValidateConnection(); err != nil {
				fmt.Printf(errorStyle.Render("❌ Connection failed: %v\n"), err)
				return err
			}

			fmt.Println(successStyle.Render("✅ Connection successful!"))
			return nil
		},
	}

	return cmd
}

func NewMCPChatCmd() *cobra.Command {
	var (
		provider string
		model    string
		temp     float64
		maxTokens int
	)

	cmd := &cobra.Command{
		Use:   "chat <message>",
		Short: "Send a chat message to an AI provider",
		Long: `Send a chat message to an AI provider and get a response.

Useful for testing configurations and direct interaction with AI models.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			message := strings.Join(args, " ")

			mcpMgr := mcp.NewMCPManager(".")
			if err := mcpMgr.LoadConfig(); err != nil {
				return fmt.Errorf("failed to load MCP config: %w", err)
			}

			// Use specified provider or default
			providerName := ""
			if provider != "" {
				providerName = provider
			}

			client, err := mcpMgr.GetClient(providerName)
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}

			// Override model if specified
			if model != "" {
				client.Model = model
			}

			messages := []mcp.Message{
				{Role: "user", Content: message},
			}

			options := make(map[string]interface{})
			if temp > 0 {
				options["temperature"] = temp
			}
			if maxTokens > 0 {
				options["max_tokens"] = maxTokens
			}

			fmt.Printf("🤖 %s (%s)\n", mcp.GetProviderDisplayName(client.Provider), client.Model)
			fmt.Println("Thinking...")

			response, err := client.Chat(messages, options)
			if err != nil {
				return fmt.Errorf("chat failed: %w", err)
			}

			if len(response.Choices) > 0 {
				fmt.Println(successStyle.Render("Response:"))
				fmt.Println(response.Choices[0].Message.Content)
				fmt.Printf(infoStyle.Render("\nUsage: %d tokens (%d prompt, %d completion)\n"),
					response.Usage.TotalTokens,
					response.Usage.PromptTokens,
					response.Usage.CompletionTokens)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&provider, "provider", "p", "", "Provider name (uses default if not specified)")
	cmd.Flags().StringVarP(&model, "model", "m", "", "Model name (overrides provider default)")
	cmd.Flags().Float64VarP(&temp, "temperature", "t", 0.7, "Temperature for response randomness")
	cmd.Flags().IntVarP(&maxTokens, "max-tokens", "x", 1000, "Maximum tokens in response")

	return cmd
}

// readPassword reads a password from stdin without echoing
func readPassword() (string, error) {
	// For demo purposes, we'll just read from stdin
	// In a real implementation, you'd use terminal.ReadPassword
	var password string
	fmt.Scanln(&password)
	return password, nil
}