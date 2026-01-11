package tools

import (
	"context"
)

// Tool represents an AI-callable tool
type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, args map[string]interface{}) (interface{}, error)
}

// ToolRegistry manages available tools
type ToolRegistry struct {
	tools map[string]Tool
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool to the registry
func (r *ToolRegistry) Register(tool Tool) {
	r.tools[tool.Name()] = tool
}

// Get retrieves a tool by name
func (r *ToolRegistry) Get(name string) (Tool, bool) {
	tool, ok := r.tools[name]
	return tool, ok
}

// List returns all registered tools
func (r *ToolRegistry) List() []Tool {
	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// ToolSchema describes a tool for AI consumption
type ToolSchema struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// GetSchemas returns JSON schemas for all tools
func (r *ToolRegistry) GetSchemas() []ToolSchema {
	schemas := make([]ToolSchema, 0, len(r.tools))
	for _, tool := range r.tools {
		schemas = append(schemas, ToolSchema{
			Name:        tool.Name(),
			Description: tool.Description(),
		})
	}
	return schemas
}

// DefaultRegistry creates a registry with all default tools
func DefaultRegistry(workDir string) *ToolRegistry {
	registry := NewToolRegistry()

	// Add all tools
	bashTool := DefaultBashTool()
	bashTool.WorkDir = workDir

	registry.Register(&BashToolWrapper{bash: bashTool})
	registry.Register(&ViewToolWrapper{view: DefaultViewTool()})
	registry.Register(&WriteToolWrapper{write: &WriteTool{}})
	registry.Register(&GrepToolWrapper{grep: DefaultGrepTool()})
	registry.Register(&GlobToolWrapper{glob: DefaultGlobTool()})
	registry.Register(&EditToolWrapper{edit: &EditTool{}})
	registry.Register(&PatchToolWrapper{patch: &PatchTool{}})
	registry.Register(&LsToolWrapper{ls: DefaultLsTool()})
	registry.Register(&FetchToolWrapper{fetch: DefaultFetchTool()})

	return registry
}

// Tool wrappers implementing the Tool interface

type BashToolWrapper struct {
	bash *BashTool
}

func (w *BashToolWrapper) Name() string        { return "bash" }
func (w *BashToolWrapper) Description() string { return "Execute shell commands" }
func (w *BashToolWrapper) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	command, _ := args["command"].(string)
	return w.bash.Execute(ctx, command)
}

type ViewToolWrapper struct {
	view *ViewTool
}

func (w *ViewToolWrapper) Name() string        { return "view" }
func (w *ViewToolWrapper) Description() string { return "View file contents" }
func (w *ViewToolWrapper) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	filePath, _ := args["file_path"].(string)
	offset, _ := args["offset"].(int)
	limit, _ := args["limit"].(int)
	return w.view.View(filePath, offset, limit)
}

type WriteToolWrapper struct {
	write *WriteTool
}

func (w *WriteToolWrapper) Name() string        { return "write" }
func (w *WriteToolWrapper) Description() string { return "Write content to a file" }
func (w *WriteToolWrapper) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	filePath, _ := args["file_path"].(string)
	content, _ := args["content"].(string)
	overwrite, _ := args["overwrite"].(bool)
	return w.write.Write(filePath, content, overwrite)
}

type GrepToolWrapper struct {
	grep *GrepTool
}

func (w *GrepToolWrapper) Name() string        { return "grep" }
func (w *GrepToolWrapper) Description() string { return "Search file contents for patterns" }
func (w *GrepToolWrapper) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	pattern, _ := args["pattern"].(string)
	path, _ := args["path"].(string)
	includes, _ := args["include"].([]string)
	literal, _ := args["literal_text"].(bool)
	return w.grep.Search(pattern, path, includes, literal)
}

type GlobToolWrapper struct {
	glob *GlobTool
}

func (w *GlobToolWrapper) Name() string        { return "glob" }
func (w *GlobToolWrapper) Description() string { return "Find files by glob pattern" }
func (w *GlobToolWrapper) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	pattern, _ := args["pattern"].(string)
	path, _ := args["path"].(string)
	return w.glob.Find(pattern, path)
}

type EditToolWrapper struct {
	edit *EditTool
}

func (w *EditToolWrapper) Name() string        { return "edit" }
func (w *EditToolWrapper) Description() string { return "Edit file contents" }
func (w *EditToolWrapper) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	filePath, _ := args["file_path"].(string)
	op := &EditOperation{}
	if v, ok := args["start_line"].(int); ok {
		op.StartLine = v
	}
	if v, ok := args["end_line"].(int); ok {
		op.EndLine = v
	}
	if v, ok := args["old_content"].(string); ok {
		op.OldContent = v
	}
	if v, ok := args["new_content"].(string); ok {
		op.NewContent = v
	}
	return w.edit.Edit(filePath, op)
}

type PatchToolWrapper struct {
	patch *PatchTool
}

func (w *PatchToolWrapper) Name() string        { return "patch" }
func (w *PatchToolWrapper) Description() string { return "Apply a diff patch to a file" }
func (w *PatchToolWrapper) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	filePath, _ := args["file_path"].(string)
	diff, _ := args["diff"].(string)
	return w.patch.Apply(filePath, diff)
}

type LsToolWrapper struct {
	ls *LsTool
}

func (w *LsToolWrapper) Name() string        { return "ls" }
func (w *LsToolWrapper) Description() string { return "List directory contents" }
func (w *LsToolWrapper) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	path, _ := args["path"].(string)
	if path == "" {
		path = "."
	}
	ignores, _ := args["ignore"].([]string)
	return w.ls.List(path, ignores)
}

type FetchToolWrapper struct {
	fetch *FetchTool
}

func (w *FetchToolWrapper) Name() string        { return "fetch" }
func (w *FetchToolWrapper) Description() string { return "Fetch content from a URL" }
func (w *FetchToolWrapper) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	url, _ := args["url"].(string)
	format := FetchFormat(args["format"].(string))
	if format == "" {
		format = FetchFormatText
	}
	return w.fetch.Fetch(ctx, url, format)
}
