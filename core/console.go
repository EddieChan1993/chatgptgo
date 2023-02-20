package core

import (
	"chatgptgo/gpt"
	goRuntime "chatgptgo/util"
	"context"
	"fmt"
	"strings"
)

func InitConsole() {
	var builderAsk strings.Builder
	goRuntime.GoRun(func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				builderAsk.Reset()
				return
			default:
				fmt.Println("输入问题?")
				ask := ""
				fmt.Scanf("%s", &ask)
				if ask == "" {
					continue
				}
				ask = strings.TrimSpace(ask)
				builderAsk.WriteString("Human:" + ask)
				//fmt.Printf("问题:%s\n", builderAsk.String())
				resp := gpt.GptApi.AskGpt(ask)
				answer := strings.TrimSpace(resp)
				builderAsk.WriteString("\nAI:" + answer)
				fmt.Printf("%s\n", answer)
				ask = ""
				fmt.Println("-----------------------")
			}
		}
	})
}
