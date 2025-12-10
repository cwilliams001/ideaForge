package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kilo40/idea-forge/internal/models"
)

// createNote handles POST /api/notes
func (s *Server) createNote(c *gin.Context) {
	var input models.NoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
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

	// TODO: Implement full flow
	// 1. Call LLM to expand note
	// 2. Call search API for links
	// 3. Save to database
	// 4. Write to Obsidian vault

	// Placeholder response
	c.JSON(http.StatusOK, gin.H{
		"message": "Note processing not yet implemented",
		"received": input.Content,
	})
}

// listNotes handles GET /api/notes
func (s *Server) listNotes(c *gin.Context) {
	category := c.Query("category")
	// limit := c.DefaultQuery("limit", "50")
	// offset := c.DefaultQuery("offset", "0")

	// TODO: Implement database query
	_ = category

	c.JSON(http.StatusOK, gin.H{
		"notes": []models.ProcessedNote{},
		"total": 0,
	})
}

// getNote handles GET /api/notes/:id
func (s *Server) getNote(c *gin.Context) {
	id := c.Param("id")

	// TODO: Implement database lookup
	_ = id

	c.JSON(http.StatusNotFound, gin.H{
		"error": "Note not found",
	})
}

// deleteNote handles DELETE /api/notes/:id
func (s *Server) deleteNote(c *gin.Context) {
	id := c.Param("id")

	// TODO: Implement deletion from database and Obsidian vault
	_ = id

	c.JSON(http.StatusNotFound, gin.H{
		"error": "Note not found",
	})
}

// listCategories handles GET /api/categories
func (s *Server) listCategories(c *gin.Context) {
	// TODO: Get counts from database
	categories := make([]gin.H, len(models.ValidCategories))
	for i, cat := range models.ValidCategories {
		categories[i] = gin.H{
			"name":  cat,
			"count": 0, // TODO: actual count
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}
