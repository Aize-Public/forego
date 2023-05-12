package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Aize-Public/forego/ctx"
)

func NewErrorf(c ctx.C, code int, f string, args ...any) Error {
	return Error{
		Code: code,
		Err:  ctx.NewErrorf(c, f, args...),
	}
}

type Error struct {
	Code int
	Err  error
}

func (this Error) Error() string {
	if this.Err == nil {
		return fmt.Sprintf("%d %s", this.Code, http.StatusText(this.Code))
	}
	return fmt.Sprintf("%d %v", this.Code, this.Err)
}

func (this Error) Unwrap() error {
	return this.Err
}

func ErrorCode(err error, def int) int {
	var e Error
	if errors.As(err, &e) {
		return e.Code
	}
	return def
}
