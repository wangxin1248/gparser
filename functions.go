package gparser

import (
	"errors"
	"go/ast"
)

// 注册可执行函数
var funcNameMap = map[string]func(args []ast.Expr, data map[string]interface{}) interface{}{}

func init() {
	funcNameMap = map[string]func(args []ast.Expr, data map[string]interface{}) interface{}{
		"in_array": inArray,
	}
}

// inArray 判断变量是否存在在数组中
func inArray(args []ast.Expr, data map[string]interface{}) interface{} {
	// 规则表达式中的变量
	param := eval(args[0], data)
	vRange, ok := args[1].(*ast.CompositeLit)
	if !ok {
		return errors.New("func in_array 2ed params is not a composite lit")
	}

	// 规则表达式中数组里的元素
	eltNodes := make([]interface{}, 0, len(vRange.Elts))
	for _, p := range vRange.Elts {
		elt := eval(p, data)
		eltNodes = append(eltNodes, elt)
	}

	for _, node := range eltNodes {
		switch node.(type) {
		case int64:
			param, err := castType(param, TypeInt64)
			if err != nil {
				return false
			}
			paramInt64, paramOk := param.(int64)
			nodeInt64, nodeOk := node.(int64)
			if !paramOk || !nodeOk {
				return false
			}
			if nodeInt64 == paramInt64 {
				return true
			}
		case string:
			param, err := castType(param, TypeString)
			if err != nil {
				return false
			}
			nodeString, nodeOk := node.(string)
			paramString, paramOk := param.(string)
			if !paramOk || !nodeOk {
				return false
			}
			if nodeString == paramString {
				return true
			}
		}
	}
	return false
}
