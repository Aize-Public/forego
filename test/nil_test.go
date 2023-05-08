package test

import "testing"

func TestNil(t *testing.T) {
	notNil(nil).fail(t)
	var noErr error
	notNil(noErr).fail(t)
	notNil((any)(noErr)).fail(t)
	func(x any) {
		notNil(x).fail(t)
	}(nil)

	notNil(1).ok(t)
	notNil(false).ok(t)
	notNil(0).ok(t)
}
