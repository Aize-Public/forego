package enc_test

import (
	"testing"

	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

func TestMap(t *testing.T) {
	c := test.Context(t)
	p := enc.Pairs{
		{JSON: "one", Value: enc.Integer(1)},
		{JSON: "none", Value: enc.Nil{}},
	}
	test.EqualsJSON(c, 1, p.Find("one"))
	test.EqualsJSON(c, nil, p.AsMap()["none"])
}
