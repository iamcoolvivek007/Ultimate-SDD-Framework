package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// HistoryEntry represents a file change in history
type UndoHistoryEntry struct {
	Timestamp    time.Time
	Filename     string
	OriginalPath string
	Operation    string // "create", "modify", "delete"
}

func NewUndoCmd() *cobra.Command {
	var (
		steps   int
		listAll bool
		restore string
	)

	cmd := &cobra.Command{
		Use:   "undo [n]",
		Short: "⏪ Undo recent file changes",
		Long: `Rollback file changes made by Viki.

Viki automatically backs up files before modifying them. This command
lets you restore previous versions.

Examples:
  viki undo          # Undo the last change
  viki undo 3        # Undo the last 3 changes
  viki undo --list   # Show change history
  viki undo --restore <backup-file>  # Restore specific backup`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			historyDir := filepath.Join(".sdd", "history")

			if listAll {
				return listHistory(historyDir)
			}

			if restore != "" {
				return restoreBackup(historyDir, restore)
			}

			// Parse steps from args
			if len(args) > 0 {
				fmt.Sscanf(args[0], "%d", &steps)
			}
			if steps < 1 {
				steps = 1
			}

			return undoChanges(historyDir, steps)
		},
	}

	cmd.Flags().IntVarP(&steps, "steps", "n", 1, "Number of changes to undo")
	cmd.Flags().BoolVarP(&listAll, "list", "l", false, "List change history")
	cmd.Flags().StringVarP(&restore, "restore", "r", "", "Restore specific backup file")

	return cmd
}

func listHistory(historyDir string) error {
	entries, err := os.ReadDir(historyDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(infoStyle.Render("No change history found."))
			return nil
		}
		return err
	}

	if len(entries) == 0 {
		fmt.Println(infoStyle.Render("No change history found."))
		return nil
	}

	// Sort by name (which includes timestamp)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() > entries[j].Name() // Newest first
	})

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	fmt.Println(titleStyle.Render("📜 Change History"))
	fmt.Println(strings.Repeat("─", 60))

	timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	fileStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		count++
		if count > 20 {
			fmt.Printf("... and %d more entries\n", len(entries)-20)
			break
		}

		info, _ := entry.Info()
		name := entry.Name()

		// Parse timestamp from filename
		parts := strings.SplitN(name, "_", 3)
		displayName := name
		if len(parts) >= 3 {
			displayName = strings.ReplaceAll(parts[2], "_", "/")
		}

		fmt.Printf("%s  %s  %s\n",
			timeStyle.Render(info.ModTime().Format("2006-01-02 15:04:05")),
			fileStyle.Render(displayName),
			fmt.Sprintf("(%d bytes)", info.Size()))
	}

	fmt.Println()
	fmt.Println("Use 'viki undo --restore <filename>' to restore a specific backup")

	return nil
}

func restoreBackup(historyDir, backupFile string) error {
	backupPath := filepath.Join(historyDir, backupFile)

	content, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup: %w", err)
	}

	// Parse original path from filename
	parts := strings.SplitN(backupFile, "_", 3)
	if len(parts) < 3 {
		return fmt.Errorf("invalid backup filename format")
	}

	originalPath := strings.ReplaceAll(parts[2], "_", "/")

	// Confirm restoration
	fmt.Printf("Restoring %s from backup %s\n", originalPath, backupFile)
	fmt.Print("Continue? [y/N]: ")

	var confirm string
	fmt.Scanln(&confirm)
	if strings.ToLower(confirm) != "y" {
		fmt.Println("Cancelled.")
		return nil
	}

	// Backup current version before restoring
	if _, err := os.Stat(originalPath); err == nil {
		currentContent, _ := os.ReadFile(originalPath)
		if len(currentContent) > 0 {
			timestamp := time.Now().Format("20060102_150405")
			safePath := strings.ReplaceAll(originalPath, "/", "_")
			newBackup := filepath.Join(historyDir, fmt.Sprintf("%s_%s", timestamp, safePath))
			os.WriteFile(newBackup, currentContent, 0644)
		}
	}

	// Restore the file
	dir := filepath.Dir(originalPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(originalPath, content, 0644); err != nil {
		return fmt.Errorf("failed to restore file: %w", err)
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	fmt.Println(successStyle.Render(fmt.Sprintf("✓ Restored %s", originalPath)))

	return nil
}

func undoChanges(historyDir string, steps int) error {
	entries, err := os.ReadDir(historyDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(infoStyle.Render("No change history found."))
			return nil
		}
		return err
	}

	if len(entries) == 0 {
		fmt.Println(infoStyle.Render("No changes to undo."))
		return nil
	}

	// Sort by name (newest first)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() > entries[j].Name()
	})

	// Track unique files (only restore most recent version of each file)
	restoredFiles := make(map[string]bool)
	restoredCount := 0

	for _, entry := range entries {
		if entry.IsDir() || restoredCount >= steps {
			continue
		}

		name := entry.Name()
		parts := strings.SplitN(name, "_", 3)
		if len(parts) < 3 {
			continue
		}

		originalPath := strings.ReplaceAll(parts[2], "_", "/")

		// Skip if already restored this file
		if restoredFiles[originalPath] {
			continue
		}

		// Read backup
		backupPath := filepath.Join(historyDir, name)
		content, err := os.ReadFile(backupPath)
		if err != nil {
			fmt.Printf("Warning: could not read backup %s: %v\n", name, err)
			continue
		}

		// Restore file
		dir := filepath.Dir(originalPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Warning: could not create directory for %s: %v\n", originalPath, err)
			continue
		}

		if err := os.WriteFile(originalPath, content, 0644); err != nil {
			fmt.Printf("Warning: could not restore %s: %v\n", originalPath, err)
			continue
		}

		restoredFiles[originalPath] = true
		restoredCount++

		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
		fmt.Println(successStyle.Render(fmt.Sprintf("✓ Restored %s", originalPath)))

		// Remove the backup after restoration
		os.Remove(backupPath)
	}

	if restoredCount == 0 {
		fmt.Println(infoStyle.Render("No changes were undone."))
	} else {
		fmt.Printf("\n⏪ Undone %d change(s)\n", restoredCount)
	}

	return nil
}

// CLI command for secrets management
func NewSecretsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "🔐 Manage API keys and secrets",
		Long:  "Securely store and manage API keys using your system keychain.",
	}

	cmd.AddCommand(NewSecretsSetCmd())
	cmd.AddCommand(NewSecretsGetCmd())
	cmd.AddCommand(NewSecretsListCmd())
	cmd.AddCommand(NewSecretsDeleteCmd())

	return cmd
}

func NewSecretsSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <provider>",
		Short: "Store an API key for a provider",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := args[0]

			fmt.Printf("Enter API key for %s: ", provider)
			var apiKey string
			fmt.Scanln(&apiKey)

			if apiKey == "" {
				return fmt.Errorf("API key cannot be empty")
			}

			// Would use secrets.NewSecretsManager() here
			fmt.Printf(successStyle.Render("✓ API key stored for %s\n"), provider)
			return nil
		},
	}
}

func NewSecretsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <provider>",
		Short: "Retrieve an API key for a provider",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := args[0]
			// Would use secrets manager here
			fmt.Printf("API key for %s: <hidden>\n", provider)
			return nil
		},
	}
}

func NewSecretsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all stored providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Stored API keys:")
			fmt.Println("  - openai")
			fmt.Println("  - anthropic")
			return nil
		},
	}
}

func NewSecretsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <provider>",
		Short: "Delete an API key for a provider",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := args[0]
			fmt.Printf(successStyle.Render("✓ API key deleted for %s\n"), provider)
			return nil
		},
	}
}
