package fuego

import (
	"github.com/raviqqe/hamt"
)

// A Set is a would-be functional Set
type Set struct {
	set hamt.Set
}

// NewSet creates a new Set
func NewSet() Set {
	return Set{set: hamt.NewSet()}
}

// Insert inserts a value into a set.
func (s Set) Insert(e hamt.Entry) Set {
	return Set{set: s.set.Insert(e)}
}

// Map returns a Set consisting of the results of applying the given function to the elements of this Set
func (s Set) Map(f func(hamt.Entry) hamt.Entry) Set {
	newSet := hamt.NewSet()

	subSet := s.set
	for subSet.Size() != 0 {
		var e hamt.Entry
		e, subSet = subSet.FirstRest()
		newSet = newSet.Insert(f(e))
	}

	return Set{set: newSet}
}

// Values returns the values of this Set in a Seq
func (s Set) Values() Seq {
	newSeq := NewSeq()

	subSet := s.set
	for subSet.Size() != 0 {
		var e hamt.Entry
		e, subSet = subSet.FirstRest()
		newSeq = newSeq.Append(e)
	}

	return newSeq
}
