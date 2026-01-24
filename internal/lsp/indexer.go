package lsp

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// Symbol represents a code symbol (function, class, etc.)
type Symbol struct {
	Name       string
	Kind       string // "function", "class", "method", "variable", "type"
	File       string
	Line       int
	Signature  string
	DocComment string
	Parent     string // Parent class/interface for methods
}

// FileIndex represents indexed data for a single file
type FileIndex struct {
	Path     string
	Language string
	Symbols  []Symbol
	Imports  []string
	Size     int64
	Modified int64
}

// ProjectIndex represents the entire project index
type ProjectIndex struct {
	Root       string
	Files      map[string]*FileIndex
	SymbolMap  map[string][]Symbol // symbol name -> locations
	mu         sync.RWMutex
}

// Indexer handles codebase indexing
type Indexer struct {
	projectRoot string
	index       *ProjectIndex
	ignorePatterns []string
}

// NewIndexer creates a new codebase indexer
func NewIndexer(projectRoot string) *Indexer {
	return &Indexer{
		projectRoot: projectRoot,
		index: &ProjectIndex{
			Root:      projectRoot,
			Files:     make(map[string]*FileIndex),
			SymbolMap: make(map[string][]Symbol),
		},
		ignorePatterns: []string{
			"node_modules", "vendor", ".git", ".sdd",
			"__pycache__", ".venv", "venv", "dist", "build",
		},
	}
}

// Index indexes the entire project
func (i *Indexer) Index() error {
	return filepath.Walk(i.projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		
		// Skip ignored directories
		if info.IsDir() {
			for _, pattern := range i.ignorePatterns {
				if info.Name() == pattern {
					return filepath.SkipDir
				}
			}
			return nil
		}
		
		// Only index source files
		lang := detectLanguage(path)
		if lang == "" {
			return nil
		}
		
		return i.indexFile(path, lang)
	})
}

// indexFile indexes a single file
func (i *Indexer) indexFile(path, lang string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	
	relPath, _ := filepath.Rel(i.projectRoot, path)
	
	fileIndex := &FileIndex{
		Path:     relPath,
		Language: lang,
		Size:     info.Size(),
		Modified: info.ModTime().Unix(),
	}
	
	// Parse symbols based on language
	switch lang {
	case "go":
		fileIndex.Symbols = parseGoSymbols(string(content), relPath)
		fileIndex.Imports = parseGoImports(string(content))
	case "javascript", "typescript":
		fileIndex.Symbols = parseJSSymbols(string(content), relPath)
		fileIndex.Imports = parseJSImports(string(content))
	case "python":
		fileIndex.Symbols = parsePythonSymbols(string(content), relPath)
		fileIndex.Imports = parsePythonImports(string(content))
	case "rust":
		fileIndex.Symbols = parseRustSymbols(string(content), relPath)
	}
	
	// Update index
	i.index.mu.Lock()
	i.index.Files[relPath] = fileIndex
	for _, sym := range fileIndex.Symbols {
		i.index.SymbolMap[sym.Name] = append(i.index.SymbolMap[sym.Name], sym)
	}
	i.index.mu.Unlock()
	
	return nil
}

// Search finds symbols matching the query
func (i *Indexer) Search(query string) []Symbol {
	i.index.mu.RLock()
	defer i.index.mu.RUnlock()
	
	var results []Symbol
	query = strings.ToLower(query)
	
	for name, symbols := range i.index.SymbolMap {
		if strings.Contains(strings.ToLower(name), query) {
			results = append(results, symbols...)
		}
	}
	
	return results
}

// GetFileSymbols returns all symbols in a file
func (i *Indexer) GetFileSymbols(path string) []Symbol {
	i.index.mu.RLock()
	defer i.index.mu.RUnlock()
	
	if file, ok := i.index.Files[path]; ok {
		return file.Symbols
	}
	return nil
}

// GetContext generates context string for AI
func (i *Indexer) GetContext(maxSize int) string {
	i.index.mu.RLock()
	defer i.index.mu.RUnlock()
	
	var context strings.Builder
	context.WriteString("## Project Structure\n\n")
	
	// Group by directory
	dirs := make(map[string][]string)
	for path := range i.index.Files {
		dir := filepath.Dir(path)
		dirs[dir] = append(dirs[dir], filepath.Base(path))
	}
	
	for dir, files := range dirs {
		context.WriteString(fmt.Sprintf("### %s/\n", dir))
		for _, f := range files {
			context.WriteString(fmt.Sprintf("- %s\n", f))
		}
		context.WriteString("\n")
	}
	
	context.WriteString("## Key Symbols\n\n")
	
	// Add top symbols
	count := 0
	for name, symbols := range i.index.SymbolMap {
		if count > 50 { // Limit symbols
			break
		}
		for _, sym := range symbols {
			if sym.Kind == "function" || sym.Kind == "class" || sym.Kind == "type" {
				if sym.Signature != "" {
					context.WriteString(fmt.Sprintf("- `%s` (%s) in %s\n", sym.Signature, sym.Kind, sym.File))
				} else {
					context.WriteString(fmt.Sprintf("- `%s` (%s) in %s\n", name, sym.Kind, sym.File))
				}
				count++
			}
		}
	}
	
	result := context.String()
	if len(result) > maxSize {
		result = result[:maxSize] + "\n...(truncated)"
	}
	
	return result
}

// GetStats returns indexing statistics
func (i *Indexer) GetStats() map[string]int {
	i.index.mu.RLock()
	defer i.index.mu.RUnlock()
	
	stats := map[string]int{
		"files":   len(i.index.Files),
		"symbols": 0,
	}
	
	for _, syms := range i.index.SymbolMap {
		stats["symbols"] += len(syms)
	}
	
	// Count by language
	for _, file := range i.index.Files {
		key := "lang_" + file.Language
		stats[key]++
	}
	
	return stats
}

// Helper functions

func detectLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go":
		return "go"
	case ".js", ".jsx", ".mjs":
		return "javascript"
	case ".ts", ".tsx":
		return "typescript"
	case ".py":
		return "python"
	case ".rs":
		return "rust"
	case ".java":
		return "java"
	case ".c", ".h":
		return "c"
	case ".cpp", ".hpp", ".cc":
		return "cpp"
	default:
		return ""
	}
}

// Regex patterns compiled once for performance
var (
	// Go patterns
	goFuncPattern   = regexp.MustCompile(`^func\s+(?:\((\w+)\s+\*?(\w+)\)\s+)?(\w+)\s*\(([^)]*)\)`)
	goTypePattern   = regexp.MustCompile(`^type\s+(\w+)\s+(struct|interface)`)
	goImportPattern = regexp.MustCompile(`import\s+(?:\(\s*([\s\S]*?)\s*\)|"([^"]+)")`)

	// JS/TS patterns
	jsFuncPattern   = regexp.MustCompile(`(?:function|const|let|var)\s+(\w+)\s*(?:=\s*(?:async\s*)?\([^)]*\)\s*=>|\([^)]*\))`)
	jsClassPattern  = regexp.MustCompile(`class\s+(\w+)`)
	jsExportPattern = regexp.MustCompile(`export\s+(?:default\s+)?(?:function|class|const|let|var)\s+(\w+)`)
	jsImportPattern = regexp.MustCompile(`(?:import|require)\s*\(?['"]([^'"]+)['"]`)

	// Python patterns
	pyFuncPattern  = regexp.MustCompile(`^(\s*)def\s+(\w+)\s*\(([^)]*)\)`)
	pyClassPattern = regexp.MustCompile(`^class\s+(\w+)`)

	// Rust patterns
	rsFuncPattern   = regexp.MustCompile(`(?:pub\s+)?fn\s+(\w+)`)
	rsStructPattern = regexp.MustCompile(`(?:pub\s+)?struct\s+(\w+)`)
	rsImplPattern   = regexp.MustCompile(`impl(?:<[^>]+>)?\s+(\w+)`)
)

func parseGoSymbols(content, file string) []Symbol {
	var symbols []Symbol
	lines := strings.Split(content, "\n")
	
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		
		if matches := goFuncPattern.FindStringSubmatch(line); matches != nil {
			sym := Symbol{
				Name:      matches[3],
				Kind:      "function",
				File:      file,
				Line:      lineNum + 1,
				Signature: line,
			}
			if matches[2] != "" {
				sym.Kind = "method"
				sym.Parent = matches[2]
			}
			symbols = append(symbols, sym)
		}
		
		if matches := goTypePattern.FindStringSubmatch(line); matches != nil {
			symbols = append(symbols, Symbol{
				Name: matches[1],
				Kind: matches[2],
				File: file,
				Line: lineNum + 1,
			})
		}
	}
	
	return symbols
}

func parseGoImports(content string) []string {
	var imports []string
	
	matches := goImportPattern.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if match[1] != "" {
			// Multi-import
			lines := strings.Split(match[1], "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "//") {
					// Extract path from quotes
					if idx := strings.Index(line, `"`); idx >= 0 {
						end := strings.LastIndex(line, `"`)
						if end > idx {
							imports = append(imports, line[idx+1:end])
						}
					}
				}
			}
		} else if match[2] != "" {
			imports = append(imports, match[2])
		}
	}
	
	return imports
}

func parseJSSymbols(content, file string) []Symbol {
	var symbols []Symbol
	lines := strings.Split(content, "\n")
	
	for lineNum, line := range lines {
		if matches := jsFuncPattern.FindStringSubmatch(line); matches != nil {
			symbols = append(symbols, Symbol{
				Name: matches[1],
				Kind: "function",
				File: file,
				Line: lineNum + 1,
			})
		}
		
		if matches := jsClassPattern.FindStringSubmatch(line); matches != nil {
			symbols = append(symbols, Symbol{
				Name: matches[1],
				Kind: "class",
				File: file,
				Line: lineNum + 1,
			})
		}
		
		if matches := jsExportPattern.FindStringSubmatch(line); matches != nil {
			// Already captured above, but mark as exported
			for i := range symbols {
				if symbols[i].Name == matches[1] && symbols[i].File == file {
					symbols[i].DocComment = "exported"
				}
			}
		}
	}
	
	return symbols
}

func parseJSImports(content string) []string {
	var imports []string
	
	matches := jsImportPattern.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		imports = append(imports, match[1])
	}
	
	return imports
}

func parsePythonSymbols(content, file string) []Symbol {
	var symbols []Symbol
	lines := strings.Split(content, "\n")
	
	var currentClass string
	
	for lineNum, line := range lines {
		if matches := pyClassPattern.FindStringSubmatch(line); matches != nil {
			currentClass = matches[1]
			symbols = append(symbols, Symbol{
				Name: matches[1],
				Kind: "class",
				File: file,
				Line: lineNum + 1,
			})
		}
		
		if matches := pyFuncPattern.FindStringSubmatch(line); matches != nil {
			indent := len(matches[1])
			sym := Symbol{
				Name:      matches[2],
				Kind:      "function",
				File:      file,
				Line:      lineNum + 1,
				Signature: fmt.Sprintf("def %s(%s)", matches[2], matches[3]),
			}
			
			// If indented under a class, it's a method
			if indent > 0 && currentClass != "" {
				sym.Kind = "method"
				sym.Parent = currentClass
			} else {
				currentClass = ""
			}
			
			symbols = append(symbols, sym)
		}
	}
	
	return symbols
}

func parsePythonImports(content string) []string {
	var imports []string
	
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "import ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				imports = append(imports, parts[1])
			}
		} else if strings.HasPrefix(line, "from ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				imports = append(imports, parts[1])
			}
		}
	}
	
	return imports
}

func parseRustSymbols(content, file string) []Symbol {
	var symbols []Symbol
	lines := strings.Split(content, "\n")
	
	var currentImpl string
	
	for lineNum, line := range lines {
		if matches := rsImplPattern.FindStringSubmatch(line); matches != nil {
			currentImpl = matches[1]
		}
		
		if matches := rsStructPattern.FindStringSubmatch(line); matches != nil {
			symbols = append(symbols, Symbol{
				Name: matches[1],
				Kind: "struct",
				File: file,
				Line: lineNum + 1,
			})
		}
		
		if matches := rsFuncPattern.FindStringSubmatch(line); matches != nil {
			sym := Symbol{
				Name: matches[1],
				Kind: "function",
				File: file,
				Line: lineNum + 1,
			}
			if currentImpl != "" && strings.Contains(line, "&self") {
				sym.Kind = "method"
				sym.Parent = currentImpl
			}
			symbols = append(symbols, sym)
		}
	}
	
	return symbols
}
