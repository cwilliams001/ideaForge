package models

import "time"

// NoteInput represents the raw input from the user
type NoteInput struct {
	Content string `json:"content" binding:"required"`
}

// ProcessedNote represents a fully processed note with expanded content
type ProcessedNote struct {
	ID        string     `json:"id"`
	Original  string     `json:"original"`
	Title     string     `json:"title"`
	Category  string     `json:"category"`
	Markdown  string     `json:"markdown"`
	Links     []Link     `json:"links"`
	CreatedAt time.Time  `json:"created_at"`
	SyncedAt  *time.Time `json:"synced_at,omitempty"`
}

// Link represents a resource link associated with a note
type Link struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Type        string `json:"type"` // github, docs, youtube, article
	Description string `json:"description,omitempty"`
}

// LLMResponse represents the structured response from the LLM
type LLMResponse struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	Markdown string `json:"markdown"`
}

// Valid categories for notes
var ValidCategories = []string{
	"homelab",
	"coding",
	"personal",
	"learning",
	"creative",
}

// IsValidCategory checks if a category is valid
func IsValidCategory(category string) bool {
	for _, c := range ValidCategories {
		if c == category {
			return true
		}
	}
	return false
}
