package enc_test

import (
	"testing"

	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

func TestPairs(t *testing.T) {
	c := test.Context(t)
	var x struct {
		Foo int `json:"foo"`
		Bar any `json:"bar"`
	}
	p := enc.Pairs{
		{Name: "Foo", JSON: "foo", Value: enc.Integer(3)},
		{Name: "Bar", JSON: "bar", Value: enc.Integer(7)},
	}
	enc.MustUnmarshal(c, p, &x)
	test.EqualsJSON(t, 3, x.Foo)
	test.EqualsJSON(t, 7, x.Bar)
}
