package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// FileChange represents a change to be applied to a file
type FileChange struct {
	Path      string
	Action    string // "create", "modify", "delete"
	Content   string
	Backup    string
	Timestamp time.Time
}

// Editor handles file modifications from AI responses
type Editor struct {
	projectRoot string
	historyDir  string
	changes     []FileChange
}

// NewEditor creates a new file editor
func NewEditor(projectRoot string) *Editor {
	historyDir := filepath.Join(projectRoot, ".sdd", "history")
	os.MkdirAll(historyDir, 0755)

	return &Editor{
		projectRoot: projectRoot,
		historyDir:  historyDir,
		changes:     []FileChange{},
	}
}

// ParseCodeBlocks extracts code blocks from AI response
func (e *Editor) ParseCodeBlocks(response string) []CodeBlock {
	var blocks []CodeBlock

	// Pattern: ```language:filename or ```language filename
	pattern := regexp.MustCompile("(?s)```(\\w+)?(?::|\\s+)?([^\\n]*)?\\n(.*?)```")
	matches := pattern.FindAllStringSubmatch(response, -1)

	for _, match := range matches {
		block := CodeBlock{
			Language: match[1],
			Filename: strings.TrimSpace(match[2]),
			Content:  match[3],
		}
		blocks = append(blocks, block)
	}

	return blocks
}

// CodeBlock represents a code block from AI response
type CodeBlock struct {
	Language string
	Filename string
	Content  string
}

// ApplyChange applies a file change with backup
func (e *Editor) ApplyChange(change FileChange) error {
	fullPath := filepath.Join(e.projectRoot, change.Path)

	// Create backup if file exists
	if _, err := os.Stat(fullPath); err == nil {
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("failed to read file for backup: %w", err)
		}
		change.Backup = string(content)

		// Save to history
		if err := e.saveToHistory(change.Path, content); err != nil {
			return fmt.Errorf("failed to save backup: %w", err)
		}
	}

	switch change.Action {
	case "create", "modify":
		// Ensure directory exists
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		if err := os.WriteFile(fullPath, []byte(change.Content), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

	case "delete":
		if err := os.Remove(fullPath); err != nil {
			return fmt.Errorf("failed to delete file: %w", err)
		}
	}

	change.Timestamp = time.Now()
	e.changes = append(e.changes, change)

	return nil
}

// saveToHistory saves a file version to history
func (e *Editor) saveToHistory(path string, content []byte) error {
	timestamp := time.Now().Format("20060102_150405")
	safePath := strings.ReplaceAll(path, "/", "_")
	historyFile := filepath.Join(e.historyDir, fmt.Sprintf("%s_%s", timestamp, safePath))

	return os.WriteFile(historyFile, content, 0644)
}

// Undo reverts the last n changes
func (e *Editor) Undo(n int) error {
	if n > len(e.changes) {
		n = len(e.changes)
	}

	for i := 0; i < n; i++ {
		change := e.changes[len(e.changes)-1-i]
		fullPath := filepath.Join(e.projectRoot, change.Path)

		if change.Backup != "" {
			// Restore from backup
			if err := os.WriteFile(fullPath, []byte(change.Backup), 0644); err != nil {
				return fmt.Errorf("failed to restore %s: %w", change.Path, err)
			}
		} else if change.Action == "create" {
			// Remove created file
			os.Remove(fullPath)
		}
	}

	// Remove undone changes from history
	e.changes = e.changes[:len(e.changes)-n]

	return nil
}

// GetHistory returns the list of history files
func (e *Editor) GetHistory() ([]HistoryEntry, error) {
	entries, err := os.ReadDir(e.historyDir)
	if err != nil {
		return nil, err
	}

	var history []HistoryEntry
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		history = append(history, HistoryEntry{
			Filename:  entry.Name(),
			Timestamp: info.ModTime(),
			Size:      info.Size(),
		})
	}

	return history, nil
}

// HistoryEntry represents a file in history
type HistoryEntry struct {
	Filename  string
	Timestamp time.Time
	Size      int64
}

// RestoreFromHistory restores a file from history
func (e *Editor) RestoreFromHistory(historyFile, targetPath string) error {
	historyPath := filepath.Join(e.historyDir, historyFile)
	content, err := os.ReadFile(historyPath)
	if err != nil {
		return fmt.Errorf("failed to read history file: %w", err)
	}

	fullPath := filepath.Join(e.projectRoot, targetPath)
	return os.WriteFile(fullPath, content, 0644)
}

// ApplyDiff applies a unified diff to a file
func (e *Editor) ApplyDiff(path string, diff string) error {
	fullPath := filepath.Join(e.projectRoot, path)

	// Read current file
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse and apply diff
	newContent := applyUnifiedDiff(string(content), diff)

	// Create change
	change := FileChange{
		Path:    path,
		Action:  "modify",
		Content: newContent,
	}

	return e.ApplyChange(change)
}

// applyUnifiedDiff applies a unified diff to content
func applyUnifiedDiff(content, diff string) string {
	lines := strings.Split(content, "\n")
	diffLines := strings.Split(diff, "\n")

	var result []string
	lineIdx := 0

	for _, dline := range diffLines {
		if strings.HasPrefix(dline, "@@") {
			// Parse hunk header
			continue
		} else if strings.HasPrefix(dline, "-") && !strings.HasPrefix(dline, "---") {
			// Remove line - skip current line
			lineIdx++
		} else if strings.HasPrefix(dline, "+") && !strings.HasPrefix(dline, "+++") {
			// Add line
			result = append(result, dline[1:])
		} else if strings.HasPrefix(dline, " ") {
			// Context line - keep
			if lineIdx < len(lines) {
				result = append(result, lines[lineIdx])
				lineIdx++
			}
		}
	}

	// Add remaining lines
	for ; lineIdx < len(lines); lineIdx++ {
		result = append(result, lines[lineIdx])
	}

	return strings.Join(result, "\n")
}

// CreateFromBlocks creates files from parsed code blocks
func (e *Editor) CreateFromBlocks(blocks []CodeBlock) ([]string, error) {
	var created []string

	for _, block := range blocks {
		if block.Filename == "" {
			continue
		}

		change := FileChange{
			Path:    block.Filename,
			Action:  "create",
			Content: block.Content,
		}

		if err := e.ApplyChange(change); err != nil {
			return created, fmt.Errorf("failed to create %s: %w", block.Filename, err)
		}

		created = append(created, block.Filename)
	}

	return created, nil
}

// GetChanges returns all changes made in this session
func (e *Editor) GetChanges() []FileChange {
	return e.changes
}
