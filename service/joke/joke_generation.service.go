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
	arthurSystemCommand :=
		fmt.Sprint("You are working as a content creator of my company, and your task is to create jokes on the topic which I give you. Now the joke can be one liner joke or it can be little bit longer but not longer than 50 words. We also need to make sure that we don't repeat any joke.And now think that we also we need to post this joke on twitter what would be the best hashtag for joke you need to think for that too and let me know... Are you comfortable with this ? And please please return the output in JSON format")
	// fmt.Println(arthurCommand)
	arthurUserCommand :=
		fmt.Sprint("Tell me a joke of Genre, ", genre, " , with proper twitter hashtag and please return the output in JSON format")

	resp, err := arthurClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: arthurSystemCommand,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: arthurUserCommand,
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
