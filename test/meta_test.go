package test_test

import (
	"testing"

	"github.com/Aize-Public/forego/test"
)

// make sure we can test the tests
func TestTest(t *testing.T) {
	test.RunOk(t, "OK", func(t test.T) {
		t.Logf("all good: %s", t.Name())
	})
	test.RunFail(t, "Fatal", func(t test.T) {
		t.Fatalf("failing")
	})
	test.RunFail(t, "panic", func(t test.T) {
		panic("panics")
	})
	test.RunOk(t, "skip", func(t test.T) {
		t.SkipNow()
		panic("oops")
	})

	// META-META!!! check that test.RunOk() fails correctly
	test.RunFail(t, "meta", func(t test.T) {
		test.RunOk(t, "meta_failing", func(t test.T) {
			t.Fatalf("meta_fail")
		})
	})
	test.RunFail(t, "meta2", func(t test.T) {
		test.RunFail(t, "meta_ok", func(t test.T) {
			t.Logf("ok")
		})
	})
}
