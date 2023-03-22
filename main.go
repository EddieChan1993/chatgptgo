package main

import (
	"chatgptgo/openai"
	goRuntime "chatgptgo/util"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	yeheiTTF := "/Users/duanchengwen/go/bin/font/yahei.ttf"
	os.Setenv("FYNE_FONT", yeheiTTF)
}

func main() {
	goRuntime.InitGoRuntime()
	openai.InitOpenAi()
	openai.CreateImgUrl("海边的落日，eric clapton站在悬崖边上，弹奏着fender电吉他，画面风格是梵高风格")
	//core.InitGui()
	//core.InitConsole()
	//core.InitWxChat()
}

func waitExit() {
	c := make(chan os.Signal, syscall.SIGKILL) // 定义一个信号的通道
	signal.Notify(c, syscall.SIGINT)           // 转发键盘中断信号到c
	<-c
	goRuntime.CloseGoRuntime()
}
