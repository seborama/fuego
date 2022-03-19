package fuego

// Tuple1 is a tuple with 1 element.
type Tuple1[E Entry[E]] struct {
	E1 E
}

// Hash returns the hash of this tuple.
// The result is unpredictable when the value is nil.
func (t Tuple1[E]) Hash() uint32 {
	var tHash1 uint32

	tHash1 = t.E1.Hash()

	result := uint32(1)
	result = 31*result + tHash1

	return result
}

// Equal returns true if 'o' and 't' are equal.
// The result is unpredictable when the value is nil.
func (t Tuple1[E]) Equal(o Tuple1[E]) bool {
	return t.Hash() == o.Hash() // TODO: return t.E1.Equal(o.E1) && t.E2.Equal(o.E2)
}

// Arity is the number of elements in this tuple.
func (t Tuple1[E]) Arity() int {
	return 1
}

// Map applies the supplied mapper to the element of this Tuple
// and returns a new Tuple.
func (t Tuple1[E]) Map(mapper Function[E, E]) Tuple1[E] {
	return Tuple1[E]{
		E1: mapper(t.E1),
	}
}

// MapMulti applies the supplied mappers one for each element
// of this Tuple and returns a new Tuple.
func (t Tuple1[E]) MapMulti(mapper1 Function[E, E]) Tuple1[E] {
	return t.Map(mapper1)
}

// ToSlice returns the elements of this tuple as a Go slice.
func (t Tuple1[E]) ToSlice() EntrySlice[E] {
	return EntrySlice[E]{t.E1}
}
