package fuego

// Tuple1 is a tuple with 1 element.
type Tuple1 struct {
	E1 Entry
}

// Hash returns the hash of this tuple.
func (t Tuple1) Hash() uint32 {
	if t.E1 == nil {
		return 0
	}
	return t.E1.Hash()
}

// Equal returns true if 'o' and 't' are equal.
func (t Tuple1) Equal(o Entry) bool {
	if oT, ok := o.(Tuple1); ok {
		return EntriesEqual(t.E1, oT.E1)
	}
	return false
}

// Arity is the number of elements in this tuple.
func (t Tuple1) Arity() int {
	return 1
}

// Map applies the supplied mapper to the element of this Tuple
// and returns a new Tuple.
func (t Tuple1) Map(mapper Function) Tuple1 {
	return Tuple1{
		E1: mapper(t.E1),
	}
}

// MapMulti applies the supplied mappers one for each element
// of this Tuple and returns a new Tuple.
func (t Tuple1) MapMulti(mapper1 Function) Tuple1 {
	return t.Map(mapper1)
}

// ToSlice returns the elements of this tuple as a Go slice.
func (t Tuple1) ToSlice() EntrySlice {
	return EntrySlice{t.E1}
}
