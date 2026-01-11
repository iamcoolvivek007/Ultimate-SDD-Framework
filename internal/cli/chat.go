package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/mcp"
)

var (
	chatUserStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	chatAssistantStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("46"))

	chatSystemStyle = lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("241"))

	chatCodeStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("252"))
)

// ChatSession represents an interactive chat session
type ChatSession struct {
	client   *mcp.ModelClient
	messages []mcp.Message
	context  string
}

// NewChatSession creates a new chat session
func NewChatSession(client *mcp.ModelClient) *ChatSession {
	return &ChatSession{
		client:   client,
		messages: []mcp.Message{},
	}
}

// AddContext adds file or project context to the session
func (s *ChatSession) AddContext(ctx string) {
	s.context = ctx
}

// SendMessage sends a message and gets a response
func (s *ChatSession) SendMessage(content string) (string, error) {
	// Add context to first message if available
	if s.context != "" && len(s.messages) == 0 {
		content = fmt.Sprintf("Context:\n%s\n\nUser request: %s", s.context, content)
	}

	s.messages = append(s.messages, mcp.Message{
		Role:    "user",
		Content: content,
	})

	response, err := s.client.Chat(s.messages, map[string]interface{}{
		"temperature": 0.7,
		"max_tokens":  4000,
	})
	if err != nil {
		return "", err
	}

	if len(response.Choices) > 0 {
		assistantMsg := response.Choices[0].Message.Content
		s.messages = append(s.messages, mcp.Message{
			Role:    "assistant",
			Content: assistantMsg,
		})
		return assistantMsg, nil
	}

	return "", fmt.Errorf("no response from AI")
}

// Clear clears the conversation history
func (s *ChatSession) Clear() {
	s.messages = []mcp.Message{}
	s.context = ""
}

// Save saves the conversation to a file
func (s *ChatSession) Save(filename string) error {
	var content strings.Builder
	for _, msg := range s.messages {
		content.WriteString(fmt.Sprintf("## %s\n\n%s\n\n---\n\n", 
			strings.Title(msg.Role), msg.Content))
	}
	return os.WriteFile(filename, []byte(content.String()), 0644)
}

func NewChatCmd() *cobra.Command {
	var (
		provider   string
		model      string
		contextDir string
	)

	cmd := &cobra.Command{
		Use:   "chat",
		Short: "💬 Start an interactive chat with AI",
		Long: `Start a continuous conversation with your AI assistant.

Unlike single-command interactions, chat mode maintains context 
across multiple messages, making it perfect for:
• Exploring ideas and iterating on solutions
• Getting progressive help with complex problems
• Pair programming sessions

Slash Commands:
  /clear    - Clear conversation history
  /save     - Save conversation to file
  /file     - Add a file to context
  /context  - Show current context
  /help     - Show available commands
  /exit     - Exit chat mode

Example:
  viki chat
  viki chat --provider my-openai`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize MCP manager
			mcpMgr := mcp.NewMCPManager(".")
			if err := mcpMgr.LoadConfig(); err != nil {
				return fmt.Errorf("failed to load MCP config: %w", err)
			}

			// Get AI client
			client, err := mcpMgr.GetClient(provider)
			if err != nil {
				return fmt.Errorf("failed to get AI client: %w", err)
			}

			if model != "" {
				client.Model = model
			}

			// Create chat session
			session := NewChatSession(client)

			// Add context if specified
			if contextDir != "" {
				ctx, err := loadDirectoryContext(contextDir)
				if err != nil {
					fmt.Printf(chatSystemStyle.Render("⚠ Could not load context: %v\n"), err)
				} else {
					session.AddContext(ctx)
					fmt.Printf(chatSystemStyle.Render("📁 Loaded context from: %s\n"), contextDir)
				}
			}

			// Print welcome message
			printChatWelcome(client)

			// Start interactive loop
			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Print(chatUserStyle.Render("\n💭 You: "))
				input, err := reader.ReadString('\n')
				if err != nil {
					break
				}

				input = strings.TrimSpace(input)
				if input == "" {
					continue
				}

				// Handle slash commands
				if strings.HasPrefix(input, "/") {
					if handleSlashCommand(input, session) {
						continue
					}
					if input == "/exit" || input == "/quit" || input == "/q" {
						fmt.Println(chatSystemStyle.Render("\n👋 Goodbye! Happy coding!"))
						break
					}
					continue
				}

				// Send message to AI
				fmt.Println(chatSystemStyle.Render("🤔 Thinking..."))
				response, err := session.SendMessage(input)
				if err != nil {
					fmt.Printf(errorStyle.Render("❌ Error: %v\n"), err)
					continue
				}

				fmt.Println()
				fmt.Println(chatAssistantStyle.Render("🤖 Viki:"))
				fmt.Println(formatResponse(response))
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&provider, "provider", "p", "", "AI provider to use")
	cmd.Flags().StringVarP(&model, "model", "m", "", "Model to use")
	cmd.Flags().StringVarP(&contextDir, "context", "c", "", "Directory to use as context")

	return cmd
}

func printChatWelcome(client *mcp.ModelClient) {
	welcomeBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1, 2).
		Render(fmt.Sprintf(`💬 Viki Chat Mode
━━━━━━━━━━━━━━━━━━━━━━
Provider: %s
Model: %s

Type your message and press Enter.
Use /help for commands, /exit to quit.`, 
		mcp.GetProviderDisplayName(client.Provider), client.Model))

	fmt.Println(welcomeBox)
}

func handleSlashCommand(input string, session *ChatSession) bool {
	parts := strings.Fields(input)
	cmd := parts[0]

	switch cmd {
	case "/clear":
		session.Clear()
		fmt.Println(chatSystemStyle.Render("🧹 Conversation cleared!"))
		return true

	case "/save":
		filename := "chat_history.md"
		if len(parts) > 1 {
			filename = parts[1]
		}
		if err := session.Save(filename); err != nil {
			fmt.Printf(errorStyle.Render("❌ Failed to save: %v\n"), err)
		} else {
			fmt.Printf(chatSystemStyle.Render("💾 Saved to: %s\n"), filename)
		}
		return true

	case "/file":
		if len(parts) < 2 {
			fmt.Println(chatSystemStyle.Render("Usage: /file <path>"))
			return true
		}
		content, err := os.ReadFile(parts[1])
		if err != nil {
			fmt.Printf(errorStyle.Render("❌ Could not read file: %v\n"), err)
		} else {
			session.AddContext(fmt.Sprintf("File: %s\n```\n%s\n```", parts[1], string(content)))
			fmt.Printf(chatSystemStyle.Render("📄 Added file to context: %s\n"), parts[1])
		}
		return true

	case "/context":
		if session.context == "" {
			fmt.Println(chatSystemStyle.Render("No context loaded."))
		} else {
			fmt.Println(chatSystemStyle.Render("Current context:"))
			// Show first 500 chars
			ctx := session.context
			if len(ctx) > 500 {
				ctx = ctx[:500] + "..."
			}
			fmt.Println(ctx)
		}
		return true

	case "/help":
		fmt.Println(chatSystemStyle.Render(`
Available Commands:
  /clear    - Clear conversation history
  /save     - Save conversation to file
  /file     - Add a file to context
  /context  - Show current context
  /help     - Show this help
  /exit     - Exit chat mode
`))
		return true

	case "/exit", "/quit", "/q":
		return false // Let main loop handle exit

	default:
		fmt.Printf(chatSystemStyle.Render("Unknown command: %s. Use /help for available commands.\n"), cmd)
		return true
	}
}

func loadDirectoryContext(dir string) (string, error) {
	var context strings.Builder
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Only include source files
		if strings.HasSuffix(name, ".go") || strings.HasSuffix(name, ".js") ||
			strings.HasSuffix(name, ".ts") || strings.HasSuffix(name, ".py") ||
			strings.HasSuffix(name, ".rs") || strings.HasSuffix(name, ".md") {
			content, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, name))
			if err != nil {
				continue
			}
			context.WriteString(fmt.Sprintf("\n--- %s ---\n%s\n", name, string(content)))
		}
	}

	return context.String(), nil
}

func formatResponse(response string) string {
	// Simple formatting - could be enhanced with markdown parsing
	lines := strings.Split(response, "\n")
	var formatted strings.Builder
	inCodeBlock := false

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			formatted.WriteString(line + "\n")
		} else if inCodeBlock {
			formatted.WriteString(chatCodeStyle.Render(line) + "\n")
		} else {
			formatted.WriteString(line + "\n")
		}
	}

	return formatted.String()
}
