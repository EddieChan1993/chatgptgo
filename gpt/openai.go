package gpt

import (
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/sashabaranov/go-openai"
)

type gpt struct {
	client       *gogpt.Client
	ctx          context.Context
	openAiClient *openai.Client
}

var openAiIns *gpt

func InitOpenAi() {
	c := gogpt.NewClient(token)
	ctx := context.Background()
	openAiClient := openai.NewClient(token)
	openAiIns = &gpt{
		client:       c,
		ctx:          ctx,
		openAiClient: openAiClient,
	}
}
