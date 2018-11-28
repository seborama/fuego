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

// Equal returns true when me and e are equal.
// Note that MapEntry defines equality as equality in Hash and
// the Hash of a MapEntry is its MapEntry.K hash (MapEntry.V is
// not considered)
func (me MapEntry) Equal(e hamt.Entry) bool {
	return me.Hash() == e.Hash()
}

// EqualMapEntry compares the key of this MapEntry with the
// given MapEntry.
// Note that MapEntry defines equality as equality in Hash and
// the Hash of a MapEntry is its MapEntry.K hash (MapEntry.V is
// not considered)
func (me MapEntry) EqualMapEntry(ome MapEntry) bool {
	return me.Hash() == ome.Hash()
}

// DeepEqual compares the key and value of this MapEntry with those
// of the given MapEntry.
func (me MapEntry) DeepEqual(o MapEntry) bool {
	return me.K.Equal(o.K) && me.V == o.V
}
