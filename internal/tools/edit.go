package tools

import (
	"fmt"
	"os"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// EditTool provides file editing capabilities
type EditTool struct{}

// EditOperation represents a single edit operation
type EditOperation struct {
	StartLine   int    `json:"start_line"`
	EndLine     int    `json:"end_line"`
	OldContent  string `json:"old_content"`
	NewContent  string `json:"new_content"`
	Description string `json:"description,omitempty"`
}

// EditResult represents the result of an edit
type EditResult struct {
	FilePath      string `json:"file_path"`
	Success       bool   `json:"success"`
	LinesChanged  int    `json:"lines_changed"`
	BeforeContent string `json:"before_content,omitempty"`
	AfterContent  string `json:"after_content,omitempty"`
	Diff          string `json:"diff"`
}

// Edit edits a file by replacing content between lines or by search/replace
func (e *EditTool) Edit(filePath string, op *EditOperation) (*EditResult, error) {
	result := &EditResult{
		FilePath: filePath,
	}

	// Read original content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	result.BeforeContent = string(content)
	lines := strings.Split(string(content), "\n")

	var newContent string

	if op.StartLine > 0 && op.EndLine > 0 {
		// Line-range based edit
		if op.StartLine > len(lines) || op.EndLine > len(lines) {
			return nil, fmt.Errorf("line range out of bounds: %d-%d (file has %d lines)", op.StartLine, op.EndLine, len(lines))
		}

		// Verify old content matches if specified
		if op.OldContent != "" {
			actualContent := strings.Join(lines[op.StartLine-1:op.EndLine], "\n")
			if strings.TrimSpace(actualContent) != strings.TrimSpace(op.OldContent) {
				return nil, fmt.Errorf("old content does not match at lines %d-%d", op.StartLine, op.EndLine)
			}
		}

		// Replace lines
		newLines := append(
			lines[:op.StartLine-1],
			strings.Split(op.NewContent, "\n")...,
		)
		newLines = append(newLines, lines[op.EndLine:]...)
		newContent = strings.Join(newLines, "\n")

		result.LinesChanged = op.EndLine - op.StartLine + 1
	} else if op.OldContent != "" {
		// Search/replace based edit
		count := strings.Count(string(content), op.OldContent)
		if count == 0 {
			return nil, fmt.Errorf("old content not found in file")
		}
		if count > 1 {
			return nil, fmt.Errorf("old content found %d times (expected exactly 1)", count)
		}

		newContent = strings.Replace(string(content), op.OldContent, op.NewContent, 1)
		result.LinesChanged = len(strings.Split(op.NewContent, "\n"))
	} else {
		return nil, fmt.Errorf("either line range or old_content must be specified")
	}

	result.AfterContent = newContent

	// Generate diff
	result.Diff = generateDiff(result.BeforeContent, newContent)

	// Write new content
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	result.Success = true
	return result, nil
}

// ReplaceAll replaces all occurrences of a pattern
func (e *EditTool) ReplaceAll(filePath, oldContent, newContent string) (*EditResult, error) {
	result := &EditResult{
		FilePath: filePath,
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	result.BeforeContent = string(content)
	count := strings.Count(string(content), oldContent)

	if count == 0 {
		return nil, fmt.Errorf("pattern not found in file")
	}

	newContentStr := strings.ReplaceAll(string(content), oldContent, newContent)
	result.AfterContent = newContentStr
	result.LinesChanged = count
	result.Diff = generateDiff(result.BeforeContent, newContentStr)

	if err := os.WriteFile(filePath, []byte(newContentStr), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	result.Success = true
	return result, nil
}

// generateDiff creates a unified diff between two strings
func generateDiff(before, after string) string {
	dmp := diffmatchpatch.New()

	// Compute line-mode diff for efficiency
	a, b, lineArray := dmp.DiffLinesToChars(before, after)
	diffs := dmp.DiffMain(a, b, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	return dmp.DiffPrettyText(diffs)
}

// PatchTool applies unified diffs/patches
type PatchTool struct{}

// PatchResult represents the result of applying a patch
type PatchResult struct {
	FilePath     string `json:"file_path"`
	Success      bool   `json:"success"`
	HunksApplied int    `json:"hunks_applied"`
	HunksFailed  int    `json:"hunks_failed"`
}

// Apply applies a unified diff patch to a file
func (p *PatchTool) Apply(filePath, diffContent string) (*PatchResult, error) {
	result := &PatchResult{
		FilePath: filePath,
	}

	// Read original content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	dmp := diffmatchpatch.New()

	// Parse patches from diff content
	patches, err := dmp.PatchFromText(diffContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse patch: %w", err)
	}

	// Apply patches
	newContent, applied := dmp.PatchApply(patches, string(content))

	for _, app := range applied {
		if app {
			result.HunksApplied++
		} else {
			result.HunksFailed++
		}
	}

	if result.HunksFailed > 0 {
		return result, fmt.Errorf("%d hunks failed to apply", result.HunksFailed)
	}

	// Write new content
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	result.Success = true
	return result, nil
}

// CreatePatch creates a patch from two file contents
func (p *PatchTool) CreatePatch(fileName, before, after string) string {
	dmp := diffmatchpatch.New()
	patches := dmp.PatchMake(before, after)
	return dmp.PatchToText(patches)
}

// InsertAt inserts content at a specific line
func (e *EditTool) InsertAt(filePath string, lineNum int, content string) (*EditResult, error) {
	result := &EditResult{
		FilePath: filePath,
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	result.BeforeContent = string(fileContent)
	lines := strings.Split(string(fileContent), "\n")

	if lineNum < 1 {
		lineNum = 1
	}
	if lineNum > len(lines)+1 {
		lineNum = len(lines) + 1
	}

	// Insert content
	insertLines := strings.Split(content, "\n")
	newLines := append(
		lines[:lineNum-1],
		append(insertLines, lines[lineNum-1:]...)...,
	)

	newContent := strings.Join(newLines, "\n")
	result.AfterContent = newContent
	result.LinesChanged = len(insertLines)
	result.Diff = generateDiff(result.BeforeContent, newContent)

	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	result.Success = true
	return result, nil
}

// DeleteLines deletes lines from a file
func (e *EditTool) DeleteLines(filePath string, startLine, endLine int) (*EditResult, error) {
	result := &EditResult{
		FilePath: filePath,
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	result.BeforeContent = string(content)
	lines := strings.Split(string(content), "\n")

	if startLine < 1 || endLine < 1 || startLine > len(lines) || endLine > len(lines) {
		return nil, fmt.Errorf("line range out of bounds")
	}

	newLines := append(lines[:startLine-1], lines[endLine:]...)
	newContent := strings.Join(newLines, "\n")

	result.AfterContent = newContent
	result.LinesChanged = endLine - startLine + 1
	result.Diff = generateDiff(result.BeforeContent, newContent)

	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	result.Success = true
	return result, nil
}
