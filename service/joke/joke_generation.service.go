package jokegenerationsvc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"twitterjokeposting/service/arthur"

	"github.com/sashabaranov/go-openai"
)

type JokeGenerator struct {
	Genre string
}

type JokePayload struct {
	Joke string `json:"joke"`
}

func NewJokeGenerator(genre string) *JokeGenerator {
	return &JokeGenerator{
		Genre: genre,
	}
}

func (jg *JokeGenerator) GenerateJoke() string {
	arthurConfig := arthur.NewArthurConfigType(os.Getenv("CHAT_GPT_API_SECRET"))
	arthurClient := arthurConfig.GetClient()

	genre := jg.Genre
	arthurCommand := fmt.Sprint("Tell me a joke of Genre, ", genre, " , with proper twitter hashtag and please return the output in JSON format")
	fmt.Println(arthurCommand)

	resp, err := arthurClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: arthurCommand,
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "error"
	}

	fmt.Println(resp.Choices[0].Message.Content)
	arthurResponse := resp.Choices[0].Message.Content

	jokePayload := &JokePayload{}
	json.Unmarshal([]byte(arthurResponse), jokePayload)

	return jokePayload.Joke
}
