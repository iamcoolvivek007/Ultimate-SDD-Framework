package lsp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CodebaseContext provides LSP-like context analysis for the project
type CodebaseContext struct {
	RootPath     string
	Files        []FileInfo
	Dependencies map[string][]string
	Structure    ProjectStructure
}

// FileInfo represents information about a file in the codebase
type FileInfo struct {
	Path     string
	Type     FileType
	Language string
	Content  string
	Size     int64
	Imports  []string
}

// FileType represents the type of file
type FileType string

const (
	FileTypeGo      FileType = "go"
	FileTypeTypeScript FileType = "typescript"
	FileTypeJavaScript FileType = "javascript"
	FileTypePython  FileType = "python"
	FileTypeRust    FileType = "rust"
	FileTypeConfig  FileType = "config"
	FileTypeDoc     FileType = "documentation"
	FileTypeOther   FileType = "other"
)

// ProjectStructure represents the overall project structure
type ProjectStructure struct {
	MainLanguage    string
	Framework       string
	HasDatabase     bool
	HasAPI          bool
	HasFrontend     bool
	HasTests        bool
	EntryPoints     []string
	ConfigFiles     []string
}

// NewCodebaseContext creates a new codebase context analyzer
func NewCodebaseContext(rootPath string) *CodebaseContext {
	return &CodebaseContext{
		RootPath:     rootPath,
		Files:        []FileInfo{},
		Dependencies: make(map[string][]string),
	}
}

// AnalyzeProject analyzes the entire project structure
func (cc *CodebaseContext) AnalyzeProject() error {
	// Walk through all files
	err := filepath.Walk(cc.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and certain files
		if strings.HasPrefix(info.Name(), ".") && info.IsDir() {
			if info.Name() == ".sdd" || info.Name() == ".agents" {
				return nil // Don't skip our own directories
			}
			return filepath.SkipDir
		}

		// Skip common directories
		if info.IsDir() && (info.Name() == "node_modules" || info.Name() == "vendor" || info.Name() == ".git") {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			fileInfo, err := cc.analyzeFile(path, info)
			if err != nil {
				return err
			}
			if fileInfo != nil {
				cc.Files = append(cc.Files, *fileInfo)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to analyze project: %w", err)
	}

	// Analyze project structure
	cc.analyzeStructure()

	return nil
}

// analyzeFile analyzes a single file
func (cc *CodebaseContext) analyzeFile(path string, info os.FileInfo) (*FileInfo, error) {
	ext := strings.ToLower(filepath.Ext(path))
	fileType := getFileType(path, ext)

	// Only analyze relevant files
	if fileType == FileTypeOther {
		return nil, nil
	}

	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fileInfo := &FileInfo{
		Path:     strings.TrimPrefix(path, cc.RootPath+"/"),
		Type:     fileType,
		Language: getLanguage(ext),
		Content:  string(content),
		Size:     info.Size(),
		Imports:  extractImports(string(content), fileType),
	}

	return fileInfo, nil
}

// analyzeStructure determines the project structure
func (cc *CodebaseContext) analyzeStructure() {
	structure := ProjectStructure{}

	// Count files by type
	typeCounts := make(map[FileType]int)
	for _, file := range cc.Files {
		typeCounts[file.Type]++
	}

	// Determine main language
	maxCount := 0
	for fileType, count := range typeCounts {
		if count > maxCount {
			maxCount = count
			switch fileType {
			case FileTypeGo:
				structure.MainLanguage = "Go"
			case FileTypeTypeScript, FileTypeJavaScript:
				structure.MainLanguage = "TypeScript/JavaScript"
			case FileTypePython:
				structure.MainLanguage = "Python"
			case FileTypeRust:
				structure.MainLanguage = "Rust"
			}
		}
	}

	// Detect framework and features
	for _, file := range cc.Files {
		content := strings.ToLower(file.Content)

		// Framework detection
		if strings.Contains(content, "gin") || strings.Contains(content, "echo") || strings.Contains(content, "fiber") {
			structure.Framework = "Go Web Framework"
		}
		if strings.Contains(content, "react") {
			structure.Framework = "React"
		}
		if strings.Contains(content, "vue") {
			structure.Framework = "Vue.js"
		}
		if strings.Contains(content, "express") {
			structure.Framework = "Express.js"
		}

		// Feature detection
		if strings.Contains(content, "database") || strings.Contains(content, "sql") || strings.Contains(content, "gorm") || strings.Contains(content, "mongo") {
			structure.HasDatabase = true
		}
		if strings.Contains(content, "http") || strings.Contains(content, "router") || strings.Contains(content, "api") {
			structure.HasAPI = true
		}
		if strings.Contains(content, "component") || strings.Contains(content, "render") || strings.Contains(content, "jsx") || strings.Contains(content, "tsx") {
			structure.HasFrontend = true
		}
		if strings.Contains(file.Path, "test") || strings.Contains(file.Path, "_test.go") || strings.Contains(file.Path, ".spec.") {
			structure.HasTests = true
		}

		// Entry points and config
		if isEntryPoint(file.Path, file.Type) {
			structure.EntryPoints = append(structure.EntryPoints, file.Path)
		}
		if isConfigFile(file.Path) {
			structure.ConfigFiles = append(structure.ConfigFiles, file.Path)
		}
	}

	cc.Structure = structure
}

// GetContextForPhase returns relevant context for a specific SDD phase
func (cc *CodebaseContext) GetContextForPhase(phase string) string {
	var context strings.Builder

	switch phase {
	case "specify":
		context.WriteString(cc.getSpecificationContext())
	case "plan":
		context.WriteString(cc.getPlanningContext())
	case "task":
		context.WriteString(cc.getTaskContext())
	case "execute":
		context.WriteString(cc.getExecutionContext())
	case "review":
		context.WriteString(cc.getReviewContext())
	}

	return context.String()
}

// getSpecificationContext provides context for requirement specification
func (cc *CodebaseContext) getSpecificationContext() string {
	var ctx strings.Builder

	ctx.WriteString("## Existing Codebase Context\n\n")

	if cc.Structure.MainLanguage != "" {
		ctx.WriteString(fmt.Sprintf("**Primary Language:** %s\n", cc.Structure.MainLanguage))
	}

	if cc.Structure.Framework != "" {
		ctx.WriteString(fmt.Sprintf("**Framework:** %s\n", cc.Structure.Framework))
	}

	ctx.WriteString("\n**Project Features:**\n")
	if cc.Structure.HasAPI {
		ctx.WriteString("- Has API endpoints\n")
	}
	if cc.Structure.HasDatabase {
		ctx.WriteString("- Uses database\n")
	}
	if cc.Structure.HasFrontend {
		ctx.WriteString("- Has frontend components\n")
	}
	if cc.Structure.HasTests {
		ctx.WriteString("- Includes tests\n")
	}

	if len(cc.Structure.EntryPoints) > 0 {
		ctx.WriteString("\n**Entry Points:**\n")
		for _, entry := range cc.Structure.EntryPoints {
			ctx.WriteString(fmt.Sprintf("- %s\n", entry))
		}
	}

	ctx.WriteString("\n**Key Files to Consider:**\n")
	for _, file := range cc.Files {
		if len(file.Content) < 5000 { // Only include smaller files
			ctx.WriteString(fmt.Sprintf("- %s (%s)\n", file.Path, file.Language))
		}
	}

	return ctx.String()
}

// getPlanningContext provides context for architecture planning
func (cc *CodebaseContext) getPlanningContext() string {
	var ctx strings.Builder

	ctx.WriteString("## Architecture Context\n\n")

	ctx.WriteString("**Current Structure:**\n")
	ctx.WriteString(fmt.Sprintf("- Main Language: %s\n", cc.Structure.MainLanguage))
	ctx.WriteString(fmt.Sprintf("- Framework: %s\n", cc.Structure.Framework))

	ctx.WriteString("\n**Technology Stack:**\n")

	// Analyze dependencies
	techStack := cc.analyzeTechStack()
	for category, technologies := range techStack {
		if len(technologies) > 0 {
			ctx.WriteString(fmt.Sprintf("- **%s:** %s\n", category, strings.Join(technologies, ", ")))
		}
	}

	ctx.WriteString("\n**Existing Patterns:**\n")
	patterns := cc.analyzePatterns()
	for _, pattern := range patterns {
		ctx.WriteString(fmt.Sprintf("- %s\n", pattern))
	}

	return ctx.String()
}

// getTaskContext provides context for task breakdown
func (cc *CodebaseContext) getTaskContext() string {
	var ctx strings.Builder

	ctx.WriteString("## Implementation Context\n\n")

	ctx.WriteString("**Codebase Size:**\n")
	ctx.WriteString(fmt.Sprintf("- Total files: %d\n", len(cc.Files)))

	fileTypes := make(map[FileType]int)
	totalSize := int64(0)
	for _, file := range cc.Files {
		fileTypes[file.Type]++
		totalSize += file.Size
	}

	ctx.WriteString(fmt.Sprintf("- Total size: %d bytes\n", totalSize))
	ctx.WriteString("- File types:\n")
	for fileType, count := range fileTypes {
		ctx.WriteString(fmt.Sprintf("  - %s: %d files\n", fileType, count))
	}

	ctx.WriteString("\n**Integration Points:**\n")
	for _, file := range cc.Files {
		if strings.Contains(strings.ToLower(file.Content), "api") ||
		   strings.Contains(strings.ToLower(file.Content), "database") ||
		   strings.Contains(strings.ToLower(file.Content), "external") {
			ctx.WriteString(fmt.Sprintf("- %s (potential integration)\n", file.Path))
		}
	}

	return ctx.String()
}

// getExecutionContext provides context for implementation
func (cc *CodebaseContext) getExecutionContext() string {
	var ctx strings.Builder

	ctx.WriteString("## Development Context\n\n")

	ctx.WriteString("**Coding Standards Found:**\n")

	// Analyze coding patterns
	if len(cc.Files) > 0 {
		sampleFile := cc.Files[0]
		if strings.Contains(sampleFile.Content, "error") {
			ctx.WriteString("- Uses Go-style error handling\n")
		}
		if strings.Contains(sampleFile.Content, "interface{}") {
			ctx.WriteString("- Uses empty interfaces for flexibility\n")
		}
	}

	ctx.WriteString("\n**Import Patterns:**\n")
	importPatterns := make(map[string]int)
	for _, file := range cc.Files {
		for _, imp := range file.Imports {
			if strings.Contains(imp, "/") {
				parts := strings.Split(imp, "/")
				if len(parts) > 1 {
					importPatterns[parts[len(parts)-2]]++
				}
			}
		}
	}

	for pattern, count := range importPatterns {
		if count > 1 {
			ctx.WriteString(fmt.Sprintf("- %s (used in %d files)\n", pattern, count))
		}
	}

	return ctx.String()
}

// getReviewContext provides context for QA review
func (cc *CodebaseContext) getReviewContext() string {
	var ctx strings.Builder

	ctx.WriteString("## Quality Assurance Context\n\n")

	ctx.WriteString("**Testing Status:**\n")
	if cc.Structure.HasTests {
		ctx.WriteString("- Testing framework detected\n")
	} else {
		ctx.WriteString("- No tests found\n")
	}

	ctx.WriteString("\n**Code Quality Metrics:**\n")
	totalLines := 0
	totalFiles := len(cc.Files)

	for _, file := range cc.Files {
		lines := strings.Split(file.Content, "\n")
		totalLines += len(lines)
	}

	if totalFiles > 0 {
		avgLines := totalLines / totalFiles
		ctx.WriteString(fmt.Sprintf("- Average lines per file: %d\n", avgLines))
	}

	ctx.WriteString(fmt.Sprintf("- Total files: %d\n", totalFiles))
	ctx.WriteString(fmt.Sprintf("- Total lines: %d\n", totalLines))

	return ctx.String()
}

// Helper functions

func getFileType(path, ext string) FileType {
	switch ext {
	case ".go":
		return FileTypeGo
	case ".ts":
		return FileTypeTypeScript
	case ".tsx":
		return FileTypeTypeScript
	case ".js":
		return FileTypeJavaScript
	case ".jsx":
		return FileTypeJavaScript
	case ".py":
		return FileTypePython
	case ".rs":
		return FileTypeRust
	case ".json", ".yaml", ".yml", ".toml", ".ini", ".cfg":
		return FileTypeConfig
	case ".md", ".txt", ".rst":
		return FileTypeDoc
	default:
		return FileTypeOther
	}
}

func getLanguage(ext string) string {
	switch ext {
	case ".go":
		return "Go"
	case ".ts", ".tsx":
		return "TypeScript"
	case ".js", ".jsx":
		return "JavaScript"
	case ".py":
		return "Python"
	case ".rs":
		return "Rust"
	default:
		return "Unknown"
	}
}

func extractImports(content string, fileType FileType) []string {
	var imports []string

	switch fileType {
	case FileTypeGo:
		// Extract Go imports
		lines := strings.Split(content, "\n")
		inImport := false
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "import") {
				if strings.Contains(line, "(") {
					inImport = true
					continue
				} else {
					// Single import
					imp := extractGoImport(line)
					if imp != "" {
						imports = append(imports, imp)
					}
				}
			} else if inImport {
				if strings.Contains(line, ")") {
					inImport = false
				} else {
					imp := extractGoImport(line)
					if imp != "" {
						imports = append(imports, imp)
					}
				}
			}
		}
	case FileTypeTypeScript, FileTypeJavaScript:
		// Extract JS/TS imports
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "import") || strings.HasPrefix(line, "require") {
				// Basic import extraction - can be enhanced
				imports = append(imports, "external_dependency")
			}
		}
	}

	return imports
}

func extractGoImport(line string) string {
	line = strings.TrimSpace(line)
	line = strings.Trim(line, `"`)
	if strings.Contains(line, `"`) {
		parts := strings.Split(line, `"`)
		if len(parts) >= 2 {
			return parts[1]
		}
	}
	return ""
}

func isEntryPoint(path string, fileType FileType) bool {
	switch fileType {
	case FileTypeGo:
		return strings.HasSuffix(path, "main.go") || strings.HasSuffix(path, "cmd/server/main.go")
	case FileTypeJavaScript, FileTypeTypeScript:
		return strings.HasSuffix(path, "index.js") || strings.HasSuffix(path, "app.js") ||
			   strings.HasSuffix(path, "server.js") || strings.HasSuffix(path, "main.ts")
	case FileTypePython:
		return strings.HasSuffix(path, "__main__.py") || path == "main.py" || path == "app.py"
	}
	return false
}

func isConfigFile(path string) bool {
	configFiles := []string{
		"go.mod", "go.sum", "package.json", "tsconfig.json", "pyproject.toml",
		"Cargo.toml", "requirements.txt", "Pipfile", ".env", "config.yaml",
		"config.json", "docker-compose.yml", "Dockerfile",
	}
	for _, config := range configFiles {
		if strings.Contains(path, config) {
			return true
		}
	}
	return false
}

func (cc *CodebaseContext) analyzeTechStack() map[string][]string {
	stack := make(map[string][]string)

	technologies := make(map[string]bool)

	for _, file := range cc.Files {
		content := strings.ToLower(file.Content)

		// Web frameworks
		if strings.Contains(content, "gin") {
			technologies["gin"] = true
		}
		if strings.Contains(content, "echo") {
			technologies["echo"] = true
		}
		if strings.Contains(content, "fiber") {
			technologies["fiber"] = true
		}
		if strings.Contains(content, "express") {
			technologies["express"] = true
		}
		if strings.Contains(content, "react") {
			technologies["react"] = true
		}
		if strings.Contains(content, "vue") {
			technologies["vue"] = true
		}

		// Databases
		if strings.Contains(content, "postgres") || strings.Contains(content, "postgresql") {
			technologies["postgresql"] = true
		}
		if strings.Contains(content, "mysql") {
			technologies["mysql"] = true
		}
		if strings.Contains(content, "mongodb") || strings.Contains(content, "mongo") {
			technologies["mongodb"] = true
		}
		if strings.Contains(content, "redis") {
			technologies["redis"] = true
		}

		// Languages and runtimes
		if file.Type == FileTypeGo {
			technologies["go"] = true
		}
		if file.Type == FileTypeTypeScript {
			technologies["typescript"] = true
		}
		if file.Type == FileTypeJavaScript {
			technologies["javascript"] = true
		}
		if file.Type == FileTypePython {
			technologies["python"] = true
		}
	}

	// Categorize technologies
	stack["Languages"] = []string{}
	stack["Frameworks"] = []string{}
	stack["Databases"] = []string{}
	stack["Tools"] = []string{}

	for tech := range technologies {
		switch tech {
		case "go", "typescript", "javascript", "python", "rust":
			stack["Languages"] = append(stack["Languages"], tech)
		case "gin", "echo", "fiber", "express", "react", "vue":
			stack["Frameworks"] = append(stack["Frameworks"], tech)
		case "postgresql", "mysql", "mongodb", "redis":
			stack["Databases"] = append(stack["Databases"], tech)
		default:
			stack["Tools"] = append(stack["Tools"], tech)
		}
	}

	return stack
}

func (cc *CodebaseContext) analyzePatterns() []string {
	var patterns []string

	hasInterfaces := false
	hasStructs := false
	hasErrorHandling := false
	hasConcurrency := false
	hasTesting := false

	for _, file := range cc.Files {
		content := file.Content

		if strings.Contains(content, "type ") && strings.Contains(content, "interface") {
			hasInterfaces = true
		}
		if strings.Contains(content, "type ") && strings.Contains(content, "struct") {
			hasStructs = true
		}
		if strings.Contains(content, "if err != nil") {
			hasErrorHandling = true
		}
		if strings.Contains(content, "go func") || strings.Contains(content, "goroutine") {
			hasConcurrency = true
		}
		if strings.Contains(file.Path, "test") || strings.Contains(file.Path, "_test.go") {
			hasTesting = true
		}
	}

	if hasInterfaces {
		patterns = append(patterns, "Uses interfaces for abstraction")
	}
	if hasStructs {
		patterns = append(patterns, "Uses structs for data modeling")
	}
	if hasErrorHandling {
		patterns = append(patterns, "Explicit error handling patterns")
	}
	if hasConcurrency {
		patterns = append(patterns, "Concurrent programming patterns")
	}
	if hasTesting {
		patterns = append(patterns, "Unit testing practices")
	}

	return patterns
}