# `forego/test`

Very similar to any other test libraries, but parses the original test source code to generate helpful log messages:

```go
func TestReadme(t *testing.T) {
	err := foo()
	test.NoError(t, err)
	test.EqualsGo(t, 2*2, add(2, 2))
	test.EqualsGo(t, 2*2, add(2, 3))
}
func foo() error { return nil }
func add(a, b int) int { return a + b }
```

Output:
```
    readme_test.go:17:   ✔ no error: `foo()`
    readme_test.go:19:   ✔ EqualsGo(2*2, add(2, 2)) ⮕  4
    readme_test.go:20: ❌  EqualsGo(2*2, add(2, 3)) ⮕  4 != 5
```

As opposed to testify:

```go
func TestReadme(t *testing.T) {
	err := foo()
	assert.NoError(t, err, "foo() has no error")

	assert.Equal(t, 2*2, add(2, 2), "add(2,2) expect 2*2")
	assert.Equal(t, 2*2, add(2, 3), "add(2,3) expect 2*2")
}
func foo() error { return nil }
func add(a, b int) int { return a + b }
```

Output:
```
    main_test.go:17: 
        	Error Trace:	x/main_test.go:17
        	Error:      	Not equal: 
        	            	expected: 4
        	            	actual  : 5
        	Test:       	TestIfy
```

Which does not provide any log for success, and wants an extra argument on each test to be easy to read.


## Functions

### `OK(t, fmt, args...)` and `Fail(…)`

Provoke a log message as the above, just to make it easier to make custom checks


### `Assert(t, cond)`

The simplest test, will log the argument source of the condition

```
  test.Assert(t, 7 > 2*2)
  test.Assert(t, 7 <= 2*2)
```

```
    readme_test.go:22:   ✔ `7 > 2*2`
    readme_test.go:23: ❌  `3 <= 2`
```


### `Equal…()` and `NotEqual…()`

```go
  test.EqualsGo(t, 123, sum(81,42)) // compare using fmt.Sprintf("%#v")
  test.EqualsStr(t, "123", "12"+"3") // compare strings
  test.EqualsJSON(c, []any{1}, []int{1})  // compare using enc.MarshalJSON()
```


### `Nil()` and `NotNil()`

```go
  test.Nil(t, (*int)(nil)) // check if nil, even if it has a type
```

Note: go does not consider `(*int)(nil)` as nil, but this does


### `Empty()` and `NotEmpty()`

Similar to `Nil()`, but has special meaning based on the type:

* pointer: check if nil
* slice, map, chan: check if len() == 0
* string: check if ""
* struct: check if it is the zero value


### `Error()` and `NoError()`

```go
  err := foo()
  test.NoError(t, err)
```


## JSON and `jsonify(any) string`

Internal function that converts any object to a JSON string, used by all the `…JSON(t, …)` functions.

Note: If the given is a `string` or a `[]byte` which is a valid JSON, it won't double-encode it.


## `ctx.C`

```go
  c := test.Context(t)
```

Provides a `ctx.C` which logs to `t.Logf()` and cancel when the tests are over


## TODO

* parse comments and use them to improve logging
* add more test functions for common use cases
