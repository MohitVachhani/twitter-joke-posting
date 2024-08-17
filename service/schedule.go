package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"twitterjokeposting/repo"
	"twitterjokeposting/util"

	"github.com/dghubble/oauth1"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/samber/lo"
	openai "github.com/sashabaranov/go-openai"
)

func getJokesToPublishToday() int {
	return util.RandomNumberGenerator(10)
}

func ScheduleJokeForToday() {
	// Create background context
	ctx := context.Background()

	// Clear Schedule Job database for today
	repo.ClearAllScheduledJobs()

	// Get How many jokes we want to publish today
	jokesToPublish := getJokesToPublishToday()

	// Get that much random times today
	randomTimesToPublishJoke := []int{}
	for iterator := 0; iterator < jokesToPublish; iterator++ {
		randomTime := util.RandomNumberGenerator(24)
		randomTimesToPublishJoke = append(randomTimesToPublishJoke, randomTime)
	}
	sort.Ints(randomTimesToPublishJoke)
	randomTimesToPublishJoke = lo.Uniq[int](randomTimesToPublishJoke)

	// Insert that much documents into database today
	repo.CreateScheduledJobs(randomTimesToPublishJoke, ctx)
}

func GetAllScheduledJokes() []repo.ScheduleJobSchema {
	return repo.GetScheuledJobs()
}

func GenerateJoke() string {
	client := openai.NewClient(os.Getenv("CHAT_GPT_API_SECRET"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Tell me a joke! and return the output in JSON format",
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

	return resp.Choices[0].Message.Content
}

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func TweetIt() string {
	ctx := context.Background()

	// setting up the oauth1 library with our api key and secret
	config := oauth1.NewConfig(os.Getenv("TWITTER_API_KEY"), os.Getenv("TWITTER_API_SECRET"))
	// setting up the token
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))

	// creating a HTTP Client with the token
	httpClient := config.Client(context.Background(), token)

	// creates a twitter client from the installed package with a httpClient
	// token1 := flag.String("token", os.Getenv("TWITTER_BEARER_TOKEN"), "twitter API token")
	token1 := os.Getenv("TWITTER_BEARER_TOKEN")

	client := twitter.Client{
		Client: httpClient,
		Host:   "https://api.twitter.com",
		Authorizer: authorize{
			Token: token1,
		},
	}

	payload, err := client.CreateTweet(ctx, twitter.CreateTweetRequest{
		Text: "Hello, World! This is the first message from the go twitter bot-3.",
	})

	if err != nil {
		log.Fatal(err)
	}

	return payload.Tweet.Text
}
