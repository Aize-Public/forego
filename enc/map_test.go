package enc_test

import (
	"testing"

	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

func TestMap(t *testing.T) {
	p := enc.Pairs{
		{JSON: "one", Value: enc.Integer(1)},
		{JSON: "none", Value: enc.Nil{}},
	}
	test.EqualsJSON(t, 1, p.Find("one"))
	test.EqualsJSON(t, nil, p.AsMap()["none"])
}
