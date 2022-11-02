# goroutinen 静态检查recover 插件
### 用来检查某些场景下goroutine内需要捕获panic,防止不必要的程序异常停止

### 使用:

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
## 集成到[golangci-lint](https://golangci-lint.run)使用

1. 创建文件 golangci-lint/pkg/golinters/gorecover.go

```go
package golinters

import (
"github.com/wangxudong123/gorecover/analyzer"
"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
"golang.org/x/tools/go/analysis"
)

func NewGoRecoverCheck() *goanalysis.Linter {
    return goanalysis.NewLinter(
            analyzer.Analyzer.Name,
            analyzer.Analyzer.Doc,
            []*analysis.Analyzer{analyzer.Analyzer},
            nil,
        ).WithLoadMode(goanalysis.LoadModeSyntax)
}

   ```
2. 导入配置 golangci-lint/pkg/lint/lintersdb/manager.go

```go
func (m Manager) GetAllSupportedLinterConfigs() []*linter.Config {

	// 一坨代码
	...
	
	// 导入NewGoRecoverCheck配置 
	linter.NewConfig(golinters.NewGoRecoverCheck()).
		WithSince("v1.0.0").
		WithPresets(linter.PresetStyle, linter.PresetBugs).
		WithURL("https://github.com/wangxudong123/gorecover"),
	
}
```

3. 在你的项目中`.golangci.yml`中添加`gorecover` 静态检查项
```yaml
linters:
  disable-all: true # 关闭全部检查
  enable: # 打开下面的检查选项
    - gorecover
```

[.golangci.yml](https://golangci-lint.run/usage/configuration) 相关配置