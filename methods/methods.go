package messages

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
	"reflect"
	"runtime"
	"strings"
	"unicode"
)

// 结构体对象方法
type structMethod struct {
	tagName     string
	onlyCopyTag bool
	deepCopy    bool
	input       interface{}
}

const (
	_defaultTag = "json"
)

// 获取struct method
func CreateStructMethod() *structMethod {
	return &structMethod{}
}

// 仅简单拷贝
func (s *structMethod) SimpleMap(ctx ...interface{}) (map[string]interface{}, error) {
	if len(ctx) <= 0 {
		ctx = append(ctx, s.getInput())
	}
	var out = make(map[string]interface{})
	v := reflect.ValueOf(ctx[0])
	// 指针类型
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}
	typ := v.Type()
	tagName := s.getTagName()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		// 拷贝 Name
		if tagName == "" && !s.onlyCopyTag {
			out[fi.Name] = v.Field(i).Interface()
			continue
		}
		if tag := fi.Tag.Get(tagName); tag != "" {
			tag = s.parseTagKey(tag)
			out[tag] = v.Field(i).Interface()
		}
	}
	return out, nil
}

// 强制转换
func (s *structMethod) MapMust(ctx ...interface{}) map[string]interface{} {
	if len(ctx) <= 0 {
		ctx = append(ctx, s.getInput())
	}
	var out = make(map[string]interface{})
	v := reflect.ValueOf(ctx[0])
	// 指针类型
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// we only accept structs
	if v.Kind() != reflect.Struct {
		return out
	}
	typ := v.Type()
	tagName := s.getTagName()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		// 拷贝 Name 小写 字段排除(私有)
		if tagName == "" && !s.onlyCopyTag && unicode.IsUpper([]rune(fi.Name)[0]) {
			vl := v.Field(i).Interface()
			if s.deepCopy {
				vfi := reflect.ValueOf(vl)
				if vfi.Kind() == reflect.Struct {
					vl = s.deepCopyStruct(vfi)
				}
				if vfi.Kind() == reflect.Ptr && vfi.Elem().Kind() == reflect.Struct {
					vl = s.deepCopyStruct(vfi.Elem())
				}
			}
			out[fi.Name] = vl
			continue
		}
		tag := s.parseTagKey(fi.Tag.Get(tagName))
		if tag != "" {
			vl := v.Field(i).Interface()
			// 深度拷贝
			if s.deepCopy {
				vfi := reflect.ValueOf(vl)
				if vfi.Kind() == reflect.Struct {
					vl = s.deepCopyStruct(vfi)
				}
				if vfi.Kind() == reflect.Ptr && vfi.Elem().Kind() == reflect.Struct {
					vl = s.deepCopyStruct(vfi.Elem())
				}
			}
			out[tag] = vl
		}
		// 非仅靠被tag 小写 字段排除(私有)
		if !s.onlyCopyTag && tag == "" && unicode.IsUpper([]rune(fi.Name)[0]) {
			vl := v.Field(i).Interface()
			// 深度拷贝
			if s.deepCopy {
				vfi := reflect.ValueOf(vl)
				if vfi.Kind() == reflect.Struct {
					vl = s.deepCopyStruct(vfi)
				}
				if vfi.Kind() == reflect.Ptr && vfi.Elem().Kind() == reflect.Struct {
					vl = s.deepCopyStruct(vfi.Elem())
				}
			}
			out[fi.Name] = vl
		}

	}
	return out
}

// 深度拷贝
func (s *structMethod) deepCopyStruct(v reflect.Value) interface{} {
	var out = make(map[string]interface{})
	typ := v.Type()
	tagName := s.getTagName()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		// 拷贝 Name
		if tagName == "" && !s.onlyCopyTag && unicode.IsUpper([]rune(fi.Name)[0]) {
			vl := v.Field(i).Interface()
			if s.deepCopy {
				vfi := reflect.ValueOf(vl)
				if vfi.Kind() == reflect.Struct {
					vl = s.deepCopyStruct(vfi)
				}
				if vfi.Kind() == reflect.Ptr && vfi.Elem().Kind() == reflect.Struct {
					vl = s.deepCopyStruct(vfi.Elem())
				}
			}
			out[fi.Name] = vl
			continue
		}
		tag := fi.Tag.Get(tagName)
		tag = s.parseTagKey(tag)
		if tag != "" {
			vl := v.Field(i).Interface()
			// 深度拷贝
			if s.deepCopy {
				vfi := reflect.ValueOf(vl)
				if vfi.Kind() == reflect.Struct {
					vl = s.deepCopyStruct(vfi)
				}
				if vfi.Kind() == reflect.Ptr && vfi.Elem().Kind() == reflect.Struct {
					vl = s.deepCopyStruct(vfi.Elem())
				}
			}
			out[tag] = vl
		}
		if !s.onlyCopyTag && tag == "" && unicode.IsUpper([]rune(fi.Name)[0]) {
			vl := v.Field(i).Interface()
			// 深度拷贝
			if s.deepCopy {
				vfi := reflect.ValueOf(vl)
				if vfi.Kind() == reflect.Struct {
					vl = s.deepCopyStruct(vfi)
				}
				if vfi.Kind() == reflect.Ptr && vfi.Elem().Kind() == reflect.Struct {
					vl = s.deepCopyStruct(vfi.Elem())
				}
			}
			out[fi.Name] = vl
		}
	}
	return out
}

// 设置 拷贝tag
func (s *structMethod) getTagName() string {
	if s.tagName == "" && s.onlyCopyTag {
		s.tagName = _defaultTag
	}
	return s.tagName
}

// 解析key
func (s *structMethod) parseTagKey(tag string) string {
	if strings.Contains(tag, ",") {
		arr := strings.Split(tag, ",")
		return arr[0]
	}
	return tag
}

// 设置是否深度拷贝
func (s *structMethod) SetDeepCopy(on bool) {
	s.deepCopy = on
}

// 设置是否仅拷贝tag 字段
func (s *structMethod) SetOnlyCopyTag(on bool) {
	s.onlyCopyTag = on
}

// 设置tag
func (s *structMethod) SetTag(name string) {
	s.tagName = name
}

func (s *structMethod) SetInput(v interface{}) {
	s.input = v
}

func (s *structMethod) getInput() interface{} {
	if s.input == nil {
		return s
	}
	return s.input
}

func (s *structMethod) Bind(self interface{}, _data interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			pc, file, line, ok := runtime.Caller(1)
			log.Println("error:", err, pc, file, line, ok)
		}
	}()
	tagName := s.getTagName()
	if tagName == "json" || tagName == "" {
		v, err := json.Marshal(_data)
		if err != nil {
			log.Println("structMethod GetMap Error1:", err.Error())
			return err
		}
		if err = json.Unmarshal(v, self); err != nil {
			log.Println("structMethod GetMap Error2:", err.Error())
			return err
		}
		return nil
	}
	switch _data.(type) {
	case map[string]interface{}:
		var m = _data.(map[string]interface{})
		err := mapstructure.Decode(&m, self)
		if err == nil {
			return err
		}
	case *map[string]interface{}:
		var m = _data.(*map[string]interface{})
		err := mapstructure.Decode(m, self)
		if err != nil {
			return err
		}
	default:
		v, err := json.Marshal(_data)
		if err != nil {
			log.Println("structMethod GetMap Error1:", err.Error())
			return err
		}
		if err = json.Unmarshal(v, self); err != nil {
			log.Println("structMethod GetMap Error2:", err.Error())
			return err
		}
	}
	return nil
}

func (s *structMethod) GetTypeOf(v interface{}) *Type {
	return TypesOf(v)
}
