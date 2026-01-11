package mcp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// StreamCallback is called for each token/chunk received
type StreamCallback func(token string, done bool)

// StreamDelta represents a streaming response delta
type StreamDelta struct {
	Content      string
	FinishReason string
	Error        error
}

// ChatStream sends a chat request with streaming response
func (mc *ModelClient) ChatStream(messages []Message, options map[string]interface{}, callback StreamCallback) error {
	switch mc.Provider {
	case ProviderOpenAI, ProviderAzure:
		return mc.streamOpenAI(messages, options, callback)
	case ProviderAnthropic:
		return mc.streamAnthropic(messages, options, callback)
	case ProviderGoogle:
		return mc.streamGoogle(messages, options, callback)
	case ProviderOllama:
		return mc.streamOllama(messages, options, callback)
	default:
		return fmt.Errorf("streaming not supported for provider: %s", mc.Provider)
	}
}

// streamOpenAI handles OpenAI streaming
func (mc *ModelClient) streamOpenAI(messages []Message, options map[string]interface{}, callback StreamCallback) error {
	request := map[string]interface{}{
		"model":    mc.Model,
		"messages": messages,
		"stream":   true,
	}

	if temp, ok := options["temperature"].(float64); ok {
		request["temperature"] = temp
	}
	if maxTokens, ok := options["max_tokens"].(int); ok {
		request["max_tokens"] = maxTokens
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := mc.BaseURL + "/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+mc.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" || line == "data: [DONE]" {
			if line == "data: [DONE]" {
				callback("", true)
			}
			continue
		}

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
				FinishReason string `json:"finish_reason"`
			} `json:"choices"`
		}

		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			done := chunk.Choices[0].FinishReason != ""
			callback(content, done)
		}
	}

	return nil
}

// streamAnthropic handles Anthropic streaming
func (mc *ModelClient) streamAnthropic(messages []Message, options map[string]interface{}, callback StreamCallback) error {
	systemMessage := ""
	userMessages := []Message{}

	for _, msg := range messages {
		if msg.Role == "system" {
			systemMessage = msg.Content
		} else {
			userMessages = append(userMessages, msg)
		}
	}

	request := map[string]interface{}{
		"model":      mc.Model,
		"messages":   userMessages,
		"max_tokens": 4096,
		"stream":     true,
	}
	if systemMessage != "" {
		request["system"] = systemMessage
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := mc.BaseURL + "/messages"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", mc.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		var event struct {
			Type  string `json:"type"`
			Delta struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"delta"`
		}

		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		switch event.Type {
		case "content_block_delta":
			if event.Delta.Type == "text_delta" {
				callback(event.Delta.Text, false)
			}
		case "message_stop":
			callback("", true)
		}
	}

	return nil
}

// streamGoogle handles Google Gemini streaming
func (mc *ModelClient) streamGoogle(messages []Message, options map[string]interface{}, callback StreamCallback) error {
	request := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]interface{}{
				{"text": messages[len(messages)-1].Content},
			}},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:streamGenerateContent?key=%s",
		mc.BaseURL, mc.Model, mc.APIKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Gemini returns line-delimited JSON
	reader := bufio.NewReader(resp.Body)
	decoder := json.NewDecoder(reader)

	// Skip opening bracket
	if _, err := decoder.Token(); err != nil {
		return err
	}

	for decoder.More() {
		var chunk struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
				FinishReason string `json:"finishReason"`
			} `json:"candidates"`
		}

		if err := decoder.Decode(&chunk); err != nil {
			break
		}

		if len(chunk.Candidates) > 0 && len(chunk.Candidates[0].Content.Parts) > 0 {
			text := chunk.Candidates[0].Content.Parts[0].Text
			done := chunk.Candidates[0].FinishReason != ""
			callback(text, done)
		}
	}

	return nil
}

// streamOllama handles Ollama streaming
func (mc *ModelClient) streamOllama(messages []Message, options map[string]interface{}, callback StreamCallback) error {
	request := map[string]interface{}{
		"model":    mc.Model,
		"messages": messages,
		"stream":   true,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := mc.BaseURL + "/api/chat"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		var chunk struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			Done bool `json:"done"`
		}

		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			continue
		}

		callback(chunk.Message.Content, chunk.Done)
		if chunk.Done {
			break
		}
	}

	return nil
}
