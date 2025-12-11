package api

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kilo40/idea-forge/internal/models"
)

// createNote handles POST /api/notes
// This is the main flow: input -> LLM expansion -> search links -> save -> sync to Obsidian
func (s *Server) createNote(c *gin.Context) {
	var input models.NoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if input.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Content is required",
		})
		return
	}

	// Check if LLM is available
	if s.llm == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "LLM service not configured. Set ANTHROPIC_API_KEY.",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	// Step 1: Expand note with LLM
	log.Printf("Expanding note: %s", input.Content)
	llmResponse, err := s.llm.ExpandNote(ctx, input.Content)
	if err != nil {
		log.Printf("LLM expansion failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to expand note",
			"details": err.Error(),
		})
		return
	}

	// Step 2: Search for relevant links (optional - don't fail if search unavailable)
	links := make([]models.Link, 0) // Initialize as empty slice, not nil (nil serializes to null in JSON)
	if s.search != nil {
		log.Printf("Searching for links: %s", llmResponse.Title)
		searchLinks, err := s.search.SearchForLinks(ctx, llmResponse.Title)
		if err != nil {
			log.Printf("Search failed (continuing without links): %v", err)
		} else {
			links = searchLinks
		}
	}

	// Build the processed note
	now := time.Now()
	note := &models.ProcessedNote{
		Original:  input.Content,
		Title:     llmResponse.Title,
		Category:  llmResponse.Category,
		Markdown:  llmResponse.Markdown,
		Links:     links,
		CreatedAt: now,
	}

	// Step 3: Save to database (optional - don't fail if db unavailable)
	if s.db != nil {
		if err := s.db.CreateNote(note); err != nil {
			log.Printf("Database save failed (continuing): %v", err)
		} else {
			log.Printf("Note saved to database: %s", note.ID)
		}
	}

	// Step 4: Write to Obsidian vault (optional - don't fail if not configured)
	if s.obsidian != nil {
		if err := s.obsidian.WriteNote(note); err != nil {
			log.Printf("Obsidian write failed (continuing): %v", err)
		} else {
			syncTime := time.Now()
			note.SyncedAt = &syncTime
			log.Printf("Note written to Obsidian: %s/%s", note.Category, note.Title)

			// Update sync time in database
			if s.db != nil {
				s.db.UpdateSyncedAt(note.ID, syncTime)
			}
		}
	}

	c.JSON(http.StatusOK, note)
}

// listNotes handles GET /api/notes
func (s *Server) listNotes(c *gin.Context) {
	if s.db == nil {
		c.JSON(http.StatusOK, gin.H{
			"notes": []models.ProcessedNote{},
			"total": 0,
		})
		return
	}

	category := c.Query("category")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notes, total, err := s.db.ListNotes(category, limit, offset)
	if err != nil {
		log.Printf("Failed to list notes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve notes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notes": notes,
		"total": total,
	})
}

// getNote handles GET /api/notes/:id
func (s *Server) getNote(c *gin.Context) {
	if s.db == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Note not found",
		})
		return
	}

	id := c.Param("id")
	note, err := s.db.GetNote(id)
	if err != nil {
		log.Printf("Failed to get note: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve note",
		})
		return
	}

	if note == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Note not found",
		})
		return
	}

	c.JSON(http.StatusOK, note)
}

// deleteNote handles DELETE /api/notes/:id
func (s *Server) deleteNote(c *gin.Context) {
	if s.db == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Note not found",
		})
		return
	}

	id := c.Param("id")

	// Get note first for Obsidian deletion
	note, err := s.db.GetNote(id)
	if err != nil {
		log.Printf("Failed to get note for deletion: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete note",
		})
		return
	}

	if note == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Note not found",
		})
		return
	}

	// Delete from database
	if err := s.db.DeleteNote(id); err != nil {
		log.Printf("Failed to delete note from database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete note",
		})
		return
	}

	// Delete from Obsidian vault
	if s.obsidian != nil {
		if err := s.obsidian.DeleteNote(note); err != nil {
			log.Printf("Failed to delete from Obsidian (note removed from db): %v", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Note deleted",
		"id":      id,
	})
}

// listCategories handles GET /api/categories
func (s *Server) listCategories(c *gin.Context) {
	var counts map[string]int
	if s.db != nil {
		var err error
		counts, err = s.db.GetCategoryCounts()
		if err != nil {
			log.Printf("Failed to get category counts: %v", err)
			counts = make(map[string]int)
		}
	} else {
		counts = make(map[string]int)
	}

	categories := make([]gin.H, len(models.ValidCategories))
	for i, cat := range models.ValidCategories {
		categories[i] = gin.H{
			"name":  cat,
			"count": counts[cat],
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}
