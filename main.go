package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Define your routes
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/about", aboutHandler).Methods("GET")

	// Start the server
	http.ListenAndServe(":6907", r)
}

// Handler for the home route
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the home page!"))
}

// Handler for the about route
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("About us page"))
}
