package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/maxBRT/todo/api/internal/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

// TaskResponse is the response struct for the /tasks endpoint
type TaskResponse struct {
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Priority    int32     `json:"priority"`
	DueDate     time.Time `json:"due_date"`
}

// Use the CreateTask function to handle POST requests to the /tasks endpoint with a query string parameter of email.
func (c *ApiCongig) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Get the email from the query string and validate it
	userEmail := r.URL.Query().Get("email")
	if userEmail == "" {
		log.Println("Error getting user by email")
		http.Error(w, "Error getting user by email", http.StatusBadRequest)
		return
	}

	// Get the user from the database
	user, err := c.DB.GetUserByEmail(r.Context(), userEmail)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error getting user by email", http.StatusInternalServerError)
		return
	}

	// Declare a struct to unmarshal the request body into
	var params struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
		DueDate     string `json:"due_date"`
	}

	// Decode the request body into the struct
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		return
	}

	// Parse the due date
	parsedDate, err := time.Parse("2006-01-02", params.DueDate)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error parsing due date", http.StatusInternalServerError)
		return
	}

	// Parse the priority
	parsedPriority, err := strconv.Atoi(params.Priority)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error parsing priority", http.StatusInternalServerError)
		return
	}

	// Create a new task in the database
	t := database.CreateTaskParams{
		ID:          uuid.New(),
		UserID:      user.ID,
		Title:       params.Title,
		Description: sql.NullString{String: params.Description, Valid: true},
		Priority:    int32(parsedPriority),
		DueDate:     sql.NullTime{Time: parsedDate, Valid: true},
	}
	task, err := c.DB.CreateTask(r.Context(), t)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		return
	}

	// Set the response struct with the task
	resp := TaskResponse{
		Title:       task.Title,
		Description: task.Description.String,
		Priority:    task.Priority,
		DueDate:     task.DueDate.Time,
	}

	// Send the task back to the client as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Use the GetTasksByUser function to handle GET requests to the /tasks endpoint with a query string parameter of email.
func (c *ApiCongig) GetTasksByUser(w http.ResponseWriter, r *http.Request) {
	// Get the email from the query string and validate it
	userEmail := r.URL.Query().Get("email")
	if userEmail == "" {
		log.Println("Error getting user by email")
		http.Error(w, "Error getting user by email", http.StatusBadRequest)
		return
	}

	// Get the user from the database
	user, err := c.DB.GetUserByEmail(r.Context(), userEmail)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error getting user by email", http.StatusInternalServerError)
		return
	}

	// Get the tasks from the database
	tasks, err := c.DB.GetTasksByUser(r.Context(), user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error getting tasks by user", http.StatusInternalServerError)
		return
	}

	// Send the tasks back to the client as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
