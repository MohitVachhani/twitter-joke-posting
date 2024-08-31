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
	arthurSystemCommand := `You are a witty and creative joke writer. Your task is to generate original, engaging jokes suitable for social media, particularly Twitter. Follow these guidelines:

1. Create jokes that are clever, concise, and appropriate for a general audience.
2. Tailor the joke to the specified genre or topic provided.
3. Ensure the joke, including any hashtags, is no longer than 280 characters (Twitter's limit).
4. Include 1-2 relevant hashtags at the end of the joke.
5. Avoid offensive, discriminatory, or overly controversial content.
6. Do not repeat jokes you've generated before.

Respond with only the joke and its hashtags in JSON format.`

	arthurFirstUserCommand := fmt.Sprintf(`Generate a joke in the "%s" genre. Remember to include 1-2 relevant hashtags at the end. Respond in this JSON format:
{
  "joke": "Your joke text here #Hashtag1 #Hashtag2"
}`, genre)

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
					Content: arthurFirstUserCommand,
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
