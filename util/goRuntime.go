package goRuntime

import (
	"context"
	"log"
	"runtime"
	"sync"
)

type tGoRuntime struct {
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

var goRuntime *tGoRuntime

func InitGoRuntime() {
	rootCtx, cancel := context.WithCancel(context.Background())
	goRuntime = &tGoRuntime{
		wg:     &sync.WaitGroup{},
		ctx:    rootCtx,
		cancel: cancel,
	}
}

func GoRun(fn func(ctx context.Context)) {
	goRuntime.wg.Add(1)
	go func() {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				log.Println(panicErr)
				PanicStack()
			}
			goRuntime.wg.Done()
		}()
		fn(goRuntime.ctx)
	}()
}

func CloseGoRuntime() {
	goRuntime.cancel()
	goRuntime.wg.Wait()
}

//PanicStack 捕获recover的panic堆栈
func PanicStack() {
	buf := make([]byte, 1<<10)
	runtime.Stack(buf, true)
	log.Println(string(buf))
}
