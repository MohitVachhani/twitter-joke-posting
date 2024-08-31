package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

func authenticateRequest(w http.ResponseWriter, r *http.Request) error {
	token := extractToken(r)
	if token == "" {
		sendErrorResponse(w, "No token provided", http.StatusUnauthorized, "NO_TOKEN")
		return fmt.Errorf("no token provided")
	}

	claims, err := validateToken(token)
	if err != nil {
		sendErrorResponse(w, "Invalid token", http.StatusUnauthorized, "INVALID_TOKEN")
		return fmt.Errorf("invalid token")
	}

	if email, ok := claims["email"].(string); !ok || email != "mohitvachhani55@gmail.com" {
		sendErrorResponse(w, "Unauthorized user", http.StatusUnauthorized, "UNAUTHORIZED_USER")
		return fmt.Errorf("unauthorized user")
	}

	return nil
}

func extractToken(r *http.Request) string {
	// Check in cookie
	if cookie, err := r.Cookie("accessToken"); err == nil {
		return cookie.Value
	}

	// Check in Authorization header
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	// Check in query parameters
	token := r.URL.Query().Get("accessToken")
	if token != "" {
		return token
	}

	// Check in request body
	var body struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
		return body.AccessToken
	}

	return ""
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
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
	if err := authenticateRequest(w, r); err != nil {
		return
	}

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
	if err := authenticateRequest(w, r); err != nil {
		return
	}

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
	if err := authenticateRequest(w, r); err != nil {
		return
	}

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
