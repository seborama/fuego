package fuego

var _ Entry[EntryInt] = EntryInt(0)

// EntryInt is an Entry for 'int'.
type EntryInt int

// Hash returns a hash for 'i'.
func (i EntryInt) Hash() uint32 {
	return uint32(i)
}

// Equal returns true if 'e' and 'i' are equal.
func (i EntryInt) Equal(e EntryInt) bool {
	return i == e
}
