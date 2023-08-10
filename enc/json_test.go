package enc_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

// Empty json.RawMessage should cause error in encoding/json, but not in enc.
func TestEmptyRawMessage(t *testing.T) {
	c := test.Context(t)
	empty := json.RawMessage{}

	_, err := json.Marshal(empty)
	test.Error(t, err)

	node, err := enc.Marshal(c, empty)
	test.NoError(t, err)
	b := enc.JSON{}.Encode(c, node)
	test.EqualsStr(t, "null", string(b)) // it started out as nothing (invalid json), but enc converts it to "null"

	err = enc.UnmarshalJSON(c, b, &empty)
	test.NoError(t, err)
	test.EqualsGo(t, json.RawMessage("null"), empty)
}

func TestStruct(t *testing.T) {
	c := test.Context(t)
	type X struct {
		S string  `json:"s"`
		I int     `json:"i"`
		F float64 `json:"f"`
		A any     `json:"a"`
	}
	x := X{S: "str", I: 42, F: 3.14, A: 3}
	n, err := enc.Marshal(c, x)
	test.NoError(t, err)
	test.ContainsGo(t, n, `42`)
	j := enc.JSON{}.Encode(c, n)
	test.EqualsStr(t, `{"s":"str","i":42,"f":3.14,"a":3}`, string(j))
}

func TestJSON(t *testing.T) {
	c := test.Context(t)
	codec := &enc.JSON{}
	check := func(t *testing.T, j string, nodeIn enc.Node) {
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
	checkLeft := func(t *testing.T, j string, nodeIn enc.Node) {
		t.Helper()
		t.Logf("%s <== %#v", j, nodeIn)
		jIn := []byte(j)
		jOut := codec.Encode(c, nodeIn)
		test.EqualsGo(t, string(jIn), string(jOut))
	}

	t.Run("scalars", func(t *testing.T) {
		check(t, `null`, enc.Nil{})
		checkLeft(t, `1`, enc.Integer(1))
		checkLeft(t, `3.14`, enc.Float(3.14))
		check(t, `3.14`, enc.Digits(`3.14`))
		check(t, `3`, enc.Digits(`3`))
		check(t, `true`, enc.Bool(true))
		check(t, `"foo"`, enc.String("foo"))
		check(t, `"\""`, enc.String(`"`))
		check(t, `"\\"`, enc.String(`\`))
		check(t, `"\\\""`, enc.String(`\"`))
	})

	t.Run("maps", func(t *testing.T) {
		check(t, `{}`, enc.Map{})
		checkLeft(t, `{"one":3.14}`, enc.Map{"one": enc.Float(3.14)})
		m := enc.Map{"one": enc.Integer(1), "nil": enc.Nil{}, "foo": enc.String("bar")}
		j := codec.Encode(c, m)
		test.ContainsJSON(t, j, `"nil":null`)
		test.ContainsJSON(t, j, `"foo":"bar"`)
		test.ContainsJSON(t, j, `"one":1`)
	})

	t.Run("pairs", func(t *testing.T) {
		checkLeft(t, `{}`, enc.Pairs{})
		checkLeft(t, `{"b":1,"a":2,"":null}`, enc.Pairs{{"b", "b", enc.Integer(1)}, {"a", "a", enc.Integer(2)}, {"", "", enc.Nil{}}})
	})

	t.Run("lists", func(t *testing.T) {
		check(t, `[]`, enc.List{})
		check(t, `[null]`, enc.List{enc.Nil{}})
		check(t, `[1,"two",false]`, enc.List{enc.Integer(1), enc.String("two"), enc.Bool(false)})
	})

	t.Run("deep", func(t *testing.T) {
		check(t, `[{}]`, enc.List{enc.Map{}})
		check(t, `[[],null]`, enc.List{enc.List{}, enc.Nil{}})
		check(t, `{"l":[]}`, enc.Map{"l": enc.List{}})
	})
}

func TestTime(t *testing.T) {
	c := test.Context(t)
	h := enc.Handler{
		Debugf: log.Debugf,
	}
	type X struct {
		T time.Time `json:"time"`
	}
	in := X{
		T: time.Now(),
	}
	n, err := h.Marshal(c, in)
	test.NoError(t, err)
	test.Contains(t, n.GoString(), fmt.Sprint(in.T.Year())) // enc.Pairs
	j := enc.JSON{}.Encode(c, n)
	test.ContainsJSON(t, j, fmt.Sprint(in.T.Year()))
	n2, err := enc.JSON{}.Decode(c, j)
	test.NoError(t, err)
	test.Contains(t, n2.GoString(), fmt.Sprint(in.T.Year())) // enc.Map
	var out X
	err = h.Unmarshal(c, n2, &out)
	test.NoError(t, err)
	test.EqualsGo(t, in, out)
}

func TestRawNode(t *testing.T) {
	c := test.Context(t)

	t.Run("enc.Node", func(t *testing.T) {
		var x struct {
			S string   `json:"s"`
			X enc.Node `json:"x"` // you can use enc.Node to have access to the interstitial similarly to json.RawMessage
		}

		in := enc.Map{
			"s": enc.String("foo"),
			"x": enc.Map{
				"ok": enc.Bool(true),
			},
		}

		err := enc.Unmarshal(c, in, &x)
		test.NoError(t, err)
		test.EqualsGo(t, enc.Map{"ok": enc.Bool(true)}, x.X)
		t.Logf("x: %+v", x)
	})

	t.Run("enc.Map", func(t *testing.T) {
		var x struct {
			S string  `json:"s"`
			X enc.Map `json:"x"` // you can use enc.Map to force the interstitial to be an object
		}

		in := enc.Map{
			"s": enc.String("foo"),
			"x": enc.Map{
				"ok": enc.Bool(true),
			},
		}

		err := enc.Unmarshal(c, in, &x)
		test.NoError(t, err)
		test.EqualsGo(t, enc.Map{"ok": enc.Bool(true)}, x.X)
		t.Logf("x: %+v", x)
	})

	t.Run("enc.Map FAIL", func(t *testing.T) {
		var x struct {
			S string  `json:"s"`
			X enc.Map `json:"x"` // you can use enc.Map to force the interstitial to be an object
		}

		in := enc.Map{
			"s": enc.String("foo"),
			"x": enc.String("bar"), // will fail to unmarshal into map
		}

		err := enc.Unmarshal(c, in, &x)
		test.Error(t, err)
		t.Logf("OK!")
	})
}
