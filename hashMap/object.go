package hashMap

import (
	"fmt"
	"github.com/weblfe/types/stringx"
	"runtime"
	"sort"
	"sync"
)
type (

	SortAnyObject interface {
		Keys() []string
		Value(key string) interface{}
		sort.Interface
		Sort()
		Remove(key string) bool
	}

	StrObject interface {
		SortAnyObject
		ValueOf(key string) string
		MapMethodProvider
	}

	MapMethodProvider interface {
		KvAnyMap() map[string]interface{}
		KvMaps() map[string]string
	}

	Decoder interface {
		Decode(v interface{}) error
	}

	Encoder interface {
		Encode() ([]byte, error)
	}

	Getter interface {
		Get(key string) interface{}
	}

	strObject struct {
		sync.RWMutex
		entries    []*Entry
		size       uint32
		dirtyCache map[string]*Index
		compare    func(string, string) bool
	}

	Entry struct {
		Key   string
		Value interface{}
	}

	Index struct {
		Pos     int
		Pointer interface{}
	}

)

func NewStrObject(size ...uint32) *strObject {
	size = append(size, 0)
	return &strObject{
		size:       size[0],
		entries:    make([]*Entry, size[0]),
		dirtyCache: make(map[string]*Index, size[0]),
	}
}

func (object *strObject) CheckType(v interface{}) bool {
	switch v.(type) {
	case string:
		return true
	case fmt.Stringer:
		return true
	}
	str := fmt.Sprintf("%v", v)
	if str != "" {
		return true
	}
	return false
}

func (object *strObject) Set(key string, value interface{}) *strObject {
	var (
		isEmpty bool
		str     string
	)
	switch value.(type) {
	case string:
		str = value.(string)
	case fmt.Stringer:
		str = value.(fmt.Stringer).String()
	case nil:
		isEmpty = true
	default:
		isEmpty = false
	}
	if !isEmpty {
		str = stringx.Stringer(value)
		if str == "" {
			_, file, line, _ := runtime.Caller(0)
			panic(fmt.Sprintf("%v, %s, at line %d ,error : %s", value, file, line, "type error: value must be type of string or fmt.Stringer"))
		}
	}
	object.Lock()
	defer object.Unlock()
	index, ok := object.dirtyCache[key]
	if ok {
		entry := object.entries[index.Pos]
		entry.Value = str
		return object
	}
	pointer := &Entry{Key: key, Value: str}
	object.entries = append(object.entries, pointer)
	object.size = uint32(len(object.entries))
	object.dirtyCache[key] = &Index{
		Pos:     int(object.size - 1),
		Pointer: pointer,
	}
	return object
}

func (object *strObject) Get(key string) interface{} {
	object.RLock()
	defer object.RUnlock()
	if object.size == 0 {
		return nil
	}
	if v, ok := object.dirtyCache[key]; ok && v != nil {
		_entry := v.Pointer
		if _entry == nil {
			return nil
		}
		switch _entry.(type) {
		case *Entry:
			return _entry.(*Entry)
		}
		return _entry
	}
	return nil
}

func (object *strObject) Keys() []string {
	object.RLock()
	defer object.RUnlock()
	if object.size == 0 {
		return []string{}
	}
	var keys []string
	for _, v := range object.entries {
		if v != nil {
			keys = append(keys, v.Key)
		}
	}
	return keys
}

func (object *strObject) Value(key string) interface{} {
	return object.Get(key)
}

func (object *strObject) ValueOf(key string) string {
	var v = object.Get(key)
	if v == nil {
		return ""
	}
	switch v.(type) {
	case string:
		return v.(string)
	case fmt.Stringer:
		return v.(fmt.Stringer).String()
	case *Entry:
		vs := v.(*Entry)
		return stringx.Stringer(vs.Value)
	}
	return fmt.Sprintf("%v", v)
}

func (object *strObject) Len() int {
	return int(object.size)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (object *strObject) Less(i, j int) bool {
	return object.GetCompare()(object.IndexOf(i), object.IndexOf(j))
}

// Swap swaps the elements with indexes i and j.
func (object *strObject) Swap(i, j int) {
	if i < 0 || j < 0 || i > int(object.size) || j > int(object.size) {
		return
	}
	object.Lock()
	defer object.Unlock()
	var (
		_entry1, _entry2 = object.entries[i], object.entries[j]
		k1, k2           = _entry1.Key, _entry2.Key
	)
	object.entries[i], object.entries[j] = _entry2, _entry1
	object.dirtyCache[k1].Pos, object.dirtyCache[k2].Pos = j, i
}

func (object *strObject) Exists(key string) bool {
	object.RLock()
	defer object.RUnlock()
	if _, ok := object.dirtyCache[key]; ok {
		return ok
	}
	return false
}

func (object *strObject) Remove(key string) bool {
	if !object.Exists(key) {
		return false
	}
	object.Lock()
	defer object.Unlock()
	if v, ok := object.dirtyCache[key]; ok {
		object.entries = append(object.entries[:v.Pos], object.entries[v.Pos+1:]...)
		object.size--
		delete(object.dirtyCache, key)
		// 更新索引号
		if v.Pos < int(object.size)-1 && object.size >= 1 {
			i := v.Pos + 1
			for {
				if i >= int(object.size) {
					break
				}
				v := object.entries[i]
				if index, _ok := object.dirtyCache[v.Key]; _ok {
					index.Pos = i
				}
				i++
			}
		}
		return true
	}
	return false
}

func (object *strObject) Sort() {
	var (
		keys       = object.Keys()
		comparator = object.GetCompare()
	)
	object.Lock()
	defer object.Unlock()
	sort.Slice(keys, func(i, j int) bool {
		return comparator(keys[i], keys[j])
	})
	object.size = uint32(len(object.entries))
	for i, k := range keys {
		// old post
		index, ok := object.dirtyCache[k]
		if !ok {
			continue
		}
		_pos := index.Pos
		// new post
		index.Pos = i
		if i >= int(object.size) || _pos >= int(object.size) {
			continue
		}
		_k := object.entries[i].Key
		object.dirtyCache[_k].Pos = _pos
		object.entries[_pos], object.entries[i] = object.entries[i], object.entries[_pos]
	}
}

func (object *strObject) IndexOf(i int) string {
	object.RLock()
	defer object.RUnlock()
	if int(object.size) < i || i < 0 {
		return ""
	}
	v := object.entries[i].Value
	switch v.(type) {
	case string:
		return v.(string)
	case fmt.Stringer:
		return v.(fmt.Stringer).String()
	}
	return ""
}

func (object *strObject) SetCompare(compare func(s1, s2 string) bool) *strObject {
	if compare == nil {
		return object
	}
	object.Lock()
	defer object.RLock()
	object.compare = compare
	return object
}

func (object *strObject) GetCompare() func(s1, s2 string) bool {
	object.Lock()
	defer object.Unlock()
	if object.compare == nil {
		object.compare = NatSortCompareAsc
	}
	return object.compare
}

func (object *strObject) String() string {
	var str = "{"
	for i, v := range object.entries {
		str += `"` + v.Key + `":"` + v.Value.(string) + `"`
		if i < int(object.size)-1 {
			str += ","
		}
		str += ""
	}
	str += "}"
	return str
}

func (object *strObject) KvAnyMap() map[string]interface{} {
	var data = make(map[string]interface{})
	object.RLock()
	defer object.RUnlock()
	for _, v := range object.entries {
		data[v.Key] = v.Value
	}
	return data
}

func (object *strObject) KvMaps() map[string]string {
	var data = make(map[string]string)
	object.RLock()
	defer object.RUnlock()
	for _, v := range object.entries {
		data[v.Key] = v.Value.(string)
	}
	return data
}

