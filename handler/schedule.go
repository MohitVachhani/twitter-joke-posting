package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"twitterjokeposting/service"
	jokegenerationsvc "twitterjokeposting/service/joke"
	"twitterjokeposting/service/tweet"

	"github.com/golang-jwt/jwt/v5"
)

type APIResponse struct {
	Success      bool
	Payload      any
	ErrorMessage *string
}

type GenerateJokeAndTweetItInput struct {
	Genre string `json:"genre"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success   bool       `json:"success"`
	ErrorCode string     `json:"errorCode,omitempty"`
	Token     *TokenData `json:"token,omitempty"`
}

type TokenData struct {
	AccessToken string `json:"accessToken"`
}

var jwtKey = []byte("twitterJokePosting6907")

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int, errorCode string) {
	response := LoginResponse{
		Success:   false,
		ErrorCode: errorCode,
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
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

func Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	// Check if the email and password match the expected values
	if loginReq.Email != "mohitvachhani55@gmail.com" || loginReq.Password != "play@1234" {
		sendErrorResponse(w, "Invalid email or password", http.StatusUnauthorized, "EMAIL_PASSWORD_WRONG")
		return
	}

	claims := jwt.MapClaims{
		"email": loginReq.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		sendErrorResponse(w, "Could not generate token", http.StatusInternalServerError, "INTERNAL_SERVER_ERROR")
		return
	}

	response := LoginResponse{
		Success: true,
		Token: &TokenData{
			AccessToken: tokenString,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
