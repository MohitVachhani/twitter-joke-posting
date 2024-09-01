package router

import (
	"net/http"
	"twitterjokeposting/handler"
)

type RouteDefinition struct {
	RouteName  string
	MethodType string
	Handler    func(http.ResponseWriter, *http.Request)
}

func SetRouteDefinitions() *[]RouteDefinition {
	routeDefinitions := []RouteDefinition{}

	routeDefinitions = append(routeDefinitions, RouteDefinition{
		RouteName:  "/",
		MethodType: "GET",
		Handler:    homeHandler,
	})
	routeDefinitions = append(routeDefinitions, RouteDefinition{
		RouteName:  "/about",
		MethodType: "GET",
		Handler:    aboutHandler,
	})
	routeDefinitions = append(routeDefinitions, RouteDefinition{
		RouteName:  "/generateJoke",
		MethodType: "POST",
		Handler:    handler.GenerateJoke,
	})
	routeDefinitions = append(routeDefinitions, RouteDefinition{
		RouteName:  "/tweetIt",
		MethodType: "POST",
		Handler:    handler.TweetIt,
	})
	routeDefinitions = append(routeDefinitions, RouteDefinition{
		RouteName:  "/tweetJoke",
		MethodType: "POST",
		Handler:    handler.GenerateJokeAndTweetIt,
	})
	routeDefinitions = append(routeDefinitions, RouteDefinition{
		RouteName:  "/login",
		MethodType: "POST",
		Handler:    handler.Login,
	})
	routeDefinitions = append(routeDefinitions, RouteDefinition{
		RouteName:  "/generate-quiz",
		MethodType: "POST",
		Handler:    handler.GenerateQuiz,
	})

	return &routeDefinitions
}

// Handler for the home route
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

// Handler for the about route
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Scheduling jokes"))
}
