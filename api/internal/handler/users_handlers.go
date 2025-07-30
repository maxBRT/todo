package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/maxBRT/todo/api/internal/database"
)

type UserParams struct {
	Email string `json:"email"`
}

// Use the GetUserWithEmail function to handle GET requests to the /users endpoint with a query string parameter of email.
func (c *ApiCongig) GetUserWithEmail(w http.ResponseWriter, r *http.Request) {
	// Get the email from the query string
	email := r.URL.Query().Get("email")

	// Get the user from the database
	users, err := c.DB.GetUserByEmail(r.Context(), email)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error getting user by email", http.StatusInternalServerError)
		return
	}

	// Send the user back to the client as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Use the CreateUser function to handle POST requests to the /users endpoint.
func (c *ApiCongig) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Initialize a struct to unmarshal the request body into
	params := UserParams{}

	// Decode the request body into the struct
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		return
	}

	// Create a new user with the email and ID
	u := database.CreateUserParams{
		Email: params.Email,
		ID:    uuid.New(),
	}
	user, err := c.DB.CreateUser(r.Context(), u)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Send the user back to the client as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (c *ApiCongig) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get the email from the query string
	email := r.URL.Query().Get("email")
	if email == "" {
		log.Println("Error getting user by email")
		http.Error(w, "Error getting user by email", http.StatusBadRequest)
		return
	}

	// Get the user from the database
	user, err := c.DB.GetUserByEmail(r.Context(), email)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error getting user by email", http.StatusInternalServerError)
		return
	}

	// Initialize a struct to unmarshal the request body into
	params := UserParams{}

	// Decode the request body into the struct
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		return
	}

	// Update the user in the database
	u := database.UpdateUserParams{
		ID:    user.ID,
		Email: params.Email,
	}
	user, err = c.DB.UpdateUser(r.Context(), u)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	// Send the user back to the client as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Use the DeleteUser function to handle DELETE requests to the /users endpoint with a query string parameter of email.
func (c *ApiCongig) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get the email from the query string
	email := r.URL.Query().Get("email")
	if email == "" {
		log.Println("Error getting user by email")
		http.Error(w, "Error getting user by email", http.StatusBadRequest)
		return
	}

	// Delete the user from the database
	err := c.DB.DeleteUser(r.Context(), email)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	// Send confirmation response
	w.WriteHeader(http.StatusNoContent)
}
