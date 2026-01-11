package tools

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// GrepTool searches file contents
type GrepTool struct {
	MaxResults int
	MaxContext int
}

// GrepMatch represents a single match
type GrepMatch struct {
	FilePath   string `json:"file_path"`
	LineNumber int    `json:"line_number"`
	Line       string `json:"line"`
	Context    struct {
		Before []string `json:"before,omitempty"`
		After  []string `json:"after,omitempty"`
	} `json:"context,omitempty"`
}

// GrepResult represents grep search results
type GrepResult struct {
	Pattern   string       `json:"pattern"`
	Path      string       `json:"path"`
	Matches   []*GrepMatch `json:"matches"`
	Total     int          `json:"total"`
	Truncated bool         `json:"truncated"`
}

// DefaultGrepTool creates a GrepTool with sensible defaults
func DefaultGrepTool() *GrepTool {
	return &GrepTool{
		MaxResults: 100,
		MaxContext: 2,
	}
}

// Search searches for a pattern in files
func (g *GrepTool) Search(pattern, path string, includePatterns []string, literalText bool) (*GrepResult, error) {
	// Check if ripgrep is available
	if _, err := exec.LookPath("rg"); err == nil {
		return g.searchWithRipgrep(pattern, path, includePatterns, literalText)
	}

	// Fallback to native Go implementation
	return g.searchNative(pattern, path, includePatterns, literalText)
}

// searchWithRipgrep uses ripgrep for fast searching
func (g *GrepTool) searchWithRipgrep(pattern, path string, includePatterns []string, literalText bool) (*GrepResult, error) {
	result := &GrepResult{
		Pattern: pattern,
		Path:    path,
	}

	args := []string{
		"-n",                             // Line numbers
		"--color=never",                  // No color codes
		"-m", strconv.Itoa(g.MaxResults), // Max count
		"-C", strconv.Itoa(g.MaxContext), // Context lines
	}

	if literalText {
		args = append(args, "-F") // Fixed string (literal)
	}

	// Add include patterns
	for _, include := range includePatterns {
		args = append(args, "-g", include)
	}

	args = append(args, pattern, path)

	cmd := exec.Command("rg", args...)
	output, err := cmd.Output()

	if err != nil {
		// Exit code 1 means no matches (not an error)
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return result, nil
		}
		// Other errors, try native search
		return g.searchNative(pattern, path, includePatterns, literalText)
	}

	// Parse ripgrep output
	lines := strings.Split(string(output), "\n")
	var currentMatch *GrepMatch

	for _, line := range lines {
		if line == "" || line == "--" {
			if currentMatch != nil {
				result.Matches = append(result.Matches, currentMatch)
				currentMatch = nil
			}
			continue
		}

		// Parse format: file:line:content or file-line-content (context)
		parts := strings.SplitN(line, ":", 3)
		if len(parts) < 3 {
			parts = strings.SplitN(line, "-", 3)
		}

		if len(parts) >= 3 {
			lineNum, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}

			match := &GrepMatch{
				FilePath:   parts[0],
				LineNumber: lineNum,
				Line:       parts[2],
			}
			result.Matches = append(result.Matches, match)
		}

		if len(result.Matches) >= g.MaxResults {
			result.Truncated = true
			break
		}
	}

	result.Total = len(result.Matches)
	return result, nil
}

// searchNative performs search using native Go
func (g *GrepTool) searchNative(pattern, path string, includePatterns []string, literalText bool) (*GrepResult, error) {
	result := &GrepResult{
		Pattern: pattern,
		Path:    path,
	}

	var re *regexp.Regexp
	var err error

	if literalText {
		re = regexp.MustCompile(regexp.QuoteMeta(pattern))
	} else {
		re, err = regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern: %w", err)
		}
	}

	err = filepath.WalkDir(path, func(filePath string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			// Skip hidden directories
			if strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Check include patterns
		if len(includePatterns) > 0 {
			matched := false
			for _, pattern := range includePatterns {
				if m, _ := filepath.Match(pattern, d.Name()); m {
					matched = true
					break
				}
			}
			if !matched {
				return nil
			}
		}

		// Search file
		matches, err := g.searchFile(filePath, re)
		if err != nil {
			return nil
		}

		result.Matches = append(result.Matches, matches...)

		if len(result.Matches) >= g.MaxResults {
			result.Truncated = true
			return filepath.SkipAll
		}

		return nil
	})

	if err != nil && err != filepath.SkipAll {
		return nil, err
	}

	result.Total = len(result.Matches)
	return result, nil
}

func (g *GrepTool) searchFile(filePath string, re *regexp.Regexp) ([]*GrepMatch, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var matches []*GrepMatch
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if re.MatchString(line) {
			match := &GrepMatch{
				FilePath:   filePath,
				LineNumber: lineNum,
				Line:       line,
			}
			matches = append(matches, match)
		}
	}

	return matches, scanner.Err()
}

// GlobTool finds files by pattern
type GlobTool struct {
	MaxResults int
}

// GlobResult represents glob search results
type GlobResult struct {
	Pattern   string   `json:"pattern"`
	Path      string   `json:"path"`
	Files     []string `json:"files"`
	Total     int      `json:"total"`
	Truncated bool     `json:"truncated"`
}

// DefaultGlobTool creates a GlobTool with sensible defaults
func DefaultGlobTool() *GlobTool {
	return &GlobTool{
		MaxResults: 200,
	}
}

// Find finds files matching a glob pattern
func (gl *GlobTool) Find(pattern, path string) (*GlobResult, error) {
	// Check if fd is available for faster search
	if _, err := exec.LookPath("fd"); err == nil {
		return gl.findWithFd(pattern, path)
	}

	// Fallback to native implementation
	return gl.findNative(pattern, path)
}

func (gl *GlobTool) findWithFd(pattern, path string) (*GlobResult, error) {
	result := &GlobResult{
		Pattern: pattern,
		Path:    path,
	}

	args := []string{
		"-g", pattern,
		"--color=never",
		"-L", // Follow symlinks
		"-a", // Absolute paths
	}

	if path != "" && path != "." {
		args = append(args, path)
	}

	cmd := exec.Command("fd", args...)
	output, err := cmd.Output()
	if err != nil {
		return gl.findNative(pattern, path)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		result.Files = append(result.Files, line)
		if len(result.Files) >= gl.MaxResults {
			result.Truncated = true
			break
		}
	}

	result.Total = len(result.Files)
	return result, nil
}

func (gl *GlobTool) findNative(pattern, path string) (*GlobResult, error) {
	result := &GlobResult{
		Pattern: pattern,
		Path:    path,
	}

	if path == "" {
		path = "."
	}

	err := filepath.WalkDir(path, func(filePath string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Match against pattern
		if matched, _ := filepath.Match(pattern, d.Name()); matched {
			result.Files = append(result.Files, filePath)
			if len(result.Files) >= gl.MaxResults {
				result.Truncated = true
				return filepath.SkipAll
			}
		}

		return nil
	})

	if err != nil && err != filepath.SkipAll {
		return nil, err
	}

	result.Total = len(result.Files)
	return result, nil
}
