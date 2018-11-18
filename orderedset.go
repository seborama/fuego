package fuego

import (
	"github.com/raviqqe/hamt"
)

// An OrderedSet is an ordered set
type OrderedSet struct {
	slice []hamt.Entry
}

// NewOrderedSet creates a new OrderedSet
func NewOrderedSet() OrderedSet {
	return OrderedSet{
		slice: []hamt.Entry{},
	}
}

// Stream returns a sequential Stream with this collection as its source.
func (s OrderedSet) Stream() Stream {
	return NewStream(
		NewEntrySliceIterator(s.slice))
}

// Insert a value into a set.
func (s OrderedSet) Insert(e hamt.Entry) Set {
	for _, entry := range s.slice {
		if e.Equal(entry) {
			return s
		}
	}
	return OrderedSet{
		slice: append(s.slice, e),
	}
}

// Delete a value from a set.
func (s OrderedSet) Delete(e hamt.Entry) Set {
	for idx, val := range s.slice {
		if val.Equal(e) {
			var slice []hamt.Entry
			if idx == 0 {
				slice = s.slice[1:]
			} else if idx == s.Size()-1 {
				slice = s.slice[:idx]
			} else {
				slice = append(s.slice[:idx], s.slice[idx+1:]...)
			}
			return OrderedSet{
				slice: slice,
			}
		}
	}
	// 'e' not found (includes the case where s.slice is empty)
	return s
}

// Size of the OrderedSet.
func (s OrderedSet) Size() int {
	return len(s.slice)
}

// FirstRest returns a value in a set and a rest of the set.
// This method is useful for iteration.
func (s OrderedSet) FirstRest() (hamt.Entry, Set) {
	return s.slice[0], OrderedSet{slice: s.slice[1:]}
}

// Merge merges 2 sets into one.
func (s OrderedSet) Merge(t Set) Set {
	merge := s
	for _, entry := range t.(OrderedSet).slice {
		merge = merge.Insert(entry).(OrderedSet)
	}
	return OrderedSet{
		slice: merge.slice,
	}
}
