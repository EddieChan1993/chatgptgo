package main

import (
	"chatgptgo/core"
	"chatgptgo/openai"
	goRuntime "chatgptgo/util"
	_ "embed"
	"os"
	"os/signal"
	"syscall"
)

//go:embed ttfPath
var ttfPath string

func init() {
	os.Setenv("FYNE_FONT", ttfPath)
}

func main() {
	goRuntime.InitGoRuntime()
	openai.InitOpenAi()
	//openai.CreateImgUrl("海边的落日，eric clapton站在悬崖边上，弹奏着fender电吉他，画面风格是梵高风格")
	core.InitGui()
	//core.InitConsole()
	//core.InitWxChat()
}

func waitExit() {
	c := make(chan os.Signal, syscall.SIGKILL) // 定义一个信号的通道
	signal.Notify(c, syscall.SIGINT)           // 转发键盘中断信号到c
	<-c
	goRuntime.CloseGoRuntime()
}
