package test_test

import (
	"io"
	"testing"

	"github.com/Aize-Public/forego/test"
)

func TestErrors(t *testing.T) {
	test.RunOk(t, "nil", func(t test.T) {
		var err error
		test.NoError(t, err)
	})
	test.RunFail(t, "EOF", func(t test.T) {
		err := io.EOF
		test.NoError(t, err)
	})
}
