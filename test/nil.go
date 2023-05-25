package test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Aize-Public/forego/utils/ast"
)

func Nil(t *testing.T, obj any) {
	t.Helper()
	isNil(obj).argument(0, 1).true(t)
}

func NotNil(t *testing.T, obj any) {
	t.Helper()
	isNil(obj).argument(0, 1).false(t)
}

func NoError(t *testing.T, err error) {
	t.Helper()
	if isNil(err).succeed {
		OK(t, "no error: %s", ast.Assignment(0, 1))
	} else {
		Fail(t, "%v", err)
	}
}

func Error(t *testing.T, err error) {
	t.Helper()
	isNil(err).assignment(0, 1).false(t)
}

func isNil(a any) res {
	switch a := a.(type) {
	case nil:
		return res{true, "nil"}
	default:
		v := reflect.ValueOf(a)
		switch v.Kind() {
		case reflect.Slice, reflect.Map, reflect.Chan, reflect.Pointer:
			return res{v.IsNil(), fmt.Sprintf("%#v", a)}
		default:
			return res{false, fmt.Sprint(a)}
		}
	}
}
