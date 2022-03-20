# gparser

基于 golang 原生语法解析器（[parser](https://pkg.go.dev/go/parser)）实现的轻量级规则引擎。支持操作：

- 规则匹配：`gparser.Match(ruleStr, params)`

下载方式：

```shell
go get github.com/wangxin1248/gparser
```

使用方式：

```go
import "github.com/wangxin1248/gparser"

ruleStr := "!(a == 1 && b == 2 && c == "test" && d == false)"

// 匹配变量
params := map[string]interface{}{
    "a": 1,
    "b": 2,
    "c": "test",
    "d": true,
}

result, err := gparser.Match(ruleStr, params)
fmt.Println(result)
```

支持类型

- int
- int64
- string
- bool

支持操作

- `!表达式`：支持一元表达式
- `&&`：支持多个表达式逻辑与
- `||`：支持多个表达式逻辑或
- `()`：支持表达式括号包裹
- `==`：int、int64、string、bool支持
- `!=`：int、int64、string、bool支持
- `>`：int、int64支持
- `<`：int、int64支持
- `>=`：int、int64支持
- `<=`：int、int64支持
- `+`：int、int64支持
- `-`：int、int64支持
- `*`：int、int64支持
- `/`：int、int64支持

性能对比

```shell
BenchmarkGParser_Match-8              127189          8912   ns/op     // gparser
BenchmarkGval_Match-8                 63584           18358  ns/op     // gval
BenchmarkGovaluateParser_Match-8      13628           86955  ns/op     // govaluate
BenchmarkYqlParser_Match-8            10364           112481 ns/op     // yql
```