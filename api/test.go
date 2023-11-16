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
	test.NoError(t, err)

	j := &JSON{}
	err = h.Client().Send(c, op, j)
	if err != nil {
		test.Fail(t, "%+v", err)
	}
	t.Logf("Test[%T](%s)...", op, j.Data)
	req, err := h.Server().Recv(c, j)
	if err != nil {
		test.Fail(t, "%+v", err)
	}

	err = req.Do(test.Context(t))
	test.NoError(t, err)

	j = &JSON{}
	err = h.Server().Send(c, req, j)
	if err != nil {
		test.Fail(t, "%+v", err)
	}
	err = h.Client().Recv(c, j, op)
	if err != nil {
		test.Fail(t, "%+v", err)
	}
	return op
}
