package caches_test

import (
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils/caches"
)

func TestLRU(t *testing.T) {
	c := test.Context(t)
	lru := caches.NewLRU(c, 5, func(c ctx.C, k string) (int, int, error) {
		return len(k), len(k), nil
	})
	t.Logf("%v", lru)

	v, hm, err := lru.Get(c, "abc")
	t.Logf("%v", lru)
	test.NoError(t, err)
	test.EqualsGo(t, 3, v)
	test.EqualsGo(t, caches.MISS, hm)

	v, hm, err = lru.Get(c, "abc")
	t.Logf("%v", lru)
	test.NoError(t, err)
	test.EqualsGo(t, 3, v)
	test.EqualsGo(t, caches.HIT, hm)

	v, hm, err = lru.Get(c, "xyz")
	t.Logf("%v", lru)
	test.NoError(t, err)
	test.EqualsGo(t, 3, v)
	test.EqualsGo(t, caches.MISS, hm)

	v, hm, err = lru.Get(c, "xy")
	t.Logf("%v", lru)
	test.NoError(t, err)
	test.EqualsGo(t, 2, v)
	test.EqualsGo(t, caches.MISS, hm)

	v, hm, err = lru.Get(c, "xyz")
	t.Logf("%v", lru)
	test.NoError(t, err)
	test.EqualsGo(t, 3, v)
	test.EqualsGo(t, caches.HIT, hm)

	v, hm, err = lru.Get(c, "abc")
	t.Logf("%v", lru)
	test.NoError(t, err)
	test.EqualsGo(t, 3, v)
	test.EqualsGo(t, caches.MISS, hm)

	v, hm, err = lru.Get(c, "a")
	t.Logf("%v", lru)
	test.NoError(t, err)
	test.EqualsGo(t, 1, v)
	test.EqualsGo(t, caches.MISS, hm)
}
