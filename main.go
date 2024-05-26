package main

import (
	"net/http"
	"twitterjokeposting/handler"

	"github.com/gorilla/mux"
)

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Define your routes
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/about", aboutHandler).Methods("GET")
	r.HandleFunc("/scheduleJokes", handler.ScheduleJokeForTodayController).Methods("POST")
	r.HandleFunc("/getScheduleJokes", handler.GetAllScheduledJokes).Methods("GET")

	// Start the server
	http.ListenAndServe(":6907", r)
}

// Handler for the home route
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

// Handler for the about route
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Scheduling jokes"))
}
