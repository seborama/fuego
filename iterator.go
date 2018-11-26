package fuego

// PanicNoSuchElement signifies that the iterator does not have
// the requested element.
const PanicNoSuchElement = "No such element"

// An Iterator over a collection.
type Iterator interface {
	Forward() Iterator
	Value() Entry
	Reverse() Iterator
	Size() int
}

// SetIterator is an Iterator over a Set.
type SetIterator struct {
	set Set
}

// NewSetIterator creates a new NewSetIterator.
func NewSetIterator(s Set) Iterator {
	if s == nil || s.Size() == 0 {
		return nil
	}
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
	if si.Size() == 0 {
		return nil
	}

	values := []Entry{}

	subSet := si.set
	for subSet.Size() != 0 {
		var e Entry
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
func (si SetIterator) Value() Entry {
	if si.Size() == 0 {
		panic(PanicNoSuchElement) // TODO return 'nil' instead?
	}
	e, _ := si.set.FirstRest()
	return e
}

// Size returns the total number of elements in the iterator.
func (si SetIterator) Size() int {
	if si.set == nil {
		return 0
	}
	return si.set.Size()
}

// SliceIterator is an Iterator over a slice.
type SliceIterator struct {
	slice []Entry
}

// NewSliceIterator creates a new SliceIterator.
func NewSliceIterator(s []Entry) Iterator {
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
	reverse := make([]Entry, len(si.slice))
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
func (si SliceIterator) Value() Entry {
	if si.Size() == 0 {
		panic(PanicNoSuchElement) // TODO return 'nil' instead?
	}
	return si.slice[0]
}

// Size returns the total number of elements in the iterator.
func (si SliceIterator) Size() int {
	return len(si.slice)
}
