package test

import (
	"fmt"
	"testing"
)

func Nil(t *testing.T, obj any) {
	notNil(obj).fail(t)
}

func NotNil(t *testing.T, obj any) {
	notNil(obj).ok(t)
}

func NoError(t *testing.T, err error) {
	notNil(err).fail(t)
}

func Error(t *testing.T, err error) {
	notNil(err).ok(t)
}

func notNil(a any) res {
	switch a := a.(type) {
	case nil:
		return res{false, "nil"}
	default:
		return res{true, fmt.Sprint(a)}
	}
}
