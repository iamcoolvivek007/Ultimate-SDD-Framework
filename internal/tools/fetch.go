package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// FetchTool fetches content from URLs
type FetchTool struct {
	Timeout   time.Duration
	UserAgent string
	MaxSize   int64
}

// FetchFormat specifies the output format
type FetchFormat string

const (
	FetchFormatText     FetchFormat = "text"
	FetchFormatHTML     FetchFormat = "html"
	FetchFormatMarkdown FetchFormat = "markdown"
	FetchFormatJSON     FetchFormat = "json"
)

// FetchResult represents the result of fetching a URL
type FetchResult struct {
	URL         string            `json:"url"`
	StatusCode  int               `json:"status_code"`
	ContentType string            `json:"content_type"`
	Content     string            `json:"content"`
	Headers     map[string]string `json:"headers,omitempty"`
	Size        int               `json:"size"`
	Duration    time.Duration     `json:"duration"`
}

// DefaultFetchTool creates a FetchTool with sensible defaults
func DefaultFetchTool() *FetchTool {
	return &FetchTool{
		Timeout:   30 * time.Second,
		UserAgent: "Viki/2.0 (AI Development Assistant)",
		MaxSize:   10 * 1024 * 1024, // 10MB
	}
}

// Fetch retrieves content from a URL
func (f *FetchTool) Fetch(ctx context.Context, url string, format FetchFormat) (*FetchResult, error) {
	result := &FetchResult{
		URL:     url,
		Headers: make(map[string]string),
	}

	start := time.Now()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", f.UserAgent)

	// Create client with timeout
	client := &http.Client{
		Timeout: f.Timeout,
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.ContentType = resp.Header.Get("Content-Type")
	result.Duration = time.Since(start)

	// Store relevant headers
	for _, key := range []string{"Content-Type", "Content-Length", "Last-Modified", "ETag"} {
		if val := resp.Header.Get(key); val != "" {
			result.Headers[key] = val
		}
	}

	// Read body with size limit
	limitReader := io.LimitReader(resp.Body, f.MaxSize)
	body, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	result.Size = len(body)

	// Process content based on format
	switch format {
	case FetchFormatText:
		result.Content = extractText(string(body))
	case FetchFormatHTML:
		result.Content = string(body)
	case FetchFormatMarkdown:
		result.Content = htmlToMarkdown(string(body))
	case FetchFormatJSON:
		// Validate and pretty-print JSON
		var obj interface{}
		if err := json.Unmarshal(body, &obj); err == nil {
			prettyJSON, _ := json.MarshalIndent(obj, "", "  ")
			result.Content = string(prettyJSON)
		} else {
			result.Content = string(body)
		}
	default:
		result.Content = string(body)
	}

	return result, nil
}

// extractText extracts plain text from HTML
func extractText(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent
	}

	var sb strings.Builder
	var extractTextNode func(*html.Node)

	extractTextNode = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				sb.WriteString(text)
				sb.WriteString(" ")
			}
		}

		// Skip script and style tags
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			return
		}

		// Add newlines for block elements
		if n.Type == html.ElementNode {
			switch n.Data {
			case "p", "div", "br", "li", "h1", "h2", "h3", "h4", "h5", "h6", "tr":
				sb.WriteString("\n")
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractTextNode(c)
		}
	}

	extractTextNode(doc)
	return strings.TrimSpace(sb.String())
}

// htmlToMarkdown converts HTML to basic markdown
func htmlToMarkdown(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent
	}

	var sb strings.Builder
	var convertNode func(*html.Node)

	convertNode = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				sb.WriteString(text)
			}
		}

		// Skip script and style
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			return
		}

		if n.Type == html.ElementNode {
			switch n.Data {
			case "h1":
				sb.WriteString("\n# ")
			case "h2":
				sb.WriteString("\n## ")
			case "h3":
				sb.WriteString("\n### ")
			case "p":
				sb.WriteString("\n\n")
			case "br":
				sb.WriteString("\n")
			case "li":
				sb.WriteString("\n- ")
			case "strong", "b":
				sb.WriteString("**")
			case "em", "i":
				sb.WriteString("*")
			case "code":
				sb.WriteString("`")
			case "pre":
				sb.WriteString("\n```\n")
			case "a":
				sb.WriteString("[")
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			convertNode(c)
		}

		if n.Type == html.ElementNode {
			switch n.Data {
			case "strong", "b":
				sb.WriteString("**")
			case "em", "i":
				sb.WriteString("*")
			case "code":
				sb.WriteString("`")
			case "pre":
				sb.WriteString("\n```\n")
			case "a":
				href := ""
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						href = attr.Val
						break
					}
				}
				sb.WriteString("](")
				sb.WriteString(href)
				sb.WriteString(")")
			}
		}
	}

	convertNode(doc)
	return strings.TrimSpace(sb.String())
}

// Head performs a HEAD request to get headers without body
func (f *FetchTool) Head(ctx context.Context, url string) (*FetchResult, error) {
	result := &FetchResult{
		URL:     url,
		Headers: make(map[string]string),
	}

	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", f.UserAgent)

	client := &http.Client{Timeout: f.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.ContentType = resp.Header.Get("Content-Type")
	result.Duration = time.Since(start)

	for key, vals := range resp.Header {
		if len(vals) > 0 {
			result.Headers[key] = vals[0]
		}
	}

	return result, nil
}
