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

// BrownfieldContext provides comprehensive analysis for existing codebases
type BrownfieldContext struct {
	CodebaseContext
	LegacyPatterns     []LegacyPattern
	ForbiddenPatterns  []ForbiddenPattern
	IntegrationPoints  []IntegrationPoint
	TechnicalDebt      []TechnicalDebtItem
	Constitution       Constitution
}

// LegacyPattern represents established patterns in the codebase
type LegacyPattern struct {
	Pattern     string
	Description string
	Files       []string
	Examples    []string
}

// ForbiddenPattern represents anti-patterns that should be avoided
type ForbiddenPattern struct {
	Pattern     string
	Description string
	Severity    string
	Occurrences []string
	Recommended string
}

// IntegrationPoint represents key integration points in the system
type IntegrationPoint struct {
	Name        string
	Type        string // api, database, external_service, etc.
	Description string
	Files       []string
	Dependencies []string
}

// TechnicalDebtItem represents identified technical debt
type TechnicalDebtItem struct {
	Issue       string
	Severity    string
	Files       []string
	Description string
	Recommendation string
}

// Constitution represents the system's architectural rules
type Constitution struct {
	TechStack        []string
	ArchitecturalRules []string
	CodingStandards []string
	IntegrationRules []string
	QualityGates    []string
}

// NewCodebaseContext creates a new codebase context analyzer
func NewCodebaseContext(rootPath string) *CodebaseContext {
	return &CodebaseContext{
		RootPath:     rootPath,
		Files:        []FileInfo{},
		Dependencies: make(map[string][]string),
	}
}

// NewBrownfieldContext creates a comprehensive brownfield analysis context
func NewBrownfieldContext(rootPath string) *BrownfieldContext {
	return &BrownfieldContext{
		CodebaseContext: CodebaseContext{
			RootPath:     rootPath,
			Files:        []FileInfo{},
			Dependencies: make(map[string][]string),
		},
		LegacyPatterns:    []LegacyPattern{},
		ForbiddenPatterns: []ForbiddenPattern{},
		IntegrationPoints: []IntegrationPoint{},
		TechnicalDebt:     []TechnicalDebtItem{},
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

// AnalyzeBrownfield performs comprehensive brownfield analysis
func (bfc *BrownfieldContext) AnalyzeBrownfield() error {
	// First perform basic codebase analysis
	if err := bfc.AnalyzeProject(); err != nil {
		return fmt.Errorf("failed to analyze codebase: %w", err)
	}

	// Analyze legacy patterns
	if err := bfc.analyzeLegacyPatterns(); err != nil {
		return fmt.Errorf("failed to analyze legacy patterns: %w", err)
	}

	// Identify forbidden patterns
	if err := bfc.identifyForbiddenPatterns(); err != nil {
		return fmt.Errorf("failed to identify forbidden patterns: %w", err)
	}

	// Map integration points
	if err := bfc.mapIntegrationPoints(); err != nil {
		return fmt.Errorf("failed to map integration points: %w", err)
	}

	// Assess technical debt
	if err := bfc.assessTechnicalDebt(); err != nil {
		return fmt.Errorf("failed to assess technical debt: %w", err)
	}

	// Load constitution if it exists
	if err := bfc.loadConstitution(); err != nil {
		// Constitution doesn't exist yet, create default
		bfc.createDefaultConstitution()
	}

	return nil
}

// analyzeLegacyPatterns identifies established patterns in the codebase
func (bfc *BrownfieldContext) analyzeLegacyPatterns() error {
	patterns := []LegacyPattern{}

	// Analyze error handling patterns
	errorPatterns := bfc.analyzeErrorHandlingPatterns()
	patterns = append(patterns, errorPatterns...)

	// Analyze data access patterns
	dataPatterns := bfc.analyzeDataAccessPatterns()
	patterns = append(patterns, dataPatterns...)

	// Analyze architectural patterns
	archPatterns := bfc.analyzeArchitecturalPatterns()
	patterns = append(patterns, archPatterns...)

	bfc.LegacyPatterns = patterns
	return nil
}

// identifyForbiddenPatterns finds anti-patterns that should be avoided
func (bfc *BrownfieldContext) identifyForbiddenPatterns() error {
	forbidden := []ForbiddenPattern{}

	// Check for deprecated patterns
	deprecatedPatterns := bfc.checkDeprecatedPatterns()
	forbidden = append(forbidden, deprecatedPatterns...)

	// Check for security issues
	securityPatterns := bfc.checkSecurityPatterns()
	forbidden = append(forbidden, securityPatterns...)

	// Check for performance issues
	performancePatterns := bfc.checkPerformancePatterns()
	forbidden = append(forbidden, performancePatterns...)

	bfc.ForbiddenPatterns = forbidden
	return nil
}

// mapIntegrationPoints identifies key integration points
func (bfc *BrownfieldContext) mapIntegrationPoints() error {
	points := []IntegrationPoint{}

	// API endpoints
	apiPoints := bfc.mapAPIPoints()
	points = append(points, apiPoints...)

	// Database connections
	dbPoints := bfc.mapDatabasePoints()
	points = append(points, dbPoints...)

	// External services
	externalPoints := bfc.mapExternalServicePoints()
	points = append(points, externalPoints...)

	// File system interactions
	filePoints := bfc.mapFileSystemPoints()
	points = append(points, filePoints...)

	bfc.IntegrationPoints = points
	return nil
}

// assessTechnicalDebt identifies technical debt items
func (bfc *BrownfieldContext) assessTechnicalDebt() error {
	debt := []TechnicalDebtItem{}

	// Code quality debt
	qualityDebt := bfc.assessCodeQualityDebt()
	debt = append(debt, qualityDebt...)

	// Architecture debt
	archDebt := bfc.assessArchitectureDebt()
	debt = append(debt, archDebt...)

	// Test coverage debt
	testDebt := bfc.assessTestCoverageDebt()
	debt = append(debt, testDebt...)

	bfc.TechnicalDebt = debt
	return nil
}

// GenerateCONTEXTFile creates the CONTEXT.md file for brownfield development
func (bfc *BrownfieldContext) GenerateCONTEXTFile() string {
	var ctx strings.Builder

	ctx.WriteString("---\n")
	ctx.WriteString("title: System Context & Constitution\n")
	ctx.WriteString("generated_by: Ultimate SDD Framework - Brownfield Discovery\n")
	ctx.WriteString("last_updated: Generated automatically\n")
	ctx.WriteString("---\n\n")

	ctx.WriteString("# System Context & Constitution\n\n")
	ctx.WriteString("This document serves as the **Source of Truth** for the current system state.\n")
	ctx.WriteString("All development decisions should reference this context to maintain consistency.\n\n")

	// Constitution
	ctx.WriteString("## 🏛️ System Constitution\n\n")
	ctx.WriteString("### Technology Stack\n")
	for _, tech := range bfc.Constitution.TechStack {
		ctx.WriteString(fmt.Sprintf("- %s\n", tech))
	}
	ctx.WriteString("\n")

	ctx.WriteString("### Architectural Rules\n")
	for _, rule := range bfc.Constitution.ArchitecturalRules {
		ctx.WriteString(fmt.Sprintf("- %s\n", rule))
	}
	ctx.WriteString("\n")

	ctx.WriteString("### Coding Standards\n")
	for _, standard := range bfc.Constitution.CodingStandards {
		ctx.WriteString(fmt.Sprintf("- %s\n", standard))
	}
	ctx.WriteString("\n")

	// Legacy Patterns
	ctx.WriteString("## 📚 Legacy Patterns\n\n")
	ctx.WriteString("These are established patterns that should be followed for consistency:\n\n")

	for i, pattern := range bfc.LegacyPatterns {
		ctx.WriteString(fmt.Sprintf("### %d. %s\n", i+1, pattern.Pattern))
		ctx.WriteString(fmt.Sprintf("**Description:** %s\n\n", pattern.Description))

		if len(pattern.Files) > 0 {
			ctx.WriteString("**Found in files:**\n")
			for _, file := range pattern.Files {
				ctx.WriteString(fmt.Sprintf("- %s\n", file))
			}
			ctx.WriteString("\n")
		}

		if len(pattern.Examples) > 0 {
			ctx.WriteString("**Examples:**\n")
			for _, example := range pattern.Examples {
				ctx.WriteString(fmt.Sprintf("```go\n%s\n```\n", example))
			}
			ctx.WriteString("\n")
		}
	}

	// Forbidden Patterns
	ctx.WriteString("## 🚫 Forbidden Patterns\n\n")
	ctx.WriteString("These anti-patterns should be avoided:\n\n")

	for i, pattern := range bfc.ForbiddenPatterns {
		ctx.WriteString(fmt.Sprintf("### %d. %s\n", i+1, pattern.Pattern))
		ctx.WriteString(fmt.Sprintf("**Severity:** %s\n", pattern.Severity))
		ctx.WriteString(fmt.Sprintf("**Description:** %s\n\n", pattern.Description))

		if len(pattern.Occurrences) > 0 {
			ctx.WriteString("**Found in:**\n")
			for _, occurrence := range pattern.Occurrences {
				ctx.WriteString(fmt.Sprintf("- %s\n", occurrence))
			}
			ctx.WriteString("\n")
		}

		if pattern.Recommended != "" {
			ctx.WriteString(fmt.Sprintf("**Recommended:** %s\n\n", pattern.Recommended))
		}
	}

	// Integration Points
	ctx.WriteString("## 🔗 Integration Points\n\n")
	ctx.WriteString("Critical integration points in the system:\n\n")

	for i, point := range bfc.IntegrationPoints {
		ctx.WriteString(fmt.Sprintf("### %d. %s (%s)\n", i+1, point.Name, point.Type))
		ctx.WriteString(fmt.Sprintf("**Description:** %s\n\n", point.Description))

		if len(point.Files) > 0 {
			ctx.WriteString("**Files:**\n")
			for _, file := range point.Files {
				ctx.WriteString(fmt.Sprintf("- %s\n", file))
			}
			ctx.WriteString("\n")
		}

		if len(point.Dependencies) > 0 {
			ctx.WriteString("**Dependencies:**\n")
			for _, dep := range point.Dependencies {
				ctx.WriteString(fmt.Sprintf("- %s\n", dep))
			}
			ctx.WriteString("\n")
		}
	}

	// Technical Debt
	ctx.WriteString("## 💸 Technical Debt\n\n")
	ctx.WriteString("Identified technical debt that should be considered:\n\n")

	for i, debt := range bfc.TechnicalDebt {
		ctx.WriteString(fmt.Sprintf("### %d. %s\n", i+1, debt.Issue))
		ctx.WriteString(fmt.Sprintf("**Severity:** %s\n", debt.Severity))
		ctx.WriteString(fmt.Sprintf("**Description:** %s\n\n", debt.Description))

		if len(debt.Files) > 0 {
			ctx.WriteString("**Affected files:**\n")
			for _, file := range debt.Files {
				ctx.WriteString(fmt.Sprintf("- %s\n", file))
			}
			ctx.WriteString("\n")
		}

		ctx.WriteString(fmt.Sprintf("**Recommendation:** %s\n\n", debt.Recommendation))
	}

	// System Overview
	ctx.WriteString("## 📊 System Overview\n\n")
	ctx.WriteString(fmt.Sprintf("**Primary Language:** %s\n", bfc.Structure.MainLanguage))
	ctx.WriteString(fmt.Sprintf("**Framework:** %s\n", bfc.Structure.Framework))
	ctx.WriteString(fmt.Sprintf("**Total Files:** %d\n", len(bfc.Files)))

	features := []string{}
	if bfc.Structure.HasAPI {
		features = append(features, "API")
	}
	if bfc.Structure.HasDatabase {
		features = append(features, "Database")
	}
	if bfc.Structure.HasFrontend {
		features = append(features, "Frontend")
	}
	if bfc.Structure.HasTests {
		features = append(features, "Tests")
	}

	if len(features) > 0 {
		ctx.WriteString(fmt.Sprintf("**Features:** %s\n", strings.Join(features, ", ")))
	}

	ctx.WriteString("\n**Entry Points:**\n")
	for _, entry := range bfc.Structure.EntryPoints {
		ctx.WriteString(fmt.Sprintf("- %s\n", entry))
	}

	ctx.WriteString("\n---\n*Generated by Ultimate SDD Framework - Brownfield Discovery Phase*\n")
	ctx.WriteString("*This document should be updated whenever the system architecture changes.*\n")

	return ctx.String()
}

// Helper methods for brownfield analysis

func (bfc *BrownfieldContext) analyzeErrorHandlingPatterns() []LegacyPattern {
	patterns := []LegacyPattern{}

	errorPattern := LegacyPattern{
		Pattern:     "Consistent Error Handling",
		Description: "Standard error handling patterns used throughout the codebase",
		Files:       []string{},
		Examples:    []string{},
	}

	for _, file := range bfc.Files {
		if strings.Contains(file.Content, "error") {
			errorPattern.Files = append(errorPattern.Files, file.Path)

			// Extract error handling examples
			if strings.Contains(file.Content, "if err != nil") {
				errorPattern.Examples = append(errorPattern.Examples,
					"if err != nil {\n    return fmt.Errorf(\"operation failed: %w\", err)\n}")
			}
		}
	}

	if len(errorPattern.Files) > 0 {
		patterns = append(patterns, errorPattern)
	}

	return patterns
}

func (bfc *BrownfieldContext) analyzeDataAccessPatterns() []LegacyPattern {
	patterns := []LegacyPattern{}

	if bfc.Structure.HasDatabase {
		dbPattern := LegacyPattern{
			Pattern:     "Database Access Layer",
			Description: "Established patterns for database interactions",
			Files:       []string{},
			Examples:    []string{},
		}

		for _, file := range bfc.Files {
			content := strings.ToLower(file.Content)
			if strings.Contains(content, "database") || strings.Contains(content, "sql") ||
			   strings.Contains(content, "query") || strings.Contains(content, "gorm") {
				dbPattern.Files = append(dbPattern.Files, file.Path)
			}
		}

		if len(dbPattern.Files) > 0 {
			patterns = append(patterns, dbPattern)
		}
	}

	return patterns
}

func (bfc *BrownfieldContext) analyzeArchitecturalPatterns() []LegacyPattern {
	patterns := []LegacyPattern{}

	// MVC pattern detection
	mvcPattern := LegacyPattern{
		Pattern:     "MVC Architecture",
		Description: "Model-View-Controller pattern implementation",
		Files:       []string{},
	}

	for _, file := range bfc.Files {
		path := strings.ToLower(file.Path)
		if strings.Contains(path, "model") || strings.Contains(path, "view") || strings.Contains(path, "controller") {
			mvcPattern.Files = append(mvcPattern.Files, file.Path)
		}
	}

	if len(mvcPattern.Files) > 3 { // At least one of each type
		patterns = append(patterns, mvcPattern)
	}

	return patterns
}

func (bfc *BrownfieldContext) checkDeprecatedPatterns() []ForbiddenPattern {
	forbidden := []ForbiddenPattern{}

	// Check for deprecated libraries or patterns
	for _, file := range bfc.Files {
		if strings.Contains(file.Content, "deprecated") ||
		   strings.Contains(strings.ToLower(file.Content), "todo: remove") {
			forbidden = append(forbidden, ForbiddenPattern{
				Pattern:     "Deprecated Code Usage",
				Description: "Code marked as deprecated should be avoided",
				Severity:    "Medium",
				Occurrences: []string{file.Path},
				Recommended: "Replace with current recommended alternatives",
			})
		}
	}

	return forbidden
}

func (bfc *BrownfieldContext) checkSecurityPatterns() []ForbiddenPattern {
	forbidden := []ForbiddenPattern{}

	// Check for potential security issues
	for _, file := range bfc.Files {
		content := file.Content

		// SQL injection risks
		if strings.Contains(content, "sprintf") && strings.Contains(content, "query") {
			forbidden = append(forbidden, ForbiddenPattern{
				Pattern:     "Potential SQL Injection",
				Description: "String formatting in SQL queries can lead to injection attacks",
				Severity:    "High",
				Occurrences: []string{file.Path},
				Recommended: "Use parameterized queries or prepared statements",
			})
		}

		// Hardcoded secrets
		if strings.Contains(content, "password") && (strings.Contains(content, "=") || strings.Contains(content, ":")) {
			forbidden = append(forbidden, ForbiddenPattern{
				Pattern:     "Hardcoded Credentials",
				Description: "Credentials should not be hardcoded in source code",
				Severity:    "Critical",
				Occurrences: []string{file.Path},
				Recommended: "Use environment variables or secure credential storage",
			})
		}
	}

	return forbidden
}

func (bfc *BrownfieldContext) checkPerformancePatterns() []ForbiddenPattern {
	forbidden := []ForbiddenPattern{}

	// Check for N+1 query patterns
	for _, file := range bfc.Files {
		content := file.Content
		if strings.Contains(content, "for") && strings.Contains(content, "query") &&
		   strings.Contains(content, "range") {
			forbidden = append(forbidden, ForbiddenPattern{
				Pattern:     "Potential N+1 Query",
				Description: "Looping and querying in loops can cause performance issues",
				Severity:    "Medium",
				Occurrences: []string{file.Path},
				Recommended: "Use batch queries or eager loading",
			})
		}
	}

	return forbidden
}

func (bfc *BrownfieldContext) mapAPIPoints() []IntegrationPoint {
	points := []IntegrationPoint{}

	for _, file := range bfc.Files {
		if strings.Contains(file.Content, "router") || strings.Contains(file.Content, "route") ||
		   strings.Contains(file.Content, "api") || strings.Contains(file.Content, "endpoint") {

			// Extract route definitions
			lines := strings.Split(file.Content, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.Contains(line, "GET") || strings.Contains(line, "POST") ||
				   strings.Contains(line, "PUT") || strings.Contains(line, "DELETE") {
					points = append(points, IntegrationPoint{
						Name:        fmt.Sprintf("API Route in %s", filepath.Base(file.Path)),
						Type:        "api",
						Description: fmt.Sprintf("API endpoint defined in %s", file.Path),
						Files:       []string{file.Path},
						Dependencies: []string{"HTTP framework"},
					})
					break
				}
			}
		}
	}

	return points
}

func (bfc *BrownfieldContext) mapDatabasePoints() []IntegrationPoint {
	points := []IntegrationPoint{}

	if bfc.Structure.HasDatabase {
		dbPoint := IntegrationPoint{
			Name:        "Database Connection",
			Type:        "database",
			Description: "Primary database connection and configuration",
			Files:       []string{},
			Dependencies: []string{"Database driver"},
		}

		for _, file := range bfc.Files {
			content := strings.ToLower(file.Content)
			if strings.Contains(content, "database") || strings.Contains(content, "sql") ||
			   strings.Contains(content, "gorm") || strings.Contains(content, "connection") {
				dbPoint.Files = append(dbPoint.Files, file.Path)
			}
		}

		if len(dbPoint.Files) > 0 {
			points = append(points, dbPoint)
		}
	}

	return points
}

func (bfc *BrownfieldContext) mapExternalServicePoints() []IntegrationPoint {
	points := []IntegrationPoint{}

	for _, file := range bfc.Files {
		content := strings.ToLower(file.Content)

		// Check for external API calls
		if strings.Contains(content, "http.client") || strings.Contains(content, "fetch") ||
		   strings.Contains(content, "axios") || strings.Contains(content, "requests") {
			points = append(points, IntegrationPoint{
				Name:        fmt.Sprintf("External API Call in %s", filepath.Base(file.Path)),
				Type:        "external_service",
				Description: "External service integration",
				Files:       []string{file.Path},
				Dependencies: []string{"HTTP client library"},
			})
		}
	}

	return points
}

func (bfc *BrownfieldContext) mapFileSystemPoints() []IntegrationPoint {
	points := []IntegrationPoint{}

	for _, file := range bfc.Files {
		if strings.Contains(file.Content, "os.Open") || strings.Contains(file.Content, "fs") ||
		   strings.Contains(file.Content, "filepath") {
			points = append(points, IntegrationPoint{
				Name:        fmt.Sprintf("File System Access in %s", filepath.Base(file.Path)),
				Type:        "filesystem",
				Description: "File system operations",
				Files:       []string{file.Path},
				Dependencies: []string{"OS filesystem"},
			})
		}
	}

	return points
}

func (bfc *BrownfieldContext) assessCodeQualityDebt() []TechnicalDebtItem {
	debt := []TechnicalDebtItem{}

	// Check for long files
	for _, file := range bfc.Files {
		lines := strings.Split(file.Content, "\n")
		if len(lines) > 500 {
			debt = append(debt, TechnicalDebtItem{
				Issue:        "Long File",
				Severity:     "Low",
				Files:        []string{file.Path},
				Description:  fmt.Sprintf("File has %d lines, making it hard to maintain", len(lines)),
				Recommendation: "Consider splitting into smaller, focused modules",
			})
		}
	}

	// Check for complex functions
	for _, file := range bfc.Files {
		lines := strings.Split(file.Content, "\n")
		for i, line := range lines {
			if strings.Contains(line, "func ") && strings.Contains(line, "{") {
				// Count lines in function (simplified)
				funcStart := i
				braceCount := 0
				funcLength := 0

				for j := funcStart; j < len(lines); j++ {
					funcLength++
					if strings.Contains(lines[j], "{") {
						braceCount++
					}
					if strings.Contains(lines[j], "}") {
						braceCount--
						if braceCount == 0 {
							break
						}
					}
				}

				if funcLength > 50 {
					debt = append(debt, TechnicalDebtItem{
						Issue:        "Complex Function",
						Severity:     "Medium",
						Files:        []string{file.Path},
						Description:  fmt.Sprintf("Function starting at line %d has %d lines", funcStart+1, funcLength),
						Recommendation: "Break down into smaller, focused functions",
					})
				}
			}
		}
	}

	return debt
}

func (bfc *BrownfieldContext) assessArchitectureDebt() []TechnicalDebtItem {
	debt := []TechnicalDebtItem{}

	// Check for circular dependencies (simplified)
	fileMap := make(map[string][]string)
	for _, file := range bfc.Files {
		fileMap[file.Path] = file.Imports
	}

	// This is a simplified check - real circular dependency detection is more complex
	for file, imports := range fileMap {
		for _, imp := range imports {
			if importedFile, exists := fileMap[imp]; exists {
				for _, reverseImp := range importedFile {
					if reverseImp == file {
						debt = append(debt, TechnicalDebtItem{
							Issue:        "Potential Circular Dependency",
							Severity:     "High",
							Files:        []string{file, imp},
							Description:  "Files may have circular import dependencies",
							Recommendation: "Refactor to break circular dependencies using interfaces or dependency injection",
						})
						break
					}
				}
			}
		}
	}

	return debt
}

func (bfc *BrownfieldContext) assessTestCoverageDebt() []TechnicalDebtItem {
	debt := []TechnicalDebtItem{}

	if !bfc.Structure.HasTests {
		debt = append(debt, TechnicalDebtItem{
			Issue:        "Missing Test Suite",
			Severity:     "High",
			Files:        []string{},
			Description:  "No test files detected in the codebase",
			Recommendation: "Implement comprehensive unit and integration tests",
		})
	} else {
		// Count test files vs source files
		sourceFiles := 0
		testFiles := 0

		for _, file := range bfc.Files {
			if strings.Contains(file.Path, "_test.go") || strings.Contains(file.Path, ".test.") {
				testFiles++
			} else if file.Type == FileTypeGo {
				sourceFiles++
			}
		}

		if testFiles < sourceFiles/2 {
			debt = append(debt, TechnicalDebtItem{
				Issue:        "Low Test Coverage",
				Severity:     "Medium",
				Files:        []string{},
				Description:  fmt.Sprintf("Only %d test files for %d source files", testFiles, sourceFiles),
				Recommendation: "Increase test coverage to at least 80%",
			})
		}
	}

	return debt
}

func (bfc *BrownfieldContext) loadConstitution() error {
	constitutionPath := filepath.Join(bfc.RootPath, "CONSTITUTION.md")

	if _, err := os.Stat(constitutionPath); os.IsNotExist(err) {
		return fmt.Errorf("constitution file not found")
	}

	_, err := os.ReadFile(constitutionPath)
	if err != nil {
		return err
	}

	// Parse constitution (simplified parsing - TODO: implement actual parsing)
	bfc.Constitution = Constitution{
		TechStack:         []string{"Go", "React", "PostgreSQL"},
		ArchitecturalRules: []string{"MVC pattern", "Repository pattern"},
		CodingStandards:   []string{"Go naming conventions", "Error handling"},
		IntegrationRules:  []string{"REST API", "Database transactions"},
		QualityGates:      []string{"Tests pass", "Linting passes"},
	}

	return nil
}

func (bfc *BrownfieldContext) createDefaultConstitution() {
	bfc.Constitution = Constitution{
		TechStack:         []string{bfc.Structure.MainLanguage, bfc.Structure.Framework},
		ArchitecturalRules: []string{"Follow established patterns", "Maintain separation of concerns"},
		CodingStandards:   []string{"Follow language conventions", "Consistent error handling"},
		IntegrationRules:  []string{"Use existing integration points", "Maintain API contracts"},
		QualityGates:      []string{"No regressions", "Tests pass", "Code review approved"},
	}
}