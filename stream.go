package fuego

// Stream is a sequence of elements supporting sequential and parallel aggregate
// operations (TODO: not yet supported).
type Stream interface {
	Map(mapper Function) Stream
	Filter(predicate Predicate) Stream
	// FlatMap(mapper Function) []interface{}
	// None(predicate Predicate) bool // i.e. NoneMatch
	// All(predicate Predicate) bool // i.e. AllMatch
	// FindFirst() Maybe
	// FindAny() Maybe
	// OfOne(i interface{}) Stream
	// Of(i ...interface{}) Stream
	ForEach(consumer Consumer)
	Reduce(f2 BiFunction) interface{}
	LeftReduce(f2 BiFunction) interface{}
	RightReduce(f2 BiFunction) interface{}
	Intersperse(e Entry) Stream
	GroupBy(classifier Function) Map
}

// ReferenceStream is a simple implementation of a Stream.
type ReferenceStream struct {
	iterator Iterator
}

// NewStream creates a new Stream.
func NewStream(it Iterator) Stream {
	return ReferenceStream{
		iterator: it,
	}
}

// Map returns a stream consisting of the results of applying the given
// function to the elements of this stream.
func (rp ReferenceStream) Map(mapper Function) Stream {
	s := []Entry{}
	for it := rp.iterator; it != nil; it = it.Forward() {
		s = append(s, mapper(it.Value()).(Entry))
	}

	return NewStream(NewSliceIterator(s))
}

// Filter returns a stream consisting of the elements of this stream that match
// the given predicate.
func (rp ReferenceStream) Filter(predicate Predicate) Stream {
	s := []Entry{}
	for it := rp.iterator; it != nil; it = it.Forward() {
		if predicate(it.Value()) {
			s = append(s, it.Value().(Entry))
		}
	}

	return NewStream(NewSliceIterator(s))
}

// ForEach executes the given function for each entry in this stream.
func (rp ReferenceStream) ForEach(consumer Consumer) {
	for it := rp.iterator; it != nil && it.Size() != 0; it = it.Forward() {
		consumer(it.Value())
	}
}

// LeftReduce accumulates the elements of this Set by
// applying the given function.
func (rp ReferenceStream) LeftReduce(f2 BiFunction) interface{} {
	it := rp.iterator
	if it == nil || it.Size() == 0 {
		return nil
	}
	res := it.Value()
	for it = it.Forward(); it != nil; it = it.Forward() {
		res = f2(res, it.Value())
	}
	return res
}

// Reduce is an alias for LeftReduce.
func (rp ReferenceStream) Reduce(f2 BiFunction) interface{} {
	return rp.LeftReduce(f2)
}

// RightReduce accumulates the elements of this Set by
// applying the given function.
func (rp ReferenceStream) RightReduce(f2 BiFunction) interface{} {
	if rp.iterator == nil || rp.iterator.Size() == 0 {
		return nil
	}
	reverse := NewStream(rp.iterator.Reverse())
	return reverse.LeftReduce(f2)
}

// Intersperse inserts an element between all elements of this Stream.
func (rp ReferenceStream) Intersperse(e Entry) Stream {
	if rp.iterator == nil || rp.iterator.Size() == 0 {
		return NewStream(NewSliceIterator([]Entry{}))
	}

	s := make([]Entry, rp.iterator.Size()*2-1)

	for it, idx := rp.iterator, 0; it != nil; it, idx = it.Forward(), idx+1 {
		s[2*idx] = it.Value().(Entry)
		if idx > 0 {
			s[2*idx-1] = e
		}
	}

	return NewStream(NewSliceIterator(s))
}

// GroupBy groups the elements of this Stream by classifying them.
func (rp ReferenceStream) GroupBy(classifier Function) Map {
	groups := map[Entry]OrderedSet{}
	for it := rp.iterator; it != nil; it = it.Forward() {
		k := classifier(it.Value())
		v := it.Value()
		groups[k] = groups[k].Insert(v).(OrderedSet)
	}

	newMap := NewOrderedMap()
	for k, v := range groups {
		newMap = newMap.Insert(k, v).(OrderedMap)
	}

	return newMap
}
