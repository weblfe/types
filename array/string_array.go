package array

import (
	"fmt"
)

type StrArr []string

func ToStrArr(arr []string) StrArr {
  return arr
}

func (arr StrArr) Len() int {
	return len(arr)
}

func (arr StrArr) IsNil() bool {
	return arr == nil
}

func (arr StrArr) Empty() bool {
	return len(arr) == 0
}

func (arr StrArr) Index(i int) string {
	if i >= 0 && i < arr.Len() {
		return arr[i]
	}
	return ""
}

func (arr StrArr) Include(v string) bool {
	for _, _v := range arr {
		if v == _v {
			return true
		}
	}
	return false
}

func (arr StrArr) Find(v string) int {
	for i, _v := range arr {
		if v == _v {
			return i
		}
	}
	return -1
}

func (arr StrArr) Append(v string) StrArr {
	return append(arr, v)
}

func (arr StrArr) Head(v string) StrArr {
	return append([]string{v}, arr...)
}

func (arr StrArr) String() string {
	return fmt.Sprintf("%v", []string(arr))
}

func (arr StrArr) First() string {
	if arr.Len() > 0 {
		return arr[0]
	}
	return ""
}

func (arr StrArr) Last() string {
	var size = arr.Len()
	if size > 0 {
		return arr[size]
	}
	return ""
}

func (arr StrArr) Merge(source []string) StrArr {
	return append(arr, source...)
}

func (arr StrArr) Unique() StrArr {
	var (
		newArr = StrArr{}
		_cache = map[string]bool{}
	)
	for _, v := range arr {
		if _, ok := _cache[v]; ok {
			continue
		}
		_cache[v] = true
		newArr = append(newArr, v)
	}
	return newArr
}

func (arr StrArr) Diff(newArr []string) StrArr {
	var (
		_arr = StrArr{}
	)
	for _, v := range arr {
		if StrArr(newArr).Find(v) < 0 {
			_arr = append(_arr, v)
		}
	}
	return _arr
}

func (arr StrArr) RowArr() []string {
	return arr
}
