package fuego

// A Map is an interface for all Map kinds.
type Map interface {
	// func (m HamtMap) BiStream() BiStream
	EntrySet() Set
	KeySet() Set
	Insert(k Entry, v interface{}) Map
	Delete(k Entry) Map
	Size() int
	FirstRest() (k Entry, v interface{}, rest Map)
	Merge(n Map) Map
	Get(k Entry) interface{}         // TODO return Maybe instead of Entry
	Has(k Entry, v interface{}) bool // TODO return EntryBool????
	HasKey(k Entry) bool             // TODO return EntryBool????
	HasValue(v interface{}) bool     // TODO return EntryBool????
}
