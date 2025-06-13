package routes

import (
	"database/sql"

	"github.com/blackzarifa/vertice-back/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	api := router.Group("/api/v1")

	// Health check
	api.GET("/health", handlers.HealthCheck(db))

	// Funcionario routes
	setupFuncionarioRoutes(api, db)
}

func setupFuncionarioRoutes(api *gin.RouterGroup, db *sql.DB) {
	api.GET("/funcionarios", handlers.ListFuncionarios(db))
	api.POST("/funcionarios", handlers.CreateFuncionario(db))
}
