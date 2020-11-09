package fuego

// TODO: consider two types of streams: CIntStreams and SIntStreams. See streams.go for more info.

// IntStream is a sequence of EntryInt elements supporting sequential
// and (in the future?) parallel operations.
//
// The current implementation is based on `Stream` and an intermediary
// channel that converts incoming `EntryInt` elements to `Entry`. This
// approach offers programming conciseness but the use of an
// intermediary channel likely decreases performance. This also
// means that type checking is weak on methods "borrowed" from
// `Stream` that expect `Entry` (instead of `EntryInt`).
type IntStream struct {
	Stream
}

// NewIntStream creates a new IntStream.
// This function leaves the provided channel is the same state
// of openness.
func NewIntStream(c chan EntryInt) IntStream {
	return NewConcurrentIntStream(c, 0)
}

// NewConcurrentIntStream creates a new IntStream with a degree
// of concurrency of n.
// This function leaves the provided channel is the same state
// of openness.
func NewConcurrentIntStream(c chan EntryInt, n int) IntStream {
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

	return IntStream{
		Stream: NewConcurrentStream(ic, n),
	}
}

// NewIntStreamFromSlice creates a new IntStream from a Go slice of EntryInt.
// The stream will be closed once all the slice data has been
// published.
func NewIntStreamFromSlice(is []EntryInt, bufsize int) IntStream {
	c := make(chan EntryInt, bufsize)

	go func() {
		defer close(c)
		for _, element := range is {
			c <- element
		}
	}()

	return NewIntStream(c)
}

// Max returns the largest number in the stream.
// Panics if the channel is nil or the stream is empty.
// This is a special case of a reduction and is equivalent to:
//   is.Reduce(max) // where max is a BiFunction that returns
//                  // the largest of two integers.
// This is a terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
func (is IntStream) Max() EntryInt {
	if is.stream == nil {
		panic(PanicMissingChannel)
	}

	val, ok := <-is.stream
	if !ok {
		panic(PanicNoSuchElement)
	}
	max := val.(EntryInt)
	for e := range is.stream {
		ei := e.(EntryInt)
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
func (is IntStream) Min() EntryInt {
	if is.stream == nil {
		panic(PanicMissingChannel)
	}

	val, ok := <-is.stream
	if !ok {
		panic(PanicNoSuchElement)
	}
	min := val.(EntryInt)
	for e := range is.stream {
		ei := e.(EntryInt)
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
func (is IntStream) Sum() EntryInt {
	if is.stream == nil {
		panic(PanicMissingChannel)
	}

	val, ok := <-is.stream
	if !ok {
		panic(PanicNoSuchElement)
	}
	sum := val.(EntryInt)
	for e := range is.stream {
		ei := e.(EntryInt)
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
func (is IntStream) Average() EntryInt {
	if is.stream == nil {
		panic(PanicMissingChannel)
	}

	val, ok := <-is.stream
	if !ok {
		panic(PanicNoSuchElement)
	}
	cnt := 1
	sum := val.(EntryInt)
	for e := range is.stream {
		ei := e.(EntryInt)
		sum += ei
		cnt++
	}
	return EntryInt(int(sum) / (cnt))
}
