package test

import "testing"

func TestNil(t *testing.T) {
	isNil(nil).true(t)
	var noErr error
	isNil(noErr).true(t)
	isNil((any)(noErr)).true(t)
	func(x any) {
		isNil(x).true(t)
	}(nil)

	isNil(1).false(t)
	isNil(true).false(t)
	isNil(0).false(t)
}
