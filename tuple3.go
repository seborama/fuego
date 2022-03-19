package fuego

// Tuple3 is a tuple with 3 elements.
type Tuple3[E Entry[E], F Entry[F], G Entry[G]] struct {
	E1 E
	E2 F
	E3 G
}

// Hash returns the hash of this tuple.
// The result is unpredictable when at least one of the values is nil.
func (t Tuple3[E, F, G]) Hash() uint32 {
	var tHash [3]uint32

	tHash[0] = t.E1.Hash()
	tHash[1] = t.E2.Hash()
	tHash[2] = t.E3.Hash()

	result := uint32(1)

	for i := range tHash {
		result = 31*result + tHash[i]
	}

	return result
}

// Equal returns true if 'o' and 't' are equal.
// The result is unpredictable when at least one of the values is nil.
func (t Tuple3[E, F, G]) Equal(o Tuple3[E, F, G]) bool {
	return t.Hash() == o.Hash() // TODO: return t.E1.Equal(o.E1) && t.E2.Equal(o.E2) && t.E3.Equal(o.E3)
}

// Arity is the number of elements in this tuple.
func (t Tuple3[E, F, G]) Arity() int {
	return 3
}

// Map applies the supplied mapper to all elements of this Tuple
// and returns a new Tuple.
func (t Tuple3[E, F, G]) Map(mapperE Function[E, E], mapperF Function[F, F], mapperG Function[G, G]) Tuple3[E, F, G] {
	return Tuple3[E, F, G]{
		E1: mapperE(t.E1),
		E2: mapperF(t.E2),
		E3: mapperG(t.E3),
	}
}

// ToESlice returns the elements of this tuple as a Go slice.
func (t Tuple3[E, F, G]) ToESlice(mapperF Function[F, E], mapperG Function[G, E]) EntrySlice[E] {
	return EntrySlice[E]{
		t.E1,
		mapperF(t.E2),
		mapperG(t.E3),
	}
}

// ToFSlice returns the elements of this tuple as a Go slice.
func (t Tuple3[E, F, G]) ToFSlice(mapperE Function[E, F], mapperG Function[G, F]) EntrySlice[F] {
	return EntrySlice[F]{
		mapperE(t.E1),
		t.E2,
		mapperG(t.E3),
	}
}

// ToGSlice returns the elements of this tuple as a Go slice.
func (t Tuple3[E, F, G]) ToGSlice(mapperE Function[E, G], mapperF Function[F, G]) EntrySlice[G] {
	return EntrySlice[G]{
		mapperE(t.E1),
		mapperF(t.E2),
		t.E3,
	}
}
