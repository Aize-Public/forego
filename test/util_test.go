package test

import (
	"testing"
)

func TestJsonify(t *testing.T) {
	// make sure we canonicalize before comparing
	c := Context(t)
	EqualsJSON(c, "[ 1 ]", []int{1})
	EqualsJSON(c, `{   "a":3,"c":1    ,"b":2}`, map[string]int{"c": 1, "b": 2, "a": 3})
}
