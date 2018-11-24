package fuego

// A Set is an interface for all Set kinds.
type Set interface {
	Stream() Stream
	Insert(e Entry) Set
	Delete(e Entry) Set
	Size() int
	FirstRest() (Entry, Set)
	Merge(t Set) Set
}
