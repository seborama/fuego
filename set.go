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

// Map returns a Set consisting of the results of applying the given function to the elements of this Set
func (s Set) Map(f func(hamt.Entry) hamt.Entry) Set {
	newSet := hamt.NewSet()

	subSet := s.set
	for subSet.Size() != 0 {
		var e hamt.Entry
		e, subSet = subSet.FirstRest()
		newSet = newSet.Insert(f(e))
	}

	return Set{
		set: newSet,
	}
}

// MapC returns a slice of channels of interface{} consisting of the results of applying the given function to the elements of this Set
func (s Set) MapC(f func(hamt.Entry) interface{}) []chan interface{} {
	stream := make(chan interface{})

	go func() {
		subSet := s.set
		for subSet.Size() != 0 {
			var e hamt.Entry
			e, subSet = subSet.FirstRest()
			stream <- f(e)
		}
		close(stream)
	}()

	return []chan interface{}{
		stream,
	}
}

// Insert inserts a value into a set.
func (s Set) Insert(e hamt.Entry) Set {
	return Set{
		set: s.set.Insert(e),
	}
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
