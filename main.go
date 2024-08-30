package main

import (
	"log"
	"net/http"
	"twitterjokeposting/router"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	godotenv.Load()

	// Create a new router
	r := mux.NewRouter()

	// Define your routes
	routeDefinitions := router.SetRouteDefinitions()
	for _, routeDefinition := range *routeDefinitions {
		r.HandleFunc(routeDefinition.RouteName, routeDefinition.Handler).Methods(routeDefinition.MethodType)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Consider specifying exact origins in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	handler := c.Handler(r)

	// Start the server
	log.Fatal(http.ListenAndServe(":4000", handler))
}
