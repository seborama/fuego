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
	oT, ok := o.(Tuple2)
	return t == oT ||
		(ok &&
			(t.E1 != nil && t.E1.Equal(oT.E1)) &&
			(t.E2 != nil && t.E2.Equal(oT.E2)))
}

// Arity is the number of elements in this tuple.
func (t Tuple2) Arity() int {
	return 2
}

// ToSlice returns the elements of this tuple as a Go slice.
func (t Tuple2) ToSlice() EntrySlice {
	return EntrySlice{
		t.E1,
		t.E2,
	}
}
