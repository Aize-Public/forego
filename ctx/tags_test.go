package ctx_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/test"
)

func TestTags(t *testing.T) {
	c := context.Background()
	c = test.WithTester(c, t)

	fetch := func(c ctx.C) []any {
		t.Helper()
		var list []any
		err := ctx.RangeTag(c, func(key string, j ctx.JSON) error {
			t.Logf("tag[%s] = %s", key, string(j))
			var v any
			err := json.Unmarshal(j, &v)
			if err != nil {
				return err
			}
			list = append(list, v)
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return list
	}

	c = ctx.WithTag(c, "a", "one")
	{
		list := fetch(c)
		test.EqualsJSON(c, []any{"one"}, list)
	}

	c = ctx.WithTag(c, "b", "two")
	{
		list := fetch(c)
		test.EqualsJSON(c, []any{"one", "two"}, list)
	}

	c = ctx.WithTag(c, "a", "typo")
	{
		list := fetch(c)
		test.EqualsJSON(c, []any{"one", "two", "typo"}, list) // NOTE(oha): we range over all the assignments, no check for duplications
	}
}
