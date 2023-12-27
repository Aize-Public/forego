package caches_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils/caches"
)

func TestRequest(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	c := test.Context(t)
	answers := map[string]chan int{}
	var waiting atomic.Int32
	k := caches.NewRequestLRU(c, 4, func(c ctx.C, k string) (int, time.Duration, error) {
		waiting.Add(1)
		defer waiting.Add(-1)
		t.Logf("waiting for %q", k)
		v := <-answers[k]
		t.Logf("got %v for %q", v, k)
		return v, time.Second, nil
	})
	test.EqualsGo(t, 0, waiting.Load())

	answers["one"] = make(chan int)
	answers["two"] = make(chan int)

	send := func(key string, expected int, flag caches.HitMiss) {
		t0 := time.Now()
		v, hit, err := k.Request(c, key)
		t.Logf("req(%q) =>  %v (%v) %v", key, v, hit, err)
		test.NoError(t, err)
		test.EqualsGo(t, expected, v)
		test.EqualsGo(t, flag, hit)
		if time.Since(t0) < 500*time.Microsecond {
			test.Fail(t, "didn't block")
		}
	}

	go send("one", 1, caches.MISS)
	time.Sleep(time.Millisecond)

	go send("one", 1, caches.HIT)
	time.Sleep(time.Millisecond)

	test.EqualsGo(t, 1, waiting.Load()) // only one allowed in, the other should be waiting for the same response

	go send("two", 2, caches.MISS)
	time.Sleep(time.Millisecond)

	test.EqualsGo(t, 2, waiting.Load()) // the second for a different key should go thru, counting 2 waiting now

	answers["one"] <- 1

	time.Sleep(time.Millisecond)
	test.EqualsGo(t, 1, waiting.Load()) // one got a response, and it's used for both requests, two is still waiting

	answers["two"] <- 2

	time.Sleep(time.Millisecond)
	test.EqualsGo(t, 0, waiting.Load()) // no more waiting

	send("one", 1, caches.HIT) // cache hits after won't block neither
}
