package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"ultimate-sdd-framework/internal/db"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// NewSessionCmd creates the session management command
func NewSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "💬 Manage chat sessions",
		Long: `Manage AI chat sessions with persistent history.

Sessions are stored in SQLite and can be:
- Listed to see all conversations
- Switched between
- Exported to markdown
- Deleted when no longer needed`,
	}

	cmd.AddCommand(newSessionListCmd())
	cmd.AddCommand(newSessionSwitchCmd())
	cmd.AddCommand(newSessionDeleteCmd())
	cmd.AddCommand(newSessionExportCmd())
	cmd.AddCommand(newSessionNewCmd())

	return cmd
}

func newSessionListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all sessions",
		Run: func(cmd *cobra.Command, args []string) {
			database, err := getDatabase()
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				return
			}
			defer database.Close()

			store := db.NewSessionStore(database)
			sessions, err := store.List("", 20)
			if err != nil {
				fmt.Printf("❌ Error listing sessions: %v\n", err)
				return
			}

			if len(sessions) == 0 {
				fmt.Println("📭 No sessions found. Start one with: viki chat")
				return
			}

			titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
			activeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
			dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

			fmt.Println(titleStyle.Render("\n💬 Chat Sessions"))
			fmt.Println(dimStyle.Render("─────────────────────────────────────────────────"))

			for _, s := range sessions {
				activeMarker := "  "
				style := dimStyle
				if s.IsActive {
					activeMarker = "▶ "
					style = activeStyle
				}

				age := formatAge(s.UpdatedAt)
				msgInfo := fmt.Sprintf("%d msgs", s.MessageCount)

				fmt.Printf("%s%s %s %s\n",
					style.Render(activeMarker),
					s.Title,
					dimStyle.Render(fmt.Sprintf("(%s)", age)),
					dimStyle.Render(msgInfo))
				fmt.Printf("   ID: %s\n", dimStyle.Render(s.ID))
			}

			fmt.Println()
		},
	}
}

func newSessionSwitchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "switch <session-id>",
		Short: "Switch to a different session",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sessionID := args[0]

			database, err := getDatabase()
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				return
			}
			defer database.Close()

			store := db.NewSessionStore(database)

			// Check if session exists
			session, err := store.GetByID(sessionID)
			if err != nil || session == nil {
				fmt.Printf("❌ Session not found: %s\n", sessionID)
				return
			}

			// Set as active
			if err := store.SetActive(sessionID); err != nil {
				fmt.Printf("❌ Error switching session: %v\n", err)
				return
			}

			fmt.Printf("✅ Switched to session: %s\n", session.Title)
		},
	}
}

func newSessionDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <session-id>",
		Short: "Delete a session",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sessionID := args[0]

			database, err := getDatabase()
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				return
			}
			defer database.Close()

			store := db.NewSessionStore(database)

			if err := store.Delete(sessionID); err != nil {
				fmt.Printf("❌ Error deleting session: %v\n", err)
				return
			}

			fmt.Printf("✅ Session deleted: %s\n", sessionID)
		},
	}
}

func newSessionExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export <session-id>",
		Short: "Export a session to markdown",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sessionID := args[0]

			database, err := getDatabase()
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				return
			}
			defer database.Close()

			store := db.NewSessionStore(database)
			msgStore := db.NewMessageStore(database)

			session, err := store.GetByID(sessionID)
			if err != nil || session == nil {
				fmt.Printf("❌ Session not found: %s\n", sessionID)
				return
			}

			messages, err := msgStore.ListBySession(sessionID, 0)
			if err != nil {
				fmt.Printf("❌ Error loading messages: %v\n", err)
				return
			}

			// Generate markdown
			md := fmt.Sprintf("# %s\n\n", session.Title)
			md += fmt.Sprintf("**Date**: %s\n", session.CreatedAt.Format("2006-01-02 15:04"))
			md += fmt.Sprintf("**Model**: %s\n\n", session.Model)
			md += "---\n\n"

			for _, msg := range messages {
				role := "**User**"
				if msg.Role == "assistant" {
					role = "**Assistant**"
				}
				md += fmt.Sprintf("%s:\n\n%s\n\n---\n\n", role, msg.Content)
			}

			// Save to file
			filename := fmt.Sprintf("session_%s.md", sessionID)
			if err := os.WriteFile(filename, []byte(md), 0644); err != nil {
				fmt.Printf("❌ Error saving: %v\n", err)
				return
			}

			fmt.Printf("✅ Exported to: %s\n", filename)
		},
	}
}

func newSessionNewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new [title]",
		Short: "Create a new session",
		Run: func(cmd *cobra.Command, args []string) {
			title := "New Session"
			if len(args) > 0 {
				title = args[0]
			}

			database, err := getDatabase()
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				return
			}
			defer database.Close()

			store := db.NewSessionStore(database)

			wd, _ := os.Getwd()
			session := &db.Session{
				Title:       title,
				ProjectPath: wd,
				IsActive:    true,
			}

			if err := store.Create(session); err != nil {
				fmt.Printf("❌ Error creating session: %v\n", err)
				return
			}

			// Set as active
			store.SetActive(session.ID)

			fmt.Printf("✅ Created new session: %s\n", title)
			fmt.Printf("   ID: %s\n", session.ID)
		},
	}
}

// NewWorkflowCmd creates the workflow command
func NewWorkflowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "🔄 Run development workflows",
		Long: `Run structured development workflows.

Workflows guide you through multi-step development processes:
- Quick Flow: Bug fixes, small changes (~5 min)
- Standard: Products and platforms (~15 min)
- Enterprise: Compliance-heavy systems (~30 min)`,
	}

	cmd.AddCommand(newWorkflowInitCmd())
	cmd.AddCommand(newWorkflowStatusCmd())
	cmd.AddCommand(newWorkflowNextCmd())
	cmd.AddCommand(newWorkflowListCmd())

	return cmd
}

func newWorkflowInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize workflow for current project",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🔍 Analyzing project to recommend workflow track...")

			wd, _ := os.Getwd()

			// Simple track detection
			fileCount := 0
			filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				ext := filepath.Ext(path)
				if ext == ".go" || ext == ".py" || ext == ".js" || ext == ".ts" {
					fileCount++
				}
				return nil
			})

			var track string
			var timeEstimate string

			switch {
			case fileCount > 50:
				track = "Enterprise"
				timeEstimate = "~30 minutes"
			case fileCount > 10:
				track = "Standard Method"
				timeEstimate = "~15 minutes"
			default:
				track = "Quick Flow"
				timeEstimate = "~5 minutes"
			}

			titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
			highlightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))

			fmt.Println()
			fmt.Println(titleStyle.Render("📋 Workflow Recommendation"))
			fmt.Println("─────────────────────────────────────────────────")
			fmt.Printf("Recommended Track: %s\n", highlightStyle.Render(track))
			fmt.Printf("Time to First Story: %s\n", timeEstimate)
			fmt.Printf("Files Analyzed: %d\n", fileCount)
			fmt.Println()
			fmt.Println("💡 Run 'viki workflow next' to start the first step")
		},
	}
}

func newWorkflowStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current workflow progress",
		Run: func(cmd *cobra.Command, args []string) {
			// Check .sdd/state.yaml for current phase
			statePath := filepath.Join(".sdd", "state.yaml")
			if _, err := os.Stat(statePath); os.IsNotExist(err) {
				fmt.Println("❌ No workflow in progress. Run 'viki workflow init' first.")
				return
			}

			titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
			completedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
			pendingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

			fmt.Println()
			fmt.Println(titleStyle.Render("📊 Workflow Progress"))
			fmt.Println("─────────────────────────────────────────────────")

			steps := []struct {
				name string
				file string
				cmd  string
			}{
				{"Initialize", ".sdd/state.yaml", "viki init"},
				{"Specify", ".sdd/spec.md", "viki specify"},
				{"Plan", ".sdd/plan.md", "viki plan"},
				{"Tasks", ".sdd/tasks.md", "viki task"},
				{"Execute", ".sdd/implementation.md", "viki execute"},
				{"Review", ".sdd/review.md", "viki review"},
			}

			for _, step := range steps {
				_, err := os.Stat(step.file)
				if err == nil {
					fmt.Printf("%s %s\n", completedStyle.Render("✅"), step.name)
				} else {
					fmt.Printf("%s %s\n", pendingStyle.Render("⬜"), step.name)
				}
			}
			fmt.Println()
		},
	}
}

func newWorkflowNextCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "next",
		Short: "Execute next workflow step",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🔄 Determining next step...")

			// Check current progress
			steps := []struct {
				name string
				file string
				cmd  string
			}{
				{"Initialize", ".sdd/state.yaml", "viki init"},
				{"Specify", ".sdd/spec.md", "viki specify"},
				{"Plan", ".sdd/plan.md", "viki plan"},
				{"Tasks", ".sdd/tasks.md", "viki task"},
				{"Execute", ".sdd/implementation.md", "viki execute"},
				{"Review", ".sdd/review.md", "viki review"},
			}

			for _, step := range steps {
				if _, err := os.Stat(step.file); os.IsNotExist(err) {
					fmt.Printf("\n➡️  Next Step: %s\n", step.name)
					fmt.Printf("   Command: %s\n\n", step.cmd)
					return
				}
			}

			fmt.Println("\n🎉 All workflow steps completed!")
		},
	}
}

func newWorkflowListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available workflows",
		Run: func(cmd *cobra.Command, args []string) {
			titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
			highlightStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
			dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

			fmt.Println()
			fmt.Println(titleStyle.Render("📋 Available Workflow Tracks"))
			fmt.Println("─────────────────────────────────────────────────")

			tracks := []struct {
				name string
				desc string
				time string
			}{
				{"Quick Flow", "Bug fixes, small features", "~5 minutes"},
				{"Standard Method", "Products and platforms", "~15 minutes"},
				{"Enterprise", "Compliance-heavy systems", "~30 minutes"},
			}

			for _, t := range tracks {
				fmt.Printf("\n%s\n", highlightStyle.Render(t.name))
				fmt.Printf("   %s\n", dimStyle.Render(t.desc))
				fmt.Printf("   Time to First Story: %s\n", t.time)
			}
			fmt.Println()
		},
	}
}

func getDatabase() (*db.DB, error) {
	homeDir, _ := os.UserHomeDir()
	cfg := db.Config{
		Path: filepath.Join(homeDir, ".viki", "viki.db"),
	}
	return db.New(cfg)
}

func formatAge(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}
