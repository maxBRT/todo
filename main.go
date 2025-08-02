package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/maxBRT/todo/internals/views"
)

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	// Serve static files from the "./assets" directory
	fs := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		views.Index().Render(r.Context(), w)
	})
	mux.HandleFunc("POST /api/tasks", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing the form: %s", err)
		}

		title := r.FormValue("title")
		description := r.FormValue("description")
		dueDate := r.FormValue("due-date")

		fmt.Println(title)
		fmt.Println(description)
		fmt.Println(dueDate)

		w.WriteHeader(http.StatusAccepted)
		views.TaskForm().Render(r.Context(), w)
	})

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
