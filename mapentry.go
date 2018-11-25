package fuego

import "github.com/raviqqe/hamt"

// MapEntry holds a hamt.Entry key and a value that form an entry in a Map.
type MapEntry struct {
	K Entry
	V interface{}
}

// NewMapEntry creates a new MapEntry.
func NewMapEntry(k Entry, v interface{}) MapEntry {
	return MapEntry{
		K: k,
		V: v,
	}
}

// Hash of the MapEntry - will panic if me.k is nil.
func (me MapEntry) Hash() uint32 {
	return me.K.Hash()
}

// Equal compares the key of this MapEntry with given Entry.
func (me MapEntry) Equal(e hamt.Entry) bool {
	return me.Hash() == e.Hash()
}

// Equal compares the key of this MapEntry with the given MapEntry.
func (me MapEntry) EqualMapEntry(ome MapEntry) bool {
	return me.Hash() == ome.Hash()
}

// DeepEqual compares the key and value of this MapEntry with those
// of the given MapEntry.
func (me MapEntry) DeepEqual(o MapEntry) bool {
	return me.K.Hash() == o.K.Hash() && me.V == o.V
}
