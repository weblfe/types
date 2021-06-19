package hashMap

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/weblfe/types/stringx"
	"io"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
)

// 字符(ASCII)自然排序 b > a ? true : false
func NatSortCompareAsc(a, b string) bool {
	var arr = []string{a, b}
	sort.Strings(arr)
	if arr[0] == b && arr[1] == a {
		return false
	}
	return true
}

// 字符(ASCII)自然排序 a > b ? true : false
func NatSortCompareDesc(a, b string) bool {
	var arr = []string{a, b}
	sort.Strings(arr)
	if arr[0] == b && arr[1] == a {
		return true
	}
	return false
}

// md5 签名
func SignatureMd5(v StrObject, _salt ...string) string {
	var (
		str  = ""
		salt = os.Getenv("SIGNATURE_SALT")
	)
	if v == nil {
		return ""
	}
	_salt = append(_salt, "")
	// 外部salt
	if len(_salt) > 0 && _salt[0] != "" {
		salt = _salt[0]
	}
	v.Sort()
	for _, key := range v.Keys() {
		str += strings.ToLower(key) + "=" + v.ValueOf(key) + "&"
	}
	str += salt
	sign := Hash(str)
	log.Println("SignatureMd5.Info sign-data", str, "sign:", sign)
	return sign
}

// md5
func Hash(data string) string {
	var _md5 = md5.New()
	if n, err := io.WriteString(_md5, data); err != nil || n == 0 {
		return ""
	}
	return fmt.Sprintf("%x", _md5.Sum(nil))
}

func GetKvMapValueOf(data map[string]interface{}, key string, def ...interface{}) interface{} {
	def = append(def, nil)
	for k, v := range data {
		if k == key || strings.EqualFold(k, key) {
			return v
		}
	}
	return def[0]
}

func GetKvStrMapValueOf(data map[string]string, key string, def ...string) string {
	def = append(def, "")
	for k, v := range data {
		if k == key || strings.EqualFold(k, key) {
			return v
		}
	}
	return def[0]
}

func ConvertMap(v interface{}) map[string]interface{} {
	switch v.(type) {
	case map[string]interface{}:
		return v.(map[string]interface{})
	case *map[string]interface{}:
		return *v.(*map[string]interface{})
	case map[string]string:
		var _data = make(map[string]interface{})
		for k, v := range v.(map[string]string) {
			_data[k] = v
		}
		return _data
	case *map[string]string:
		var _data = make(map[string]interface{})
		for k, v := range *v.(*map[string]string) {
			_data[k] = v
		}
		return _data
	case map[string]int:
		var _data = make(map[string]interface{})
		for k, v := range v.(map[string]int) {
			_data[k] = v
		}
		return _data
	case *map[string]int:
		var _data = make(map[string]interface{})
		for k, v := range *v.(*map[string]int) {
			_data[k] = v
		}
		return _data
	case map[string]float64:
		var _data = make(map[string]interface{})
		for k, v := range *v.(*map[string]float64) {
			_data[k] = v
		}
		return _data
	case *map[string]float64:
		var _data = make(map[string]interface{})
		for k, v := range v.(map[string]float64) {
			_data[k] = v
		}
		return _data
	case string:
		var (
			vs    = v.(string)
			_data = make(map[string]interface{})
		)
		if vs == "" {
			return nil
		}
		if err := json.Unmarshal([]byte(vs), &_data); err != nil {
			log.Println("ConvertMap String ERROR", err.Error())
			return nil
		}
		return _data
	case []byte:
		var (
			vs    = v.([]byte)
			_data = make(map[string]interface{})
		)
		if vs == nil || len(vs) == 0 {
			return nil
		}
		if err := json.Unmarshal(vs, &_data); err != nil {
			log.Println("ConvertMap []byte ERROR", err.Error())
			return nil
		}
		return _data
	case fmt.Stringer:
		var (
			vs    = v.(fmt.Stringer).String()
			_data = make(map[string]interface{})
		)
		if vs == "" {
			return nil
		}
		if err := json.Unmarshal([]byte(vs), &_data); err != nil {
			log.Println("ConvertMap String ERROR", err.Error())
			return nil
		}
		return _data
	}
	var (
		el  reflect.Value
		vfi = reflect.ValueOf(v)
	)
	if vfi.Kind() == reflect.Struct {
		el = vfi
	}
	if vfi.Kind() == reflect.Ptr {
		el = vfi.Elem()
		if el.Kind() != reflect.Struct {
			return nil
		}
	}
	_data := make(map[string]interface{})
	err := mapstructure.Decode(el.Interface(), &_data)
	if err != nil {
		log.Println("ConvertMap Error", err.Error())
		return nil
	}
	return _data
}

func ConvertStrMap(v interface{}) map[string]string {
	switch v.(type) {
	case map[string]string:
		return v.(map[string]string)
	case *map[string]string:
		return *v.(*map[string]string)
	case map[string]int:
		var _data = make(map[string]string)
		for k, v := range v.(map[string]int) {
			_data[k] = stringx.Stringer(v)
		}
		return _data
	case *map[string]int:
		var _data = make(map[string]string)
		for k, v := range *v.(*map[string]int) {
			_data[k] = stringx.Stringer(v)
		}
		return _data
	case map[string]float64:
		var _data = make(map[string]string)
		for k, v := range v.(map[string]float64) {
			_data[k] = stringx.Stringer(v)
		}
		return _data
	case *map[string]float64:
		var _data = make(map[string]string)
		for k, v := range *v.(*map[string]float64) {
			_data[k] = stringx.Stringer(v)
		}
		return _data
	case *map[string]interface{}:
		var _data = make(map[string]string)
		for k, v := range *v.(*map[string]interface{}) {
			_data[k] = stringx.Stringer(v)
		}
		return _data
	case map[string]interface{}:
		var _data = make(map[string]string)
		for k, v := range v.(map[string]interface{}) {
			_data[k] = stringx.Stringer(v)
		}
		return _data
	case string:
		var (
			vs    = v.(string)
			_data = make(map[string]string)
		)
		if vs == "" {
			return nil
		}
		if err := json.Unmarshal([]byte(vs), &_data); err != nil {
			log.Println("ConvertMap String ERROR", err.Error())
			return nil
		}
		return _data
	case []byte:
		var (
			vs    = v.([]byte)
			_data = make(map[string]string)
		)
		if vs == nil || len(vs) == 0 {
			return nil
		}
		if err := json.Unmarshal(vs, &_data); err != nil {
			log.Println("ConvertMap []byte ERROR", err.Error())
			return nil
		}
		return _data
	case fmt.Stringer:
		var (
			vs    = v.(fmt.Stringer).String()
			_data = make(map[string]string)
		)
		if vs == "" {
			return nil
		}
		if err := json.Unmarshal([]byte(vs), &_data); err != nil {
			log.Println("ConvertMap String ERROR", err.Error())
			return nil
		}
		return _data
	}
	var (
		el  reflect.Value
		vfi = reflect.ValueOf(v)
	)
	if vfi.Kind() == reflect.Struct {
		el = vfi
	}
	if vfi.Kind() == reflect.Ptr {
		el = vfi.Elem()
		if el.Kind() != reflect.Struct {
			return nil
		}
	}
	_data := make(map[string]string)
	err := mapstructure.Decode(el.Interface(), &_data)
	if err != nil {
		log.Println("ConvertStrMap Error", err.Error())
	}
	return _data
}
