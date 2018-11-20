package fuego

type Tuple0 struct{}

func (t Tuple0) Hash() uint32 {
	return 1
}

func (t Tuple0) Equal(o Tuple) bool {
	// Tuple0 is considered to meet equality when o and t are the same object (in memory)
	// However, we pass objects by value, not reference, hence they can never be the same.
	return false
}

func (t Tuple0) Arity() int {
	return 0
}

func (t Tuple0) ToSet() Set {
	return NewOrderedSet()
}
