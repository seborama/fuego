package fuego

type Tuple1 struct {
	E1 Entry
}

func (t Tuple1) Hash() uint32 {
	if t.E1 == nil {
		return 0
	}
	return t.E1.Hash()
}

func (t Tuple1) Equal(o Tuple) bool {
	oT, ok := o.(Tuple1)
	return t == o ||
		(ok &&
			(t.E1 != nil && t.E1.Equal(oT.E1)))
}

func (t Tuple1) Arity() int {
	return 1
}

func (t Tuple1) ToSet() Set {
	return NewOrderedSet().
		Insert(t.E1)
}
