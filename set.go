package fuego

import (
	"github.com/raviqqe/hamt"
)

// A Set is a Set
type Set struct {
	set hamt.Set
}

// NewSet creates a new Set
func NewSet() Set {
	return Set{
		set: hamt.NewSet(),
	}
}

// Stream returns a sequential Stream with this collection as its source.
func (s Set) Stream() Stream {
	return NewStream(
		NewBaseIterable(NewSetIterator(s)))
}

// Insert inserts a value into a set.
func (s Set) Insert(e hamt.Entry) Set {
	return Set{
		set: s.set.Insert(e),
	}
}

// Size of the Set.
func (s Set) Size() int {
	return s.set.Size()
}

// FirstRest returns a value in a set and a rest of the set.
// This method is useful for iteration.
func (s Set) FirstRest() (hamt.Entry, Set) {
	e, s2 := s.set.FirstRest()
	return e, Set{set: s2}
}

// Values returns the values of this Set
func (s Set) Values() []hamt.Entry {
	values := []hamt.Entry{}

	subSet := s.set
	for subSet.Size() != 0 {
		var e hamt.Entry
		e, subSet = subSet.FirstRest()
		values = append(values, e)
	}

	return values
}
