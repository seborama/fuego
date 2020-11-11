package fuego

// Entry is the simplest behaviour that functional types must adhere to.
type Entry interface {
	Hash() uint32 // TODO: remove Hash() since the project no longer includes collections? Hashes suffer from collision.
	Equal(Entry) bool
}

// EntriesEqual checks the equality of 2 Entry objects.
// Note:
// EntriesEqual(&entry1, &entry2) will not produce the desired
// outcome with this method.
func EntriesEqual(e1, e2 Entry) bool {
	return (e1 == nil && e2 == nil) || (e1 != nil && e1.Equal(e2))
}
