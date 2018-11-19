package fuego

import (
	"github.com/raviqqe/hamt"
)

// PanicNoSuchElement signifies that the iterator does not have
// the requested element.
const PanicNoSuchElement = "No such element"

// An Iterator over a collection.
type Iterator interface {
	Forward() Iterator
	Value() interface{}
	Reverse() Iterator
	Size() int
}

// SetIterator is an Iterator over a Set.
type SetIterator struct {
	set Set
}

// NewSetIterator creates a new NewSetIterator.
func NewSetIterator(s Set) Iterator {
	return SetIterator{set: s}
}

// Forward to the next element in the collection.
// This consumes an element, unlike Reverse().
func (si SetIterator) Forward() Iterator {
	_, si2 := si.set.FirstRest()
	if si2.Size() == 0 {
		return nil
	}

	return SetIterator{
		set: si2,
	}
}

// Reverse this iterator - useful for RightReduce amongst other things.
// This does not consume any element unlike Forward().
// IMPORTANT NOTE: currently, this function uses OrderedSet for the reverse!
func (si SetIterator) Reverse() Iterator {
	values := []hamt.Entry{}

	subSet := si.set
	for subSet.Size() != 0 {
		var e hamt.Entry
		e, subSet = subSet.FirstRest()
		values = append(values, e)
	}

	reverse := NewOrderedSet()
	for i := len(values) - 1; i >= 0; i-- {
		reverse = reverse.Insert(values[i]).(OrderedSet)
	}

	return SetIterator{
		set: reverse,
	}
}

// Value returns the element of the collection currently pointed
// to by the Iterator.
func (si SetIterator) Value() interface{} {
	if si.Size() == 0 {
		panic(PanicNoSuchElement)
	}
	e, _ := si.set.FirstRest()
	return e
}

// Size returns the total number of elements in the iterator.
func (si SetIterator) Size() int {
	return si.set.Size()
}

// SliceIterator is an Iterator over a slice.
type SliceIterator struct {
	slice []interface{}
}

// NewSliceIterator creates a new SliceIterator.
func NewSliceIterator(s []interface{}) Iterator {
	return SliceIterator{
		slice: s,
	}
}

// Forward to the next element in the collection.
// This consumes an element, unlike Reverse().
func (si SliceIterator) Forward() Iterator {
	if si.Size() <= 1 {
		return nil
	}

	return SliceIterator{
		slice: si.slice[1:],
	}
}

// Reverse this iterator - useful for RightReduce amongst other things.
// This does not consume any element unlike Forward().
func (si SliceIterator) Reverse() Iterator {
	reverse := make([]interface{}, len(si.slice))
	copy(reverse, si.slice)

	for left, right := 0, len(reverse)-1; left < right; left, right = left+1, right-1 {
		reverse[left], reverse[right] = reverse[right], reverse[left]
	}

	return SliceIterator{
		slice: reverse,
	}
}

// Value returns the element of the collection currently pointed
// to by the Iterator.
func (si SliceIterator) Value() interface{} {
	if si.slice == nil || len(si.slice) == 0 {
		panic(PanicNoSuchElement)
	}
	return si.slice[0]
}

// Size returns the total number of elements in the iterator.
func (si SliceIterator) Size() int {
	return len(si.slice)
}

// EntrySliceIterator is an Iterator over a slice.
type EntrySliceIterator struct {
	slice []hamt.Entry
}

// NewEntrySliceIterator creates a new EntrySliceIterator.
func NewEntrySliceIterator(s []hamt.Entry) Iterator {
	return EntrySliceIterator{
		slice: s,
	}
}

// Forward to the next element in the collection.
// This consumes an element, unlike Reverse().
func (si EntrySliceIterator) Forward() Iterator {
	if si.Size() <= 1 {
		return nil
	}

	return EntrySliceIterator{
		slice: si.slice[1:],
	}
}

// Reverse this iterator - useful for RightReduce amongst other things.
// This does not consume any element unlike Forward().
func (si EntrySliceIterator) Reverse() Iterator {
	reverse := make([]hamt.Entry, len(si.slice))
	copy(reverse, si.slice)

	for left, right := 0, len(reverse)-1; left < right; left, right = left+1, right-1 {
		reverse[left], reverse[right] = reverse[right], reverse[left]
	}

	return EntrySliceIterator{
		slice: reverse,
	}
}

// Value returns the element of the collection currently pointed
// to by the Iterator.
func (si EntrySliceIterator) Value() interface{} {
	if si.slice == nil || len(si.slice) == 0 {
		panic(PanicNoSuchElement)
	}
	return si.slice[0]
}

// Size returns the total number of elements in the iterator.
func (si EntrySliceIterator) Size() int {
	return len(si.slice)
}
