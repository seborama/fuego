package fuego

// Tuple3 is a tuple with 3 elements.
type Tuple3 struct {
	E1 Entry
	E2 Entry
	E3 Entry
}

// Hash returns the hash of this tuple.
func (t Tuple3) Hash() uint32 {
	var tHash [3]uint32
	if t.E1 != nil {
		tHash[0] = t.E1.Hash()
	}
	if t.E2 != nil {
		tHash[1] = t.E2.Hash()
	}
	if t.E3 != nil {
		tHash[2] = t.E3.Hash()
	}

	result := uint32(1)
	for i := range tHash {
		result = 31*result + tHash[i]
	}
	return result
}

// Equal returns true if 'o' and 't' are equal.
func (t Tuple3) Equal(o Entry) bool {
	if oT, ok := o.(Tuple3); ok {
		return EntriesEqual(t.E1, oT.E1) &&
			EntriesEqual(t.E2, oT.E2) &&
			EntriesEqual(t.E3, oT.E3)
	}
	return false
}

// Arity is the number of elements in this tuple.
func (t Tuple3) Arity() int {
	return 3
}

// Map applies the supplied mapper to all elements of this Tuple
// and returns a new Tuple.
func (t Tuple3) Map(mapper Function) Tuple3 {
	return Tuple3{
		E1: mapper(t.E1),
		E2: mapper(t.E2),
		E3: mapper(t.E3),
	}
}

// ToSlice returns the elements of this tuple as a Go slice.
func (t Tuple3) ToSlice() EntrySlice {
	return EntrySlice{
		t.E1,
		t.E2,
		t.E3,
	}
}
