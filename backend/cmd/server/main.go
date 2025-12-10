package main

import (
	"log"
	"os"

	"github.com/kilo40/idea-forge/internal/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := api.NewServer()

	log.Printf("Starting Idea Forge API server on :%s", port)
	if err := server.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
