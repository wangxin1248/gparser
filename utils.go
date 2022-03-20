package gparser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	TypeString = "string"
	TypeInt64  = "int64"
	TypeBool   = "bool"
	TypeFloat  = "float"
	TypeObject = "object"
)

// castType 基础类型转换，支持 string int64 bool float object 几种类型
func castType(data interface{}, typeName string) (interface{}, error) {
	if typeName == "" {
		return data, nil
	}
	switch strings.ToLower(typeName) {
	case TypeString:
		return castToString(data)
	case TypeInt64:
		return castToInt64(data)
	case TypeBool:
		return castToBoolean(data)
	case TypeFloat:
		return castToFloat(data)
	case TypeObject:
		return data, nil
	default:
		return nil, fmt.Errorf("type cast failure, unexpected type: %s", typeName)
	}
}

func castToBoolean(data interface{}) (bool, error) {
	if data == nil {
		return false, nil
	}
	switch t := data.(type) {
	case bool:
		return t, nil
	case string:
		if t == "" {
			return false, nil
		}
		return strconv.ParseBool(t)
	default:
		return false, fmt.Errorf("type cast failure, unexpected boolean value: %v", data)
	}
}

func castToString(data interface{}) (string, error) {
	if data == nil {
		return "", nil
	}
	return fmt.Sprint(data), nil
}

func castToInt64(data interface{}) (interface{}, error) {
	if data == nil {
		return 0, nil
	}

	switch t := data.(type) {
	case int:
		return int64(t), nil
	case float32:
		return int64(t), nil
	case float64:
		return int64(t), nil
	case json.Number:
		return t.Int64()
	case string:
		if t == "" {
			return 0, nil
		}
		return strconv.ParseInt(t, 10, 64)
	}
	return strconv.ParseInt(fmt.Sprint(data), 10, 64)
}

func castToFloat(data interface{}) (interface{}, error) {
	if data == nil {
		return 0, nil
	}

	switch t := data.(type) {
	case int:
		return float64(t), nil
	case int16:
		return float32(t), nil
	case int32:
		return float32(t), nil
	case int64:
		return float64(t), nil
	case uint:
		return float64(t), nil
	case uint16:
		return float32(t), nil
	case uint32:
		return float32(t), nil
	case uint64:
		return float64(t), nil
	case float32:
		return t, nil
	case float64:
		return t, nil
	case json.Number:
		return t.Float64()
	case string:
		if t == "" {
			return 0., nil
		}
		return strconv.ParseFloat(t, 64)
	}
	return strconv.ParseFloat(fmt.Sprint(data), 64)
}
