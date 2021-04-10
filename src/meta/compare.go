package meta

import (
	d "github.com/teja2010/golconda/src/debug"
	"reflect"
)

// ALessThanB - Is A less than B ?
func ALessThanB(a, b interface{}, elem string) bool {
	d.DebugLog("elem", elem)
	_a := reflect.ValueOf(a)
	_b := reflect.ValueOf(b)

	_aEl := _a.FieldByName(elem)
	_bEl := _b.FieldByName(elem)

	if _aEl.IsZero() || _bEl.IsZero() {
		d.Warning("Elements are Zero")
	}

	return _ALessThanB(_aEl, _bEl)
}

// handle each type seperately
func _ALessThanB(_a, _b reflect.Value) bool {
	switch _a.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		inta := _a.Int()
		intb := _b.Int()

		return inta < intb

	case reflect.String:
		stra := _a.String()
		strb := _b.String()

		return stra < strb

	default:
		d.Bug("Unhandled Kind", _a.Kind())
	}

	return false
}

// Contains - check if a has the elem within
func Contains(a interface{}, elem string) bool {
	_a := reflect.ValueOf(a)
	ta := _a.Type()
	_, found := ta.FieldByName(elem)
	return found
}
