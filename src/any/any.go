package any

import (
	. "github.com/teja2010/golconda/src/debug"
)

const (
	NoneType valType = -1
	IntType valType = 0
	StrType valType = 1
	MapType valType = 2
)

type valType int

type AnyValue struct {
	Type valType
	_val interface{}
}

// Map which has maps within
type TreeMap map[string]AnyValue

func NoneValue() AnyValue {
	return AnyValue{NoneType, nil}
}

func IntValue(i int) AnyValue {
	return AnyValue{IntType, i}
}

func StrValue(s string) AnyValue {
	return AnyValue{StrType, s}
}

func MapValue(m TreeMap) AnyValue {
	return AnyValue{MapType, m}
}

func (a AnyValue) Int() int {
	if DebugCheck(a.Type != IntType) {
		Bug("Not Int type");
	}

	return a._val.(int)
}

func (a AnyValue) Str() string {
	if DebugCheck(a.Type != StrType) {
		Bug("Not Str Type");
	}

	return a._val.(string)
}

func (a AnyValue) Map() TreeMap {
	if DebugCheck(a.Type != MapType) {
		Bug("Not Map Type");
	}

	return a._val.(TreeMap)
}

func (t TreeMap) Get(key string) AnyValue {
	val, ok := t[key]
	if DebugCheck(!ok) {
		Bug("Value Not found")
	}

	return val
}

func (t TreeMap) Set(key string, val AnyValue) {
	t[key] = val
}

func ReadMap(file_name string) TreeMap {
	t := make(map[string]AnyValue)
	return t
}

// values in t2 will be overwrite values in t1.
// finally return t1
func RightUpdate(t1 TreeMap, t2 TreeMap) TreeMap {
	// TODO
	return t1
}
