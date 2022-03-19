package fuego

// Entry is the simplest behaviour that functional types must adhere to.
type entry[E any] interface {
	Hash() uint32 // TODO: remove Hash() since the project no longer includes collections? Hashes suffer from collision.
	Equal(E) bool
}
type Entry[E entry[E]] interface {
	// Hash() uint32 // TODO: remove Hash() since the project no longer includes collections? Hashes suffer from collision.
	// Equal(E) bool
	entry[E]
}
