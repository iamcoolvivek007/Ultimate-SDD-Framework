package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ModelProvider represents different AI model providers
type ModelProvider string

const (
	ProviderOpenAI     ModelProvider = "openai"
	ProviderAnthropic  ModelProvider = "anthropic"
	ProviderGoogle     ModelProvider = "google"
	ProviderOllama     ModelProvider = "ollama"
	ProviderAzure      ModelProvider = "azure"
)

// ModelClient handles communication with AI models via API keys
type ModelClient struct {
	Provider    ModelProvider
	APIKey      string
	BaseURL     string
	Model       string
	httpClient  *http.Client
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a request to the AI model
type ChatRequest struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// ChatResponse represents a response from the AI model
type ChatResponse struct {
	Choices []struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// NewModelClient creates a new AI model client
func NewModelClient(provider ModelProvider, apiKey, model string) *ModelClient {
	client := &ModelClient{
		Provider:   provider,
		APIKey:     apiKey,
		Model:      model,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}

	// Set default base URLs
	switch provider {
	case ProviderOpenAI:
		client.BaseURL = "https://api.openai.com/v1"
	case ProviderAnthropic:
		client.BaseURL = "https://api.anthropic.com/v1"
	case ProviderGoogle:
		client.BaseURL = "https://generativelanguage.googleapis.com/v1beta"
	case ProviderOllama:
		client.BaseURL = "http://localhost:11434"
	case ProviderAzure:
		// Azure OpenAI requires custom base URL
		client.BaseURL = ""
	}

	return client
}

// SetBaseURL allows overriding the default base URL
func (mc *ModelClient) SetBaseURL(url string) {
	mc.BaseURL = url
}

// Chat sends a chat request to the AI model
func (mc *ModelClient) Chat(messages []Message, options map[string]interface{}) (*ChatResponse, error) {
	var request ChatRequest
	var endpoint string
	var headers map[string]string

	// Prepare request based on provider
	switch mc.Provider {
	case ProviderOpenAI, ProviderAzure:
		request = ChatRequest{
			Messages: messages,
			Model:    mc.Model,
		}
		endpoint = "/chat/completions"
		headers = map[string]string{
			"Authorization": "Bearer " + mc.APIKey,
			"Content-Type":  "application/json",
		}

	case ProviderAnthropic:
		// Anthropic uses a different message format
		systemMessage := ""
		userMessages := []Message{}

		for _, msg := range messages {
			if msg.Role == "system" {
				systemMessage = msg.Content
			} else {
				userMessages = append(userMessages, msg)
			}
		}

		requestBody := map[string]interface{}{
			"model":      mc.Model,
			"messages":   userMessages,
			"max_tokens": 4096,
		}
		if systemMessage != "" {
			requestBody["system"] = systemMessage
		}

		return mc.sendAnthropicRequest(requestBody, headers)

	case ProviderGoogle:
		// Google Gemini format
		requestBody := map[string]interface{}{
			"contents": []map[string]interface{}{
				{"parts": []map[string]interface{}{
					{"text": messages[len(messages)-1].Content},
				}},
			},
		}
		endpoint = fmt.Sprintf("/models/%s:generateContent", mc.Model)
		headers = map[string]string{
			"Content-Type": "application/json",
		}
		if strings.Contains(mc.BaseURL, "generativelanguage.googleapis.com") {
			headers["x-goog-api-key"] = mc.APIKey
		}

		return mc.sendGoogleRequest(requestBody, endpoint, headers)

	case ProviderOllama:
		request = ChatRequest{
			Messages: messages,
			Model:    mc.Model,
			Stream:   false,
		}
		endpoint = "/api/chat"
		headers = map[string]string{
			"Content-Type": "application/json",
		}

	default:
		return nil, fmt.Errorf("unsupported provider: %s", mc.Provider)
	}

	// Apply options
	if temp, ok := options["temperature"].(float64); ok {
		request.Temperature = temp
	}
	if maxTokens, ok := options["max_tokens"].(int); ok {
		request.MaxTokens = maxTokens
	}

	return mc.sendRequest(request, endpoint, headers)
}

// sendRequest sends a generic HTTP request
func (mc *ModelClient) sendRequest(request interface{}, endpoint string, headers map[string]string) (*ChatResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := mc.BaseURL + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// sendAnthropicRequest handles Anthropic's specific API format
func (mc *ModelClient) sendAnthropicRequest(requestBody map[string]interface{}, headers map[string]string) (*ChatResponse, error) {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := mc.BaseURL + "/messages"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", mc.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse Anthropic response format
	var anthropicResp struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&anthropicResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to standard format
	response := &ChatResponse{
		Choices: []struct {
			Message      Message `json:"message"`
			FinishReason string  `json:"finish_reason"`
		}{
			{
				Message: Message{
					Role:    "assistant",
					Content: anthropicResp.Content[0].Text,
				},
				FinishReason: "stop",
			},
		},
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     anthropicResp.Usage.InputTokens,
			CompletionTokens: anthropicResp.Usage.OutputTokens,
			TotalTokens:      anthropicResp.Usage.InputTokens + anthropicResp.Usage.OutputTokens,
		},
	}

	return response, nil
}

// sendGoogleRequest handles Google's Gemini API format
func (mc *ModelClient) sendGoogleRequest(requestBody map[string]interface{}, endpoint string, headers map[string]string) (*ChatResponse, error) {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := mc.BaseURL + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to standard format
	content := ""
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		content = geminiResp.Candidates[0].Content.Parts[0].Text
	}

	response := &ChatResponse{
		Choices: []struct {
			Message      Message `json:"message"`
			FinishReason string  `json:"finish_reason"`
		}{
			{
				Message: Message{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: "stop",
			},
		},
	}

	return response, nil
}

// ValidateConnection tests the API key and connection
func (mc *ModelClient) ValidateConnection() error {
	// Send a simple test message
	testMessages := []Message{
		{Role: "user", Content: "Hello, this is a test message. Please respond with 'OK'."},
	}

	options := map[string]interface{}{
		"temperature": 0.1,
		"max_tokens":  10,
	}

	_, err := mc.Chat(testMessages, options)
	return err
}

// GetAvailableModels returns available models for the provider
func (mc *ModelClient) GetAvailableModels() ([]string, error) {
	switch mc.Provider {
	case ProviderOpenAI, ProviderAzure:
		return []string{
			"gpt-4",
			"gpt-4-turbo",
			"gpt-4-turbo-preview",
			"gpt-3.5-turbo",
			"gpt-3.5-turbo-16k",
		}, nil
	case ProviderAnthropic:
		return []string{
			"claude-3-opus-20240229",
			"claude-3-sonnet-20240229",
			"claude-3-haiku-20240307",
			"claude-2.1",
			"claude-2",
		}, nil
	case ProviderGoogle:
		return []string{
			"gemini-pro",
			"gemini-pro-vision",
			"gemini-1.5-pro-latest",
		}, nil
	case ProviderOllama:
		// For Ollama, we'd need to query the local instance
		return []string{
			"llama2",
			"codellama",
			"mistral",
			"vicuna",
		}, nil
	default:
		return []string{}, fmt.Errorf("unsupported provider: %s", mc.Provider)
	}
}