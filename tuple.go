package fuego

import "github.com/raviqqe/hamt"

type Tuple interface {
	Hash() bool
	Equal(o interface{}) uint32
	Arity() int
	ToSet() Set
}

type Tuple0 struct{}

func (t Tuple0) Hash() uint32 {
	return 1
}

func (t Tuple0) Equal(o interface{}) bool {
	return &t == o
}

func (t Tuple0) Arity() int {
	return 0
}

func (t Tuple0) ToSet() Set {
	return NewOrderedSet()
}

type Tuple1 struct {
	E1 hamt.Entry
}

func (t Tuple1) Hash() uint32 {
	return t.E1.Hash()
}

func (t Tuple1) Equal(o interface{}) bool {
	oT, ok := o.(Tuple1)
	return &t == o || (ok && t.E1.Equal(oT.E1))
}

func (t Tuple1) Arity() int {
	return 1
}

func (t Tuple1) ToSet() Set {
	return NewOrderedSet().
		Insert(t.E1)
}

type Tuple2 struct {
	E1 hamt.Entry
	E2 hamt.Entry
}

type Tuple3 struct {
	E1 hamt.Entry
	E2 hamt.Entry
	E3 hamt.Entry
}

type Tuple4 struct {
	E1 hamt.Entry
	E2 hamt.Entry
	E3 hamt.Entry
	E4 hamt.Entry
}

type Tuple5 struct {
	E1 hamt.Entry
	E2 hamt.Entry
	E3 hamt.Entry
	E4 hamt.Entry
	E5 hamt.Entry
}

type Tuple6 struct {
	E1 hamt.Entry
	E2 hamt.Entry
	E3 hamt.Entry
	E4 hamt.Entry
	E5 hamt.Entry
	E6 hamt.Entry
}

type Tuple7 struct {
	E1 hamt.Entry
	E2 hamt.Entry
	E3 hamt.Entry
	E4 hamt.Entry
	E5 hamt.Entry
	E6 hamt.Entry
	E7 hamt.Entry
}

type Tuple8 struct {
	E1 hamt.Entry
	E2 hamt.Entry
	E3 hamt.Entry
	E4 hamt.Entry
	E5 hamt.Entry
	E6 hamt.Entry
	E7 hamt.Entry
	E8 hamt.Entry
}

type Tuple9 struct {
	E1 hamt.Entry
	E2 hamt.Entry
	E3 hamt.Entry
	E4 hamt.Entry
	E5 hamt.Entry
	E6 hamt.Entry
	E7 hamt.Entry
	E8 hamt.Entry
	E9 hamt.Entry
}
