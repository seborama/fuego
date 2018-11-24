package fuego

import "github.com/raviqqe/hamt"

// EntryNone is an Entry for no value.
// TODO: confirm this pattern - Should use MaybeNone()
type EntryNone struct{}

// Hash returns a hash for 'i'.
func (i EntryNone) Hash() uint32 {
	return 0
}

// Equal returns true if 'e' and 'i' are equal.
func (i EntryNone) Equal(e hamt.Entry) bool {
	j, ok := e.(EntryNone)
	if !ok {
		return false
	}

	return i == j
}
