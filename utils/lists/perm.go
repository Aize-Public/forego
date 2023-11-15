package lists

import (
	"fmt"
	"strings"
)

/*
Creates a permute generator which will permute *inplace*

Note: the first call of the next function won't change anything
this is to make it easier to do loops:

		next := lists.Permute(inplace)
		for next() {
	        use(inplace)
		}

uses Johnson and Trotter algorithm
*/
func Permute[T any](inplace []T) (next func() bool) {
	this := permute(inplace)
	first := true
	return func() bool {
		if first {
			first = false
			return true // allow for 1 extra call before starting
		}
		return this.next()
	}
}

type permElem struct {
	right bool
	orig  int
}

type permGen[T any] struct {
	inplace []T
	elems   []permElem
}

func permute[T any](inplace []T) *permGen[T] {
	this := &permGen[T]{
		inplace: inplace,
		elems:   make([]permElem, len(inplace)),
	}
	for i := range inplace {
		this.elems[i] = permElem{
			orig: i,
		}
	}
	return this
}

func (this *permGen[T]) next() bool {
	best_orig := -1
	best := -1

	for i, el := range this.elems {
		// skip the elements that cannot move
		if el.right {
			if i == len(this.elems)-1 {
				continue
			}
			if this.elems[i+1].orig > this.elems[i].orig {
				continue
			}
		} else {
			if i == 0 {
				continue
			}
			if this.elems[i-1].orig > this.elems[i].orig {
				continue
			}
		}

		// remember the best one
		if el.orig > best_orig {
			best_orig = el.orig
			best = i
		}
	}

	// no moves? we are done
	if best < 0 {
		return false
	}

	// swap the best towards it's direction
	if this.elems[best].right {
		this.elems[best], this.elems[best+1] = this.elems[best+1], this.elems[best]
		this.inplace[best], this.inplace[best+1] = this.inplace[best+1], this.inplace[best]
	} else {
		this.elems[best], this.elems[best-1] = this.elems[best-1], this.elems[best]
		this.inplace[best], this.inplace[best-1] = this.inplace[best-1], this.inplace[best]
	}

	// switch direction of anything bigger
	for i, el := range this.elems {
		if el.orig > best_orig {
			this.elems[i].right = !el.right
		}
	}
	return true
}

func (this *permGen[T]) String() string {
	out := []string{}
	for i := 0; i < len(this.elems); i++ {
		if this.elems[i].right {
			out = append(out, fmt.Sprintf("%d>", this.elems[i].orig))
		} else {
			out = append(out, fmt.Sprintf("<%d", this.elems[i].orig))
		}
	}
	return strings.Join(out, " ")
}
