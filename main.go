package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/blackzarifa/vertice-back/config"
	"github.com/blackzarifa/vertice-back/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := config.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Setup routes
	router := gin.Default()
	setupRoutes(router, db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(router *gin.Engine, db *sql.DB) {
	// Create handler
	funcHandler := handlers.NewFuncionarioHandler(db)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "database": "connected"})
		})

		// Funcionario routes
		v1.POST("/funcionarios", funcHandler.Create)
		v1.GET("/funcionarios", funcHandler.List)
	}
}
