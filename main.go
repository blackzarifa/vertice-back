package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/blackzarifa/vertice-back/config"
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Database connection failed",
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"message": "Server is running",
		})
	}).Methods("GET")

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
