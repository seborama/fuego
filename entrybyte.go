package fuego

// EntryByte is an Entry for 'byte'.
type EntryByte byte

// Hash returns a hash for 'i'.
func (i EntryByte) Hash() uint32 {
	return uint32(i)
}

// Equal returns true if 'e' and 'i' are equal.
func (i EntryByte) Equal(e Entry) bool {
	return i == e
}
