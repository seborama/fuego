package fuego

import (
	"encoding/binary"
	"hash/fnv"

	"github.com/raviqqe/hamt"
)

// A Seq is a sequence
type Seq struct {
	seq []hamt.Entry
}

// NewSeq creates a new Seq
func NewSeq() Seq {
	return Seq{
		seq: []hamt.Entry{},
	}
}

// Head returns the first element or panics if none exists
func (s Seq) Head() hamt.Entry {
	return s.seq[0]
}

// HashCode returns the hash of the Seq
func (s Seq) HashCode() uint32 {
	h := fnv.New32a()
	for _, e := range s.seq {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, e.Hash())
		h.Write(bs)
	}
	return h.Sum32()
}

// Last returns the last element or panics if none exists
func (s Seq) Last() hamt.Entry {
	return s.seq[len(s.seq)-1]
}

// Length returns the number of elements in the Seq
func (s Seq) Length() int {
	return len(s.seq)
}

// Tail returns a new Seq made of all but the first element of the Seq or
// panics if no element exists.
func (s Seq) Tail() Traversable {
	return Seq{
		seq: s.seq[1:],
	}
}

// Get returns the first element or panics if none exists
func (s Seq) Get() hamt.Entry {
	return s.Head()
}

// IsEmpty checks if the Seq is empty
func (s Seq) IsEmpty() bool {
	return s.Length() == 0
}

// NonEmpty checks if the Seq is not empty
func (s Seq) NonEmpty() bool {
	return !s.IsEmpty()
}

// Size computes the number of elements of the Seq
func (s Seq) Size() int {
	return s.Length()
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
