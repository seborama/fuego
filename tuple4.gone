package fuego

import "github.com/raviqqe/hamt"

type Tuple4 struct {
	E1 hamt.Entry
	E2 hamt.Entry
	E3 hamt.Entry
	E4 hamt.Entry
}

func (t Tuple4) Hash() uint32 {
	result := uint32(1)
	result = 31*result + t.E1.Hash()
	result = 31*result + t.E2.Hash()
	result = 31*result + t.E3.Hash()
	result = 31*result + t.E4.Hash()
	return result
}

func (t Tuple4) Equal(o Tuple) bool {
	oT, ok := o.(Tuple4)
	return &t == o ||
		(ok &&
			t.E1.Equal(oT.E1) &&
			t.E1.Equal(oT.E2) &&
			t.E1.Equal(oT.E3) &&
			t.E1.Equal(oT.E4))
}

func (t Tuple4) Arity() int {
	return 4
}

func (t Tuple4) ToSet() Set {
	return NewOrderedSet().
		Insert(t.E1).
		Insert(t.E2).
		Insert(t.E3).
		Insert(t.E4)
}
