package fuego

import "math"

// EntryFloat is an Entry for 'float32'.
type EntryFloat float32

// Hash returns a hash for 'f'.
func (f EntryFloat) Hash() uint32 {
	return math.Float32bits(float32(f))
}

// Equal returns true if 'e' and 'f' are equal.
func (f EntryFloat) Equal(e Entry) bool {
	return f == e
}
