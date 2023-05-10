package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Aize-Public/forego/utils/ast"
)

// helper, returns the jsonish value as string, or an error as string
// just to make tests easier to manage
// NOTE: we assume the error message is never a valid jsonish, so there is no ambiguity
func jsonish(v any) string {
	switch v := v.(type) {
	case json.RawMessage:
		return string(v)
	case []byte:
		if json.Valid(v) {
			return string(v)
		}
	}
	j, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(j)
}

type res struct {
	succeed bool
	msg     string
}

func (res res) argument(above, argNum int) res {
	call, _ := ast.Caller(above + 1)
	res.msg = call.Args[argNum].Src + ": " + res.msg
	return res
}

func (res res) assignment(above, argNum int) res {
	res.msg = ast.Assignment(above+1, argNum) + ": " + res.msg
	return res
}

func (res res) prefix(f string, args ...any) res {
	res.msg = fmt.Sprintf(f, args...) + ": " + res.msg
	return res
}

// expect true
func (res res) true(t *testing.T, f ...any) {
	t.Helper()
	if res.succeed {
		t.Logf("OK %s", res.msg)
	} else {
		t.Fatalf("FAIL %s", res.msg)
	}
}

// expect false
func (res res) false(t *testing.T) {
	t.Helper()
	if res.succeed {
		t.Fatalf("FAIL %s", res.msg)
	} else {
		t.Logf("OK %s", res.msg)
	}
}
