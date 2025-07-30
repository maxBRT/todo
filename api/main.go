package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/maxBRT/todo/api/internal/database"
	"github.com/maxBRT/todo/api/internal/handler"
	"log"
	"net/http"
	"os"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database and ping it
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create the API config
	config := handler.ApiCongig{
		DB: database.New(db),
	}

	// Initialize the server
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":5000",
		Handler: mux,
	}

	// Register the handlers
	mux.HandleFunc("POST /api/users", config.CreateUser)
	mux.HandleFunc("GET /api/users", config.GetUserWithEmail)
	mux.HandleFunc("PATCH /api/users", config.UpdateUser)
	mux.HandleFunc("DELETE /api/users", config.DeleteUser)

	mux.HandleFunc("POST /api/tasks", config.CreateTask)
	mux.HandleFunc("GET /api/tasks", config.GetTasksByUser)

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
