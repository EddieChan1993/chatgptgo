package openai

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/sashabaranov/go-openai"
	"io"
	"strings"
)

//go:embed token
var token string
var builderAsk strings.Builder

const Ai = "AI:"

func AskGpt(content string) string {
	resp, err := openAiIns.openAiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0301,
			Stop:  []string{"Human:", " " + Ai}, //连续发问的标志词
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}
	data := resp.Choices[0].Message.Content
	return filedContent(data)
}

func filedContent(content string) string {
	answer := strings.TrimSpace(content)
	answer = strings.Trim(answer, "\n")
	answer = strings.TrimSpace(answer)
	answer = strings.Trim(answer, "Bot:")
	answer = strings.Trim(answer, "Robot:")
	answer = strings.Trim(answer, "Computer:")
	answer = strings.Trim(answer, "回答：")
	if strings.Index(answer, Ai) == -1 {
		builderAsk.WriteString("\n" + Ai + answer)
		answer = Ai + answer
	} else {
		builderAsk.WriteString("\n" + answer)
	}
	return answer
}

func AskGptStream(content string) string {
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
	stream, err := openAiIns.client.CreateCompletionStream(openAiIns.ctx, req)
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

func ClearAsk() {
	builderAsk.Reset()
}
