package enc_test

import (
	"testing"

	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

func TestNumeric(t *testing.T) {
	c := test.Context(t)
	i := int64(1000000234567890123) // big enough to be rounded as float64
	t.Logf("i64: %d", i)
	j := enc.MustMarshalJSON(c, i)
	test.EqualsStr(t, string(j), `1000000234567890123`)
	test.NoError(t, enc.UnmarshalJSON(c, j, &i))
	test.EqualsGo(t, 1000000234567890123, i)
	var a any
	test.NoError(t, enc.UnmarshalJSON(c, j, &a))
	_ = a.(float64)
	t.Logf("a: %v", a)
}

func TestFloat(t *testing.T) {
	c := test.Context(t)
	{
		f := float64(42.42)
		j, err := enc.MarshalJSON(c, f)
		test.NoError(t, err)
		f = 0
		test.NoError(t, enc.UnmarshalJSON(c, j, &f))
		test.Assert(t, f > 42 && f < 43)
	}
	{
		// Test unmarshaling float into other stuff
		j := []byte(`42.42`)
		var i int64
		test.NoError(t, enc.UnmarshalJSON(c, j, &i))
		test.EqualsGo(t, int64(42), i)
		var u uint
		test.NoError(t, enc.UnmarshalJSON(c, j, &u))
		test.EqualsGo(t, uint(42), u)
		var s string
		test.NoError(t, enc.UnmarshalJSON(c, j, &s))
		test.EqualsGo(t, "42.42", s)
	}
}
