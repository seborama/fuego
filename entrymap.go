package fuego

import (
	"sort"
)

const (
	// PanicDuplicateKey signifies that the key already exists for this map.
	PanicDuplicateKey = "duplicate key"
)

// EntryMap is an Entry for 'map[Entry]Entry'.
type EntryMap map[Entry]Entry

// Stream returns a stream of tuples the elements of the EntryMap.
func (em EntryMap) Stream(bufsize int) Stream {
	c := make(chan Entry, bufsize)

	go func() {
		defer close(c)

		for k, v := range em {
			c <- Tuple2{E1: k, E2: v}
		}
	}()

	return NewStream(c)
}

// Hash returns a hash for this Entry.
func (em EntryMap) Hash() uint32 {
	type tuple2 struct {
		key  Entry
		hash uint32
	}

	if len(em) == 0 {
		return 0
	}

	sortedKeyHashes := []tuple2{}
	for k := range em {
		sortedKeyHashes = append(sortedKeyHashes, tuple2{
			key:  k,
			hash: k.Hash(),
		})
	}
	sort.SliceStable(sortedKeyHashes, func(i, j int) bool {
		return sortedKeyHashes[i].hash < sortedKeyHashes[j].hash
	})

	result := uint32(1)
	for _, keyHash := range sortedKeyHashes {
		result = 31*result + (keyHash.hash ^ em[keyHash.key].Hash())
	}

	return result
}

// Equal returns true if this type is equal to 'e'.
func (em EntryMap) Equal(e Entry) bool {
	if _, ok := e.(EntryMap); !ok {
		return false
	}

	return em.Hash() == e.Hash()
}

// HasKey returns true if this map has the supplied key.
func (em EntryMap) HasKey(key Entry) bool {
	for k := range em {
		if k.Equal(key) {
			return true
		}
	}

	return false
}

// Len returns the number of Entries in this EntryMap.
func (em EntryMap) Len() int {
	return len(em)
}

// Merge performs a value merge on the specified key in this EntryMap.
// If the supplied key does not exist in this EntryMap, the supplied value
// is assigned.
// If the supplied key already exists, mergeFunction is called to determine the new value
// that will replace the current.
func (em EntryMap) Merge(key, value Entry, mergeFunction BiFunction) EntryMap {
	if !em.HasKey(key) {
		em[key] = value
		return em
	}

	em[key] = mergeFunction(em[key], value)
	return em
}
