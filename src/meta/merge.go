package meta

import (
	"errors"
	"reflect"

	d "github.com/teja2010/golconda/src/debug"
)

// LeftMerge takes three values of the same type:
//    left : the default values
//    right : incomplete but fresh values from the user
//    zero: zero values, to check if the value in right is incomplete
//
// rightVal such that, the values in right val if filled are left unmodified
// but if the value is "nil", then the value in leftVal is used to fill it
func LeftMerge(leftVal, rightVal, zeroVal interface{}) error {

	_leftv := reflect.ValueOf(leftVal)
	_zerov := reflect.ValueOf(zeroVal)
	_rightv := reflect.ValueOf(rightVal)

	if _leftv.Kind() != reflect.Ptr || _leftv.IsNil() {
		return errors.New("leftVal should be a pointer")
	}

	if _rightv.Kind() != reflect.Ptr || _rightv.IsNil() {
		return errors.New("rightVal should be a non-nil pointer")
	}

	if _zerov.Kind() != reflect.Ptr || _zerov.IsNil() {
		return errors.New("zeroVal should be a non-nil pointer")
	}

	return leftMerge(_leftv, _rightv, _zerov)
}

func leftMerge(_leftv, _rightv, _zerov reflect.Value) error {

	if err := typeCheck(_leftv, _rightv, _zerov); err != nil {
		return err
	}

	lElems := _leftv.Elem()
	rElems := _rightv.Elem()
	zElems := _zerov.Elem()

	for i := 0; i < lElems.NumField(); i++ {
		lEl := lElems.Field(i)
		rEl := rElems.Field(i)
		zEl := zElems.Field(i)

		if lEl.Kind() == reflect.Struct {
			if err := iterCopy(lEl, rEl, zEl); err != nil {
				d.Error("iterCopy of type", lEl.Type().String(),
					"failed")
				return err
			}
		} else if isRightEqualToDefault(rEl, zEl) {
			if err := safeCopy(rEl, lEl, zEl); err != nil {
				d.Error("safeCopy of type", lEl.Type().String(),
					"failed")
				return err
			}
		}
	}

	return nil
}

// right is zero, copy everything in left into right
func safeCopy(lEl, rEl, zEl reflect.Value) error {

	d.DebugLog("Copy <", lEl, "> into <", rEl, "> zero <", zEl, ">")

	if err := typeCheck(lEl, rEl, zEl); err != nil {
		return err
	}

	lEl.Set(rEl) // dst is right, source is left.
	// Opposite of all other Set* calls               _/\_

	d.DebugLog("Right set to", rEl)
	return nil
}

// if this is a struct, call LeftMerge again
func iterCopy(lEl, rEl, zEl reflect.Value) error {

	if lEl.Kind() == reflect.Struct {
		return leftMerge(lEl.Addr(), rEl.Addr(), zEl.Addr())
	}

	return nil
}

func isRightEqualToDefault(_rightv, _zerov reflect.Value) bool {
	d.DebugLog(_rightv.Kind(), _zerov.Kind())
	if _rightv.Kind() != _zerov.Kind() {
		return false
	}

	// TODO: this switch does not handle all cases.
	switch _rightv.Kind() {
	case reflect.Int:
		return _rightv.Int() == _zerov.Int()

	case reflect.Float32:
		return _rightv.Float() == _zerov.Float()

	case reflect.Float64:
		return _rightv.Float() == _zerov.Float()

	case reflect.String:
		return _rightv.String() == _zerov.String()

	case reflect.Slice:
		rlen := _rightv.Len()

		// zero slice is []
		// so return true if rlen is zero.
		return rlen == 0

		//zlen := _zerov.Len()
		//if rlen != zlen {
		//	return false
		//}
		//
		//rslice := _rightv.Slice(0, rlen)
		//zslice := _zerov.Slice(0, zlen)
		//
		//for i := 0; i < rlen; i++ {
		//	if rslice.Field(i) != zslice.Field(i) {
		//		return false
		//	}
		//}
		//return true
		//
		//eq := reflect.DeepEqual(_rightv.Slice(0, rlen), _zerov.Slice(0, zlen))
		//return eq

	case reflect.Struct:
		return false // structs need to be handled by iterCopy

	default:
		d.Warning("Unhandled Type", _rightv.Kind())
	}

	return false
}

func typeCheck(left, right, zero reflect.Value) error {

	d.DebugLog("Types", left.Type().String(), right.Type().String(),
		zero.Type().String())
	if left.Type() != right.Type() || left.Type() != zero.Type() {
		return errors.New("iterCopy Types dont match" +
			left.Type().String() +
			right.Type().String() +
			zero.Type().String())
	}

	return nil
}
