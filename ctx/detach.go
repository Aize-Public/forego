package ctx

import "time"

// create a new context which inherits all the Values from the parent, but does NOT propagate cancels or timeouts
func Detach(c C) C {
	return detachedC{c}
}

type detachedC struct {
	parent C
}

func (this detachedC) Value(key any) any {
	return this.parent.Value(key)
}

func (this detachedC) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (this detachedC) Done() <-chan struct{} {
	return nil
}

func (this detachedC) Err() error {
	return nil
}
