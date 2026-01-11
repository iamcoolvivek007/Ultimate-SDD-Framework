package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Command represents a custom slash command
type Command struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Content     string            `json:"content"`
	Scope       string            `json:"scope"` // "user", "project", "builtin"
	FilePath    string            `json:"file_path,omitempty"`
	Arguments   []string          `json:"arguments,omitempty"` // Extracted $VARIABLE placeholders
	Metadata    map[string]string `json:"metadata,omitempty"`  // Parsed YAML frontmatter
}

// CommandLoader loads custom commands from filesystem
type CommandLoader struct {
	UserDir    string // ~/.viki/commands/
	ProjectDir string // .viki/commands/
}

// NewCommandLoader creates a command loader with default paths
func NewCommandLoader(projectDir string) *CommandLoader {
	homeDir, _ := os.UserHomeDir()
	return &CommandLoader{
		UserDir:    filepath.Join(homeDir, ".viki", "commands"),
		ProjectDir: filepath.Join(projectDir, ".viki", "commands"),
	}
}

// LoadAll loads all commands from user and project directories
func (l *CommandLoader) LoadAll() ([]*Command, error) {
	var commands []*Command

	// Load user commands
	userCmds, err := l.loadFromDir(l.UserDir, "user")
	if err == nil {
		commands = append(commands, userCmds...)
	}

	// Load project commands
	projectCmds, err := l.loadFromDir(l.ProjectDir, "project")
	if err == nil {
		commands = append(commands, projectCmds...)
	}

	// Add built-in commands
	commands = append(commands, getBuiltinCommands()...)

	return commands, nil
}

// loadFromDir loads commands from a directory
func (l *CommandLoader) loadFromDir(dir, scope string) ([]*Command, error) {
	var commands []*Command

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		cmd, err := l.loadCommand(path, scope)
		if err != nil {
			return nil // Skip invalid commands
		}

		commands = append(commands, cmd)
		return nil
	})

	return commands, err
}

// loadCommand loads a single command from a file
func (l *CommandLoader) loadCommand(path, scope string) (*Command, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Generate command ID from path
	baseDir := l.UserDir
	if scope == "project" {
		baseDir = l.ProjectDir
	}

	relPath, _ := filepath.Rel(baseDir, path)
	id := strings.TrimSuffix(relPath, ".md")
	id = strings.ReplaceAll(id, string(filepath.Separator), "/")
	id = scope + "/" + id

	cmd := &Command{
		ID:       id,
		Name:     filepath.Base(id),
		Scope:    scope,
		FilePath: path,
		Metadata: make(map[string]string),
	}

	// Parse content
	cmd.Content, cmd.Metadata = parseCommandFile(string(content))

	// Extract description from metadata
	if desc, ok := cmd.Metadata["description"]; ok {
		cmd.Description = desc
	}

	// Extract argument placeholders ($VARIABLE_NAME)
	cmd.Arguments = extractArguments(cmd.Content)

	return cmd, nil
}

// parseCommandFile parses YAML frontmatter and content
func parseCommandFile(content string) (string, map[string]string) {
	metadata := make(map[string]string)

	// Check for YAML frontmatter
	if !strings.HasPrefix(content, "---") {
		return content, metadata
	}

	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return content, metadata
	}

	// Parse simple YAML metadata
	lines := strings.Split(parts[1], "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 {
			continue
		}

		key := strings.TrimSpace(line[:colonIdx])
		value := strings.TrimSpace(line[colonIdx+1:])
		value = strings.Trim(value, "\"'")
		metadata[key] = value
	}

	return strings.TrimSpace(parts[2]), metadata
}

// extractArguments extracts $VARIABLE_NAME placeholders from content
func extractArguments(content string) []string {
	re := regexp.MustCompile(`\$([A-Z][A-Z0-9_]*)`)
	matches := re.FindAllStringSubmatch(content, -1)

	seen := make(map[string]bool)
	var args []string

	for _, match := range matches {
		if len(match) >= 2 && !seen[match[1]] {
			seen[match[1]] = true
			args = append(args, match[1])
		}
	}

	return args
}

// getBuiltinCommands returns the built-in commands
func getBuiltinCommands() []*Command {
	return []*Command{
		{
			ID:          "viki/help",
			Name:        "help",
			Description: "Show available commands",
			Content:     "List all available viki commands with descriptions.",
			Scope:       "builtin",
		},
		{
			ID:          "viki/constitution",
			Name:        "constitution",
			Description: "Create or update project constitution",
			Content: `Create or update the project constitution at .viki/constitution.md.

## User Input
$ARGUMENTS

If no input provided, interactively gather project principles.

The constitution should include:
1. Project name and description
2. Core principles (3-5 non-negotiables)
3. Coding standards
4. Quality requirements
5. Governance rules`,
			Scope:     "builtin",
			Arguments: []string{"ARGUMENTS"},
		},
		{
			ID:          "viki/clarify",
			Name:        "clarify",
			Description: "Clarify specifications with structured Q&A",
			Content: `Review the current specification and identify areas that need clarification.

Generate a structured set of questions to fill gaps in:
- Requirements completeness
- Edge cases
- Technical constraints
- User flows
- Error handling

For each question, explain why it needs clarification and suggest possible answers.`,
			Scope: "builtin",
		},
		{
			ID:          "viki/checklist",
			Name:        "checklist",
			Description: "Generate quality checklist",
			Content: `Generate a quality checklist for the current specification/plan.

Include checks for:
- [ ] Requirements completeness
- [ ] Acceptance criteria defined
- [ ] Edge cases identified
- [ ] Error handling planned
- [ ] Security considerations
- [ ] Performance requirements
- [ ] Testing strategy
- [ ] Documentation needs

Mark as [X] if already addressed, [ ] if needs attention.`,
			Scope: "builtin",
		},
		{
			ID:          "viki/analyze",
			Name:        "analyze",
			Description: "Cross-artifact consistency analysis",
			Content: `Perform cross-artifact consistency analysis.

Compare:
1. spec.md ↔ plan.md alignment
2. plan.md ↔ tasks.md coverage
3. Tasks ↔ Implementation status

Report any:
- Inconsistencies
- Missing coverage
- Orphaned items
- Scope creep`,
			Scope: "builtin",
		},
		{
			ID:          "viki/compact",
			Name:        "compact",
			Description: "Summarize and compact current session",
			Content: `Summarize the current conversation and create a new session with the summary.

Include in summary:
- Key decisions made
- Current progress
- Open items
- Important context`,
			Scope: "builtin",
		},
		{
			ID:          "viki/context",
			Name:        "context",
			Description: "Load project context",
			Content: `Load and analyze project context.

Read:
- .viki/constitution.md
- .sdd/CONTEXT.md
- README.md
- Key source files

Provide a summary of the project state.`,
			Scope: "builtin",
		},
	}
}

// Execute processes a command, replacing arguments with provided values
func (c *Command) Execute(args map[string]string) string {
	content := c.Content

	// Replace argument placeholders
	for argName, argValue := range args {
		placeholder := "$" + argName
		content = strings.ReplaceAll(content, placeholder, argValue)
	}

	return content
}

// HasRequiredArgs checks if all required arguments are provided
func (c *Command) HasRequiredArgs(args map[string]string) (bool, []string) {
	var missing []string

	for _, arg := range c.Arguments {
		if _, ok := args[arg]; !ok {
			missing = append(missing, arg)
		}
	}

	return len(missing) == 0, missing
}

// CommandRegistry provides access to all commands
type CommandRegistry struct {
	commands map[string]*Command
	loader   *CommandLoader
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry(projectDir string) *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]*Command),
		loader:   NewCommandLoader(projectDir),
	}
}

// Load loads all commands
func (r *CommandRegistry) Load() error {
	cmds, err := r.loader.LoadAll()
	if err != nil {
		return err
	}

	for _, cmd := range cmds {
		r.commands[cmd.ID] = cmd
	}

	return nil
}

// Get retrieves a command by ID or name
func (r *CommandRegistry) Get(id string) (*Command, bool) {
	// Try exact match first
	if cmd, ok := r.commands[id]; ok {
		return cmd, true
	}

	// Try with viki/ prefix
	if cmd, ok := r.commands["viki/"+id]; ok {
		return cmd, true
	}

	// Try with user/ prefix
	if cmd, ok := r.commands["user/"+id]; ok {
		return cmd, true
	}

	// Try with project/ prefix
	if cmd, ok := r.commands["project/"+id]; ok {
		return cmd, true
	}

	return nil, false
}

// List returns all commands
func (r *CommandRegistry) List() []*Command {
	cmds := make([]*Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// ListByScope returns commands filtered by scope
func (r *CommandRegistry) ListByScope(scope string) []*Command {
	var cmds []*Command
	for _, cmd := range r.commands {
		if cmd.Scope == scope {
			cmds = append(cmds, cmd)
		}
	}
	return cmds
}

// FormatCommandList formats commands for display
func FormatCommandList(commands []*Command) string {
	var sb strings.Builder

	// Group by scope
	groups := map[string][]*Command{
		"builtin": {},
		"user":    {},
		"project": {},
	}

	for _, cmd := range commands {
		groups[cmd.Scope] = append(groups[cmd.Scope], cmd)
	}

	// Format each group
	for _, scope := range []string{"builtin", "user", "project"} {
		cmds := groups[scope]
		if len(cmds) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("\n### %s Commands\n", capitalize(scope)))
		for _, cmd := range cmds {
			sb.WriteString(fmt.Sprintf("  /%s - %s\n", cmd.ID, cmd.Description))
		}
	}

	return sb.String()
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
