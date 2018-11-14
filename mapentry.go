package fuego

import (
	"github.com/raviqqe/hamt"
)

// MapEntry holds a hamt.Entry key and a value that form an entry in a Map.
type MapEntry struct {
	K hamt.Entry
	V interface{}
}

// NewMapEntry creates a new MapEntry.
func NewMapEntry(k hamt.Entry, v interface{}) MapEntry {
	return MapEntry{
		K: k,
		V: v,
	}
}

// Hash of the MapEntry - will panic if me.k is nil.
func (me MapEntry) Hash() uint32 {
	return me.K.Hash()
}

// Equal compares the key of this MapEntry with with given Entry.
func (me MapEntry) Equal(e hamt.Entry) bool {
	return me.Hash() == e.Hash()
}

// DeepEqual compares the key and value of this MapEntry with those of the given MapEntry.
func (me MapEntry) DeepEqual(o MapEntry) bool {
	return me.K == o.K && me.V == o.V
}
