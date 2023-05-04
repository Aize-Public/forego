package ctx

import "context"

func WithCancel(c C) (C, CancelFunc) {
	c, cf := context.WithCancelCause(c)
	return c, CancelFunc(cf)
}

func TODO() C {
	return context.TODO()
}
