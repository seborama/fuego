package fuego

// Stream is a sequence of elements supporting sequential and parallel aggregate
// operations (TODO: not yet supported).
type Stream interface {
	Map(mapper Function) Stream
	Filter(predicate Predicate) Stream
	FlatMap(mapper Function) []interface{}
	None(predicate Predicate) Stream // TODO: a special case of Filter
	All(predicate Predicate) Stream  // TODO: a special case of Filter, opposite of None()
	FindFirst() Maybe
	FindAny() Maybe
	OfOne(i interface{}) Stream
	Of(i ...interface{}) Stream
}

// ReferenceStream is a simple implementation of a Stream.
type ReferenceStream struct {
	iterator Iterator
}

// NewStream creates a new Stream.
// TODO eventually change to Spliterable?
func NewStream(it Iterator) Stream {
	return ReferenceStream{
		iterator: it,
	}
}

// Map returns a stream consisting of the results of applying the given
// function to the elements of this stream.
func (rp ReferenceStream) Map(mapper Function) Stream {
	s := []interface{}{}
	for it := rp.iterator; it != nil; it = it.Forward() {
		s = append(s, mapper.Apply(it.Value()))
	}

	return NewStream(NewSliceIterator(s))
}

func (rp ReferenceStream) Filter(predicate Predicate) Stream {
	s := []interface{}{}
	for it := rp.iterator; it != nil; it = it.Forward() {
		if predicate.Test(it.Value()) {
			s = append(s, it.Value())
		}
	}

	return NewStream(NewSliceIterator(s))
}
