package fuego

// Tuple2 is a tuple with 2 elements.
type Tuple2[E Entry[E], F Entry[F]] struct {
	E1 E
	E2 F
}

// Hash returns the hash of this tuple.
// The result is unpredictable when at least one of the values is nil.
func (t Tuple2[E, F]) Hash() uint32 {
	var tHash1, tHash2 uint32

	tHash1 = t.E1.Hash()
	tHash2 = t.E2.Hash()

	result := uint32(1)
	result = 31*result + tHash1
	result = 31*result + tHash2

	return result
}

// Equal returns true if 'o' and 't' are equal.
// The result is unpredictable when at least one of the values is nil.
func (t Tuple2[E, F]) Equal(o Tuple2[E, F]) bool {
	return t.Hash() == o.Hash() // TODO: return t.E1.Equal(o.E1) && t.E2.Equal(o.E2)
}

// Arity is the number of elements in this tuple.
func (t Tuple2[E, F]) Arity() int {
	return 2
}

// Map applies the supplied mapper to all elements of this Tuple
// and returns a new Tuple.
func (t Tuple2[E, F]) Map(mapperE Function[E, E], mapperF Function[F, F]) Tuple2[E, F] {
	return Tuple2[E, F]{
		E1: mapperE(t.E1),
		E2: mapperF(t.E2),
	}
}

// ToESlice returns the elements of this tuple as a Go slice.
func (t Tuple2[E, F]) ToESlice(mapperF Function[F, E]) EntrySlice[E] {
	return EntrySlice[E]{
		t.E1,
		mapperF(t.E2),
	}
}

// ToFSlice returns the elements of this tuple as a Go slice.
func (t Tuple2[E, F]) ToFSlice(mapperE Function[E, F]) EntrySlice[F] {
	return EntrySlice[F]{
		mapperE(t.E1),
		t.E2,
	}
}
