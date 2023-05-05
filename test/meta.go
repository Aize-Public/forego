package test

import (
	"log"
	"runtime"
	"testing"

	"github.com/Aize-Public/forego/utils/ast"
)

// used to test the tests
type T interface {
	Helper()
	Logf(f string, args ...any)
	Fatalf(f string, args ...any)
	Name() string
	SkipNow()
	Setenv(key, value string)
}

var _ T = &testing.T{}

// implementation of test.T
type testT struct {
	name    string
	Logger  func(f string, args ...any)
	failed  bool
	skipped bool
	helper  func() // Note(oha): this is hard to simplify, we need to carry the pointer to the original t.Helper function
}

// used to distinguish normal panic from SkipNow() and FailNow()
var panicVal = &struct{}{}

var _ T = &testT{}

// MAKE IT LOOK LIKE A *testing.T
func (t *testT) Fatal(args ...any) {
	t.Logf("FATAL: %v", args)
	t.failed = true
	panic(panicVal)
}

func (t *testT) Helper() {
	// we really can't do anything about this, testing.T hardcoded the call above value
}
func (t *testT) Fail() {
	t.failed = true
	panic(panicVal)
}

func (t *testT) Skipped() bool {
	return t.skipped
}
func (t *testT) FailNow() {
	t.failed = true
	panic(panicVal)
}
func (t *testT) Failed() bool {
	return t.failed
}
func (t *testT) Cleanup(f func()) {
}
func (t *testT) Error(args ...any) {
	t.Helper()
	t.Logf("ERR: %v", args)
}
func (t *testT) Errorf(f string, args ...any) {
	t.helper()
	if t.Logger == nil {
		log.Printf(f, args...)
	} else {
		t.Logger(f, args...)
	}
}
func (t *testT) Log(args ...any) {
	t.helper()
	t.Logf("%v", args)
}
func (t *testT) Logf(f string, args ...any) {
	t.helper()
	/*
		_, me, _, _ := runtime.Caller(0)
		for i := 1; i < 100; i++ {
			_, file, line, _ := runtime.Caller(i)
			if file != me {
				f = fmt.Sprintf("%s:%d %s", filepath.Base(file), line, f)
				break
			}
		}
	*/
	if t.Logger == nil {
		log.Printf(f, args...)
	} else {
		t.Logger(f, args...)
	}
}
func (t *testT) Fatalf(f string, args ...any) {
	t.helper()
	t.Logf(f, args...)
	t.failed = true
	panic(panicVal)
}
func (t *testT) Name() string {
	return t.name
}
func (t *testT) SkipNow() {
	t.skipped = true
	panic(panicVal)
}
func (t *testT) Skip(args ...any) {
	t.helper()
	t.Logf("SKIP: %+v", args)
	t.skipped = true
	panic(panicVal)
}
func (t *testT) Skipf(f string, args ...any) {
	t.helper()
	t.Logf("SKIP: "+f, args)
	t.skipped = true
	panic(panicVal)
}

func RunOk(t T, name string, f func(t T)) {
	t.Helper()
	t.Logf("RUN expect OK %q", name)
	tt := testT{
		helper: t.Helper,
		Logger: t.Logf,
	}
	if k, ok := t.(*testT); ok {
		tt.helper = k.helper
	}
	ok, _ := tt.run(name, f)
	if ok {
		t.Logf("ok: %s", ast.Assignment(1, 1))
	} else {
		t.Fatalf("FAIL: %s", ast.Assignment(1, 1))
	}
}

func RunFail(t T, name string, f func(t T)) {
	t.Helper()
	t.Logf("RUN expect FAIL %q", name)
	tt := testT{
		helper: t.Helper,
		Logger: t.Logf,
	}
	if k, ok := t.(*testT); ok {
		tt.helper = k.helper
	}
	ok, _ := tt.run(name, f)
	if !ok {
		t.Logf("failed ok: %s", ast.Assignment(1, 1))
	} else {
		t.Fatalf("FAIL EXPECTED: %s", ast.Assignment(1, 1))
	}
}

func (t *testT) Run(name string, f func(t T)) bool {
	t.helper()
	t.Logf("RUN %q", name)
	ok, code := t.run(name, f)
	if ok {
		t.Logf("OK %q %q", name, code)
	} else {
		t.Logf("FAIL %q %q", name, code)
	}
	return ok
}

func (o *testT) run(name string, f func(t T)) (bool, string) {
	o.helper()
	t := &testT{
		helper: o.helper,
		Logger: func(f string, args ...any) {
			o.helper()
			o.Logf("  "+f, args...)
		},
	}
	t.name = name
	r := func() (r any) {
		o.helper()
		defer func() {
			o.helper()
			r = recover()
			if r != nil && r != panicVal {
				t.Logf("panic: %v", r)
				for i := 1; i < 100; i++ {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					t.Logf("  %s:%d\n", file, line)
				}
			}
		}()
		f(t)
		return
	}()
	if t.skipped {
		return true, "skip"
	}
	if r != nil {
		if t.failed {
			return false, "fail"
		}
		return false, "panic"
		//f(&t) // run again
	}
	return true, "ok"
}

func (t *testT) Setenv(key, value string) {}
