//go:build linux
// +build linux

package test

import (
	"os"

	"golang.org/x/sys/unix"
)

// NOTE: this is meant to control logging mode, but it doesn't work with go test, so consider using env instead
// this ugly thing is true if the output goes to a console, false if the output is piped somewhere
var isTerminal = func() bool {
	_, err := unix.IoctlGetTermios(int(os.Stdout.Fd()), unix.TCGETS)
	return err == nil
}()
