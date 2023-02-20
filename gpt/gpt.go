package gpt

import (
	"context"
	_ "embed"
	"errors"
	gogpt "github.com/sashabaranov/go-gpt3"
	"io"
	"log"
	"strings"
)

//go:embed token
var token string

type gpt struct {
	client *gogpt.Client
	ctx    context.Context
}

var GptApi IGptApi

type IGptApi interface {
	AskGpt(content string) string
	AskGptStream(content string) string
}

func InitGpt() {
	c := gogpt.NewClient(token)
	ctx := context.Background()

	GptApi = &gpt{
		client: c,
		ctx:    ctx,
	}
}

func (g *gpt) AskGpt(content string) string {
	req := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Temperature:      0,
		MaxTokens:        1000,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0.6,
		Stop:             []string{"Human:", " AI:"}, //连续发问的标志词
		Prompt:           content,
	}
	resp, err := g.client.CreateCompletion(g.ctx, req)
	if err != nil {
		return ""
	}
	return resp.Choices[0].Text
}

func (g *gpt) AskGptStream(content string) string {
	req := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Temperature:      0,
		MaxTokens:        1000,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		BestOf:           1,
		Prompt:           content,
		Stream:           true,
	}
	stream, err := g.client.CreateCompletionStream(g.ctx, req)
	if err != nil {
		log.Fatalf("CreateCompletionStream Err %v\n", err)
		return ""
	}
	defer stream.Close()
	var builder strings.Builder
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalf("Stream Err %v\n", err)
		}
		if len(response.Choices) != 0 {
			_, _ = builder.WriteString(response.Choices[0].Text)
		}
		//fmt.Printf("Stream response: %v\n", response)
	}
	return strings.TrimSpace(builder.String())
}
