package fuego

const PanicNoSuchElement = "No such element"

// An Iterator over a collection.
type Iterator interface {
	Forward() Iterator
	Value() interface{}
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
func (si SetIterator) Forward() Iterator {
	_, si2 := si.set.FirstRest()
	if si2.Size() == 0 {
		return nil
	}

	return SetIterator{
		set: si2,
	}
}

// Value returns the element of the collection currently pointed to by the Iterator.
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
func (si SliceIterator) Forward() Iterator {
	if si.Size() <= 1 {
		return nil
	}

	return SliceIterator{
		slice: si.slice[1:],
	}
}

// Value returns the element of the collection currently pointed to by the Iterator.
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
