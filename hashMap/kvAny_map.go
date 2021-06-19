package hashMap

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/weblfe/types/stringx"
	"log"
	"strconv"
	"strings"
)

type KvAnyMap map[string]interface{}

func ToKvAnyMap(m map[string]interface{}) KvAnyMap {
	return m
}

// 将map 转对应 struct
func MapBind(m interface{}, bind interface{}) error {
	if bind == nil {
		return errors.New("nil pointer bind")
	}
	return mapstructure.Decode(m, bind)
}

func (kv KvAnyMap) Decode(v interface{}) error {
	return MapBind(kv, v)
}

func (kv KvAnyMap) Empty() bool {
	return len(kv) == 0
}

func (kv KvAnyMap) Exists(key string) bool {
	if _, ok := kv[key]; ok {
		return true
	}
	return false
}

func (kv KvAnyMap) Get(key string, def ...interface{}) interface{} {
	if len(def) == 0 {
		def = append(def, nil)
	}
	if v, ok := kv[key]; ok {
		return v
	}
	return def[0]
}

func (kv KvAnyMap) Remove(key ...string) int {
	var n = 0
	for _, k := range key {
		delete(kv, k)
		n++
	}
	return n
}

func (kv KvAnyMap) String() string {
	var data = kv.Bytes()
	if data != nil {
		return string(data)
	}
	return ""
}

func (kv KvAnyMap) Copy() KvAnyMap {
	var _new = map[string]interface{}{}
	for k, v := range kv {
		_new[k] = v
	}
	return _new
}

func (kv KvAnyMap) Bytes() []byte {
	var _bytes, err = json.Marshal(kv)
	if err != nil {
		log.Println(err)
		return nil
	}
	return _bytes
}

func (kv KvAnyMap) Size() int {
	return len(kv.Bytes())
}

func (kv KvAnyMap) Count() int {
	return len(kv)
}

func (kv KvAnyMap) Map() map[string]interface{} {
	return kv
}

func (kv KvAnyMap) MapPointer() *map[string]interface{} {
	var m = kv.Map()
	return &m
}

func (kv KvAnyMap) GetStrArrOf(key string, def ...[]string) []string {
	var v = kv.Get(key)
	def = append(def, nil)
	if v == nil {
		return def[0]
	}
	switch v.(type) {
	case []string:
		return v.([]string)
	case []interface{}:
		var (
			_strArr []string
			arr     = v.([]interface{})
		)
		for _, v := range arr {
			_strArr = append(_strArr, stringx.Stringer(v))
		}
		return _strArr
	case *[]string:
		return *v.(*[]string)
	case string:
		var str = v.(string)
		// 逗号分割
		if strings.Contains(str, ",") && !strings.Contains(str, "[") {
			return strings.Split(str, ",")
		}
		// json-array 对象
		if strings.Contains(str, "[") && strings.Contains(str, "]") {
			var _jsonArr []string
			if err := json.Unmarshal([]byte(str), &_jsonArr); err == nil {
				return _jsonArr
			}
		}
		// 空格分割
		return strings.Split(str, " ")
	}
	return def[0]
}

func (kv KvAnyMap) GetStringOf(key string, def ...string) string {
	def = append(def, "")
	var v = kv.Get(key)
	if v == nil {
		return def[0]
	}
	switch v.(type) {
	case string:
		return v.(string)
	case *string:
		return *v.(*string)
	case fmt.Stringer:
		return v.(fmt.Stringer).String()
	}
	var _json, err = json.Marshal(v)
	if err != nil {
		return def[0]
	}
	return string(_json)
}

func (kv KvAnyMap) GetIntOf(key string, def ...int) int {
	def = append(def, 0)
	var v = kv.Get(key)
	if v == nil {
		return def[0]
	}
	switch v.(type) {
	case *bool:
		var b = *v.(*bool)
		if b {
			return 1
		}
		return 0
	case bool:
		var b = v.(bool)
		if b {
			return 1
		}
		return 0
	case *int:
		return *v.(*int)
	case int:
		return v.(int)
	case *int64:
		return int(*v.(*int64))
	case int64:
		return int(v.(int64))
	case *int32:
		return int(*v.(*int32))
	case int32:
		return int(v.(int32))
	case *uint:
		return int(*v.(*uint))
	case uint:
		return int(v.(uint))
	case *uint32:
		return int(*v.(*uint32))
	case uint32:
		return int(v.(uint32))
	case *uint64:
		return int(*v.(*uint64))
	case uint64:
		return int(v.(uint64))
	case *float64:
		return int(*v.(*float64))
	case float64:
		return int(v.(float64))
	case float32:
		return int(v.(float32))
	case *float32:
		return int(*v.(*float32))
	case string:
		vs := v.(string)
		if n, err := strconv.Atoi(vs); err == nil {
			return n
		}
	case fmt.Stringer:
		vs := v.(fmt.Stringer).String()
		if n, err := strconv.Atoi(vs); err == nil {
			return n
		}
	}
	return def[0]
}

func (kv KvAnyMap) GetMapOf(key string, def ...map[string]interface{}) KvAnyMap {
	def = append(def, nil)
	var v = kv.Get(key)
	if v == nil {
		return def[0]
	}
	switch v.(type) {
	case map[string]interface{}:
		return v.(map[string]interface{})
	case *map[string]interface{}:
		return *v.(*map[string]interface{})
	case *string:
		var (
			vs   = *v.(*string)
			data = make(map[string]interface{})
		)
		if vs == "" {
			return def[0]
		}
		if err := json.Unmarshal([]byte(vs), &data); err == nil {
			return data
		}
	case string:
		var (
			vs   = v.(string)
			data = make(map[string]interface{})
		)
		if vs == "" {
			return def[0]
		}
		if err := json.Unmarshal([]byte(vs), &data); err == nil {
			return data
		}
	case fmt.Stringer:
		var (
			vs   = v.(fmt.Stringer).String()
			data = make(map[string]interface{})
		)
		if vs == "" {
			return def[0]
		}
		if err := json.Unmarshal([]byte(vs), &data); err == nil {
			return data
		}
	}
	return def[0]
}

func (kv KvAnyMap) GetMapListOf(key string, def ...[]map[string]interface{}) []KvAnyMap {
	def = append(def, nil)
	var v = kv.Get(key)
	if v == nil {
		return listMap(def[0])
	}
	switch v.(type) {
	case []interface{}:
		return anyArray2MapArr(v.([]interface{}))
	case []map[string]interface{}:
		return listMap(v.([]map[string]interface{}))
	case *[]map[string]interface{}:
		return listMap(*v.(*[]map[string]interface{}))
	case *[]map[string]string:
		return arrayMap2KvArr(*v.(*[]map[string]string))
	case []map[string]string:
		return arrayMap2KvArr(v.([]map[string]string))
	}
	return listMap(def[0])
}

func listMap(list []map[string]interface{}) []KvAnyMap {
	if list == nil {
		return nil
	}
	var _list []KvAnyMap
	for _, v := range list {
		_list = append(_list, v)
	}
	return _list
}

func anyArray2MapArr(arr []interface{}) []KvAnyMap {
	var list []KvAnyMap
	for _, v := range arr {
		if v == nil {
			continue
		}
		switch v.(type) {
		case map[string]interface{}:
			list = append(list, v.(map[string]interface{}))
		case *map[string]interface{}:
			list = append(list, *v.(*map[string]interface{}))
		}
	}
	return list
}

func arrayMap2KvArr(list []map[string]string) []KvAnyMap {
	if list == nil {
		return nil
	}
	var _list []KvAnyMap
	for _, v := range list {
		_list = append(_list, strMap2Kv(v))
	}
	return _list
}

func strMap2Kv(data map[string]string) map[string]interface{} {
	if data == nil {
		return nil
	}
	var _map = map[string]interface{}{}
	for k, v := range data {
		_map[k] = v
	}
	return _map
}

func MapDecode(input interface{}, output interface{}, tagName ...string) error {
	tagName = append(tagName, "json")
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   output,
		TagName:  tagName[0],
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}
