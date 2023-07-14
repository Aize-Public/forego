package ctx

import (
	"fmt"
	"sync/atomic"

	"github.com/google/uuid"
)

type trackingCtx struct {
	C
	ID   string
	Last int32
}

// create a new tracking id embedded in the context, the suggested id can be used if valid, otherwise a new id is created
func WithTracking(c C, suggest string) C {
	if suggest == "" {
		suggest = uuid.NewString()
	}
	return trackingCtx{
		C:    WithTag(c, "tracking-id", suggest),
		ID:   suggest,
		Last: 0,
	}
}

type trackingGet struct{}

func GetTracking(c C) string {
	v := c.Value(trackingGet{})
	if v == nil {
		return ""
	} else {
		return v.(string)
	}
}

type trackingNext struct{}

func (c trackingCtx) Value(k any) any {
	switch k.(type) {
	case trackingNext:
		step := atomic.AddInt32(&c.Last, 1)
		return fmt.Sprintf("%s.%x", c.ID, step)
	case trackingGet:
		return c.ID
	default:
		return c.C.Value(k)
	}
}

// TEMP
func Span(c C, name string) (C, CancelFunc) {
	// TODO add opentelemetry support
	c, cf := WithCancel(c)
	k := c.Value(trackingNext{})
	if k == nil {
		k = ""
	}
	return WithTracking(c, k.(string)), cf
}
