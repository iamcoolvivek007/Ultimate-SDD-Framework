package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// MCPConfig represents the Model Context Protocol configuration
type MCPConfig struct {
	Providers map[string]ProviderConfig `json:"providers"`
	DefaultProvider string               `json:"default_provider"`
}

// ProviderConfig represents configuration for a specific AI provider
type ProviderConfig struct {
	Provider ModelProvider `json:"provider"`
	APIKey   string        `json:"api_key"`
	BaseURL  string        `json:"base_url,omitempty"`
	Model    string        `json:"model"`
	Enabled  bool          `json:"enabled"`
}

// MCPManager manages MCP connections and configurations
type MCPManager struct {
	configPath string
	config     *MCPConfig
	clients    map[string]*ModelClient
}

// NewMCPManager creates a new MCP manager
func NewMCPManager(projectRoot string) *MCPManager {
	configPath := filepath.Join(projectRoot, ".sdd", "mcp.json")

	return &MCPManager{
		configPath: configPath,
		clients:    make(map[string]*ModelClient),
	}
}

// LoadConfig loads the MCP configuration from disk
func (m *MCPManager) LoadConfig() error {
	// Check if config exists
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		// Create default config
		m.config = &MCPConfig{
			Providers:       make(map[string]ProviderConfig),
			DefaultProvider: "",
		}
		return m.SaveConfig()
	}

	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return fmt.Errorf("failed to read MCP config: %w", err)
	}

	if err := json.Unmarshal(data, &m.config); err != nil {
		return fmt.Errorf("failed to parse MCP config: %w", err)
	}

	// Initialize clients for enabled providers
	for name, provider := range m.config.Providers {
		if provider.Enabled {
			client := NewModelClient(provider.Provider, provider.APIKey, provider.Model)
			if provider.BaseURL != "" {
				client.SetBaseURL(provider.BaseURL)
			}
			m.clients[name] = client
		}
	}

	return nil
}

// SaveConfig saves the MCP configuration to disk
func (m *MCPManager) SaveConfig() error {
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal MCP config: %w", err)
	}

	// Ensure .sdd directory exists
	dir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write MCP config: %w", err)
	}

	return nil
}

// AddProvider adds a new AI provider configuration
func (m *MCPManager) AddProvider(name string, provider ModelProvider, apiKey, model string, options map[string]interface{}) error {
	config := ProviderConfig{
		Provider: provider,
		APIKey:   apiKey,
		Model:    model,
		Enabled:  true,
	}

	// Apply additional options
	if baseURL, ok := options["base_url"].(string); ok {
		config.BaseURL = baseURL
	}

	m.config.Providers[name] = config

	// Create client
	client := NewModelClient(provider, apiKey, model)
	if config.BaseURL != "" {
		client.SetBaseURL(config.BaseURL)
	}
	m.clients[name] = client

	// Set as default if it's the first provider
	if m.config.DefaultProvider == "" {
		m.config.DefaultProvider = name
	}

	return m.SaveConfig()
}

// RemoveProvider removes a provider configuration
func (m *MCPManager) RemoveProvider(name string) error {
	if _, exists := m.config.Providers[name]; !exists {
		return fmt.Errorf("provider '%s' not found", name)
	}

	delete(m.config.Providers, name)
	delete(m.clients, name)

	// Update default provider if necessary
	if m.config.DefaultProvider == name {
		m.config.DefaultProvider = ""
		for name := range m.config.Providers {
			m.config.DefaultProvider = name
			break
		}
	}

	return m.SaveConfig()
}

// SetDefaultProvider sets the default provider
func (m *MCPManager) SetDefaultProvider(name string) error {
	if _, exists := m.config.Providers[name]; !exists {
		return fmt.Errorf("provider '%s' not found", name)
	}

	m.config.DefaultProvider = name
	return m.SaveConfig()
}

func (m *MCPManager) GetDefaultProvider() string {
	return m.config.DefaultProvider
}

// GetClient returns a model client for the specified provider
func (m *MCPManager) GetClient(providerName string) (*ModelClient, error) {
	if providerName == "" {
		providerName = m.config.DefaultProvider
	}

	client, exists := m.clients[providerName]
	if !exists {
		return nil, fmt.Errorf("provider '%s' not configured or disabled", providerName)
	}

	return client, nil
}

// ListProviders returns a list of configured providers
func (m *MCPManager) ListProviders() map[string]ProviderConfig {
	return m.config.Providers
}

// ValidateProvider tests a provider configuration
func (m *MCPManager) ValidateProvider(name string) error {
	client, err := m.GetClient(name)
	if err != nil {
		return err
	}

	return client.ValidateConnection()
}

// ChatWithProvider sends a chat request to a specific provider
func (m *MCPManager) ChatWithProvider(providerName string, messages []Message, options map[string]interface{}) (*ChatResponse, error) {
	client, err := m.GetClient(providerName)
	if err != nil {
		return nil, err
	}

	return client.Chat(messages, options)
}

// Chat sends a chat request to the default provider
func (m *MCPManager) Chat(messages []Message, options map[string]interface{}) (*ChatResponse, error) {
	return m.ChatWithProvider("", messages, options)
}

// GetAvailableProviders returns a list of supported providers
func GetAvailableProviders() []ModelProvider {
	return []ModelProvider{
		ProviderOpenAI,
		ProviderAnthropic,
		ProviderGoogle,
		ProviderOllama,
		ProviderAzure,
	}
}

// GetProviderDisplayName returns a human-readable name for a provider
func GetProviderDisplayName(provider ModelProvider) string {
	switch provider {
	case ProviderOpenAI:
		return "OpenAI"
	case ProviderAnthropic:
		return "Anthropic"
	case ProviderGoogle:
		return "Google Gemini"
	case ProviderOllama:
		return "Ollama (Local)"
	case ProviderAzure:
		return "Azure OpenAI"
	default:
		return string(provider)
	}
}

// GetDefaultModelForProvider returns the default model for a provider
func GetDefaultModelForProvider(provider ModelProvider) string {
	switch provider {
	case ProviderOpenAI:
		return "gpt-4"
	case ProviderAnthropic:
		return "claude-3-sonnet-20240229"
	case ProviderGoogle:
		return "gemini-pro"
	case ProviderOllama:
		return "llama2"
	case ProviderAzure:
		return "gpt-4"
	default:
		return ""
	}
}