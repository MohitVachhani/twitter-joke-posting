package arthur

import (
	"github.com/sashabaranov/go-openai"
)

type ArthurConfigType struct {
	ApiSecret string
}

func NewArthurConfigType(apiSecret string) *ArthurConfigType {
	return &ArthurConfigType{
		ApiSecret: apiSecret,
	}
}

func (a *ArthurConfigType) GetClient() *openai.Client {
	client := openai.NewClient(a.ApiSecret)
	return client
}
