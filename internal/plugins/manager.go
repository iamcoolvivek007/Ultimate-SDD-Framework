package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/goccy/go-yaml"
)

// Plugin represents a Viki plugin
type Plugin interface {
	// Name returns the plugin name
	Name() string
	// Version returns the plugin version
	Version() string
	// Description returns the plugin description
	Description() string
	// Initialize initializes the plugin
	Initialize(ctx PluginContext) error
	// Execute runs the plugin
	Execute(args []string) error
}

// PluginContext provides context to plugins
type PluginContext struct {
	WorkDir    string
	ConfigDir  string
	SDDDir     string
	AIProvider AIProviderInterface
}

// AIProviderInterface is the interface plugins use to access AI
type AIProviderInterface interface {
	Chat(messages []Message, options map[string]interface{}) (string, error)
	Stream(messages []Message, callback func(string)) error
}

// Message is a simplified message for plugins
type Message struct {
	Role    string
	Content string
}

// PluginManifest describes a plugin's metadata
type PluginManifest struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	Description string   `yaml:"description"`
	Author      string   `yaml:"author"`
	Entry       string   `yaml:"entry"`    // Entry point file
	Type        string   `yaml:"type"`     // "go", "script", "config"
	Commands    []string `yaml:"commands"` // Commands this plugin adds
	Hooks       []string `yaml:"hooks"`    // Hooks this plugin listens to
	Agents      []string `yaml:"agents"`   // Custom agents this plugin provides
}

// PluginInfo contains loaded plugin information
type PluginInfo struct {
	Manifest PluginManifest
	Path     string
	Loaded   bool
	Instance Plugin
}

// PluginManager manages Viki plugins
type PluginManager struct {
	pluginsDir string
	plugins    map[string]*PluginInfo
	context    PluginContext
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(pluginsDir string) *PluginManager {
	return &PluginManager{
		pluginsDir: pluginsDir,
		plugins:    make(map[string]*PluginInfo),
	}
}

// SetContext sets the plugin context
func (pm *PluginManager) SetContext(ctx PluginContext) {
	pm.context = ctx
}

// Discover discovers all plugins in the plugins directory
func (pm *PluginManager) Discover() error {
	if err := os.MkdirAll(pm.pluginsDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugins directory: %w", err)
	}

	entries, err := os.ReadDir(pm.pluginsDir)
	if err != nil {
		return fmt.Errorf("failed to read plugins directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginPath := filepath.Join(pm.pluginsDir, entry.Name())
		manifestPath := filepath.Join(pluginPath, "plugin.yaml")

		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			continue
		}

		manifest, err := pm.loadManifest(manifestPath)
		if err != nil {
			fmt.Printf("Warning: failed to load plugin %s: %v\n", entry.Name(), err)
			continue
		}

		pm.plugins[manifest.Name] = &PluginInfo{
			Manifest: *manifest,
			Path:     pluginPath,
			Loaded:   false,
		}
	}

	return nil
}

// loadManifest loads a plugin manifest
func (pm *PluginManager) loadManifest(path string) (*PluginManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest PluginManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// Load loads a plugin by name
func (pm *PluginManager) Load(name string) error {
	info, ok := pm.plugins[name]
	if !ok {
		return fmt.Errorf("plugin not found: %s", name)
	}

	if info.Loaded {
		return nil
	}

	switch info.Manifest.Type {
	case "go":
		return pm.loadGoPlugin(info)
	case "script":
		return pm.loadScriptPlugin(info)
	case "config":
		// Config-only plugins don't need loading
		info.Loaded = true
		return nil
	default:
		return fmt.Errorf("unknown plugin type: %s", info.Manifest.Type)
	}
}

// loadGoPlugin loads a Go plugin (.so file)
func (pm *PluginManager) loadGoPlugin(info *PluginInfo) error {
	pluginPath := filepath.Join(info.Path, info.Manifest.Entry)

	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %w", err)
	}

	sym, err := p.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin missing 'Plugin' symbol: %w", err)
	}

	plug, ok := sym.(Plugin)
	if !ok {
		return fmt.Errorf("plugin 'Plugin' symbol has wrong type")
	}

	if err := plug.Initialize(pm.context); err != nil {
		return fmt.Errorf("failed to initialize plugin: %w", err)
	}

	info.Instance = plug
	info.Loaded = true

	return nil
}

// loadScriptPlugin loads a script-based plugin
func (pm *PluginManager) loadScriptPlugin(info *PluginInfo) error {
	// Script plugins are executed on-demand
	info.Loaded = true
	return nil
}

// Execute executes a plugin command
func (pm *PluginManager) Execute(name string, args []string) error {
	info, ok := pm.plugins[name]
	if !ok {
		return fmt.Errorf("plugin not found: %s", name)
	}

	if !info.Loaded {
		if err := pm.Load(name); err != nil {
			return err
		}
	}

	if info.Instance != nil {
		return info.Instance.Execute(args)
	}

	if info.Manifest.Type == "script" {
		return pm.executeScript(info, args)
	}

	return fmt.Errorf("plugin cannot be executed")
}

// executeScript executes a script-based plugin
func (pm *PluginManager) executeScript(info *PluginInfo, args []string) error {
	scriptPath := filepath.Join(info.Path, info.Manifest.Entry)

	// Determine script runner
	ext := filepath.Ext(info.Manifest.Entry)
	var runner string
	switch ext {
	case ".sh":
		runner = "bash"
	case ".py":
		runner = "python3"
	case ".js":
		runner = "node"
	default:
		return fmt.Errorf("unsupported script type: %s", ext)
	}

	// Execute script
	cmd := fmt.Sprintf("%s %s %s", runner, scriptPath, strings.Join(args, " "))
	fmt.Printf("Executing: %s\n", cmd)
	// In real implementation, use os/exec
	return nil
}

// List returns all discovered plugins
func (pm *PluginManager) List() []*PluginInfo {
	var plugins []*PluginInfo
	for _, info := range pm.plugins {
		plugins = append(plugins, info)
	}
	return plugins
}

// GetCommands returns all commands provided by plugins
func (pm *PluginManager) GetCommands() map[string]string {
	commands := make(map[string]string)
	for _, info := range pm.plugins {
		for _, cmd := range info.Manifest.Commands {
			commands[cmd] = info.Manifest.Name
		}
	}
	return commands
}

// GetAgents returns all agents provided by plugins
func (pm *PluginManager) GetAgents() map[string]string {
	agents := make(map[string]string)
	for _, info := range pm.plugins {
		for _, agent := range info.Manifest.Agents {
			agents[agent] = filepath.Join(info.Path, "agents", agent+".md")
		}
	}
	return agents
}

// Install installs a plugin from a URL or local path
func (pm *PluginManager) Install(source string) error {
	// Placeholder for plugin installation
	// Would handle:
	// - GitHub URLs
	// - Local paths
	// - ZIP archives
	return fmt.Errorf("plugin installation not yet implemented")
}

// Uninstall removes a plugin
func (pm *PluginManager) Uninstall(name string) error {
	info, ok := pm.plugins[name]
	if !ok {
		return fmt.Errorf("plugin not found: %s", name)
	}

	if err := os.RemoveAll(info.Path); err != nil {
		return fmt.Errorf("failed to remove plugin: %w", err)
	}

	delete(pm.plugins, name)
	return nil
}

// CreatePluginTemplate creates a new plugin template
func (pm *PluginManager) CreatePluginTemplate(name, pluginType string) error {
	pluginDir := filepath.Join(pm.pluginsDir, name)

	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return err
	}

	// Create manifest
	manifest := PluginManifest{
		Name:        name,
		Version:     "1.0.0",
		Description: "A custom Viki plugin",
		Author:      "Your Name",
		Type:        pluginType,
	}

	switch pluginType {
	case "script":
		manifest.Entry = "main.sh"
		// Create shell script
		script := "#!/bin/bash\necho \"Hello from " + name + " plugin!\"\n"
		if err := os.WriteFile(filepath.Join(pluginDir, "main.sh"), []byte(script), 0755); err != nil {
			return err
		}
	case "config":
		// Create agents directory
		agentsDir := filepath.Join(pluginDir, "agents")
		if err := os.MkdirAll(agentsDir, 0755); err != nil {
			return err
		}
		manifest.Agents = []string{"custom-agent"}
		// Create custom agent
		agent := `---
role: Custom Agent
expertise: Your expertise here
personality: Helpful and friendly
---

# Custom Agent

Describe your custom agent's behavior here.
`
		if err := os.WriteFile(filepath.Join(agentsDir, "custom-agent.md"), []byte(agent), 0644); err != nil {
			return err
		}
	}

	// Write manifest
	manifestData, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}

	header := "# Plugin Manifest\n# Documentation: https://viki.dev/plugins\n\n"
	return os.WriteFile(filepath.Join(pluginDir, "plugin.yaml"), []byte(header+string(manifestData)), 0644)
}
