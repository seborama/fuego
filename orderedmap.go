package fuego

// A OrderedMap is an ordered map
type OrderedMap struct {
	entries []MapEntry
}

// NewOrderedMap creates a new Map
func NewOrderedMap() OrderedMap {
	return OrderedMap{
		entries: []MapEntry{},
	}
}

// BiStream returns a sequential stream with this collection as its source.
// TODO: implement this - it will iterate through the (k,v) pairs of the Map.
// func (m OrderedMap) BiStream() BiStream {
// 	return NewStream(NewMapIterator(m))
// }

// EntrySet returns a Set of MapEntry's from the (k, v) pairs contained
// in this map.
// Since EntrySet returns a Set, it can be streamed with Set.Stream().
func (m OrderedMap) EntrySet() Set {
	newSet := NewOrderedSet()
	for _, e := range m.entries {
		newSet = newSet.Insert(MapEntry{e.K, e.V}).(OrderedSet)
	}
	return newSet
}

// KeySet returns a Set of keys contained in this map.
// Since KeySet returns a Set, it can be streamed with Set.Stream().
// Note that ValueSet() is not implemented because Values can be present
// multiple times. This could possibly be implemented via []interface{}?
// It also could be better to use the BiStream() proposed in this file.
func (m OrderedMap) KeySet() Set {
	newSet := NewOrderedSet()
	for _, e := range m.entries {
		newSet = newSet.Insert(e.K).(OrderedSet)
	}
	return newSet
}

// Insert a value into this map.
func (m OrderedMap) Insert(k Entry, v interface{}) Map {
	newMap := make([]MapEntry, len(m.entries)+1) // keep room for the '(k,v)' if not already present
	copy(newMap, m.entries)

	foundExisting := false
	for _, e := range m.entries {
		if e.Equal(k) {
			foundExisting = true
			newMap = append(newMap, MapEntry{K: k, V: v})
			break
		}
	}
	if !foundExisting {
		newMap[len(newMap)-1] = MapEntry{K: k, V: v}
	}

	return OrderedMap{
		entries: newMap,
	}
}

// Delete a value from this map.
func (m OrderedMap) Delete(k Entry) Map {
	for idx, e := range m.entries {
		if e.K.Equal(k) {
			var sCopy []MapEntry
			if idx == 0 {
				sCopy = make([]MapEntry, len(m.entries)-1)
				copy(sCopy, m.entries[1:])
			} else if idx == m.Size()-1 {
				sCopy = make([]MapEntry, len(m.entries)-1)
				copy(sCopy, m.entries[:idx])
			} else {
				sCopy = append(m.entries[:idx], m.entries[idx+1:]...)
			}
			return OrderedMap{
				entries: sCopy,
			}
		}
	}

	// 'k' not found (includes the case where s.entries is empty)
	sCopy := make([]MapEntry, len(m.entries))
	copy(sCopy, m.entries)
	return OrderedMap{
		entries: sCopy,
	}
}

// Size of the Set.
func (m OrderedMap) Size() int {
	return len(m.entries)
}

// FirstRest returns a key-value pair in a map and a rest of the map.
// This method is useful for iteration.
// The key and value would be nil if the map is empty.
func (m OrderedMap) FirstRest() (Entry, interface{}, Map) {
	sCopy := make([]MapEntry, len(m.entries)-1)
	copy(sCopy, m.entries[1:])
	return m.entries[0].K, m.entries[0].V, OrderedMap{entries: sCopy}
}

// Merge this map and given map.
func (m OrderedMap) Merge(n Map) Map {
	merge := make([]MapEntry, len(m.entries))
	copy(merge, m.entries)

	for _, entry := range n.(OrderedMap).entries {
		merge = append(merge, entry)
	}
	return OrderedMap{
		entries: merge,
	}
}

type notFound struct{}

// Get a value in this map corresponding to a given key.
// It returns nil if no value is found.
// TODO return Maybe instead of interface{}
func (m OrderedMap) Get(k Entry) interface{} {
	for _, e := range m.entries {
		if e.K.Equal(k) {
			return e.V
		}
	}

	return notFound{}
}

// Has returns true if a key-value pair corresponding with a given key is
// included in a map, or false otherwise.
func (m OrderedMap) Has(k Entry, v interface{}) bool {
	value := m.Get(k)

	if value == (notFound{}) {
		return false
	}
	return value == v
}

// HasKey returns true if a given key exists
// in a map, or false otherwise.
func (m OrderedMap) HasKey(k Entry) bool {
	return m.Get(k) != (notFound{})
}

// HasValue returns true if a given value exists
// in a map, or false otherwise.
func (m OrderedMap) HasValue(v interface{}) bool {
	for _, e := range m.entries {
		if e.V == v {
			return true
		}
	}
	return false
}
