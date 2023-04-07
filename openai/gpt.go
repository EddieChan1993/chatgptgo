package openai

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
	"strings"
)

//go:embed token
var token string
var builderAsk strings.Builder

const Ai = "AI:"

func AskGpt(content string) (string, error) {
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
		return "", err
	}
	data := resp.Choices[0].Message.Content
	return filedContent(data), nil
}

func filedContent(content string) string {
	answer := strings.Trim(content, "Bot:")
	answer = strings.Trim(answer, "Robot:")
	answer = strings.Trim(answer, "Computer:")
	answer = strings.Trim(answer, "回答：")
	return answer
}

func AskGptStream(content string, fn func(answer string, err error)) {
	req := openai.ChatCompletionRequest{
		Model:  openai.GPT3Dot5Turbo0301,
		Stop:   []string{"Human:", " " + Ai}, //连续发问的标志词
		Stream: true,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
	}
	stream, err := openAiIns.openAiClient.CreateChatCompletionStream(openAiIns.ctx, req)
	if err != nil {
		fmt.Printf("CreateCompletionStream Err %v\n", err)
		fn("", err)
		return
	}
	defer stream.Close()
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Printf("Stream Err %v\n", err)
			fn("", err)
			break
		}
		answerRsp := response.Choices
		if len(answerRsp) != 0 {
			builderAsk.WriteString(answerRsp[0].Delta.Content)
			fn(answerRsp[0].Delta.Content, nil)
		}
	}
	return
}

func GetAskContent(ask string) string {
	ask = strings.TrimSpace(ask)
	builderAsk.WriteString("\nHuman:" + ask)
	return builderAsk.String()
}

func ClearAsk() {
	builderAsk.Reset()
}
