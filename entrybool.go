package fuego

import "github.com/raviqqe/hamt"

// EntryBool is an Entry for 'bool'.
type EntryBool bool

// Hash returns a hash for 'i'.
func (i EntryBool) Hash() uint32 {
	if bool(i) {
		return 1
	}
	return 0
}

// Equal returns true if 'e' and 'i' are equal.
func (i EntryBool) Equal(e hamt.Entry) bool {
	j, ok := e.(EntryBool)

	if !ok {
		return false
	}

	return i == j
}
