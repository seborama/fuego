package fuego

import (
	"github.com/raviqqe/hamt"
)

// A Map is an interface for all Map kinds.
type Map interface {
	// func (m HamtMap) BiStream() BiStream
	EntrySet() Set
	KeySet() Set
	Insert(k hamt.Entry, v interface{}) Map
	Delete(k hamt.Entry) Map
	Size() int
	FirstRest() (hamt.Entry, interface{}, Map)
	Merge(n Map) Map
	Get(k hamt.Entry) interface{} // TODO return Maybe instead of interface{}
	Has(k hamt.Entry, v interface{}) bool
	HasKey(k hamt.Entry) bool
	HasValue(v interface{}) bool
}
