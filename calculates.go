package gparser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
)

// 计算int类型表达式
func calculateForInt(x, y interface{}, op token.Token) interface{} {
	x, err := castType(x, TypeInt64)
	if err != nil {
		return err
	}
	y, err = castType(y, TypeInt64)
	if err != nil {
		return err
	}
	return calculateForInt64(x, y, op)
}

// 计算int64类型表达式
func calculateForInt64(x, y interface{}, op token.Token) interface{} {
	x, err := castType(x, TypeInt64)
	if err != nil {
		return err
	}
	xInt, xok := x.(int64)
	yInt, yok := y.(int64)
	if !xok || !yok {
		return errors.New(fmt.Sprintf("%v %v %v eval failed", x, op, y))
	}

	// 计算逻辑
	switch op {
	case token.EQL:
		return xInt == yInt
	case token.NEQ:
		return xInt != yInt
	case token.GTR:
		return xInt > yInt
	case token.LSS:
		return xInt < yInt
	case token.GEQ:
		return xInt >= yInt
	case token.LEQ:
		return xInt <= yInt
	case token.ADD:
		return xInt + yInt
	case token.SUB:
		return xInt - yInt
	case token.MUL:
		return xInt * yInt
	case token.QUO:
		if yInt == 0 {
			return 0
		}
		return xInt / yInt
	default:
		return errors.New(fmt.Sprintf("unsupported binary operator: %s", op.String()))
	}
}

// 计算float类型表达式
func calculateForFloat(x, y interface{}, op token.Token) interface{} {
	x, err := castType(x, TypeFloat)
	if err != nil {
		return err
	}
	y, err = castType(y, TypeFloat)
	if err != nil {
		return err
	}
	xFloat, xok := x.(float64)
	yFloat, yok := y.(float64)
	if !xok || !yok {
		return errors.New(fmt.Sprintf("%v %v %v eval failed", x, op, y))
	}

	switch op {
	case token.EQL:
		return xFloat == yFloat
	case token.NEQ:
		return xFloat != yFloat
	case token.GTR:
		return xFloat > yFloat
	case token.LSS:
		return xFloat < yFloat
	case token.GEQ:
		return xFloat >= yFloat
	case token.LEQ:
		return xFloat <= yFloat
	case token.ADD:
		return xFloat + yFloat
	case token.SUB:
		return xFloat - yFloat
	case token.MUL:
		return xFloat * yFloat
	case token.QUO:
		if yFloat == 0 {
			return 0.0
		}
		return xFloat / yFloat
	}
	return errors.New(fmt.Sprintf("unsupported binary operator: %s", op.String()))
}

// 计算string类型表达式
func calculateForString(x, y interface{}, op token.Token) interface{} {
	x, err := castType(x, TypeString)
	if err != nil {
		return err
	}
	xString, xok := x.(string)
	yString, yok := y.(string)
	if !xok || !yok {
		return errors.New(fmt.Sprintf("%v %v %v eval failed", x, op, y))
	}

	// 计算逻辑
	switch op {
	case token.EQL: // ==
		return xString == yString
	case token.NEQ: // !=
		return xString != yString
	}
	return errors.New(fmt.Sprintf("unsupported binary operator: %s", op.String()))
}

// 计算bool类型表达式
func calculateForBool(x, y interface{}, op token.Token) interface{} {
	x, err := castType(x, TypeBool)
	if err != nil {
		return err
	}
	xb, xok := x.(bool)
	yb, yok := y.(bool)
	if !xok || !yok {
		return errors.New(fmt.Sprintf("%v %v %v eval failed", x, op, y))
	}

	// 计算逻辑
	switch op {
	case token.LAND:
		return xb && yb
	case token.LOR:
		return xb || yb
	case token.EQL:
		return xb == yb
	case token.NEQ:
		return xb != yb
	}
	return errors.New(fmt.Sprintf("unsupported binary operator: %s", op.String()))
}

// calculateForFunc 计算函数表达式
func calculateForFunc(funcName string, args []ast.Expr, data map[string]interface{}) interface{} {
	// 根据funcName分发逻辑
	handler, ok := funcNameMap[funcName]
	if !ok {
		return errors.New(fmt.Sprintf("%+v func not support", funcName))
	}
	return handler(args, data)
}
