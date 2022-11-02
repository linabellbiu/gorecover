# goroutinen 静态检查recover 插件
### 用来检查某些场景下goroutine内需要捕获panic,防止不必要的程序异常停止

使用:

```
./gorecover ../../testdata/src/p/main.go

// 输出检查的异常
gorecover/testdata/src/p/main.go:9:5: goroutine should have recover in defer func

```

```go
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
```


