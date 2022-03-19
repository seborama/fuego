package fuego

var _ Entry[EntrySlice[EntryInt]] = EntrySlice[EntryInt]{0}

// EntrySlice is a slice of type E.
type EntrySlice[E Entry[E]] []E

// TODO: implement Stream() (see EntryMap)

// Hash returns a hash for this Entry.
func (es EntrySlice[E]) Hash() uint32 {
	if len(es) == 0 {
		return 0
	}

	result := uint32(1)

	for _, element := range es {
		h := element.Hash()
		result = 31*result + h
	}

	return result
}

// Equal returns true if this type is equal to 'e'.
func (es EntrySlice[E]) Equal(other EntrySlice[E]) bool {
	return es.Hash() == other.Hash()
}

// Append an Entry to this EntrySlice.
func (es EntrySlice[E]) Append(e E) EntrySlice[E] {
	return append(es, e)
}

// Len returns the number of Entries in this EntrySlice.
func (es EntrySlice[E]) Len() int {
	return len(es)
}
