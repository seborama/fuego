package fuego

// Entry is the simplest behaviour that functional types must adhere to.
type Entry interface {
	Hash() uint32 // TODO: remove Hash() since the project no longer includes collections? Hashes suffer from collision.
	Equal(Entry) bool
	// Value() Entry
}
