package fuego

import (
	"github.com/raviqqe/hamt"
)

// A Map is a Map
type Map struct {
	myMap hamt.Map
}

// NewMap creates a new Map
func NewMap() Map {
	return Map{
		myMap: hamt.NewMap(),
	}
}

// Stream returns a sequential Stream with this collection as its source.
func (m Map) Stream() Stream {
	return NewStream(NewMapIterator(m))
}

// Insert inserts a value into a set.
func (m Map) Insert(k hamt.Entry, v interface{}) Map {
	return Map{
		myMap: m.myMap.Insert(k, v),
	}
}

// Delete deletes a value from a set.
func (m Map) Delete(k hamt.Entry) Map {
	return Map{
		myMap: m.myMap.Delete(k),
	}
}

// Size of the Set.
func (m Map) Size() int {
	return m.myMap.Size()
}

// FirstRest returns a value in a set and a rest of the set.
// This method is useful for iteration.
func (m Map) FirstRest() (hamt.Entry, Set) {
	e, s2 := s.set.FirstRest()
	return e, Map{myMap: s2}
}

// Merge merges 2 maps into one.
func (m Map) Merge(n Map) Map {
}

// Find finds a value corresponding to a given key from a map.
// It returns nil if no value is found.
func (m Map) Find(k hamt.Entry) interface{} {
	e := m.myMap.Find(k)

	if e == nil {
		return nil
	}

	return e.(keyValue).value
}

// FindKey finds a value corresponding to a given key from a map.
// It returns nil if no value is found.
func (m Map) FindKey(interface{}) hamt.Entry {
}

// FindValue finds a value corresponding to a given key from a map.
// It returns nil if no value is found.
func (m Map) FindValue(k hamt.Entry) interface{} {
	// TODO implement.
	// NOTE this would have to return a Set of values, but Set only accepts Entry so the values will have to be wrapped
}

// Has returns true if a key-value pair corresponding with a given key is
// included in a map, or false otherwise.
func (m Map) Has(k hamt.Entry) bool {
	return m.Find(k) != nil // TODO build from HasKey?
}

// Has returns true if a given key is
// included in a map, or false otherwise.
func (m Map) HasKey(k hamt.Entry) bool {
}

// Has returns true if a given value is
// included in a map, or false otherwise.
func (m Map) HasValue(k hamt.Entry) bool {
}

// Values returns the values of thim Map
func (m Map) Values() []hamt.Entry {
	values := []hamt.Entry{}

	subSet := m.myMap
	for subSet.Size() != 0 {
		var e hamt.Entry
		e, subSet = subSet.FirstRest()
		values = append(values, e)
	}

	return values
}
