package fuego

import "github.com/raviqqe/hamt"

type Seq struct {
	seq []hamt.Entry
}

func NewSeq() Seq {
	return Seq{}
}

func (s Seq) Size() int {
	return len(s.seq)
}

func (s Seq) Append(e hamt.Entry) Seq {
	return Seq{
		seq: append(s.seq, e),
	}
}

func (s Seq) AsGo() []hamt.Entry {
	cpy := make([]hamt.Entry, len(s.seq))
	copy(cpy, s.seq)
	return cpy
}
