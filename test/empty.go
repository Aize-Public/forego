package test

import (
	"fmt"
	"reflect"
	"testing"
)

func NotEmpty(t *testing.T, obj any) {
	t.Helper()
	empty(obj).argument(0, 1).false(t)
}

func Empty(t *testing.T, obj any) {
	t.Helper()
	empty(obj).argument(0, 1).true(t)
}

func empty(obj any) res {
	if obj == nil {
		return res{true, fmt.Sprintf("is nil")}
	}
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Pointer:
		// TODO(oha): is this a good idea? if we pass a io.Reader it will say true even if it contains no data...
		if v.IsNil() {
			return res{true, fmt.Sprintf("%T is nil", obj)}
		} else {
			return res{false, fmt.Sprintf("%T %s", obj, obj)}
		}
	case reflect.Slice, reflect.Map, reflect.Chan:
		if v.IsNil() {
			return res{true, fmt.Sprintf("%T is nil", obj)}
		} else if v.Len() == 0 {
			return res{true, fmt.Sprintf("%T is empty", obj)}
		} else {
			return res{false, fmt.Sprintf("%T has %d elements", obj, v.Len())}
		}
	case reflect.String:
		if v.String() == "" {
			return res{true, `is ""`}
		} else {
			return res{false, fmt.Sprintf("is %q", v.String())}
		}
	case reflect.Struct:
		if v.IsZero() {
			return res{true, fmt.Sprintf("%T is zero", obj)}
		} else {
			return res{false, fmt.Sprintf("%T %+v", obj, obj)}
		}
	default:
		panic(fmt.Sprintf("can't test emptyness for %T", obj))
	}
}
