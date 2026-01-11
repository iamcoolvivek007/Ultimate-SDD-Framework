package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/goccy/go-yaml"
)

// Config represents the global Viki configuration
type Config struct {
	// Default AI provider
	DefaultProvider string `yaml:"default_provider"`

	// Theme settings
	Theme ThemeConfig `yaml:"theme"`

	// Editor settings
	Editor EditorConfig `yaml:"editor"`

	// AI settings
	AI AIConfig `yaml:"ai"`

	// Project defaults
	ProjectDefaults ProjectConfig `yaml:"project_defaults"`

	// Telemetry settings
	Telemetry TelemetryConfig `yaml:"telemetry"`
}

// ThemeConfig represents theme settings
type ThemeConfig struct {
	ColorScheme string `yaml:"color_scheme"` // "dark", "light", "auto"
	Accent      string `yaml:"accent"`       // accent color
	Emoji       bool   `yaml:"emoji"`        // use emojis
}

// EditorConfig represents editor settings
type EditorConfig struct {
	Command    string `yaml:"command"`     // editor command (vim, code, etc.)
	AutoFormat bool   `yaml:"auto_format"` // auto-format on save
	TabSize    int    `yaml:"tab_size"`
}

// AIConfig represents AI settings
type AIConfig struct {
	Temperature     float64 `yaml:"temperature"`
	MaxTokens       int     `yaml:"max_tokens"`
	StreamResponses bool    `yaml:"stream_responses"`
	AutoApprove     bool    `yaml:"auto_approve"` // Skip approval gates
}

// ProjectConfig represents project defaults
type ProjectConfig struct {
	Language   string   `yaml:"language"`    // Default language (go, python, etc.)
	Framework  string   `yaml:"framework"`   // Default framework
	TestRunner string   `yaml:"test_runner"` // Default test runner
	Agents     []string `yaml:"agents"`      // Default agents to load
}

// TelemetryConfig represents telemetry settings
type TelemetryConfig struct {
	Enabled   bool `yaml:"enabled"`
	Anonymous bool `yaml:"anonymous"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		DefaultProvider: "",
		Theme: ThemeConfig{
			ColorScheme: "dark",
			Accent:      "39", // cyan
			Emoji:       true,
		},
		Editor: EditorConfig{
			Command:    getDefaultEditor(),
			AutoFormat: true,
			TabSize:    4,
		},
		AI: AIConfig{
			Temperature:     0.7,
			MaxTokens:       4000,
			StreamResponses: true,
			AutoApprove:     false,
		},
		ProjectDefaults: ProjectConfig{
			Language:   "go",
			Framework:  "",
			TestRunner: "go test",
			Agents:     []string{"strategist", "architect", "developer", "qa"},
		},
		Telemetry: TelemetryConfig{
			Enabled:   false,
			Anonymous: true,
		},
	}
}

// ConfigManager handles loading and saving configuration
type ConfigManager struct {
	configDir  string
	configFile string
	config     *Config
}

// NewConfigManager creates a new config manager
func NewConfigManager() *ConfigManager {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".viki")

	return &ConfigManager{
		configDir:  configDir,
		configFile: filepath.Join(configDir, "config.yaml"),
		config:     DefaultConfig(),
	}
}

// GetConfigDir returns the config directory path
func (cm *ConfigManager) GetConfigDir() string {
	return cm.configDir
}

// Load loads the configuration from disk
func (cm *ConfigManager) Load() error {
	// Ensure config directory exists
	if err := os.MkdirAll(cm.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(cm.configFile); os.IsNotExist(err) {
		// Create default config
		return cm.Save()
	}

	// Read config file
	data, err := os.ReadFile(cm.configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, cm.config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}

// Save saves the configuration to disk
func (cm *ConfigManager) Save() error {
	// Ensure config directory exists
	if err := os.MkdirAll(cm.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(cm.config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Add header comment
	header := `# Viki Global Configuration
# Location: ~/.viki/config.yaml
# Documentation: https://github.com/viki-dev/viki#configuration

`

	// Write to file
	if err := os.WriteFile(cm.configFile, []byte(header+string(data)), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Get returns the current configuration
func (cm *ConfigManager) Get() *Config {
	return cm.config
}

// Set updates the configuration
func (cm *ConfigManager) Set(config *Config) {
	cm.config = config
}

// GetValue gets a configuration value by key path
func (cm *ConfigManager) GetValue(key string) (interface{}, error) {
	switch key {
	case "default_provider":
		return cm.config.DefaultProvider, nil
	case "theme.color_scheme":
		return cm.config.Theme.ColorScheme, nil
	case "theme.accent":
		return cm.config.Theme.Accent, nil
	case "theme.emoji":
		return cm.config.Theme.Emoji, nil
	case "editor.command":
		return cm.config.Editor.Command, nil
	case "ai.temperature":
		return cm.config.AI.Temperature, nil
	case "ai.max_tokens":
		return cm.config.AI.MaxTokens, nil
	case "ai.stream_responses":
		return cm.config.AI.StreamResponses, nil
	default:
		return nil, fmt.Errorf("unknown config key: %s", key)
	}
}

// SetValue sets a configuration value by key path
func (cm *ConfigManager) SetValue(key string, value interface{}) error {
	switch key {
	case "default_provider":
		cm.config.DefaultProvider = value.(string)
	case "theme.color_scheme":
		cm.config.Theme.ColorScheme = value.(string)
	case "theme.accent":
		cm.config.Theme.Accent = value.(string)
	case "theme.emoji":
		cm.config.Theme.Emoji = value.(bool)
	case "editor.command":
		cm.config.Editor.Command = value.(string)
	case "ai.temperature":
		cm.config.AI.Temperature = value.(float64)
	case "ai.max_tokens":
		cm.config.AI.MaxTokens = value.(int)
	case "ai.stream_responses":
		cm.config.AI.StreamResponses = value.(bool)
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	return cm.Save()
}

// Reset resets the configuration to defaults
func (cm *ConfigManager) Reset() error {
	cm.config = DefaultConfig()
	return cm.Save()
}

// MergeWithProject merges global config with project-level config
func (cm *ConfigManager) MergeWithProject(projectConfig map[string]interface{}) *Config {
	merged := *cm.config

	// Override with project-level settings
	if provider, ok := projectConfig["default_provider"].(string); ok && provider != "" {
		merged.DefaultProvider = provider
	}

	if ai, ok := projectConfig["ai"].(map[string]interface{}); ok {
		if temp, ok := ai["temperature"].(float64); ok {
			merged.AI.Temperature = temp
		}
		if tokens, ok := ai["max_tokens"].(int); ok {
			merged.AI.MaxTokens = tokens
		}
	}

	return &merged
}

// getDefaultEditor returns the default editor based on OS
func getDefaultEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}

	switch runtime.GOOS {
	case "windows":
		return "notepad"
	case "darwin":
		return "nano"
	default:
		return "nano"
	}
}

// ListAllKeys returns all available config keys
func ListAllKeys() []string {
	return []string{
		"default_provider",
		"theme.color_scheme",
		"theme.accent",
		"theme.emoji",
		"editor.command",
		"editor.auto_format",
		"editor.tab_size",
		"ai.temperature",
		"ai.max_tokens",
		"ai.stream_responses",
		"ai.auto_approve",
		"project_defaults.language",
		"project_defaults.framework",
		"project_defaults.test_runner",
		"telemetry.enabled",
		"telemetry.anonymous",
	}
}
