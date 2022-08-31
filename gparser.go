package gparser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

// Match 利用原生parser完成表达式与输入数据匹配任务
func Match(expr string, data map[string]interface{}) (bool, error) {
	// 空表达式默认匹配成功
	if expr == "" {
		return true, nil
	}
	// 空数据默认匹配失败
	if data == nil {
		return false, nil
	}
	// 解析表达式
	parseExpr, err := parser.ParseExpr(expr)
	if err != nil {
		return false, err
	}
	// 表达式中的 true、false 值被识别为变量
	data["true"] = true
	data["false"] = false

	// 匹配表达式与输入数据
	result := eval(parseExpr, data)

	// 返回匹配结果
	if _, ok := result.(error); ok {
		return false, result.(error)
	}
	return result.(bool), nil
}

func eval(expr ast.Expr, data map[string]interface{}) interface{} {
	switch expr := expr.(type) {
	case *ast.BasicLit: // 匹配到数据
		return getlitValue(expr)
	case *ast.BinaryExpr: // 匹配到子树
		x := eval(expr.X, data)
		y := eval(expr.Y, data)
		if x == nil || y == nil {
			return errors.New(fmt.Sprintf("%+v, %+v is nil", x, y))
		}
		op := expr.Op
		// 规则计算（按照规则表达式中变量的类型进行匹配）
		switch y.(type) {
		case int:
			return calculateForInt(x, y, op)
		case int64:
			return calculateForInt64(x, y, op)
		case string:
			return calculateForString(x, y, op)
		case bool:
			return calculateForBool(x, y, op)
		case error:
			return errors.New(fmt.Sprintf("%+v %+v %+v eval failed", x, op, y))
		default:
			return errors.New(fmt.Sprintf("%+v op is not support", op))
		}
	case *ast.CallExpr: // 匹配到函数
		return calculateForFunc(expr.Fun.(*ast.Ident).Name, expr.Args, data)
	case *ast.ParenExpr: // 匹配到括号
		return eval(expr.X, data)
	case *ast.UnaryExpr: // 匹配到一元表达式
		x := eval(expr.X, data)
		if x == nil {
			return errors.New(fmt.Sprintf("%+v is nil", x))
		}
		op := expr.Op
		switch op {
		case token.NOT:
			switch x.(type) {
			case bool:
				xb := x.(bool)
				return !xb
			}
		}
		return errors.New(fmt.Sprintf("%x type is not support", expr))
	case *ast.Ident: // 匹配到变量
		return data[expr.Name]
	default:
		return errors.New(fmt.Sprintf("%x type is not support", expr))
	}
}

// 获取AST中变量的数据（表达式中的数字为int，转为int64）
func getlitValue(basicLit *ast.BasicLit) interface{} {
	switch basicLit.Kind {
	case token.INT:
		value, err := strconv.ParseInt(basicLit.Value, 10, 64)
		if err != nil {
			return err
		}
		return value
	case token.STRING:
		value, err := strconv.Unquote(basicLit.Value)
		if err != nil {
			return err
		}
		return value
	}
	return errors.New(fmt.Sprintf("%s is not support type", basicLit.Kind))
}
