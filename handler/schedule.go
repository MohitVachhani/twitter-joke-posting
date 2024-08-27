package handler

import (
	"encoding/json"
	"net/http"
	"twitterjokeposting/service"
	jokegenerationsvc "twitterjokeposting/service/joke"
	"twitterjokeposting/service/tweet"
)

type APIResponse struct {
	Success      bool
	Payload      any
	ErrorMessage *string
}

type GenerateJokeAndTweetItInput struct {
	Genre string `json:"genre"`
}

func ScheduleJokeForTodayController(w http.ResponseWriter, r *http.Request) {
	service.ScheduleJokeForToday()
	apiResponse := APIResponse{
		Success: true,
	}

	jsonData, _ := json.Marshal(apiResponse)

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetAllScheduledJokes(w http.ResponseWriter, r *http.Request) {
	scheduleJobs := service.GetAllScheduledJokes()
	apiResponse := APIResponse{
		Success: true,
		Payload: scheduleJobs,
	}

	jsonData, _ := json.Marshal(apiResponse)

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GenerateJoke(w http.ResponseWriter, r *http.Request) {
	var body GenerateJokeAndTweetItInput
	json.NewDecoder(r.Body).Decode(&body)

	jokeGenerator := jokegenerationsvc.NewJokeGenerator(body.Genre)
	generatedJoke := jokeGenerator.GenerateJoke()

	apiResponse := APIResponse{
		Success: true,
		Payload: generatedJoke,
	}

	jsonData, _ := json.Marshal(apiResponse)

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func TweetIt(w http.ResponseWriter, r *http.Request) {
	tweet := service.TweetIt()
	apiResponse := APIResponse{
		Success: true,
		Payload: tweet,
	}

	jsonData, _ := json.Marshal(apiResponse)

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GenerateJokeAndTweetIt(w http.ResponseWriter, r *http.Request) {
	var body GenerateJokeAndTweetItInput
	json.NewDecoder(r.Body).Decode(&body)

	generatedJoke := tweet.GenerateJokeAndTweetIt(body.Genre)
	apiResponse := APIResponse{
		Success: true,
		Payload: generatedJoke,
	}

	jsonData, _ := json.Marshal(apiResponse)

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
