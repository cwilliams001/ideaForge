package api

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server represents the API server
type Server struct {
	router *gin.Engine
	// TODO: Add dependencies (db, llm client, search client)
}

// NewServer creates a new API server instance
func NewServer() *Server {
	router := gin.Default()

	// Configure CORS for frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s := &Server{
		router: router,
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	api := s.router.Group("/api")
	{
		// Health check
		api.GET("/health", s.healthCheck)

		// Notes endpoints
		api.POST("/notes", s.createNote)
		api.GET("/notes", s.listNotes)
		api.GET("/notes/:id", s.getNote)
		api.DELETE("/notes/:id", s.deleteNote)

		// Categories endpoint
		api.GET("/categories", s.listCategories)
	}
}

// Run starts the HTTP server
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// healthCheck returns server health status
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "idea-forge",
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}
