package fuego

// Function that accepts one argument and produces a result.
type Function interface {
	Apply(i interface{}) interface{}
}

// Stream is a sequence of elements supporting sequential and parallel aggregate
// operations (TODO: not yet supported).
type Stream interface {
	Map(mapper Function) Stream
}

// ReferenceStream is a simple implementation of a Stream.
type ReferenceStream struct {
	iterator Iterator
}

// NewStream creates a new Stream.
// TODO eventually change to Spliterable?
func NewStream(it Iterable) Stream {
	return ReferenceStream{
		iterator: it.Iterator(),
	}
}

// Map returns a stream consisting of the results of applying the given
// function to the elements of this stream.
func (rp ReferenceStream) Map(mapper Function) Stream {
	s := []interface{}{}
	for it := rp.iterator; it != nil; it = it.Forward() {
		val := it.Value()
		s = append(s, mapper.Apply(val))
	}

	return NewStream(
		NewBaseIterable(NewSliceIterator(s)))
}
