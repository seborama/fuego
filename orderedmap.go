package fuego

// A OrderedMap is an ordered map
type OrderedMap struct {
	entries OrderedSet
}

// NewOrderedMap creates a new Map
func NewOrderedMap() OrderedMap {
	return OrderedMap{
		entries: NewOrderedSet(),
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
	return m.entries
}

// KeySet returns a Set of keys contained in this map.
// Since KeySet returns a Set, it can be streamed with Set.Stream().
// Note that ValueSet() is not implemented because Values can be present
// multiple times. This could possibly be implemented via []interface{}
// or a Sequence?
// It also could be better to use the BiStream() proposed in this file.
func (m OrderedMap) KeySet() Set {
	keySet := NewOrderedSet()
	for it := NewSetIterator(m.entries); it != nil; it = it.Forward() {
		keySet = keySet.Insert(it.Value().(MapEntry).K).(OrderedSet)
	}
	return keySet
}

// Insert a value into this map.
func (m OrderedMap) Insert(k Entry, v interface{}) Map {
	return OrderedMap{
		entries: m.entries.
			Insert(MapEntry{
				K: k,
				V: v,
			}).(OrderedSet),
	}
}

// Delete a value from this map.
func (m OrderedMap) Delete(k Entry) Map {
	return OrderedMap{
		entries: m.entries.
			Delete(MapEntry{
				K: k,
				V: nil, // MapEntry equality is based solely on MapEntry.K
			}).(OrderedSet),
	}
}

// Size of the Set.
func (m OrderedMap) Size() int {
	return m.entries.Size()
}

// FirstRest returns a key-value pair in a map and a rest of the map.
// This method is useful for iteration.
// The key and value would be nil if the map is empty.
func (m OrderedMap) FirstRest() (Entry, interface{}, Map) {
	e, rest := m.entries.FirstRest()
	return e.(MapEntry).K, e.(MapEntry).V,
		OrderedMap{
			entries: rest.(OrderedSet)}
}

// Merge this map and given map.
func (m OrderedMap) Merge(n Map) Map {
	newMap := m
	for it := NewSetIterator(n.EntrySet()); it != nil; it = it.Forward() {
		newMap = newMap.Insert(it.Value().(MapEntry).K, it.Value().(MapEntry).V).(OrderedMap)
	}
	return newMap
}

// Get a value in this map corresponding to a given key.
// It returns nil if no value is found.
func (m OrderedMap) Get(k Entry) interface{} {
	for it := NewSetIterator(m.entries); it != nil; it = it.Forward() {
		if it.Value().(MapEntry).K.Equal(k) {
			return it.Value().(MapEntry).V
		}
	}
	return EntryNone{}
}

// Has returns true if a key-value pair corresponding with a given key is
// included in a map, or false otherwise.
func (m OrderedMap) Has(k Entry, v interface{}) bool {
	return m.Get(k) == v // TODO review this to use .Equal should we implement MapEntry(Entry, Entry)
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
	for it := NewSetIterator(m.entries); it != nil; it = it.Forward() {
		if it.Value().(MapEntry).V == v {
			return true
		}
	}
	return false
}
