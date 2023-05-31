package sync_test

import (
	"testing"

	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils/sync"
)

func TestMapNil(t *testing.T) {
	m := sync.Map[string, *string]{}

	v, exists := m.Load("no")
	test.Assert(t, !exists)
	test.Assert(t, v == nil)
}

func TestMapIntInt(t *testing.T) {
	m := &sync.Map[int, int]{}

	v, exists := m.Load(3)
	test.Assert(t, !exists)
	test.EqualsGo(t, 0, v)

	v, exists = m.LoadOrStore(1, 2)
	test.Assert(t, !exists)
	test.EqualsGo(t, 2, v)

	m.Store(3, 4)
	sum := 0
	m.Range(func(k int, v int) bool {
		t.Logf("k: %v, v: %v", k, v)
		sum += k*1000 + v
		return true
	})
	test.EqualsGo(t, 4006, sum) // k: 1+3, v: 2+4

	v, exists = m.Load(3)
	test.Assert(t, exists)
	test.EqualsGo(t, 4, v)
}
