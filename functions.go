package gparser

import (
	"errors"
	"go/ast"
)

// Register executable functions
var funcNameMap = map[string]func(args []ast.Expr, data map[string]interface{}) interface{}{}

func init() {
	funcNameMap = map[string]func(args []ast.Expr, data map[string]interface{}) interface{}{
		"in_array": inArray,
		"max":      maxFunc,
		"min":      minFunc,
	}
}

// maxFunc takes the maximum value among multiple numbers
func maxFunc(args []ast.Expr, data map[string]interface{}) interface{} {
	if len(args) == 0 {
		return errors.New("func max requires at least one argument")
	}
	var maxFloat float64
	for i, arg := range args {
		val := eval(arg, data)
		if err, ok := val.(error); ok {
			return err
		}
		f, err := castType(val, TypeFloat)
		if err != nil {
			return err
		}
		num, ok := f.(float64)
		if !ok {
			return errors.New("func max argument is not numeric")
		}
		if i == 0 || num > maxFloat {
			maxFloat = num
		}
	}
	return maxFloat
}

// minFunc takes the minimum value among multiple numbers
func minFunc(args []ast.Expr, data map[string]interface{}) interface{} {
	if len(args) == 0 {
		return errors.New("func min requires at least one argument")
	}
	var minFloat float64
	for i, arg := range args {
		val := eval(arg, data)
		if err, ok := val.(error); ok {
			return err
		}
		f, err := castType(val, TypeFloat)
		if err != nil {
			return err
		}
		num, ok := f.(float64)
		if !ok {
			return errors.New("func min argument is not numeric")
		}
		if i == 0 || num < minFloat {
			minFloat = num
		}
	}
	return minFloat
}

// inArray checks if a variable exists in an array
func inArray(args []ast.Expr, data map[string]interface{}) interface{} {
	// Variables in rule expressions
	param := eval(args[0], data)
	vRange, ok := args[1].(*ast.CompositeLit)
	if !ok {
		return errors.New("func in_array 2ed params is not a composite lit")
	}

	// Elements in arrays in rule expressions
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
