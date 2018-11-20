package fuego

import "github.com/raviqqe/hamt"

type Tuple3 struct {
	E1 hamt.Entry
	E2 hamt.Entry
	E3 hamt.Entry
}

func (t Tuple3) Hash() uint32 {
	result := uint32(1)
	result = 31*result + t.E1.Hash()
	result = 31*result + t.E2.Hash()
	result = 31*result + t.E3.Hash()
	return result
}

func (t Tuple3) Equal(o Tuple) bool {
	oT, ok := o.(Tuple3)
	return &t == o ||
		(ok &&
			t.E1.Equal(oT.E1) &&
			t.E1.Equal(oT.E2) &&
			t.E1.Equal(oT.E2))
}

func (t Tuple3) Arity() int {
	return 3
}

func (t Tuple3) ToSet() Set {
	return NewOrderedSet().
		Insert(t.E1).
		Insert(t.E2).
		Insert(t.E3)
}
