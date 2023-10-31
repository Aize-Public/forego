package api

import (
	"testing"

	"github.com/Aize-Public/forego/test"
)

// TODO(oha): should we do marshalling here to be safe?
func Test[T Op](t *testing.T, op T) T {
	c := test.Context(t)
	t.Helper()
	h, err := NewHandler(c, op)
	if err != nil {
		panic(err)
	}

	j := &JSON{}
	err = h.Client().Send(c, op, j)
	if err != nil {
		panic(err)
	}
	t.Logf("Test[%T](%s)...", op, j.Data)
	req, err := h.Server().Recv(c, j)

	err = req.Do(test.Context(t))
	test.NoError(t, err)

	j = &JSON{}
	err = h.Server().Send(c, req, j)
	if err != nil {
		panic(err)
	}
	err = h.Client().Recv(c, j, op)
	if err != nil {
		panic(err)
	}
	return op
}
