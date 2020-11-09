package fuego

import (
	"fmt"
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

// UniqueKeyMapAppender appends entry 'e' to map 'm' e.E1 is not a key
// of 'm' or it panics.
// nolint: gochecknoglobals
var UniqueKeyMapAppender = func(m, e Entry) Entry {
	if m.(EntryMap).HasKey(e.(Tuple2).E1) {
		panic(fmt.Sprintf("%s: %v", PanicDuplicateKey, e.(Tuple2).E1))
	}
	m.(EntryMap)[e.(Tuple2).E1] = e.(Tuple2).E2
	return m
}

// Append an Entry to this EntryMap.
func (em EntryMap) Append(e Entry) EntryMap {
	return em.AppendMerge(e, UniqueKeyMapAppender)
}

// AppendMerge appends an Entry to this EntryMap as with Append but
// calls the merge function when this map already contains the key.
func (em EntryMap) AppendMerge(e Entry, mergeFunction BiFunction) EntryMap {
	return mergeFunction(em, e).(EntryMap)
}

// Len returns the number of Entries in this EntryMap.
func (em EntryMap) Len() int {
	return len(em)
}
