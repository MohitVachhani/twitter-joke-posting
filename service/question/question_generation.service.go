package questiongenerationsvc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	schemaInterface "twitterjokeposting/interface"
	"twitterjokeposting/service/arthur"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/sashabaranov/go-openai"
)

type QuestionGenerator struct {
	ThemeName         string
	NumberOfQuestions int
}

type QuestionGeneratorChatGPTResponseType struct {
	QuestionText   string   `json:"questionText"`
	Options        []string `json:"options"`
	CorrectIndexes []int    `json:"correctIndexes"`
}

type QuizGeneratorChatGPTResponseType struct {
	Quiz []QuestionGeneratorChatGPTResponseType `json:"quiz"`
}

func NewQuestionGenerator(themeName string, numberOfQuestions int) *QuestionGenerator {
	return &QuestionGenerator{
		ThemeName:         themeName,
		NumberOfQuestions: numberOfQuestions,
	}
}

func (questionPromptResponse *QuestionGeneratorChatGPTResponseType) createQuestionsFromPromptResponse(themeID string) schemaInterface.QuestionSchema {
	options := []schemaInterface.OptionSchema{}
	correctOptions := questionPromptResponse.CorrectIndexes
	for i := range questionPromptResponse.Options {
		isCorrect := false
		if lo.Contains(correctOptions, i) {
			isCorrect = true
		}
		options = append(options, schemaInterface.OptionSchema{
			ID:        generateUniqueID(),
			Text:      questionPromptResponse.Options[i],
			IsCorrect: isCorrect,
		})
	}

	return schemaInterface.QuestionSchema{
		QuestionText: questionPromptResponse.QuestionText,
		Options:      options,
		ThemeID:      themeID,
		BaseSchema: schemaInterface.BaseSchema{
			ID:        generateUniqueID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		QuestionType: schemaInterface.SingleSelect,
		Attachment:   nil,
		Status:       schemaInterface.Active,
	}
}

func (qg *QuestionGenerator) GenerateQuestions() ([]schemaInterface.QuestionSchema, error) {
	arthurConfig := arthur.NewArthurConfigType(os.Getenv("CHAT_GPT_API_SECRET"))
	arthurClient := arthurConfig.GetClient()

	arthurSystemCommand :=
		fmt.Sprintf(`You are an educational content creator. Your task is to Generate a quiz with different questions on theme: %s. Follow these guidelines, generate original, engaging questions suitable for quizzes.`, qg.ThemeName)
	arthurInstructionCommands := []string{
		`Use tone and style to create more engaging and personalized responses`,
		`Incorporate humor or wit when appropriate`,
		`Create questions that are clear, concise, and appropriate for the specified theme.`,
		`Provide 4 options for each question, with one or more correct answers.`,
		`Ensure the questions are suitable for a general audience.`,
		`Avoid offensive, discriminatory, or overly controversial content.`,
		`Do not repeat questions you've generated before.`,
	}
	arthurResponseSystemCommand := `Respond with the questions and please return in array json format. Every json should have questionText, options array, and which option indexes are correct. The json format should be, key should be quiz and value should be array of questions. Multiple options can be correct so inside question json, correctIndexes should be an array.`

	arthurMessages := []openai.ChatCompletionMessage{}
	arthurMessages = append(arthurMessages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: arthurSystemCommand,
	})
	for i := range arthurInstructionCommands {
		arthurMessages = append(arthurMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: arthurInstructionCommands[i],
		})
	}
	arthurMessages = append(arthurMessages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: arthurResponseSystemCommand,
	})
	arthurMessages = append(arthurMessages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: fmt.Sprintf(`Generate a quiz on the theme "%s" with %d questions. Respond in array JSON format.`, qg.ThemeName, qg.NumberOfQuestions),
	})

	resp, err := arthurClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			MaxTokens: 1800,
			Messages:  arthurMessages,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		},
	)

	if err != nil {
		return nil, err
	}

	var questions []QuestionGeneratorChatGPTResponseType
	var quiz QuizGeneratorChatGPTResponseType
	message := resp.Choices[0].Message
	err = json.Unmarshal([]byte(message.Content), &quiz)
	if err != nil {
		return nil, err
	}
	questions = quiz.Quiz

	var requiredQuestions []schemaInterface.QuestionSchema

	themeId := generateUniqueID()
	// Set common fields for each question
	for i := range questions {
		requiredQuestions = append(requiredQuestions, questions[i].createQuestionsFromPromptResponse(themeId))
	}

	return requiredQuestions, nil
}

func generateUniqueID() string {
	return uuid.New().String()
}