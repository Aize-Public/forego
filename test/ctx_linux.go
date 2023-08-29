//go:build linux
// +build linux

package test

import (
	"os"

	"golang.org/x/sys/unix"
)

// this ugly thing is true if the output goes to a console, false if the output is piped somewhere
var isTerminal = func() bool {
	_, err := unix.IoctlGetTermios(int(os.Stdout.Fd()), unix.TCGETS)
	return err == nil
}()
