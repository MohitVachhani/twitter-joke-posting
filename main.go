package main

import (
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
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"*"},
	})

	handler := c.Handler(r)

	// Start the server
	http.ListenAndServe(":4000", handler)
}
