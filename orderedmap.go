package fuego

// A OrderedMap is an ordered map
type OrderedMap struct {
	entries []MapEntry // TODO use OrderedSet
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
	// len+1 to keep room for the '(k,v)' if not already present
	newMap := make([]MapEntry, len(m.entries)+1)
	copy(newMap, m.entries)

	foundExisting := false
	for idx, e := range m.entries {
		if e.Equal(k) {
			foundExisting = true
			newMap[idx] = MapEntry{K: k, V: v}
			newMap = newMap[:len(newMap)-1] // remove unneeded extra room
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
			switch idx {
			case 0:
				sCopy = make([]MapEntry, len(m.entries)-1)
				copy(sCopy, m.entries[1:])

			case m.Size() - 1:
				sCopy = make([]MapEntry, len(m.entries)-1)
				copy(sCopy, m.entries[:idx])

			default:
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

	sliceIndex := make(map[Entry]bool, len(m.entries))
	for _, v := range m.entries {
		sliceIndex[v.K] = true
	}

	for _, entry := range n.(OrderedMap).entries {
		if !sliceIndex[entry.K] {
			merge = append(merge, entry)
		}
	}
	return OrderedMap{
		entries: merge,
	}
}

// Get a value in this map corresponding to a given key.
// It returns nil if no value is found.
func (m OrderedMap) Get(k Entry) interface{} {
	for _, e := range m.entries {
		if e.K.Equal(k) {
			return e.V
		}
	}

	return EntryNone{}
}

// Has returns true if a key-value pair corresponding with a given key is
// included in a map, or false otherwise.
func (m OrderedMap) Has(k Entry, v interface{}) bool {
	value := m.Get(k)

	if _, ok := value.(EntryNone); ok {
		return false
	}
	return value == v
}

// HasKey returns true if a given key exists
// in a map, or false otherwise.
func (m OrderedMap) HasKey(k Entry) bool {
	_, ok := m.Get(k).(EntryNone)
	return !ok
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
