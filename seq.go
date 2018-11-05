package fuego

import "github.com/raviqqe/hamt"

// A Seq is a sequence
type Seq struct {
	seq []hamt.Entry
}

// NewSeq creates a new Seq
func NewSeq() Seq {
	return Seq{}
}

// Size returns the number of elements in the Seq
func (s Seq) Size() int {
	return len(s.seq)
}

// Append an element
func (s Seq) Append(e hamt.Entry) Seq {
	return Seq{
		seq: append(s.seq, e),
	}
}

// find an element's index in the Seq's slice
func (s Seq) find(e hamt.Entry) int {
	for p, v := range s.seq {
		if v.Equal(e) {
			return p
		}
	}
	return -1
}

// Remove the first occurrence of the given element
func (s Seq) Remove(e hamt.Entry) Seq {
	i := s.find(e)
	if i < 0 {
		cpy := make([]hamt.Entry, len(s.seq))
		copy(cpy, s.seq)
		return Seq{seq: cpy}
	}

	return Seq{
		seq: append(s.seq[:i], s.seq[i+1:]...),
	}
}

// AsGo returns the values of the Seq in a Go slice.
// The Go slice is mutable.
// The Go slice and Seq are disconnected: changes to one won't affect the other.
func (s Seq) AsGo() []hamt.Entry {
	cpy := make([]hamt.Entry, len(s.seq))
	copy(cpy, s.seq)
	return cpy
}
