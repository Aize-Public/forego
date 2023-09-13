package ctx

import (
	"context"
)

// just an alias, so you can type `c ctx.C` instead of `ctx context.Context`
type C context.Context

type CancelFunc context.CancelCauseFunc

func (f CancelFunc) Exit() {
	f(nil)
}
