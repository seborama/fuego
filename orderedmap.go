package fuego

import (
	"github.com/raviqqe/hamt"
)

// A OrderedMap is an ordered map
type OrderedMap struct {
	myMap map[hamt.Entry]interface{}
}

// NewOrderedMap creates a new Map
func NewOrderedMap() OrderedMap {
	return OrderedMap{
		myMap: map[hamt.Entry]interface{}{},
	}
}

// BiStream returns a sequential stream with this collection as its source.
// TODO: implement this - it will iterator through the (k,v) pairs of the Map.
// func (m OrderedMap) BiStream() BiStream {
// 	return NewStream(NewMapIterator(m))
// }

// EntrySet returns a Set of MapEntry's from the (k, v) pairs contained
// in this map.
// Since EntrySet returns a Set, it can be streamed with Set.Stream().
func (m OrderedMap) EntrySet() Set {
	s := NewOrderedSet()
	for k, v := range m.myMap {
		s = s.Insert(MapEntry{k, v}).(OrderedSet)
	}
	return s
}

// KeySet returns a Set of keys contained in this map.
// Since KeySet returns a Set, it can be streamed with Set.Stream().
// Note that ValueSet() is not implemented because Values can be present
// multiple times. This could possibly be implemented via []interface{}?
// It also could be better to use the BiStream() proposed in this file.
func (m OrderedMap) KeySet() Set {
	s := NewOrderedSet()
	for k := range m.myMap {
		s = s.Insert(k).(OrderedSet)
	}
	return s
}

// Insert inserts a value into this map.
func (m OrderedMap) Insert(k hamt.Entry, v interface{}) Map {
	newMap := make(map[hamt.Entry]interface{}, len(m.myMap))
	for k2, v2 := range m.myMap {
		newMap[k2] = v2
	}
	newMap[k] = v
	return OrderedMap{
		myMap: newMap,
	}
}

// Delete deletes a value from this map.
func (m OrderedMap) Delete(k hamt.Entry) Map {
	newMap := make(map[hamt.Entry]interface{}, len(m.myMap))
	for k2, v2 := range m.myMap {
		newMap[k2] = v2
	}
	delete(newMap, k)
	return OrderedMap{
		myMap: newMap,
	}
}

// Size of the Set.
func (m OrderedMap) Size() int {
	return len(m.myMap)
}

// FirstRest returns a key-value pair in a map and a rest of the map.
// This method is useful for iteration.
// The key and value would be nil if the map is empty.
func (m OrderedMap) FirstRest() (hamt.Entry, interface{}, Map) {
	k, v, m2 := m.myMap.FirstRest()
	return k, v, OrderedMap{myMap: m2}
}

// Merge merges 2 maps into one.
func (m OrderedMap) Merge(n Map) Map {
	return OrderedMap{
		myMap: m.myMap.Merge(n.(OrderedMap).myMap),
	}
}

// Find finds a value corresponding to a given key from a map.
// It returns nil if no value is found.
func (m OrderedMap) Find(k hamt.Entry) MapEntry {
	v := m.myMap.Find(k)
	if v == nil {
		return MapEntry{nil, nil}
	}

	return MapEntry{K: k, V: v}
}

// FindKey finds a value corresponding to a given key from a map.
// It returns nil if no value is found.
func (m OrderedMap) FindKey(k hamt.Entry) interface{} {
	return m.myMap.Find(k)
}

// FindValue finds a value corresponding to a given key from a map.
// It returns nil if no value is found.
// func (m OrderedMap) FindValue(k hamt.Entry) interface{} {
// TODO implement.
// NOTE this would have to return a Set of values, but Set only
// accepts Entry so the values will have to be wrapped
// }

// Has returns true if a key-value pair corresponding with a given key is
// included in a map, or false otherwise.
func (m OrderedMap) Has(k hamt.Entry, v interface{}) bool {
	value := m.FindKey(k)
	if value != nil && value == v {
		return true
	}
	return false
}

// HasKey returns true if a given key exists
// in a map, or false otherwise.
func (m OrderedMap) HasKey(k hamt.Entry) bool {
	return m.myMap.Include(k)
}

// HasValue returns true if a given value exists
// in a map, or false otherwise.
func (m OrderedMap) HasValue(v interface{}) bool {
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
