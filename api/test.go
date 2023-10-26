package api

import (
	"testing"

	"github.com/Aize-Public/forego/test"
)

// TODO(oha): should we do marshalling here to be safe?
func Test[T Op](t *testing.T, op T) T {
	t.Helper()
	t.Logf("Test(%+v)...", op)
	err := op.Do(test.Context(t))
	test.NoError(t, err)
	return op
}
