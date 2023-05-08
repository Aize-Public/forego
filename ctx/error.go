package ctx

import (
	"errors"
	"fmt"

	"github.com/Aize-Public/forego/utils"
)

func Errorf(c C, f string, args ...any) error {
	return maybeWrap(c, fmt.Errorf(f, args...))
}

func Error(c C, err error) error {
	return maybeWrap(c, err)
}

// a generic error which contains the stack trace
type Err struct {
	error
	Stack []string
	C     C
}

func (err Err) Unwrap() error {
	return err.error
}

func (this Err) Is(err error) bool {
	switch err.(type) {
	case *Err, Err:
		return true
	default:
		return errors.Is(this.error, err)
	}
}

func maybeWrap(c C, err error) error {
	if errors.Is(err, Err{}) {
		return err // already wrapped
	}
	return Err{
		error: err,
		Stack: utils.Stack(2, 100),
		C:     c,
	}
}

// just a []byte, but marshal and unmarshal like json.RawMessage and it is printed as string in logs, win win
type JSON []byte

func (this JSON) MarshalJSON() ([]byte, error) {
	return this, nil
}
func (this *JSON) UnmarshalJSON(j []byte) error {
	*this = j
	return nil
}
func (this JSON) String() string {
	return string(this)
}
