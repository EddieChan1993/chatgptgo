package openai

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"io"
	"strings"
	"time"
)

//go:embed token
var token string
var builderAsk strings.Builder

const Ai = "AI:"

func AskGpt(content string) (string, error) {
	resp, err := openAiIns.openAiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT40314,
			Stop:  []string{"Human:", " " + Ai}, //连续发问的标志词
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Name:    "Hat",
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

func AskGptStream(content string, fn func(answer string, err error)) {
	resCh := make(chan openai.ChatCompletionStreamResponse)
	errCh := make(chan error)
	stream, err := callStream(errCh, resCh, content)
	if err != nil {
		fn("", err)
		return
	}
	wait := 30 * time.Second
	ticker := time.NewTicker(wait)
	defer func() {
		stream.Close()
		ticker.Stop()
	}()
	builderAsk.WriteString("\n" + Ai)
	for {
		select {
		case <-ticker.C:
			fn("", fmt.Errorf("响应超时"))
			return
		case errRes := <-errCh:
			fn("", errRes)
		case answerRes := <-resCh:
			if answerRes.ID == "Ok" {
				return
			}
			ticker.Reset(wait)
			answerRsp := answerRes.Choices
			if len(answerRsp) != 0 {
				builderAsk.WriteString(answerRsp[0].Delta.Content)
				fn(answerRsp[0].Delta.Content, nil)
			}
		}
	}
}

func callStream(errCh chan error, resCh chan openai.ChatCompletionStreamResponse, content string) (stream *openai.ChatCompletionStream, err error) {
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
	stream, err = openAiIns.openAiClient.CreateChatCompletionStream(openAiIns.ctx, req)
	if err != nil {
		fmt.Printf("CreateCompletionStream Err %v\n", err)
		return nil, err
	}
	go func() {
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				resCh <- openai.ChatCompletionStreamResponse{
					ID: "Ok",
				}
				break
			}
			if err != nil {
				fmt.Printf("Stream Err %v\n", err)
				errCh <- err
				break
			}
			resCh <- response
		}
	}()
	return stream, err
}

func GetAskContent(ask string) string {
	ask = strings.TrimSpace(ask)
	builderAsk.WriteString("\nHuman:" + ask)
	return builderAsk.String()
}

func ClearAsk() {
	builderAsk.Reset()
}

func filedContent(content string) string {
	answer := strings.Trim(content, "Bot:")
	answer = strings.Trim(answer, "Robot:")
	answer = strings.Trim(answer, "Computer:")
	answer = strings.Trim(answer, "回答：")
	return answer
}
