package enc_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

func TestUnmarshal(t *testing.T) {
	c := test.Context(t)

	// Unmarshal into any
	t.Run("any", func(t *testing.T) {
		check := func(n enc.Node, obj any) {
			t.Logf("%+v", n)
			var x any
			err := enc.Handler{}.Unmarshal(c, n, &x)
			test.NoError(t, err)
			test.EqualsGo(t, obj, x)
		}
		check(
			enc.Map{"yes": enc.Bool(true)},
			map[string]any{"yes": true},
		)
		check(
			enc.List{enc.Number(3.14), enc.Nil{}},
			[]any{3.14, nil},
		)
		check(
			enc.Number(3.14),
			3.14,
		)
		check(
			enc.String("ok"),
			"ok",
		)
		check(
			enc.Nil{},
			nil,
		)
	})

	t.Run("map[string]any", func(t *testing.T) {
		check := func(n enc.Node, obj any) {
			t.Logf("%+v", n)
			var x map[string]any
			err := enc.Unmarshal(c, n, &x)
			test.NoError(t, err)
			test.EqualsGo(t, obj, x)
		}
		check(
			enc.Map{"yes": enc.Bool(true)},
			map[string]any{"yes": true},
		)
		check(
			enc.Map{"list": enc.List{enc.Nil{}, enc.Number(3.14)}},
			map[string]any{"list": []any{nil, 3.14}},
		)
		// TODO check for error if passing not a map
	})

	t.Run("*struct", func(t *testing.T) {
		type X struct {
			I int `json:"i"`
		}
		var y struct {
			X *X `json:"x"`
		}
		err := enc.Unmarshal(c, enc.Map{"x": enc.Map{"i": enc.Number(314)}}, &y)
		test.NoError(t, err)
		test.ContainsJSON(t, y, "314")

		y.X = nil
		err = enc.Unmarshal(c, enc.Map{}, &y)
		test.NoError(t, err)
		test.Nil(t, y.X)
		test.ContainsJSON(t, y, "null")

		y.X = nil
		err = enc.Unmarshal(c, enc.Map{"x": enc.Nil{}}, &y)
		test.NoError(t, err)
		test.Nil(t, y.X)
		test.ContainsJSON(t, y, "null")
	})
}

func TestMarshal(t *testing.T) {
	c := test.Context(t)

	x := struct {
		S string `json:"s"`
		I int    `json:"i"`
		V any    `json:"v"`
	}{
		S: "foo",
		I: 42,
		V: []any{nil, 2, true},
	}

	n, err := enc.Marshal(c, x)
	test.NoError(t, err)

	test.EqualsJSON(t, enc.Pairs{ // NOTE(oha): since we conflate a struct, we preserve the order of the fields using enc.Pairs
		{"S", "s", enc.String("foo")},
		{"I", "i", enc.Number(42)},
		{"V", "v", enc.List{
			enc.Nil{},
			enc.Number(2),
			enc.Bool(true),
		}},
	}, n)
	{
		n, err := enc.Marshal(c, []string{"foo", "bar"})
		test.NoError(t, err)
		t.Logf("%v", n)
	}
	{
		n, err := enc.Marshal(c, (map[string]any)(nil))
		test.NoError(t, err)
		switch n.(type) {
		case enc.Nil:
		default:
			test.Fail(t, "expected enc.Nil, got %T", n)
		}
	}
	{
		n, err := enc.Marshal(c, ([]any)(nil))
		test.NoError(t, err)
		switch n.(type) {
		case enc.Nil:
		default:
			test.Fail(t, "expected enc.Nil, got %T", n)
		}
	}
}

func TestCompat(t *testing.T) {
	c := test.Context(t)
	{
		n, err := enc.Marshal(c, map[int]string{3: "three"})
		test.NoError(t, err)
		t.Logf("n: %+v", n)
		j, _ := json.Marshal(map[int]string{3: "three"})
		test.EqualsStr(t, string(j), fmt.Sprint(n))
	}
	{
		var m map[int]string
		err := enc.Unmarshal(c, enc.Map{"3": enc.String("three")}, &m)
		test.NoError(t, err)
		test.EqualsGo(t, map[int]string{3: "three"}, m)
	}

	if false { // strictly json is different, but does it really matter?
		type Pair struct {
			x int
			y int
		}
		_, err := json.Marshal(map[Pair]string{{3, 4}: "three"})
		test.Error(t, err)

		_, err = enc.Marshal(c, map[Pair]string{{3, 4}: "three"})
		test.Error(t, err)
	}
}
