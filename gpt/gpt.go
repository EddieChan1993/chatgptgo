package gpt

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"io"
	"strings"
)

//go:embed token
var token string
var builderAsk strings.Builder

const Ai = "AI:"

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
		Stop:             []string{"Human:", " " + Ai}, //连续发问的标志词
		Prompt:           content,
	}
	resp, err := g.client.CreateCompletion(g.ctx, req)
	if err != nil {
		fmt.Printf("CreateCompletion %v", err)
		return ""
	}
	data := resp.Choices[0].Text
	answer := strings.TrimSpace(data)
	answer = strings.Trim(answer, "\n")
	answer = strings.Trim(answer, "Bot:")
	answer = strings.Trim(answer, "Robot:")
	answer = strings.Trim(answer, "Computer:")
	if strings.Index(answer, Ai) == -1 {
		builderAsk.WriteString("\n" + Ai + answer)
		answer = Ai + answer
	} else {
		builderAsk.WriteString("\n" + answer)
	}
	return answer
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
		fmt.Printf("CreateCompletionStream Err %v\n", err)
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
			fmt.Printf("Stream Err %v\n", err)
		}
		if len(response.Choices) != 0 {
			_, _ = builder.WriteString(response.Choices[0].Text)
		}
		//fmt.Printf("Stream response: %v\n", response)
	}
	return strings.TrimSpace(builder.String())
}

func GetAskContent(ask string) string {
	ask = strings.TrimSpace(ask)
	builderAsk.WriteString("\nHuman:" + ask)
	return builderAsk.String()
}
