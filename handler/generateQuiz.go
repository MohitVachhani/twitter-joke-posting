package handler

import (
	"encoding/json"
	"net/http"
	schemaInterface "twitterjokeposting/interface"
	questiongenerationsvc "twitterjokeposting/service/question"
)

type GenerateQuizRequest struct {
	ThemeName         string `json:"themeName"`
	NumberOfQuestions int    `json:"numberOfQuestions"`
}

type GenerateQuizResponse struct {
	Success   bool                             `json:"success"`
	Questions []schemaInterface.QuestionSchema `json:"questions,omitempty"`
	Error     string                           `json:"error,omitempty"`
}

func GenerateQuiz(w http.ResponseWriter, r *http.Request) {
	var req GenerateQuizRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ThemeName == "" || req.NumberOfQuestions <= 0 {
		http.Error(w, "Invalid input parameters", http.StatusBadRequest)
		return
	}

	questionGenerator := questiongenerationsvc.NewQuestionGenerator(req.ThemeName, req.NumberOfQuestions)
	questions, err := questionGenerator.GenerateQuestions()
	if err != nil {
		http.Error(w, "Failed to generate questions", http.StatusInternalServerError)
		return
	}

	response := GenerateQuizResponse{
		Success:   true,
		Questions: questions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
