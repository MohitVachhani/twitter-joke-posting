package main

import (
	"log"
	"net/http"
	"twitterjokeposting/middleware"
	"twitterjokeposting/router"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

// contentTypeMiddleware sets the Content-Type header to application/json for all responses
func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	// Create a new router
	r := mux.NewRouter()

	// Define your routes
	routeDefinitions := router.SetRouteDefinitions()
	for _, routeDefinition := range *routeDefinitions {
		r.HandleFunc(routeDefinition.RouteName, routeDefinition.Handler).Methods(routeDefinition.MethodType)
	}

	// Create a new rate limiter that allows 5 requests per minute with a burst of 10
	limiter := middleware.NewIPRateLimiter(rate.Limit(5/60), 10)

	// Apply middlewares
	r.Use(middleware.RateLimitMiddleware(limiter))
	r.Use(contentTypeMiddleware)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Consider specifying exact origins in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	// Apply CORS middleware
	handler := c.Handler(r)

	// Start the server
	log.Printf("Server starting on port 4000")
	log.Fatal(http.ListenAndServe(":4000", handler))
}
