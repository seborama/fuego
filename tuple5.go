package fuego

import "github.com/raviqqe/hamt"

type Tuple5 struct {
	E1 hamt.Entry
	E2 hamt.Entry
	E3 hamt.Entry
	E4 hamt.Entry
	E5 hamt.Entry
}

func (t Tuple5) Hash() uint32 {
	result := uint32(1)
	result = 31*result + t.E1.Hash()
	result = 31*result + t.E2.Hash()
	result = 31*result + t.E3.Hash()
	result = 31*result + t.E4.Hash()
	result = 31*result + t.E5.Hash()
	return result
}

func (t Tuple5) Equal(o Tuple) bool {
	oT, ok := o.(Tuple5)
	return &t == o ||
		(ok &&
			t.E1.Equal(oT.E1) &&
			t.E1.Equal(oT.E2) &&
			t.E1.Equal(oT.E3) &&
			t.E1.Equal(oT.E4) &&
			t.E1.Equal(oT.E5))
}

func (t Tuple5) Arity() int {
	return 5
}

func (t Tuple5) ToSet() Set {
	return NewOrderedSet().
		Insert(t.E1).
		Insert(t.E2).
		Insert(t.E3).
		Insert(t.E4).
		Insert(t.E5)
}
