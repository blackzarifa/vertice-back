package main

import (
	"log"
	"os"

	"github.com/blackzarifa/vertice-back/config"
	"github.com/blackzarifa/vertice-back/routes"
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

	// Setup all routes
	routes.SetupRoutes(router, db)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}
