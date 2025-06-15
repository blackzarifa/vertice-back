package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/blackzarifa/vertice-back/config"
	"github.com/blackzarifa/vertice-back/handlers"
	"github.com/blackzarifa/vertice-back/middleware"
	"github.com/gin-contrib/cors"
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

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

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
	jwtSecret := "vertice-bank-secret-2025"

	// Create handlers
	funcHandler := handlers.NewFuncionarioHandler(db)
	authHandler := handlers.NewAuthHandler(db, jwtSecret)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "database": "connected"})
		})

		// Auth routes (public)
		v1.POST("/auth/login", authHandler.Login)

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthRequired(jwtSecret))
		{
			// Funcionario routes
			protected.POST("/funcionarios", funcHandler.Create)
			protected.GET("/funcionarios", funcHandler.List)
		}
	}
}
