package messages

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Type struct {
	v interface{}
}

func TypesOf(v interface{}) *Type {
	return &Type{v: v}
}

func (_type *Type) GetInt(def ...int) int {
	def = append(def, 0)
	if _type.v == nil {
		return def[0]
	}
	switch _type.v.(type) {
	case string:
		if n, err := strconv.Atoi(_type.v.(string)); err == nil {
			return n
		}
		return def[0]
	case *string:
		if n, err := strconv.Atoi(*_type.v.(*string)); err == nil {
			return n
		}
		return def[0]
	case fmt.Stringer:
		if n, err := strconv.Atoi(_type.v.(fmt.Stringer).String()); err == nil {
			return n
		}
		return def[0]
	case int:
		return _type.v.(int)
	case *int:
		return *(_type.v.(*int))
	case int32:
		return int(_type.v.(int32))
	case *int32:
		return int(*_type.v.(*int32))
	case int8:
		return int(_type.v.(int8))
	case *int8:
		return int(*_type.v.(*int8))
	case int16:
		return int(_type.v.(int16))
	case *int16:
		return int(*_type.v.(*int16))
	case int64:
		return int(_type.v.(int64))
	case *int64:
		return int(*_type.v.(*int64))
	case uint8:
		return int(_type.v.(uint8))
	case *uint8:
		return int(*_type.v.(*uint8))
	case uint32:
		return int(_type.v.(uint32))
	case *uint32:
		return int(*_type.v.(*uint32))
	case uint16:
		return int(_type.v.(uint16))
	case *uint16:
		return int(*_type.v.(*uint16))
	case uint64:
		return int(_type.v.(uint64))
	case *uint64:
		return int(*_type.v.(*uint64))
	case float64:
		return int(_type.v.(float64))
	case *float64:
		return int(*_type.v.(*float64))
	case float32:
		return int(_type.v.(float32))
	case *float32:
		return int(*_type.v.(*float32))
	case bool:
		if _type.v.(bool) {
			return 1
		}
		return 0
	case *bool:
		if *_type.v.(*bool) {
			return 1
		}
		return 0
	}
	return def[0]
}

func (_type *Type) GetFloat(def ...float32) float32 {
	def = append(def, 0)
	if _type.v == nil {
		return def[0]
	}
	switch _type.v.(type) {
	case string:
		if n, err := strconv.ParseFloat(_type.v.(string), 32); err == nil {
			return float32(n)
		}
		return def[0]
	case *string:
		if n, err := strconv.ParseFloat(*_type.v.(*string), 32); err == nil {
			return float32(n)
		}
		return def[0]
	case fmt.Stringer:
		if n, err := strconv.ParseFloat(_type.v.(fmt.Stringer).String(), 32); err == nil {
			return float32(n)
		}
		return def[0]
	case int:
		return float32(_type.v.(int))
	case *int:
		return float32(*(_type.v.(*int)))
	case int32:
		return float32(_type.v.(int32))
	case *int32:
		return float32(*_type.v.(*int32))
	case int8:
		return float32(_type.v.(int8))
	case *int8:
		return float32(*_type.v.(*int8))
	case int16:
		return float32(_type.v.(int16))
	case *int16:
		return float32(*_type.v.(*int16))
	case int64:
		return float32(_type.v.(int64))
	case *int64:
		return float32(*_type.v.(*int64))
	case uint8:
		return float32(_type.v.(uint8))
	case *uint8:
		return float32(*_type.v.(*uint8))
	case uint32:
		return float32(_type.v.(uint32))
	case *uint32:
		return float32(*_type.v.(*uint32))
	case uint16:
		return float32(_type.v.(uint16))
	case *uint16:
		return float32(*_type.v.(*uint16))
	case uint64:
		return float32(_type.v.(uint64))
	case *uint64:
		return float32(*_type.v.(*uint64))
	case float64:
		return float32(_type.v.(float64))
	case *float64:
		return float32(*_type.v.(*float64))
	case float32:
		return _type.v.(float32)
	case *float32:
		return *_type.v.(*float32)
	case bool:
		if _type.v.(bool) {
			return 1
		}
		return 0
	case *bool:
		if *_type.v.(*bool) {
			return 1
		}
		return 0
	}
	return def[0]
}

func (_type *Type) GetString(def ...string) string {
	def = append(def, "")
	if _type.v == nil {
		return def[0]
	}
	switch _type.v.(type) {
	case string:
		return _type.v.(string)
	case fmt.Stringer:
		return _type.v.(fmt.Stringer).String()
	}
	if _json, err := json.Marshal(_type.v); err == nil {
		return string(_json)
	}
	return fmt.Sprintf("%v", _type.v)
}

func (_type *Type) GetBool(def ...bool) bool {
	def = append(def, false)
	if _type.v == nil {
		return def[0]
	}
	switch _type.v.(type) {
	case string:
		vs := _type.v.(string)
		if _type.isTrue(vs) {
			return true
		}
		if b, err := strconv.ParseBool(vs); err == nil {
			return b
		}
		return def[0]
	case *string:
		vs := *_type.v.(*string)
		if _type.isTrue(vs) {
			return true
		}
		if b, err := strconv.ParseBool(vs); err == nil {
			return b
		}
		return def[0]
	case fmt.Stringer:
		v := _type.v.(fmt.Stringer).String()
		if _type.isTrue(v) {
			return true
		}
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
		return def[0]
	case bool:
		return _type.v.(bool)
	case *bool:
		return *_type.v.(*bool)
	}

	if _type.GetInt() == 1 {
		return true
	}
	return def[0]
}

func (_type *Type) isTrue(vs string) bool {
	switch vs {
	case "true":
		return true
	case "True":
		return true
	case "yes":
		return true
	case "Yes":
		return true
	case "On":
		return true
	case "on":
		return true
	case "1":
		return true
	}
	return false
}
