package fuego

// An OrderedSet is an ordered set
type OrderedSet struct {
	slice []Entry
}

// NewOrderedSet creates a new OrderedSet
func NewOrderedSet() OrderedSet {
	return OrderedSet{
		slice: []Entry{},
	}
}

// Stream returns a sequential Stream with this collection as its source.
func (s OrderedSet) Stream() Stream {
	return NewStream(
		NewSliceIterator(s.slice))
}

// Insert a value into this set.
func (s OrderedSet) Insert(e Entry) Set {
	for _, entry := range s.slice {
		if e.Equal(entry) {
			sCopy := make([]Entry, len(s.slice))
			copy(sCopy, s.slice)
			return OrderedSet{
				slice: sCopy,
			}
		}
	}
	return OrderedSet{
		slice: append(s.slice, e),
	}
}

// Delete a value from this set.
func (s OrderedSet) Delete(e Entry) Set {
	for idx, val := range s.slice {
		if val.Equal(e) {
			var sCopy []Entry
			if idx == 0 {
				sCopy = make([]Entry, len(s.slice)-1)
				copy(sCopy, s.slice[1:])
			} else if idx == s.Size()-1 {
				sCopy = make([]Entry, len(s.slice)-1)
				copy(sCopy, s.slice[:idx])
			} else {
				sCopy = append(s.slice[:idx], s.slice[idx+1:]...)
			}
			return OrderedSet{
				slice: sCopy,
			}
		}
	}

	// 'e' not found (includes the case where s.slice is empty)
	sCopy := make([]Entry, len(s.slice))
	copy(sCopy, s.slice)
	return OrderedSet{
		slice: sCopy,
	}
}

// Size of the OrderedSet.
func (s OrderedSet) Size() int {
	return len(s.slice)
}

// FirstRest returns a value in a set and a rest of the set.
// This method is useful for iteration.
func (s OrderedSet) FirstRest() (Entry, Set) {
	sCopy := make([]Entry, len(s.slice)-1)
	copy(sCopy, s.slice[1:])
	return s.slice[0], OrderedSet{slice: sCopy}
}

// Merge 2 sets into one.
func (s OrderedSet) Merge(t Set) Set {
	merge := make([]Entry, len(s.slice))
	copy(merge, s.slice)

	sliceIndex := make(map[Entry]bool, len(s.slice))
	for _, v := range s.slice {
		sliceIndex[v] = true
	}

	for _, entry := range t.(OrderedSet).slice {
		if !sliceIndex[entry] {
			merge = append(merge, entry)
		}
	}

	return OrderedSet{
		slice: merge,
	}
}
