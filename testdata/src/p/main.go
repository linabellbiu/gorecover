package main

import (
	"fmt"
	"runtime"
)

func goFuncWithRecover() {
	go func() { // want "goroutine should have recover in defer func"
	}()

	go func() {
		defer func() {
			recover()
		}()
	}()

	go HandlerPanic(func() {
		panic("panic2")
	})
}

func HandlerPanic(f func()) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Println(r)
		}
	}()

	f()
}

func main() {
	goFuncWithRecover()
}
