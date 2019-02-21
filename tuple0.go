package fuego

// Tuple0 is a tuple with 0 element.
type Tuple0 struct{}

// Hash returns the hash of this tuple.
func (t Tuple0) Hash() uint32 {
	return 1
}

// Equal returns true if 'o' and 't' are equal.
func (t Tuple0) Equal(o Entry) bool {
	// Tuple0 is considered to meet equality when o and t are the
	// same object (in memory)
	// However, we pass objects by value, not reference, hence they
	// can never be the same.
	return false
}

// Arity is the number of elements in this tuple.
func (t Tuple0) Arity() int {
	return 0
}

// ToSlice returns the elements of this tuple as a Go slice.
func (t Tuple0) ToSlice() EntrySlice {
	return EntrySlice{}
}
