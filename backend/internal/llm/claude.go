package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kilo40/idea-forge/internal/models"
)

const (
	anthropicAPIURL = "https://api.anthropic.com/v1/messages"
	defaultModel    = "claude-sonnet-4-20250514"
)

// Client represents the Claude API client
type Client struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewClient creates a new Claude API client
func NewClient() (*Client, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable is required")
	}

	model := os.Getenv("LLM_MODEL")
	if model == "" {
		model = defaultModel
	}

	return &Client{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// SystemPrompt for note expansion
const systemPrompt = `You are a productivity assistant that transforms quick notes into structured, actionable markdown todo lists.

Given a brief note or idea, you will:
1. Determine the most appropriate category from: homelab, coding, personal, learning, creative
2. Create a clear, descriptive title
3. Expand the note into a markdown checklist with logical steps
4. Keep steps actionable and specific
5. Add brief context where helpful

Respond ONLY with valid JSON in this exact format:
{
  "title": "Clear Title Here",
  "category": "category_name",
  "markdown": "# Title\n\n## Tasks\n- [ ] First step\n- [ ] Second step\n..."
}

Keep the markdown concise but comprehensive. Each task should be completable in one sitting.
Do not include any text outside the JSON object.`

// anthropicRequest represents the API request structure
type anthropicRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system"`
	Messages  []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// anthropicResponse represents the API response structure
type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
	Error      *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// ExpandNote takes a raw note and returns structured LLM response
func (c *Client) ExpandNote(ctx context.Context, note string) (*models.LLMResponse, error) {
	reqBody := anthropicRequest{
		Model:     c.model,
		MaxTokens: 2048,
		System:    systemPrompt,
		Messages: []message{
			{
				Role:    "user",
				Content: note,
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", anthropicAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var apiResp anthropicResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if apiResp.Error != nil {
		return nil, fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	if len(apiResp.Content) == 0 || apiResp.Content[0].Type != "text" {
		return nil, fmt.Errorf("unexpected response format")
	}

	// Clean the response text (LLMs sometimes wrap JSON in markdown code blocks)
	responseText := cleanJSONResponse(apiResp.Content[0].Text)

	var llmResponse models.LLMResponse
	if err := json.Unmarshal([]byte(responseText), &llmResponse); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response as JSON: %w", err)
	}

	// Validate category
	if !models.IsValidCategory(llmResponse.Category) {
		llmResponse.Category = "personal" // Default fallback
	}

	return &llmResponse, nil
}

// cleanJSONResponse strips markdown code blocks from LLM responses
func cleanJSONResponse(text string) string {
	text = strings.TrimSpace(text)

	// Remove ```json or ``` prefix
	if strings.HasPrefix(text, "```json") {
		text = strings.TrimPrefix(text, "```json")
	} else if strings.HasPrefix(text, "```") {
		text = strings.TrimPrefix(text, "```")
	}

	// Remove trailing ```
	if strings.HasSuffix(text, "```") {
		text = strings.TrimSuffix(text, "```")
	}

	return strings.TrimSpace(text)
}
