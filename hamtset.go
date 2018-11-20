package fuego

import (
	"github.com/raviqqe/hamt"
)

// A HamtSet is an unnaturally ordered set
type HamtSet struct {
	set hamt.Set
}

// NewHamtSet creates a new HamtSet
func NewHamtSet() HamtSet {
	return HamtSet{
		set: hamt.NewSet(),
	}
}

// Stream returns a sequential Stream with this collection as its source.
func (s HamtSet) Stream() Stream {
	return NewStream(NewSetIterator(s))
}

// Insert inserts a value into a set.
func (s HamtSet) Insert(e hamt.Entry) Set {
	return HamtSet{
		set: s.set.Insert(e),
	}
}

// Delete deletes a value from a set.
func (s HamtSet) Delete(e hamt.Entry) Set {
	return HamtSet{
		set: s.set.Delete(e),
	}
}

// Size of the HamtSet.
func (s HamtSet) Size() int {
	return s.set.Size()
}

// FirstRest returns a value in a set and a rest of the set.
// This method is useful for iteration.
func (s HamtSet) FirstRest() (hamt.Entry, Set) {
	e, s2 := s.set.FirstRest()
	return e, HamtSet{set: s2}
}

// Merge merges 2 sets into one.
func (s HamtSet) Merge(t Set) Set {
	return HamtSet{
		set: s.set.Merge((t.(HamtSet).set)),
	}
}
