package test

import (
	"encoding/json"
	"testing"
)

// helper, returns the json value as string, or an error as string
// just to make tests easier to manage
// NOTE: we assume the error message is never a valid json, so there is no ambiguity
func JSON(v any) string {
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

func EqualsJSON(t testing.TB, expect, got any) {
	t.Helper()
	e := JSON(expect)
	g := JSON(got)
	if e == g {
		t.Logf("%s", e)
	} else {
		t.Fatalf("FAIL: expected %s got %s", e, g)
	}
}
