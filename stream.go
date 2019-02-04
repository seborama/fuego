package fuego

// Stream is a sequence of elements supporting sequential and parallel
// parallel operations.
type Stream struct {
	stream chan Entry
}

// NewStream creates a new Stream.
func NewStream(c chan Entry) Stream {
	return Stream{
		stream: c,
	}
}

// NewStreamFromSlice creates a new Stream from a Go slice.
func NewStreamFromSlice(s []Entry) Stream {
	c := make(chan Entry, 1e3)
	defer close(c)
	for _, element := range s {
		c <- element
	}

	return NewStream(c)
}

// TODO: implement NewStreamFromMap?

// Map returns a slice of channel of Set consisting of the results of
// applying the given function to the elements of this Set
func (s Stream) Map(mapper Function) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() { // TODO: introduce a cut-off to prevent the go func from straying
		defer close(outstream)
		if s.stream == nil {
			return
		}
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
	outstream := make(chan Entry, cap(s.stream))

	go func() { // TODO: introduce a cut-off to prevent the go func from straying
		defer close(outstream)
		if s.stream == nil {
			return
		}
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
func (s Stream) LeftReduce(f2 BiFunction) Entry {
	if s.stream == nil {
		return nil
	}

	res := <-s.stream
	for val := range s.stream {
		res = f2(res, val)
	}

	return res
}

// Reduce is an alias for LeftReduce.
func (s Stream) Reduce(f2 BiFunction) Entry {
	return s.LeftReduce(f2)
}

// Intersperse inserts an element between all elements of this Stream.
func (s Stream) Intersperse(e Entry) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() { // TODO: introduce a cut-off to prevent the go func from straying
		defer close(outstream)
		if s.stream == nil {
			return
		}
		if val := <-s.stream; val != nil {
			outstream <- val
		}
		for val := range s.stream {
			outstream <- e
			outstream <- val
		}
	}()

	return Stream{
		stream: outstream,
	}
}

// GroupBy groups the elements of this Stream by classifying them.
func (s Stream) GroupBy(classifier Function) EntryMap {
	resultMap := EntryMap{}

	if s.stream != nil {
		for val := range s.stream {
			k := classifier(val)
			resultMap[k] = append(resultMap[k], val)
		}
	}

	return resultMap
}
