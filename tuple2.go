package fuego

// Tuple2 is a tuple with 2 elements.
type Tuple2 struct {
	E1 Entry
	E2 Entry
}

// Hash returns the hash of this tuple.
func (t Tuple2) Hash() uint32 {
	var tHash1, tHash2 uint32
	if t.E1 != nil {
		tHash1 = t.E1.Hash()
	}
	if t.E2 != nil {
		tHash2 = t.E2.Hash()
	}

	result := uint32(1)
	result = 31*result + tHash1
	result = 31*result + tHash2
	return result
}

// Equal returns true if 'o' and 't' are equal.
func (t Tuple2) Equal(o Entry) bool {
	if oT, ok := o.(Tuple2); ok {
		return EntriesEqual(t.E1, oT.E1) &&
			EntriesEqual(t.E2, oT.E2)
	}
	return false
}

// Arity is the number of elements in this tuple.
func (t Tuple2) Arity() int {
	return 2
}

// Map applies the supplied mapper to all elements of this Tuple
// and returns a new Tuple.
func (t Tuple2) Map(mapper Function) Tuple2 {
	return Tuple2{
		E1: mapper(t.E1),
		E2: mapper(t.E2),
	}
}

// MapMulti applies the supplied mappers one for each element
// of this Tuple and returns a new Tuple.
func (t Tuple2) MapMulti(mapper1 Function, mapper2 Function) Tuple2 {
	return Tuple2{
		E1: mapper1(t.E1),
		E2: mapper2(t.E2),
	}
}

// ToSlice returns the elements of this tuple as a Go slice.
func (t Tuple2) ToSlice() EntrySlice {
	return EntrySlice{
		t.E1,
		t.E2,
	}
}
