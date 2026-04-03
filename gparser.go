package gparser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

// Match uses the native parser to complete expression matching with input data
func Evaluate(expr string, data map[string]interface{}) (interface{}, error) {
	// Empty expression returns nil
	if expr == "" {
		return nil, nil
	}
	if data == nil {
		data = make(map[string]interface{})
	}

	// Parse the expression
	parseExpr, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, err
	}

	// Expression true and false values are recognized as variables
	data["true"] = true
	data["false"] = false

	result := eval(parseExpr, data)
	if errVal, ok := result.(error); ok {
		return nil, errVal
	}
	return result, nil
}

func Match(expr string, data map[string]interface{}) (bool, error) {
	// Empty expression defaults to match success
	if expr == "" {
		return true, nil
	}

	// Empty data defaults to match failure
	if data == nil {
		return false, nil
	}

	result, err := Evaluate(expr, data)
	if err != nil {
		return false, err
	}

	boolResult, err := castType(result, "bool")
	if err != nil {
		return false, err
	}

	if v, ok := boolResult.(bool); ok {
		return v, nil
	}
	return false, errors.New("expression result is not bool")
}

func eval(expr ast.Expr, data map[string]interface{}) interface{} {
	switch expr := expr.(type) {
	case *ast.BasicLit: // Matched to data
		return getlitValue(expr)
	case *ast.BinaryExpr: // Matched to subtree
		x := eval(expr.X, data)
		y := eval(expr.Y, data)
		if x == nil || y == nil {
			return errors.New(fmt.Sprintf("%+v, %+v is nil", x, y))
		}
		op := expr.Op
		// Rule calculation (matching according to the type of variables in the rule expression)
		switch y.(type) {
		case float64:
			return calculateForFloat(x, y, op)
		case int:
			return calculateForInt(x, y, op)
		case int64:
			if _, ok := x.(float64); ok {
				return calculateForFloat(x, y, op)
			}
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
	case *ast.CallExpr: // Matched to function
		return calculateForFunc(expr.Fun.(*ast.Ident).Name, expr.Args, data)
	case *ast.ParenExpr: // Matched to parentheses
		return eval(expr.X, data)
	case *ast.UnaryExpr: // Matched to unary expression
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
	case *ast.Ident: // Matched to variable
		if val, ok := data[expr.Name]; ok {
			return val
		} else {
			return errors.New(fmt.Sprintf("no parameter %s", expr.Name))
		}
	default:
		return errors.New(fmt.Sprintf("%x type is not support", expr))
	}
}

// Get the data of variables in AST (numbers in expressions are int, converted to int64)
func getlitValue(basicLit *ast.BasicLit) interface{} {
	switch basicLit.Kind {
	case token.INT:
		value, err := strconv.ParseInt(basicLit.Value, 10, 64)
		if err != nil {
			return err
		}
		return value
	case token.FLOAT:
		value, err := strconv.ParseFloat(basicLit.Value, 64)
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
