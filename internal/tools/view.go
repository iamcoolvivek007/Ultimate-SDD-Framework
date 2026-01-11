package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ViewTool reads and displays file contents
type ViewTool struct {
	MaxLines    int
	MaxFileSize int64
}

// ViewResult represents the result of viewing a file
type ViewResult struct {
	FilePath    string   `json:"file_path"`
	Content     string   `json:"content"`
	Lines       []string `json:"lines,omitempty"`
	TotalLines  int      `json:"total_lines"`
	StartLine   int      `json:"start_line"`
	EndLine     int      `json:"end_line"`
	Truncated   bool     `json:"truncated"`
	FileSize    int64    `json:"file_size"`
	FileExists  bool     `json:"file_exists"`
	IsDirectory bool     `json:"is_directory"`
}

// DefaultViewTool creates a ViewTool with sensible defaults
func DefaultViewTool() *ViewTool {
	return &ViewTool{
		MaxLines:    500,
		MaxFileSize: 1024 * 1024, // 1MB
	}
}

// View reads a file and returns its contents
func (v *ViewTool) View(filePath string, offset, limit int) (*ViewResult, error) {
	result := &ViewResult{
		FilePath:  filePath,
		StartLine: offset + 1, // Convert to 1-indexed
	}

	// Check if file exists
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		result.FileExists = false
		return result, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	result.FileExists = true
	result.FileSize = info.Size()
	result.IsDirectory = info.IsDir()

	if info.IsDir() {
		return result, nil
	}

	// Check file size limit
	if info.Size() > v.MaxFileSize {
		return nil, fmt.Errorf("file too large: %d bytes (max: %d)", info.Size(), v.MaxFileSize)
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	result.TotalLines = len(lines)

	// Apply offset and limit
	if offset < 0 {
		offset = 0
	}
	if offset >= len(lines) {
		offset = len(lines) - 1
		if offset < 0 {
			offset = 0
		}
	}

	if limit <= 0 {
		limit = v.MaxLines
	}

	endOffset := offset + limit
	if endOffset > len(lines) {
		endOffset = len(lines)
	}

	result.StartLine = offset + 1
	result.EndLine = endOffset
	result.Lines = lines[offset:endOffset]
	result.Content = strings.Join(result.Lines, "\n")
	result.Truncated = endOffset < len(lines)

	return result, nil
}

// ViewRange reads specific line range from a file (1-indexed)
func (v *ViewTool) ViewRange(filePath string, startLine, endLine int) (*ViewResult, error) {
	if startLine < 1 {
		startLine = 1
	}
	offset := startLine - 1
	limit := endLine - startLine + 1
	if limit < 1 {
		limit = 1
	}
	return v.View(filePath, offset, limit)
}

// FormatWithLineNumbers formats content with line numbers
func FormatWithLineNumbers(lines []string, startLine int) string {
	var sb strings.Builder
	width := len(fmt.Sprintf("%d", startLine+len(lines)))

	for i, line := range lines {
		lineNum := startLine + i
		sb.WriteString(fmt.Sprintf("%*d │ %s\n", width, lineNum, line))
	}

	return sb.String()
}

// WriteTool writes content to files
type WriteTool struct{}

// WriteResult represents the result of a write operation
type WriteResult struct {
	FilePath    string `json:"file_path"`
	BytesWriten int    `json:"bytes_written"`
	Created     bool   `json:"created"`
	Overwritten bool   `json:"overwritten"`
}

// Write writes content to a file, creating directories if needed
func (w *WriteTool) Write(filePath, content string, overwrite bool) (*WriteResult, error) {
	result := &WriteResult{
		FilePath: filePath,
	}

	// Check if file exists
	_, err := os.Stat(filePath)
	if err == nil {
		if !overwrite {
			return nil, fmt.Errorf("file already exists: %s (use overwrite=true to replace)", filePath)
		}
		result.Overwritten = true
	} else if os.IsNotExist(err) {
		result.Created = true
	}

	// Create directory if needed
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	result.BytesWriten = len(content)
	return result, nil
}

// Append appends content to a file
func (w *WriteTool) Append(filePath, content string) (*WriteResult, error) {
	result := &WriteResult{
		FilePath: filePath,
	}

	// Check if file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return w.Write(filePath, content, false)
	}

	// Open file for appending
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	n, err := f.WriteString(content)
	if err != nil {
		return nil, fmt.Errorf("failed to append to file: %w", err)
	}

	result.BytesWriten = n
	return result, nil
}

// LsTool lists directory contents
type LsTool struct {
	ShowHidden bool
	MaxEntries int
}

// FileEntry represents a file or directory entry
type FileEntry struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	IsDir       bool   `json:"is_dir"`
	Size        int64  `json:"size"`
	Permissions string `json:"permissions"`
	ModTime     string `json:"mod_time"`
}

// LsResult represents the result of listing a directory
type LsResult struct {
	Path      string       `json:"path"`
	Entries   []*FileEntry `json:"entries"`
	Total     int          `json:"total"`
	Truncated bool         `json:"truncated"`
}

// DefaultLsTool creates an LsTool with sensible defaults
func DefaultLsTool() *LsTool {
	return &LsTool{
		ShowHidden: false,
		MaxEntries: 100,
	}
}

// List lists directory contents
func (l *LsTool) List(dirPath string, ignorePatterns []string) (*LsResult, error) {
	result := &LsResult{
		Path: dirPath,
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files if not showing them
		if !l.ShowHidden && strings.HasPrefix(name, ".") {
			continue
		}

		// Check ignore patterns
		ignored := false
		for _, pattern := range ignorePatterns {
			if matched, _ := filepath.Match(pattern, name); matched {
				ignored = true
				break
			}
		}
		if ignored {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileEntry := &FileEntry{
			Name:        name,
			Path:        filepath.Join(dirPath, name),
			IsDir:       entry.IsDir(),
			Size:        info.Size(),
			Permissions: info.Mode().String(),
			ModTime:     info.ModTime().Format("2006-01-02 15:04:05"),
		}

		result.Entries = append(result.Entries, fileEntry)

		if len(result.Entries) >= l.MaxEntries {
			result.Truncated = true
			break
		}
	}

	result.Total = len(result.Entries)
	return result, nil
}

// ListRecursive lists directory contents recursively
func (l *LsTool) ListRecursive(dirPath string, maxDepth int, ignorePatterns []string) (*LsResult, error) {
	result := &LsResult{
		Path: dirPath,
	}

	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Check depth
		rel, _ := filepath.Rel(dirPath, path)
		depth := len(strings.Split(rel, string(filepath.Separator)))
		if maxDepth > 0 && depth > maxDepth {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		name := d.Name()

		// Skip hidden files
		if !l.ShowHidden && strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Check ignore patterns
		for _, pattern := range ignorePatterns {
			if matched, _ := filepath.Match(pattern, name); matched {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		fileEntry := &FileEntry{
			Name:        name,
			Path:        path,
			IsDir:       d.IsDir(),
			Size:        info.Size(),
			Permissions: info.Mode().String(),
			ModTime:     info.ModTime().Format("2006-01-02 15:04:05"),
		}

		result.Entries = append(result.Entries, fileEntry)

		if len(result.Entries) >= l.MaxEntries {
			result.Truncated = true
			return filepath.SkipAll
		}

		return nil
	})

	if err != nil && err != filepath.SkipAll {
		return nil, err
	}

	result.Total = len(result.Entries)
	return result, nil
}
