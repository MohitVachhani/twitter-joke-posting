package tweet

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	jokegenerationsvc "twitterjokeposting/service/joke"

	"github.com/dghubble/oauth1"
	"github.com/g8rswimmer/go-twitter/v2"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func CreateTwitterClient() twitter.Client {
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

	return client
}

func GenerateJokeAndTweetIt() string {
	ctx := context.Background()

	jokeGenerator := jokegenerationsvc.NewJokeGenerator("Programming")
	generatedJoke := jokeGenerator.GenerateJoke()

	twitterClient := CreateTwitterClient()
	payload, err := twitterClient.CreateTweet(ctx, twitter.CreateTweetRequest{
		Text: generatedJoke,
	})

	if err != nil {
		log.Fatal(err)
	}

	return payload.Tweet.Text
}
