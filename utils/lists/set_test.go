package lists_test

import (
	"testing"

	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils/lists"
)

func TestSet(t *testing.T) {
	s := lists.Set[string]{}
	test.Assert(t, s.Add("foo") == true)
	test.Assert(t, s.Add("foo") == false)
	test.Assert(t, s.Add("bar") == true)
	test.Assert(t, s.Has("bar") == true)
	test.Assert(t, s.Remove("bar") == true)
	test.Assert(t, s.Has("bar") == false)
	test.Assert(t, s.Remove("bar") == false)
	test.Assert(t, s.Has("bar") == false)
	test.Assert(t, s.Add("bar") == true)
	test.Assert(t, s.Has("bar") == true)
	test.Assert(t, len(s.All()) == 2)
}
