package test

import (
	"encoding/json"
	"testing"
)

func TestEquals(t *testing.T) {
	c := Context(t)

	equalJSON(c, false, false).true(t)
	equalJSON(c, 1, 1).true(t)
	equalJSON(c, 1, 2).false(t)
	equalJSON(c, 1, "1").true(t)
	equalJSON(c, 1, 1.0).true(t)

	// json can be compared directly
	equalJSON(c, 1, json.RawMessage(`1`)).true(t)
	equalJSON(c, []byte("null"), nil).true(t)

	// array types don't matter
	equalJSON(c, []int{1, 2}, []any{1.0, 2.0}).true(t)

	// map and struct are just an object
	equalJSON(c,
		map[string]int{"one": 1},
		struct {
			One int `json:"one"`
		}{1},
	).true(t)
}
