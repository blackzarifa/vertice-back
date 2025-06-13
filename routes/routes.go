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
	
	// Authentication routes
	setupAuthRoutes(api, db)
	
	// Cliente routes
	setupClienteRoutes(api, db)
	
	// Conta routes
	setupContaRoutes(api, db)
	
	// Transacao routes
	setupTransacaoRoutes(api, db)
}

func setupAuthRoutes(api *gin.RouterGroup, db *sql.DB) {
	// api.POST("/auth/login", handlers.Login(db))
	// api.POST("/auth/logout", handlers.Logout())
	// api.POST("/auth/refresh", handlers.RefreshToken())
}

func setupClienteRoutes(api *gin.RouterGroup, db *sql.DB) {
	// api.GET("/clientes", handlers.ListClientes(db))
	// api.GET("/clientes/:id", handlers.GetCliente(db))
	// api.POST("/clientes", handlers.CreateCliente(db))
	// api.PUT("/clientes/:id", handlers.UpdateCliente(db))
	// api.DELETE("/clientes/:id", handlers.DeleteCliente(db))
}

func setupContaRoutes(api *gin.RouterGroup, db *sql.DB) {
	// api.GET("/contas", handlers.ListContas(db))
	// api.GET("/contas/:id", handlers.GetConta(db))
	// api.POST("/contas", handlers.CreateConta(db))
	// api.PUT("/contas/:id", handlers.UpdateConta(db))
	// api.DELETE("/contas/:id", handlers.DeleteConta(db))
	// api.GET("/contas/:id/saldo", handlers.GetSaldo(db))
}

func setupTransacaoRoutes(api *gin.RouterGroup, db *sql.DB) {
	// api.GET("/transacoes", handlers.ListTransacoes(db))
	// api.POST("/transacoes/transferencia", handlers.Transferencia(db))
	// api.POST("/transacoes/deposito", handlers.Deposito(db))
	// api.POST("/transacoes/saque", handlers.Saque(db))
}
