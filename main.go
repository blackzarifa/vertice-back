package main

import (
	"log"
	"net/http"
	"os"

	"github.com/blackzarifa/vertice-back/config"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	log.Println("Database connected successfully!")

	if err := config.RunMigrations(db); err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migrations completed successfully!")

	// Create Gin router with default middleware (logger + recovery)
	router := gin.Default()

	// Basic health check
	router.GET("/health", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "Database connection failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Server is running",
		})
	})

	// Example: Path parameters
	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(200, gin.H{
			"message": "You requested user",
			"id":      id,
		})
	})

	// Example: JSON body parsing
	router.POST("/login", func(c *gin.Context) {
		var loginReq struct {
			CPF   string `json:"cpf" binding:"required"`
			Senha string `json:"senha" binding:"required"`
		}

		// This automatically parses JSON and validates required fields
		if err := c.ShouldBindJSON(&loginReq); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"message": "Login attempt",
			"cpf":     loginReq.CPF,
		})
	})

	// Example: Query parameters
	router.GET("/contas", func(c *gin.Context) {
		// GET /contas?tipo=POUPANCA&limit=10
		tipo := c.Query("tipo")                // "POUPANCA"
		limit := c.DefaultQuery("limit", "20") // "10" or default "20"

		c.JSON(200, gin.H{
			"tipo":  tipo,
			"limit": limit,
		})
	})

	// Group routes with common prefix
	api := router.Group("/api/v1")
	{
		api.GET("/clientes", func(c *gin.Context) {
			c.JSON(200, gin.H{"endpoint": "list clientes"})
		})
		api.POST("/clientes", func(c *gin.Context) {
			c.JSON(201, gin.H{"endpoint": "create cliente"})
		})
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}
