package fuego

// A Map is an interface for all Map kinds.
type Map interface {
	// func (m HamtMap) BiStream() BiStream
	EntrySet() Set
	KeySet() Set
	Insert(k Entry, v interface{}) Map
	Delete(k Entry) Map
	Size() int
	FirstRest() (Entry, interface{}, Map)
	Merge(n Map) Map
	Get(k Entry) interface{} // TODO return Maybe instead of interface{}
	Has(k Entry, v interface{}) bool
	HasKey(k Entry) bool
	HasValue(v interface{}) bool
}
