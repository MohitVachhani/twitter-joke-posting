package handler

import (
	"encoding/json"
	"net/http"
	"twitterjokeposting/service"
)

type APIResponse struct {
	Success      bool
	Payload      any
	ErrorMessage *string
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
