package enc_test

import (
	"testing"

	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

func TestExpand(t *testing.T) {
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
	})
}

func TestConflate(t *testing.T) {
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
}
