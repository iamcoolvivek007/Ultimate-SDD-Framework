package tools

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// BashTool executes shell commands
type BashTool struct {
	Shell     string
	ShellArgs []string
	Timeout   time.Duration
	WorkDir   string
}

// BashResult represents the result of a bash command execution
type BashResult struct {
	Command    string        `json:"command"`
	ExitCode   int           `json:"exit_code"`
	Stdout     string        `json:"stdout"`
	Stderr     string        `json:"stderr"`
	Duration   time.Duration `json:"duration"`
	TimedOut   bool          `json:"timed_out"`
	Terminated bool          `json:"terminated"`
}

// DefaultBashTool creates a BashTool with sensible defaults
func DefaultBashTool() *BashTool {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	return &BashTool{
		Shell:     shell,
		ShellArgs: []string{"-c"},
		Timeout:   60 * time.Second,
		WorkDir:   ".",
	}
}

// Execute runs a bash command and returns the result
func (b *BashTool) Execute(ctx context.Context, command string) (*BashResult, error) {
	result := &BashResult{
		Command: command,
	}

	start := time.Now()

	// Create context with timeout if specified
	if b.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, b.Timeout)
		defer cancel()
	}

	// Build the command
	args := append(b.ShellArgs, command)
	cmd := exec.CommandContext(ctx, b.Shell, args...)

	// Set working directory
	if b.WorkDir != "" {
		cmd.Dir = b.WorkDir
	}

	// Set up process group for proper cleanup
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set environment
	cmd.Env = os.Environ()

	// Run the command
	err := cmd.Run()

	result.Duration = time.Since(start)
	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	// Handle exit status
	if ctx.Err() == context.DeadlineExceeded {
		result.TimedOut = true
		result.ExitCode = -1
	} else if ctx.Err() == context.Canceled {
		result.Terminated = true
		result.ExitCode = -1
	} else if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("failed to execute command: %w", err)
		}
	} else {
		result.ExitCode = 0
	}

	return result, nil
}

// ExecuteWithInput runs a command with stdin input
func (b *BashTool) ExecuteWithInput(ctx context.Context, command, input string) (*BashResult, error) {
	result := &BashResult{
		Command: command,
	}

	start := time.Now()

	if b.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, b.Timeout)
		defer cancel()
	}

	args := append(b.ShellArgs, command)
	cmd := exec.CommandContext(ctx, b.Shell, args...)

	if b.WorkDir != "" {
		cmd.Dir = b.WorkDir
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = strings.NewReader(input)
	cmd.Env = os.Environ()

	err := cmd.Run()

	result.Duration = time.Since(start)
	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	if ctx.Err() == context.DeadlineExceeded {
		result.TimedOut = true
		result.ExitCode = -1
	} else if ctx.Err() == context.Canceled {
		result.Terminated = true
		result.ExitCode = -1
	} else if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("failed to execute command: %w", err)
		}
	}

	return result, nil
}

// IsCommandSafe performs basic safety check on commands
// Returns false for potentially dangerous commands
func IsCommandSafe(command string) (bool, string) {
	dangerous := []struct {
		pattern string
		reason  string
	}{
		{"rm -rf /", "Attempts to delete root filesystem"},
		{"rm -rf /*", "Attempts to delete root filesystem"},
		{":(){:|:&};:", "Fork bomb detected"},
		{"dd if=", "Direct disk write detected"},
		{"mkfs", "Filesystem format command detected"},
		{"> /dev/sd", "Direct device write detected"},
		{"chmod -R 777 /", "Dangerous permission change"},
		{"wget", "Network download - review URL"},
		{"curl", "Network request - review URL"},
	}

	lowerCmd := strings.ToLower(command)
	for _, d := range dangerous {
		if strings.Contains(lowerCmd, strings.ToLower(d.pattern)) {
			return false, d.reason
		}
	}

	return true, ""
}

// FormatResult formats the bash result for display
func (r *BashResult) FormatResult() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Command: %s\n", r.Command))
	sb.WriteString(fmt.Sprintf("Exit Code: %d\n", r.ExitCode))
	sb.WriteString(fmt.Sprintf("Duration: %v\n", r.Duration))

	if r.TimedOut {
		sb.WriteString("Status: TIMED OUT\n")
	} else if r.Terminated {
		sb.WriteString("Status: TERMINATED\n")
	}

	if r.Stdout != "" {
		sb.WriteString("\n--- STDOUT ---\n")
		sb.WriteString(r.Stdout)
	}

	if r.Stderr != "" {
		sb.WriteString("\n--- STDERR ---\n")
		sb.WriteString(r.Stderr)
	}

	return sb.String()
}

// Success returns true if the command succeeded
func (r *BashResult) Success() bool {
	return r.ExitCode == 0 && !r.TimedOut && !r.Terminated
}
