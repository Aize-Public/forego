package ctx

import (
	"context"
	"time"
)

func WithCancel(c C) (C, CancelFunc) {
	c, cf := context.WithCancelCause(c)
	return c, CancelFunc(cf)
}

func WithValue(c C, key, val any) C {
	return context.WithValue(c, key, val)
}

func TODO() C {
	return context.TODO()
}

func Background() (C, CancelFunc) {
	c, cf := context.WithCancelCause(context.Background())
	return c, CancelFunc(cf)
}

func Cause(c C) error {
	return context.Cause(c)
}

func WithTimeout(c C, d time.Duration) (C, func()) {
	return context.WithTimeout(c, d)
}

func WithDeadline(c C, t time.Time) (C, func()) {
	return context.WithDeadline(c, t)
}
