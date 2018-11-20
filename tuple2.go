package fuego

import "github.com/raviqqe/hamt"

type Tuple2 struct {
	E1 hamt.Entry
	E2 hamt.Entry
}

func (t Tuple2) Hash() uint32 {
	result := uint32(1)
	result = 31*result + t.E1.Hash()
	result = 31*result + t.E2.Hash()
	return result
}

func (t Tuple2) Equal(o Tuple) bool {
	oT, ok := o.(Tuple2)
	return &t == o ||
		(ok &&
			t.E1.Equal(oT.E1) &&
			t.E1.Equal(oT.E2))
}

func (t Tuple2) Arity() int {
	return 2
}

func (t Tuple2) ToSet() Set {
	return NewOrderedSet().
		Insert(t.E1).
		Insert(t.E2)
}
