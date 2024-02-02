package utils

import "errors"

// similar to io.EOF, used to signal the end of data
// normally used in Range() functions to stop looping while still returning nil as error
type EOD struct{}

var _ error = EOD{}

func (err EOD) Error() string { return "EOD" }

// call the given function passing an error, if the error is returned back (wrapped or not) it will be ignored
func ErrorGate(f func(err error) error) error {
	err := f(errGate{})
	if errors.Is(err, errGate{}) {
		return nil
	}
	return err
}

type errGate struct{}

func (err errGate) Error() string { return "Error: should have been trapped by ErrorGate()" }
