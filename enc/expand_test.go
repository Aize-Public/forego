package enc_test

import (
	"testing"

	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

func TestExpand(t *testing.T) {
	c := test.Context(t)
	t.Run("any", func(t *testing.T) {
		check := func(n enc.Node, obj any) {
			t.Logf("%+v", n)
			var x any
			err := enc.Codec{}.Expand(c, n, &x)
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
}
