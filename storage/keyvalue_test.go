package storage_test

import (
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/storage"
	"github.com/Aize-Public/forego/test"
)

func TestKeyvalue(t *testing.T) {
	c := test.Context(t)
	kv := storage.NewMemKeyValue()
	{
		n, err := kv.Get(c, "one")
		test.NoError(t, err)
		test.Nil(t, n)
	}
	{
		err := kv.Upsert(c, "one", enc.Map{"num": enc.Integer(1)})
		test.NoError(t, err)
	}
	{
		n, err := kv.Get(c, "one")
		test.NoError(t, err)
		test.EqualsJSON(t, `{"num":1}`, n)
	}
	{
		err := kv.Upsert(c, "two", enc.Map{"num": enc.Integer(2), "type": enc.String("foo")})
		test.NoError(t, err)
	}
	{
		n, err := kv.Get(c, "one")
		test.NoError(t, err)
		test.EqualsJSON(t, `{"num":1}`, n)
	}
	{
		tot := 0
		err := kv.Range(c, func(c ctx.C, k string, v enc.Map) error {
			tot += int(v["num"].(enc.Integer))
			return nil
		})
		test.NoError(t, err)
		test.EqualsGo(t, 3, tot)
	}
	{
		vals := []any{}
		err := kv.Range(c, func(c ctx.C, k string, v enc.Map) error {
			vals = append(vals, v["num"])
			return nil
		}, storage.Filter{
			Field: "type",
			Cmp:   storage.Equal,
			Val:   enc.String("foo"),
		})
		test.NoError(t, err)
		test.EqualsJSON(t, `[2]`, vals)
	}
	{
		err := kv.Upsert(c, "three", enc.Map{"num": enc.Integer(3)})
		test.NoError(t, err)
		err = kv.Delete(c, "three")
		test.NoError(t, err)
		val, err := kv.Get(c, "three")
		test.NoError(t, err)
		test.Nil(t, val)
	}
}
