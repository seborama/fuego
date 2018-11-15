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

// BiStream returns a sequential stream with this collection as its source.
// TODO: implement this - it will iterator through the (k,v) pairs of the Map.
// func (m Map) BiStream() BiStream {
// 	return NewStream(NewMapIterator(m))
// }

// EntrySet returns a Set of MapEntry's from the (k, v) pairs contained in this map.
// Since EntrySet returns a Set, it can be streamed with Set.Stream().
func (m Map) EntrySet() Set {
	s := Set{}

	subMap := m.myMap
	for subMap.Size() != 0 {
		var k2 hamt.Entry
		var v2 interface{}
		k2, v2, subMap = subMap.FirstRest()
		s = s.Insert(MapEntry{k2, v2})
	}
	return s
}

// KeySet returns a Set of keys contained in this map.
// Since KeySet returns a Set, it can be streamed with Set.Stream().
// Note that ValueSet() is not implemented because Values can be present
// multiple times. This could possibly be implemented via []interface{}?
// It also could be better to use the BiStream() proposed in this file.
func (m Map) KeySet() Set {
	s := Set{}

	subMap := m.myMap
	for subMap.Size() != 0 {
		var k2 hamt.Entry
		k2, _, subMap = subMap.FirstRest()
		s = s.Insert(k2)
	}
	return s
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
func (m Map) Find(k hamt.Entry) MapEntry {
	v := m.myMap.Find(k)
	if v == nil {
		return MapEntry{nil, nil}
	}

	return MapEntry{K: k, V: v}
}

// FindKey finds a value corresponding to a given key from a map.
// It returns nil if no value is found.
func (m Map) FindKey(k hamt.Entry) interface{} {
	return m.myMap.Find(k)
}

// FindValue finds a value corresponding to a given key from a map.
// It returns nil if no value is found.
// func (m Map) FindValue(k hamt.Entry) interface{} {
// TODO implement.
// NOTE this would have to return a Set of values, but Set only accepts Entry so the values will have to be wrapped
// }

// Has returns true if a key-value pair corresponding with a given key is
// included in a map, or false otherwise.
func (m Map) Has(k hamt.Entry, v interface{}) bool {
	value := m.FindKey(k)
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
		_, v2, subMap = subMap.FirstRest()
		if v2 == v {
			return true
		}
	}

	return false
}
