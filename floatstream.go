package fuego

// TODO: consider two types of streams: CFloatStreams and SFloatStreams. See streams.go for more info.

// FloatStream is a sequence of EntryFloat elements supporting sequential
// and (in the future?) parallel operations.
//
// The current implementation is based on `Stream` and an intermediary
// channel that converts incoming `EntryFloat` elements to `Entry`. This
// approach offers programming conciseness but the use of an
// intermediary channel likely decreases performance. This also
// means that type checking is weak on methods "borrowed" from
// `Stream` that expect `Entry` (instead of `EntryFloat`).
type FloatStream struct {
	Stream
}

// NewFloatStream creates a new FloatStream.
// This function leaves the provided channel is the same state
// of openness.
func NewFloatStream(c chan EntryFloat) FloatStream {
	return NewConcurrentFloatStream(c, 0)
}

// NewConcurrentFloatStream creates a new FloatStream with a degree
// of concurrency of n.
// This function leaves the provided channel is the same state
// of openness.
func NewConcurrentFloatStream(c chan EntryFloat, n int) FloatStream {
	if c == nil {
		panic(PanicMissingChannel)
	}
	ic := make(chan Entry, cap(c))
	go func() {
		defer close(ic)
		for val := range c {
			ic <- val
		}
	}()

	return FloatStream{
		NewConcurrentStream(ic, n),
	}
}

// NewFloatStreamFromSlice creates a new FloatStream from a Go slice of EntryFloat.
// The stream will be closed once all the slice data has been
// published.
func NewFloatStreamFromSlice(is []EntryFloat, bufsize int) FloatStream {
	c := make(chan EntryFloat, bufsize)

	go func() {
		defer close(c)
		for _, element := range is {
			c <- element
		}
	}()

	return NewFloatStream(c)
}

// Max returns the largest number in the stream.
// Panics if the channel is nil or the stream is empty.
// This is a special case of a reduction and is equivalent to:
//   is.Reduce(max) // where max is a BiFunction that returns
//                  // the largest of two integers.
// This is a terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
func (is FloatStream) Max() EntryFloat {
	if is.stream == nil {
		panic(PanicMissingChannel)
	}

	val, ok := <-is.stream
	if !ok {
		panic(PanicNoSuchElement)
	}
	max := val.(EntryFloat)
	for e := range is.stream {
		ei := e.(EntryFloat)
		if ei > max {
			max = ei
		}
	}
	return max
}

// Min returns the smallest number in the stream.
// Panics if the channel is nil or the stream is empty.
// This is a special case of a reduction and is equivalent to:
//   is.Reduce(min) // where min is a BiFunction that returns
//                  // the smallest of two integers.
// This is a terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
func (is FloatStream) Min() EntryFloat {
	if is.stream == nil {
		panic(PanicMissingChannel)
	}

	val, ok := <-is.stream
	if !ok {
		panic(PanicNoSuchElement)
	}
	min := val.(EntryFloat)
	for e := range is.stream {
		ei := e.(EntryFloat)
		if ei < min {
			min = ei
		}
	}
	return min
}

// Sum adds the numbers in the stream.
// Panics if the channel is nil or the stream is empty.
// This is a special case of a reduction and is equivalent to:
//   is.Reduce(sum) // where max is a BiFunction that adds
//                  // two integers.
// This is a terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
func (is FloatStream) Sum() EntryFloat {
	if is.stream == nil {
		panic(PanicMissingChannel)
	}

	val, ok := <-is.stream
	if !ok {
		panic(PanicNoSuchElement)
	}
	sum := val.(EntryFloat)
	for e := range is.stream {
		ei := e.(EntryFloat)
		sum += ei
	}
	return sum
}

// Average returns the average of the numbers in the stream.
// Panics if the channel is nil or the stream is empty.
// This is a special case of a reduction.
// This is a terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
func (is FloatStream) Average() EntryFloat {
	if is.stream == nil {
		panic(PanicMissingChannel)
	}

	val, ok := <-is.stream
	if !ok {
		panic(PanicNoSuchElement)
	}
	cnt := 1
	sum := val.(EntryFloat)
	for e := range is.stream {
		ei := e.(EntryFloat)
		sum += ei
		cnt++
	}
	return EntryFloat(float32(sum) / float32(cnt))
}
