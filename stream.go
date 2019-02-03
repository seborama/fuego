package fuego

type Map map[Entry][]Entry

// Stream is a sequence of elements supporting sequential and parallel
// parallel operations.
type Stream struct {
	stream chan Entry
}

// NewStream creates a new Stream.
func NewStream(s chan Entry) Stream {
	return Stream{
		stream: s,
	}
}

// Map returns a slice of channel of Set consisting of the results of
// applying the given function to the elements of this Set
func (s Stream) Map(mapper Function) Stream {
	if s.stream == nil {
		return Stream{stream: nil}
	}

	outstream := make(chan Entry, cap(s.stream))

	go func() { // TODO: introduce a cut-off to prevent the go func leak
		defer close(outstream)
		for val := range s.stream {
			outstream <- mapper(val)
		}
	}()

	return Stream{
		stream: outstream,
	}
}

// Filter returns a stream consisting of the elements of this stream that
// match the given predicate.
func (s Stream) Filter(predicate Predicate) Stream {
	if s.stream == nil {
		return Stream{stream: nil}
	}

	outstream := make(chan Entry, cap(s.stream))

	go func() { // TODO: introduce a cut-off to prevent the go func leak
		defer close(outstream)
		for val := range s.stream {
			if predicate(val) {
				outstream <- val
			}
		}
	}()

	return Stream{
		stream: outstream,
	}
}

// ForEach executes the given function for each entry in this stream.
func (s Stream) ForEach(consumer Consumer) {
	if s.stream == nil {
		return
	}

	for val := range s.stream {
		consumer(val)
	}
}

// LeftReduce accumulates the elements of this Set by
// applying the given function.
// func (rp Stream) LeftReduce(f2 BiFunction) interface{} {
// 	it := rp.iterator
// 	if it == nil || it.Size() == 0 {
// 		return nil
// 	}
// 	res := it.Value()
// 	for it = it.Forward(); it != nil; it = it.Forward() {
// 		res = f2(res, it.Value())
// 	}
// 	return res
// }

// Reduce is an alias for LeftReduce.
// func (rp Stream) Reduce(f2 BiFunction) interface{} {
// 	return rp.LeftReduce(f2)
// }

// Intersperse inserts an element between all elements of this Stream.
// func (rp Stream) Intersperse(e Entry) Stream {
// 	if rp.iterator == nil || rp.iterator.Size() == 0 {
// 		return NewStream(NewSliceIterator([]Entry{}))
// 	}

// 	s := make([]Entry, rp.iterator.Size()*2-1)

// 	for it, idx := rp.iterator, 0; it != nil; it, idx = it.Forward(), idx+1 {
// 		s[2*idx] = it.Value().(Entry)
// 		if idx > 0 {
// 			s[2*idx-1] = e
// 		}
// 	}

// 	return NewStream(NewSliceIterator(s))
// }

// GroupBy groups the elements of this Stream by classifying them.
// func (rp Stream) GroupBy(classifier Function) Map {
// 	groups := Map{}
// 	for it := rp.iterator; it != nil; it = it.Forward() {
// 		k := classifier(it.Value())
// 		v := it.Value()
// 		groups[k] = groups[k].Insert(v).(OrderedSet)
// 	}

// 	newMap := NewOrderedMap()
// 	for k, v := range groups {
// 		newMap = newMap.Insert(k, v).(OrderedMap)
// 	}

// 	return newMap
// }
