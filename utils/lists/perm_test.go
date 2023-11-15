package lists_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils/lists"
)

func TestPerm(t *testing.T) {
	list := strings.Fields("a b c d e")
	perm := lists.Permute(list)
	out := map[string]int{}
	ct := 0
	for perm() {
		x := fmt.Sprintf("%+v", list)
		t.Logf("%+v", x)
		out[x]++
		if out[x] != 1 {
			test.Fail(t, "dup for %v", x)
		}
		ct++
	}
	test.EqualsGo(t, 5*4*3*2*1, ct) // n! permutations

	/*
		in := []int{1, 2, 3}
		ct := map[string]int{}
		tot := 0
		p := lists.Permute[int]{Slice: in}
		for p.Next() {
			tot++
			k := fmt.Sprintf("%+v", p.Slice)
			t.Logf("%d: %v", tot, k)
			if ct[k] == 1 {
				t.Fatalf("dup: %v", k)
			}
			ct[k]++
		}
	*/
}
