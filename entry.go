package fuego

import "github.com/raviqqe/hamt"

// Entry represents an entry in a collection.
type Entry interface {
	// Hash() uint32
	// Equal(Entry) bool
	hamt.Entry
	// Value() Entry
}
