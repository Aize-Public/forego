package prom

import (
	"testing"
)

func mustPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		t.Helper()
		r := recover()
		if r == nil {
			t.Fatalf("should have panicked")
		}
	}()
	f()
}

func TestLabels(t *testing.T) {
	mustPanic(t, func() {
		stringify(
			[]string{"a", "b", "c"},
			[]string{"1", "2"},
		)
	})

	mustPanic(t, func() {
		stringify(
			[]string{"a", "b", "c"},
			nil,
		)
	})

	// temporary fix for empty labels
	if `` != stringify(nil, nil) {
		t.Fatalf("no longer?")
	}
}
