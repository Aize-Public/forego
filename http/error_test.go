package http_test

import (
	"io"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/http"
	"github.com/Aize-Public/forego/test"
)

func TestError(t *testing.T) {
	c := test.Context(t)
	{
		err := io.EOF
		test.EqualsGo(t, 999, http.ErrorCode(err, 999))
	}
	{
		err := io.EOF
		err = http.Error{401, err}
		test.EqualsGo(t, 401, http.ErrorCode(err, 999))
	}
	{
		err := io.EOF
		err = http.Error{401, err}
		err = ctx.NewErrorf(c, "err: %w", err)
		test.EqualsGo(t, 401, http.ErrorCode(err, 999))
	}
	{
		err := io.EOF
		err = ctx.NewErrorf(c, "err: %w", err)
		err = http.Error{401, err}
		err = ctx.NewErrorf(c, "err: %w", err)
		test.EqualsGo(t, 401, http.ErrorCode(err, 999))
	}
	{
		err := io.EOF
		err = ctx.NewErrorf(c, "err: %w", err)
		err = http.Error{401, err}
		err = ctx.NewErrorf(c, "err: %w", err)
		err = http.Error{403, err}
		err = ctx.NewErrorf(c, "err: %w", err)
		test.EqualsGo(t, 403, http.ErrorCode(err, 999))
	}
}
