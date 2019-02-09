package fuego

// TODO: consider two types of streams: CIntStreams and SIntStreams. See streams.go for more info.

// IntStream is a sequence of EntryInt elements supporting sequential
// and (in the future?) parallel operations.
type IntStream struct {
	stream chan EntryInt
}

// NewIntStream creates a new IntStream.
// This function leaves the provided channel is the same state
// of openness.
func NewIntStream(c chan EntryInt) IntStream {
	if c == nil {
		panic(PanicMissingChannel)
	}
	return IntStream{
		stream: c,
	}
}

// NewIntStreamFromSlice creates a new IntStream from a Go slice of EntryInt.
// The stream will be closed once all the slice data has been
// published.
func NewIntStreamFromSlice(is []EntryInt) IntStream {
	c := make(chan EntryInt, 1e3)

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

	max, ok := <-is.stream
	if !ok {
		panic(PanicNoSuchElement)
	}
	for ei := range is.stream {
		if ei > max {
			max = ei
		}
	}
	return max
}

// Close the stream and returns true if success.
func (is IntStream) Close() bool {
	closed := false
	func() {
		defer func() { _ = recover() }()
		close(is.stream)
		closed = true
	}()
	return closed
}
