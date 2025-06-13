package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
)

func SetupRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()

	// api.HandleFunc("/usuarios", handlers.CreateUsuario(db)).Methods("POST")
	// api.HandleFunc("/usuarios/{id}", handlers.GetUsuario(db)).Methods("GET")
	//
	// api.HandleFunc("/login", handlers.Login(db)).Methods("POST")
	//
	// api.HandleFunc("/contas", handlers.CreateConta(db)).Methods("POST")
	// api.HandleFunc("/contas/{numero}", handlers.GetConta(db)).Methods("GET")
	//
	// api.HandleFunc("/transacoes", handlers.CreateTransacao(db)).Methods("POST")
	// api.HandleFunc("/transacoes/conta/{conta_id}", handlers.GetTransacoesByConta(db)).
	// 	Methods("GET")

	return router
}
