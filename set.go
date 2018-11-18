package fuego

import (
	"github.com/raviqqe/hamt"
)

// A Set is an interface for all Set kinds.
type Set interface {
	Stream() Stream
	Insert(e hamt.Entry) Set
	Delete(e hamt.Entry) Set
	Size() int
	FirstRest() (hamt.Entry, Set)
	Merge(t Set) Set
}
