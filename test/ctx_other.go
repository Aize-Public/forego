//go:build !linux
// +build !linux

package test

var isTerminal = func() bool {
	return true
}()
