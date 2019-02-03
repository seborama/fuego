package fuego

import (
	"sort"
)

// EntryMap is an Entry for 'map[Entry]EntrySlice'.
type EntryMap map[Entry]EntrySlice

// Stream returns a stream of tuples the elements of the EntryMap.
func (em EntryMap) Stream() Stream {
	if em == nil {
		return NewStream(nil)
	}

	c := make(chan Entry, 1e3)

	for k, v := range em {
		c <- Tuple2{k, v}
	}
	close(c)

	return NewStream(c)
}

type tuple2 struct {
	key  Entry
	hash uint32
}

// Hash returns a hash for this Entry.
func (em EntryMap) Hash() uint32 {
	if em == nil || len(em) == 0 {
		return 0
	}

	var sortedKeyHashes []tuple2
	for k := range em {
		sortedKeyHashes = append(sortedKeyHashes, tuple2{
			key:  k,
			hash: k.Hash(),
		})
	}
	sort.Slice(sortedKeyHashes, func(i, j int) bool {
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
