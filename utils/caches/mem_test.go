package caches_test

import (
	"io"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils/caches"
)

func TestMem(t *testing.T) {
	c := test.Context(t)
	cache := caches.Mem(c, func(c ctx.C, k string) (int, error) {
		t.Logf("cache miss for %q", k)
		switch k {
		case "":
			return 0, io.EOF
		default:
			return len(k), nil
		}
	})

	miss := func(k string, expect int) {
		t.Helper()
		v, hm, err := cache.Get(c, k)
		test.NoError(t, err)
		test.EqualsGo(t, expect, v)
		test.EqualsGo(t, caches.MISS, hm)
	}
	hit := func(k string, expect int) {
		t.Helper()
		v, hm, err := cache.Get(c, k)
		test.NoError(t, err)
		test.EqualsGo(t, expect, v)
		test.EqualsGo(t, caches.HIT, hm)
	}
	err := func(k string) {
		t.Helper()
		_, hm, err := cache.Get(c, k)
		test.Error(t, err)
		test.EqualsGo(t, caches.ERR, hm)
	}
	miss("one", 3)
	hit("one", 3)
	err("")
	hit("one", 3)
	miss("two", 3)
	hit("two", 3)
	hit("one", 3)
	err("")
}
