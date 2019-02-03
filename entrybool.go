package fuego

// EntryBool is an Entry for 'bool'.
type EntryBool bool

// Hash returns a hash for this Entry.
func (eb EntryBool) Hash() uint32 {
	if bool(eb) {
		return 1
	}
	return 0
}

// Equal returns true if this type is equal to 'e'.
func (eb EntryBool) Equal(e Entry) bool {
	return eb == e
}
