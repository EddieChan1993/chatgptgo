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
				fmt.Print("我:")
				ask := ""
				fmt.Scanf("%s", &ask)
				if ask == "" {
					continue
				}
				ask = strings.TrimSpace(ask)
				builderAsk.WriteString("\nHuman:" + ask)
				//fmt.Printf("问题:%s\n", builderAsk.String())
				resp := gpt.GptApi.AskGpt(builderAsk.String())
				answer := strings.TrimSpace(resp)
				answer = strings.Trim(answer, "\n")
				fmt.Printf("%s\n", answer)
				if strings.Index(answer, "AI:") == -1 {
					builderAsk.WriteString("\nAI:" + answer)
				} else {
					builderAsk.WriteString("\n" + answer)
				}
				ask = ""
				fmt.Println("-----------------------")
			}
		}
	})
}
