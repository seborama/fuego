package fuego

// TODO: consider two types of streams: CIntStreams and SIntStreams. See streams.go for more info.

// IntStream is a sequence of EntryInt elements supporting sequential
// and (in the future?) parallel operations.
type IntStream struct {
	stream chan EntryInt
}

// NewIntStream creates a new IntStream.
func NewIntStream(c chan EntryInt) IntStream {
	return IntStream{
		stream: c,
	}
}

// NewIntStreamFromSlice creates a new IntStream from a Go slice of EntryInt.
func NewIntStreamFromSlice(is []EntryInt) IntStream {
	c := make(chan EntryInt, 1e3)
	defer close(c)
	for _, element := range is {
		c <- element
	}

	return NewIntStream(c)
}

// Close the stream.
// Panics if s.stream is nil or already closed.
func (is IntStream) Close() {
	close(is.stream)
}
