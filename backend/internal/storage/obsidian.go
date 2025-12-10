package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/kilo40/idea-forge/internal/models"
)

// ObsidianWriter handles writing notes to an Obsidian vault
type ObsidianWriter struct {
	vaultPath  string
	folderName string
}

// NewObsidianWriter creates a new Obsidian writer
func NewObsidianWriter() (*ObsidianWriter, error) {
	vaultPath := os.Getenv("OBSIDIAN_VAULT_PATH")
	if vaultPath == "" {
		return nil, fmt.Errorf("OBSIDIAN_VAULT_PATH environment variable is required")
	}

	folderName := os.Getenv("OBSIDIAN_FOLDER")
	if folderName == "" {
		folderName = "IdeaForge"
	}

	// Verify vault path exists
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("vault path does not exist: %s", vaultPath)
	}

	return &ObsidianWriter{
		vaultPath:  vaultPath,
		folderName: folderName,
	}, nil
}

// WriteNote writes a processed note to the Obsidian vault
func (w *ObsidianWriter) WriteNote(note *models.ProcessedNote) error {
	// Create category folder if it doesn't exist
	categoryPath := filepath.Join(w.vaultPath, w.folderName, note.Category)
	if err := os.MkdirAll(categoryPath, 0755); err != nil {
		return fmt.Errorf("failed to create category folder: %w", err)
	}

	// Generate filename
	filename := w.generateFilename(note)
	filePath := filepath.Join(categoryPath, filename)

	// Prevent path traversal
	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(w.vaultPath)) {
		return fmt.Errorf("invalid file path: attempted path traversal")
	}

	// Generate file content
	content := w.generateContent(note)

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// DeleteNote removes a note file from the Obsidian vault
func (w *ObsidianWriter) DeleteNote(note *models.ProcessedNote) error {
	filename := w.generateFilename(note)
	filePath := filepath.Join(w.vaultPath, w.folderName, note.Category, filename)

	// Prevent path traversal
	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(w.vaultPath)) {
		return fmt.Errorf("invalid file path: attempted path traversal")
	}

	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// generateFilename creates a sanitized filename for the note
func (w *ObsidianWriter) generateFilename(note *models.ProcessedNote) string {
	// Format: YYYY-MM-DD-slugified-title.md
	date := note.CreatedAt.Format("2006-01-02")
	slug := slugify(note.Title)
	return fmt.Sprintf("%s-%s.md", date, slug)
}

// generateContent creates the full markdown content with frontmatter
func (w *ObsidianWriter) generateContent(note *models.ProcessedNote) string {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("created: %s\n", note.CreatedAt.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("category: %s\n", note.Category))
	sb.WriteString("source: idea-forge\n")
	sb.WriteString(fmt.Sprintf("id: %s\n", note.ID))
	sb.WriteString("---\n\n")

	// Main content
	sb.WriteString(note.Markdown)

	// Resources section if links exist
	if len(note.Links) > 0 {
		sb.WriteString("\n\n## Resources\n")
		for _, link := range note.Links {
			if link.Description != "" {
				sb.WriteString(fmt.Sprintf("- [%s](%s) - %s\n", link.Title, link.URL, link.Description))
			} else {
				sb.WriteString(fmt.Sprintf("- [%s](%s)\n", link.Title, link.URL))
			}
		}
	}

	return sb.String()
}

// slugify converts a title to a URL-safe slug
func slugify(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove non-alphanumeric characters except hyphens
	reg := regexp.MustCompile("[^a-z0-9-]")
	slug = reg.ReplaceAllString(slug, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Limit length
	if len(slug) > 50 {
		slug = slug[:50]
		// Don't end on a hyphen
		slug = strings.TrimSuffix(slug, "-")
	}

	return slug
}
