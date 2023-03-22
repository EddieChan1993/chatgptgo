package core

import (
	"chatgptgo/openai"
	goRuntime "chatgptgo/util"
	"context"
	"fmt"
	"strings"
)

var builderAsk strings.Builder

func InitConsole() {
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
				msg := openai.GetAskContent(ask)
				answer := openai.AskGpt(msg)
				fmt.Printf("%s\n", answer)
				fmt.Println("-----------------------")
			}
		}
	})
}
