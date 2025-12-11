package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kilo40/idea-forge/internal/llm"
	"github.com/kilo40/idea-forge/internal/search"
	"github.com/kilo40/idea-forge/internal/storage"
)

// Server represents the API server with all dependencies
type Server struct {
	router   *gin.Engine
	db       *storage.Database
	llm      *llm.Client
	search   *search.Client
	obsidian *storage.ObsidianWriter
}

// NewServer creates a new API server instance with all dependencies
func NewServer() *Server {
	router := gin.Default()

	// Configure CORS for frontend
	// Allow all origins since this runs on a private Tailscale network (no public access)
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true // Allow all origins on trusted private network
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s := &Server{
		router: router,
	}

	// Initialize dependencies (log errors but don't fail - allows partial functionality)
	if db, err := storage.NewDatabase(); err != nil {
		log.Printf("Warning: Database initialization failed: %v", err)
	} else {
		s.db = db
	}

	if llmClient, err := llm.NewClient(); err != nil {
		log.Printf("Warning: LLM client initialization failed: %v", err)
	} else {
		s.llm = llmClient
	}

	if searchClient, err := search.NewClient(); err != nil {
		log.Printf("Warning: Search client initialization failed: %v", err)
	} else {
		s.search = searchClient
	}

	if obsidianWriter, err := storage.NewObsidianWriter(); err != nil {
		log.Printf("Warning: Obsidian writer initialization failed: %v", err)
	} else {
		s.obsidian = obsidianWriter
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Routes under /api prefix (for local development and direct access)
	api := s.router.Group("/api")
	{
		api.GET("/health", s.healthCheck)
		api.POST("/notes", s.createNote)
		api.GET("/notes", s.listNotes)
		api.GET("/notes/:id", s.getNote)
		api.DELETE("/notes/:id", s.deleteNote)
		api.GET("/categories", s.listCategories)
	}

	// Same routes at root level (for Tailscale serve which strips /api/ prefix)
	s.router.GET("/health", s.healthCheck)
	s.router.POST("/notes", s.createNote)
	s.router.GET("/notes", s.listNotes)
	s.router.GET("/notes/:id", s.getNote)
	s.router.DELETE("/notes/:id", s.deleteNote)
	s.router.GET("/categories", s.listCategories)
}

// Run starts the HTTP server
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// healthCheck returns server health status
func (s *Server) healthCheck(c *gin.Context) {
	status := gin.H{
		"status":  "healthy",
		"service": "idea-forge",
		"time":    time.Now().UTC().Format(time.RFC3339),
	}

	// Add component status
	components := gin.H{
		"database": s.db != nil,
		"llm":      s.llm != nil,
		"search":   s.search != nil,
		"obsidian": s.obsidian != nil,
	}
	status["components"] = components

	c.JSON(http.StatusOK, status)
}
