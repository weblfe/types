package stringx

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Stringer(v interface{}) string {
	if v == nil {
		return ""
	}
	switch v.(type) {
	case string:
		return v.(string)
	case fmt.Stringer:
		return v.(fmt.Stringer).String()
	case []byte:
		return string(v.([]byte))
	case []rune:
		return string(v.([]rune))
	case error:
		return v.(error).Error()
	}
	vfi := reflect.ValueOf(v)
	switch vfi.Kind() {
	case reflect.Ptr:
		el := vfi.Elem()
		return Stringer(el.Interface())
	case reflect.String:
		return vfi.String()
	case reflect.Invalid:
		return ""
	}
	if _json, err := json.Marshal(v); err == nil {
		return string(_json)
	}
	return fmt.Sprintf("%v", v)
}