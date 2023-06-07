package shutdown

import (
	"errors"

	"github.com/Aize-Public/forego/ctx"
)

// Use to signal that an operation can't be executed because the system is shutting down
var Err = errors.New("shutdown")

type ReleaseFn func()

func (fn ReleaseFn) Release() { fn() }

// start a global shutdown unless already started
func Begin() {
	shutdowner.begin()
}

// return a closed channel when the shutdown has started
func Started() <-chan struct{} {
	return shutdowner.started()
}

// return a closed channel 5 seconds after the shutdown has started
func Started5Sec() <-chan struct{} {
	return shutdowner.started5Sec()
}

// returns a channel that will close when the shutdown has completed
func Done() <-chan struct{} {
	return shutdowner.done()
}

// prevent the shutdown to complete until released
func Hold() ReleaseFn {
	return shutdowner.hold()
}

// prevents the shutdown to complete until released, and also wait for the shutdown to start
func HoldAndWait() ReleaseFn {
	return shutdowner.holdAndWait()
}

// setup signals and wait for the shutdown to complete
func WaitForSignal(c ctx.C, cf ctx.CancelFunc) {
	shutdowner.waitForSignal(c, cf)
}
