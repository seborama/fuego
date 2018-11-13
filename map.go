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
// func (m Map) Stream() Stream {
// 	return NewStream(NewMapIterator(m))
// }

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

// FirstRest returns a key-value pair in a map and a rest of the map.
// This method is useful for iteration.
// The key and value would be nil if the map is empty.
func (m Map) FirstRest() (hamt.Entry, interface{}, Map) {
	k, v, m2 := m.myMap.FirstRest()
	return k, v, Map{myMap: m2}
}

// Merge merges 2 maps into one.
func (m Map) Merge(n Map) Map {
	return Map{
		myMap: m.myMap.Merge(n.myMap),
	}
}

// Find finds a value corresponding to a given key from a map.
// It returns nil if no value is found.
func (m Map) Find(k hamt.Entry) interface{} {
	return m.myMap.Find(k)
}

// FindKey finds a value corresponding to a given key from a map.
// It returns nil if no value is found.
// func (m Map) FindKey(interface{}) hamt.Entry {
// }

// FindValue finds a value corresponding to a given key from a map.
// It returns nil if no value is found.
// func (m Map) FindValue(k hamt.Entry) interface{} {
// TODO implement.
// NOTE this would have to return a Set of values, but Set only accepts Entry so the values will have to be wrapped
// }

// Has returns true if a key-value pair corresponding with a given key is
// included in a map, or false otherwise.
func (m Map) Has(k hamt.Entry, v interface{}) bool {
	value := m.Find(k)
	if value != nil && value == v {
		return true
	}
	return false
}

// HasKey returns true if a given key exists
// in a map, or false otherwise.
func (m Map) HasKey(k hamt.Entry) bool {
	return m.myMap.Include(k)
}

// HasValue returns true if a given value exists
// in a map, or false otherwise.
func (m Map) HasValue(v interface{}) bool {
	subMap := m.myMap

	for subMap.Size() != 0 {
		var v2 interface{}
		_, v2, subMap = m.myMap.FirstRest()
		if v2 == v {
			return true
		}
	}

	return false
}

// Values returns the values of thim Map
// func (m Map) Values() []hamt.Entry {
// 	values := []hamt.Entry{}

// 	subSet := m.myMap
// 	for subSet.Size() != 0 {
// 		var e hamt.Entry
// 		k, v, subSet = subSet.FirstRest()
// 		values = append(values, e)
// 	}

// 	return values
// }
