package search

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/kilo40/idea-forge/internal/models"
)

// Client represents the search client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new search client
func NewClient() (*Client, error) {
	baseURL := os.Getenv("SEARXNG_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("SEARXNG_URL environment variable is required")
	}

	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// searxngResponse represents the SearXNG API response
type searxngResponse struct {
	Results []struct {
		Title   string `json:"title"`
		URL     string `json:"url"`
		Content string `json:"content"`
	} `json:"results"`
}

// SearchForLinks searches for relevant links for a topic
func (c *Client) SearchForLinks(ctx context.Context, topic string) ([]models.Link, error) {
	var allLinks []models.Link

	// Search queries to find different types of resources
	queries := []struct {
		query    string
		linkType string
	}{
		{fmt.Sprintf("%s github", topic), "github"},
		{fmt.Sprintf("%s documentation official", topic), "docs"},
		{fmt.Sprintf("%s tutorial setup guide", topic), "article"},
	}

	for _, q := range queries {
		links, err := c.search(ctx, q.query, q.linkType)
		if err != nil {
			// Log error but continue with other queries
			continue
		}
		allLinks = append(allLinks, links...)
	}

	// Deduplicate and limit results
	seen := make(map[string]bool)
	var uniqueLinks []models.Link
	for _, link := range allLinks {
		if !seen[link.URL] {
			seen[link.URL] = true
			uniqueLinks = append(uniqueLinks, link)
		}
		if len(uniqueLinks) >= 5 {
			break
		}
	}

	return uniqueLinks, nil
}

// search performs a single search query
func (c *Client) search(ctx context.Context, query string, defaultType string) ([]models.Link, error) {
	searchURL := fmt.Sprintf("%s/search", c.baseURL)
	params := url.Values{}
	params.Set("q", query)
	params.Set("format", "json")
	params.Set("engines", "google,duckduckgo,bing")

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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
		return nil, fmt.Errorf("search API error (status %d): %s", resp.StatusCode, string(body))
	}

	var searchResp searxngResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	var links []models.Link
	for _, result := range searchResp.Results {
		if len(links) >= 2 { // Limit per query
			break
		}
		link := models.Link{
			Title:       result.Title,
			URL:         result.URL,
			Type:        categorizeURL(result.URL, defaultType),
			Description: truncate(result.Content, 150),
		}
		links = append(links, link)
	}

	return links, nil
}

// categorizeURL determines the link type based on URL
func categorizeURL(rawURL string, defaultType string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return defaultType
	}

	host := strings.ToLower(u.Host)

	switch {
	case strings.Contains(host, "github.com"):
		return "github"
	case strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be"):
		return "youtube"
	case strings.Contains(host, "docs.") || strings.HasSuffix(u.Path, "/docs") || strings.Contains(u.Path, "/docs/"):
		return "docs"
	default:
		return defaultType
	}
}

// truncate shortens a string to max length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
