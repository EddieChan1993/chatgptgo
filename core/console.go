package core

import (
	"chatgptgo/gpt"
	goRuntime "chatgptgo/util"
	"context"
	"fmt"
)

func InitConsole() {
	goRuntime.GoRun(func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Println("输入问题?")
				ask := ""
				fmt.Scanf("%s", &ask)
				if ask == "" {
					continue
				}
				resp := gpt.GptApi.AskGptStream(ask)
				fmt.Printf("robotAI:%s\n", resp)
				ask = ""
				fmt.Println("-----------------------")
			}
		}
	})
}
