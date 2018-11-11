package fuego

// Iterable performs the given action for each element
// until all elements have been processed or the action panics.
type Iterable interface {
	Iterator() Iterator
}

// BaseIterable is a basic Iterable implementation.
type BaseIterable struct {
	iterator Iterator
}

// NewBaseIterable creates a new BaseIterable.
func NewBaseIterable(i Iterator) Iterable {
	return BaseIterable{iterator: i}
}

// Iterator returns this Iterable's Iterator.
func (si BaseIterable) Iterator() Iterator {
	return si.iterator
}

// An Iterator over a collection.
type Iterator interface {
	Forward() Iterator
	Value() interface{}
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
func (si SetIterator) Forward() Iterator {
	_, si2 := si.set.FirstRest()
	// TODO: create a Iterator wrapper that manages state
	if si2.Size() == 0 {
		return nil
	}

	return SetIterator{
		set: si2,
	}
}

// Value returns the element of the collection currently pointed to by the Iterator.
func (si SetIterator) Value() interface{} {
	e, _ := si.set.FirstRest()
	return e
}

// SliceIterator is an Iterator over a slice.
type SliceIterator struct {
	slice []interface{}
	size  int
}

// NewSliceIterator creates a new SliceIterator.
func NewSliceIterator(s []interface{}) Iterator {
	return SliceIterator{
		slice: s,
		size:  len(s),
	}
}

// Forward to the next element in the collection.
func (si SliceIterator) Forward() Iterator {
	if si.size <= 1 {
		return nil
	}

	return SliceIterator{
		slice: si.slice[1:],
		size:  si.size - 1,
	}
}

// Value returns the element of the collection currently pointed to by the Iterator.
func (si SliceIterator) Value() interface{} {
	return si.slice[0]
}
