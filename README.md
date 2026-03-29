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

函数用法示例：

```go
ruleStr = `in_array(c, []int{1,2,3}) && max(a, b, 5) == 5 && min(a, b, 0) == 0 && d < 5.5`
params = map[string]interface{}{
    "a": 1,
    "b": 3,
    "c": 2,
    "d": 4.2,
}
result, err = gparser.Match(ruleStr, params)
fmt.Println(result) // true
```

支持类型

- int
- int64
- float（`float32`, `float64`, `json.Number` 形式也支持转换）
- string
- bool

支持操作

- `!表达式`：支持一元表达式
- `&&`：支持多个表达式逻辑与
- `||`：支持多个表达式逻辑或
- `()`：支持表达式括号包裹
- `==`：int、int64、float、string、bool支持
- `!=`：int、int64、float、string、bool支持
- `>`：int、int64、float支持
- `<`：int、int64、float支持
- `>=`：int、int64、float支持
- `<=`：int、int64、float支持
- `+`：int、int64、float支持
- `-`：int、int64、float支持
- `*`：int、int64、float支持
- `/`：int、int64、float支持
- `函数调用`：支持 `in_array`, `max`, `min`

常见问题/报错模式

- `func in_array 2ed params is not a composite lit`
  - 说明：`in_array` 第二个参数需为复合字面量，如 `[]int{1,2,3}`。
  - 处理：确保表达式语法是 `in_array(var, []类型{...})`。

- `func max requires at least one argument` / `func min requires at least one argument`
  - 说明：`max`/`min` 必须传入至少一个参数。
  - 处理：补充参数，如 `max(a,b)`。

- `unsupported binary operator: ...`
  - 说明：当前不支持例如 `^` 之类运算符。
  - 处理：改用支持的运算符 `+ - * / == != > < >= <= && ||`。

- `...eval failed`（常见于类型不匹配）
  - 说明：变量类型与操作类型不一致，例如对 `string` 做 `>`。
  - 处理：检查数据类型、使用类型兼容值或转换为支持类型。

- 除数为 0 时除法结果为 0（`/` 逻辑行为）
  - 说明：为了避免 panic，`x / 0` 会返回 `0`；浮点也是 `0.0`。
  - 处理：提前判断除数是否为 0，或在表达式中避免该分支。

性能对比

```shell
BenchmarkGParser_Match-8              127189          8912   ns/op     // gparser
BenchmarkGval_Match-8                 63584           18358  ns/op     // gval
BenchmarkGovaluateParser_Match-8      13628           86955  ns/op     // govaluate
BenchmarkYqlParser_Match-8            10364           112481 ns/op     // yql
```