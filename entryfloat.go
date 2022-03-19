package fuego

import "math"

var _ Entry[EntryFloat] = EntryFloat(0)

// EntryFloat is an Entry for 'float32'.
type EntryFloat float32

// Hash returns a hash for 'f'.
func (f EntryFloat) Hash() uint32 {
	return math.Float32bits(float32(f))
}

// Equal returns true if 'e' and 'f' are equal.
func (f EntryFloat) Equal(e EntryFloat) bool {
	return f == e
}
