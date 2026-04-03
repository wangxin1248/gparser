# gparser

A lightweight rule engine implemented based on Golang's native syntax parser ([parser](https://pkg.go.dev/go/parser)). Supports operations:

- Rule matching: `gparser.Match(ruleStr, params)`
- Expression evaluation: `gparser.Evaluate(expr, params)`

Installation:

```shell
go get github.com/wangxin1248/gparser
```

Usage:

```go
import "github.com/wangxin1248/gparser"

ruleStr := "!(a == 1 && b == 2 && c == \"test\" && d == false)"

// Match variables
params := map[string]interface{}{
    "a": 1,
    "b": 2,
    "c": "test",
    "d": true,
}

result, err := gparser.Match(ruleStr, params)
fmt.Println(result)
```

Function usage example:

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

Supported types

- int
- int64
- float (`float32`, `float64`, `json.Number` forms are also supported for conversion)
- string
- bool

Supported operations

- `!expression`: Supports unary expressions
- `&&`: Supports logical AND of multiple expressions
- `||`: Supports logical OR of multiple expressions
- `()`: Supports expression parentheses wrapping
- `==`: Supported for int, int64, float, string, bool
- `!=`: Supported for int, int64, float, string, bool
- `>`: Supported for int, int64, float
- `<`: Supported for int, int64, float
- `>=`: Supported for int, int64, float
- `<=`: Supported for int, int64, float
- `+`: Supported for int, int64, float
- `-`: Supported for int, int64, float
- `*`: Supported for int, int64, float
- `/`: Supported for int, int64, float
- Function calls: Supports `in_array`, `max`, `min`

Common issues/error patterns

- `func in_array 2ed params is not a composite lit`
  - Description: The second parameter of `in_array` must be a composite literal, such as `[]int{1,2,3}`.
  - Handling: Ensure the expression syntax is `in_array(var, []type{...})`.

- `func max requires at least one argument` / `func min requires at least one argument`
  - Description: `max`/`min` must have at least one argument.
  - Handling: Add arguments, such as `max(a,b)`.

- `unsupported binary operator: ...`
  - Description: Operators like `^` are not currently supported.
  - Handling: Use supported operators `+ - * / == != > < >= <= && ||`.

- `...eval failed` (commonly due to type mismatch)
  - Description: Variable types do not match the operation types, e.g., doing `>` on `string`.
  - Handling: Check data types, use type-compatible values, or convert to supported types.

- Division by zero results in 0 (`/` logic behavior)
  - Description: To avoid panic, `x / 0` returns `0`; for floats, it's `0.0`.
  - Handling: Check if the divisor is 0 in advance, or avoid that branch in the expression.

Performance comparison

```shell
BenchmarkGParser_Match-8              127189          8912   ns/op     // gparser
BenchmarkGval_Match-8                 63584           18358  ns/op     // gval
BenchmarkGovaluateParser_Match-8      13628           86955  ns/op     // govaluate
BenchmarkYqlParser_Match-8            10364           112481 ns/op     // yql
```