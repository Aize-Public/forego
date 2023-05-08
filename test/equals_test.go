package test

import (
	"encoding/json"
	"testing"
)

func TestEquals(t *testing.T) {
	equalJSON(false, false).true(t)
	equalJSON(1, 1).true(t)
	equalJSON(1, 2).false(t)
	equalJSON(1, "1").false(t)
	equalJSON(1, 1.0).true(t)

	// json can be compared directly
	equalJSON(1, json.RawMessage(`1`)).true(t)
	equalJSON([]byte("null"), nil).true(t)

	// array types don't matter
	equalJSON([]int{1, 2}, []any{1.0, 2.0}).true(t)

	// map and struct are just an object
	equalJSON(
		map[string]int{"one": 1},
		struct {
			One int `json:"one"`
		}{1},
	).true(t)
}
