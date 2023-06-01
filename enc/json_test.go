package enc_test

import (
	"testing"

	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

func TestJSON(t *testing.T) {
	c := test.Context(t)
	codec := &enc.JSON{}
	check := func(j string, nodeIn enc.Node) {
		t.Helper()
		t.Logf("%s <=> %#v", j, nodeIn)
		jIn := []byte(j)
		jOut := codec.Encode(c, nodeIn)
		test.EqualsGo(t, string(jIn), string(jOut))
		nodeOut, err := codec.Decode(c, jIn)
		test.NoError(t, err)
		test.EqualsGo(t, nodeIn, nodeOut)
	}
	// only check from node to json
	checkLeft := func(j string, nodeIn enc.Node) {
		t.Helper()
		t.Logf("%s <== %#v", j, nodeIn)
		jIn := []byte(j)
		jOut := codec.Encode(c, nodeIn)
		test.EqualsGo(t, string(jIn), string(jOut))
	}

	t.Run("scalars", func(t *testing.T) {
		check(`null`, enc.Nil{})
		check(`1`, enc.Number(1))
		check(`3.14`, enc.Number(3.14))
		check(`true`, enc.Bool(true))
		check(`"foo"`, enc.String("foo"))
		check(`"\""`, enc.String(`"`))
		check(`"\\"`, enc.String(`\`))
		check(`"\\\""`, enc.String(`\"`))
	})

	t.Run("maps", func(t *testing.T) {
		check(`{}`, enc.Map{})
		check(`{"one":3.14}`, enc.Map{"one": enc.Number(3.14)})
		m := enc.Map{"one": enc.Number(1), "nil": enc.Nil{}, "foo": enc.String("bar")}
		j := codec.Encode(c, m)
		test.ContainsJSON(t, j, `"nil":null`)
		test.ContainsJSON(t, j, `"foo":"bar"`)
		test.ContainsJSON(t, j, `"one":1`)
	})

	t.Run("pairs", func(t *testing.T) {
		checkLeft(`{}`, enc.Pairs{})
		checkLeft(`{"b":1,"a":2,"":null}`, enc.Pairs{{"b", "b", enc.Number(1)}, {"a", "a", enc.Number(2)}, {"", "", enc.Nil{}}})
	})

	t.Run("lists", func(t *testing.T) {
		check(`[]`, enc.List{})
		check(`[null]`, enc.List{enc.Nil{}})
		check(`[1,"two",false]`, enc.List{enc.Number(1), enc.String("two"), enc.Bool(false)})
	})

	t.Run("deep", func(t *testing.T) {
		check(`[{}]`, enc.List{enc.Map{}})
		check(`[[],null]`, enc.List{enc.List{}, enc.Nil{}})
		check(`{"l":[]}`, enc.Map{"l": enc.List{}})
	})
}
