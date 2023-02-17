package main

import (
	"chatgptgo/core"
	"chatgptgo/gpt"
	goRuntime "chatgptgo/util"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	goRuntime.InitGoRuntime()
	gpt.InitGpt()
	core.InitConsole()
	//core.InitWxChat()
	waitExit()
}

func waitExit() {
	c := make(chan os.Signal, syscall.SIGKILL) // 定义一个信号的通道
	signal.Notify(c, syscall.SIGINT)           // 转发键盘中断信号到c
	<-c
	goRuntime.CloseGoRuntime()
}
