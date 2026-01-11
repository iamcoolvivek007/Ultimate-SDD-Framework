package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed builtin/*
var builtinTemplates embed.FS

// Template represents a project template
type Template struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Language    string            `yaml:"language"`
	Framework   string            `yaml:"framework"`
	Files       map[string]string `yaml:"files"`
	Variables   []TemplateVar     `yaml:"variables"`
	PostCreate  []string          `yaml:"post_create"`
}

// TemplateVar represents a template variable
type TemplateVar struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
	Required    bool   `yaml:"required"`
}

// TemplateManager manages project templates
type TemplateManager struct {
	templatesDir string
	templates    map[string]*Template
}

// NewTemplateManager creates a new template manager
func NewTemplateManager(templatesDir string) *TemplateManager {
	return &TemplateManager{
		templatesDir: templatesDir,
		templates:    make(map[string]*Template),
	}
}

// LoadBuiltinTemplates loads the built-in templates
func (tm *TemplateManager) LoadBuiltinTemplates() {
	// Go API template
	tm.templates["go-api"] = &Template{
		Name:        "go-api",
		Description: "Go REST API with Fiber",
		Language:    "go",
		Framework:   "fiber",
		Files: map[string]string{
			"main.go":        goMainTemplate,
			"go.mod":         goModTemplate,
			"README.md":      readmeTemplate,
			".gitignore":     gitignoreTemplate,
			"Makefile":       makefileTemplate,
			"api/handler.go": goHandlerTemplate,
		},
		Variables: []TemplateVar{
			{Name: "ProjectName", Description: "Project name", Default: "myapp", Required: true},
			{Name: "ModulePath", Description: "Go module path", Default: "github.com/user/myapp", Required: true},
		},
		PostCreate: []string{
			"go mod tidy",
		},
	}

	// React App template
	tm.templates["react-app"] = &Template{
		Name:        "react-app",
		Description: "React + TypeScript + Vite",
		Language:    "typescript",
		Framework:   "react",
		Files: map[string]string{
			"package.json":   reactPackageTemplate,
			"vite.config.ts": viteConfigTemplate,
			"tsconfig.json":  tsconfigTemplate,
			"index.html":     indexHtmlTemplate,
			"src/App.tsx":    reactAppTemplate,
			"src/main.tsx":   reactMainTemplate,
			"src/index.css":  reactCssTemplate,
			".gitignore":     gitignoreTemplate,
			"README.md":      readmeTemplate,
		},
		Variables: []TemplateVar{
			{Name: "ProjectName", Description: "Project name", Default: "my-react-app", Required: true},
		},
		PostCreate: []string{
			"npm install",
		},
	}

	// Python FastAPI template
	tm.templates["python-api"] = &Template{
		Name:        "python-api",
		Description: "Python FastAPI",
		Language:    "python",
		Framework:   "fastapi",
		Files: map[string]string{
			"main.py":          pythonMainTemplate,
			"requirements.txt": pythonRequirementsTemplate,
			"README.md":        readmeTemplate,
			".gitignore":       pythonGitignoreTemplate,
		},
		Variables: []TemplateVar{
			{Name: "ProjectName", Description: "Project name", Default: "myapi", Required: true},
		},
		PostCreate: []string{
			"python -m venv venv",
			"pip install -r requirements.txt",
		},
	}

	// Next.js template
	tm.templates["nextjs"] = &Template{
		Name:        "nextjs",
		Description: "Next.js 14 with App Router",
		Language:    "typescript",
		Framework:   "nextjs",
		Files: map[string]string{
			"package.json":    nextPackageTemplate,
			"next.config.js":  nextConfigTemplate,
			"tsconfig.json":   tsconfigTemplate,
			"app/layout.tsx":  nextLayoutTemplate,
			"app/page.tsx":    nextPageTemplate,
			"app/globals.css": nextCssTemplate,
			".gitignore":      gitignoreTemplate,
			"README.md":       readmeTemplate,
		},
		Variables: []TemplateVar{
			{Name: "ProjectName", Description: "Project name", Default: "my-next-app", Required: true},
		},
		PostCreate: []string{
			"npm install",
		},
	}

	// CLI template
	tm.templates["go-cli"] = &Template{
		Name:        "go-cli",
		Description: "Go CLI with Cobra",
		Language:    "go",
		Framework:   "cobra",
		Files: map[string]string{
			"main.go":     goCliMainTemplate,
			"cmd/root.go": goCliRootTemplate,
			"go.mod":      goModTemplate,
			"README.md":   readmeTemplate,
			".gitignore":  gitignoreTemplate,
		},
		Variables: []TemplateVar{
			{Name: "ProjectName", Description: "CLI name", Default: "mycli", Required: true},
			{Name: "ModulePath", Description: "Go module path", Default: "github.com/user/mycli", Required: true},
		},
		PostCreate: []string{
			"go mod tidy",
		},
	}
}

// List returns all available templates
func (tm *TemplateManager) List() []*Template {
	var templates []*Template
	for _, t := range tm.templates {
		templates = append(templates, t)
	}
	return templates
}

// Get returns a template by name
func (tm *TemplateManager) Get(name string) (*Template, error) {
	t, ok := tm.templates[name]
	if !ok {
		return nil, fmt.Errorf("template not found: %s", name)
	}
	return t, nil
}

// Create creates a new project from a template
func (tm *TemplateManager) Create(templateName, targetDir string, vars map[string]string) error {
	t, err := tm.Get(templateName)
	if err != nil {
		return err
	}

	// Create target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Process each template file
	for path, content := range t.Files {
		targetPath := filepath.Join(targetDir, path)

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", path, err)
		}

		// Process template
		processed, err := processTemplate(content, vars)
		if err != nil {
			return fmt.Errorf("failed to process template %s: %w", path, err)
		}

		// Write file
		if err := os.WriteFile(targetPath, []byte(processed), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", path, err)
		}

		fmt.Printf("  Created: %s\n", path)
	}

	// Initialize .sdd directory
	sddDir := filepath.Join(targetDir, ".sdd")
	if err := os.MkdirAll(sddDir, 0755); err != nil {
		return err
	}

	return nil
}

// processTemplate processes a template string with variables
func processTemplate(content string, vars map[string]string) (string, error) {
	tmpl, err := template.New("file").Parse(content)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// LoadCustomTemplates loads custom templates from a directory
func (tm *TemplateManager) LoadCustomTemplates() error {
	if _, err := os.Stat(tm.templatesDir); os.IsNotExist(err) {
		return nil
	}

	entries, err := os.ReadDir(tm.templatesDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		templateDir := filepath.Join(tm.templatesDir, entry.Name())
		// Would load template.yaml and files here
		_ = templateDir
	}

	return nil
}

// Template content strings

const goMainTemplate = `package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "{{.ProjectName}}",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to {{.ProjectName}}!",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	log.Fatal(app.Listen(":3000"))
}
`

const goModTemplate = `module {{.ModulePath}}

go 1.21

require github.com/gofiber/fiber/v2 v2.52.0
`

const goHandlerTemplate = `package api

import "github.com/gofiber/fiber/v2"

func HelloHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}
`

const readmeTemplate = `# {{.ProjectName}}

A project created with Viki.

## Getting Started

` + "```bash\n# Run the project\nmake run\n```" + `

## License

MIT
`

const gitignoreTemplate = `.DS_Store
*.log
.env
.env.local
node_modules/
dist/
build/
.vscode/
.idea/
*.exe
*.dll
*.so
*.dylib
`

const makefileTemplate = `.PHONY: build run test clean

build:
	go build -o bin/{{.ProjectName}} .

run:
	go run .

test:
	go test -v ./...

clean:
	rm -rf bin/
`

const reactPackageTemplate = `{
  "name": "{{.ProjectName}}",
  "private": true,
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.0",
    "@types/react-dom": "^18.2.0",
    "@vitejs/plugin-react": "^4.2.0",
    "typescript": "^5.2.0",
    "vite": "^5.0.0"
  }
}
`

const viteConfigTemplate = `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
})
`

const tsconfigTemplate = `{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true
  },
  "include": ["src"]
}
`

const indexHtmlTemplate = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{.ProjectName}}</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
`

const reactAppTemplate = `import './index.css'

function App() {
  return (
    <div className="app">
      <h1>Welcome to {{.ProjectName}}</h1>
      <p>Created with Viki 🤖</p>
    </div>
  )
}

export default App
`

const reactMainTemplate = `import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
`

const reactCssTemplate = `:root {
  font-family: Inter, system-ui, sans-serif;
  line-height: 1.5;
  color: #213547;
  background-color: #ffffff;
}

.app {
  max-width: 1280px;
  margin: 0 auto;
  padding: 2rem;
  text-align: center;
}

h1 {
  font-size: 3.2em;
  line-height: 1.1;
}
`

const pythonMainTemplate = `from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(title="{{.ProjectName}}")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
async def root():
    return {"message": "Welcome to {{.ProjectName}}!"}

@app.get("/health")
async def health():
    return {"status": "healthy"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
`

const pythonRequirementsTemplate = `fastapi>=0.109.0
uvicorn>=0.27.0
pydantic>=2.0.0
`

const pythonGitignoreTemplate = `.DS_Store
*.pyc
__pycache__/
.env
.venv/
venv/
*.egg-info/
dist/
build/
.pytest_cache/
.coverage
`

const nextPackageTemplate = `{
  "name": "{{.ProjectName}}",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint"
  },
  "dependencies": {
    "next": "14.1.0",
    "react": "^18",
    "react-dom": "^18"
  },
  "devDependencies": {
    "@types/node": "^20",
    "@types/react": "^18",
    "@types/react-dom": "^18",
    "typescript": "^5"
  }
}
`

const nextConfigTemplate = `/** @type {import('next').NextConfig} */
const nextConfig = {}

module.exports = nextConfig
`

const nextLayoutTemplate = `export const metadata = {
  title: '{{.ProjectName}}',
  description: 'Created with Viki',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}
`

const nextPageTemplate = `export default function Home() {
  return (
    <main style={{ padding: '2rem', textAlign: 'center' }}>
      <h1>Welcome to {{.ProjectName}}</h1>
      <p>Created with Viki 🤖</p>
    </main>
  )
}
`

const nextCssTemplate = `* {
  box-sizing: border-box;
  padding: 0;
  margin: 0;
}

html, body {
  max-width: 100vw;
  overflow-x: hidden;
}
`

const goCliMainTemplate = `package main

import "{{.ModulePath}}/cmd"

func main() {
	cmd.Execute()
}
`

const goCliRootTemplate = `package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "{{.ProjectName}}",
	Short: "{{.ProjectName}} - A CLI application",
	Long:  "{{.ProjectName}} is a command-line application created with Viki.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from {{.ProjectName}}!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
}
`

// Suppress unused warning
var _ = builtinTemplates
var _ fs.FS = builtinTemplates
