package fuego

// MathableStream is a Stream of Mathable type.
type MathableStream[T Mathable] struct {
	Stream[T]
}

// Sum return the sum of all items on the stream.
// Panics if the channel is nil or the stream is empty.
// This is a special case of a reduction.
// This is a terminal operation and hence expects the producer to close the stream in order to complete.
func (s MathableStream[T]) Sum() T {
	if s.stream == nil {
		panic(PanicMissingChannel)
	}

	sum, ok := <-s.stream
	if !ok {
		panic(PanicNoSuchElement)
	}

	for val := range s.stream {
		sum = Sum(sum, val)
	}

	return sum
}

// Average returns the arithmetic average of the numbers in the stream.
// Panics if the channel is nil or the stream is empty.
// This is a special case of a reduction.
// This is a terminal operation and hence expects the producer to close the stream in order to complete.
func (s MathableStream[T]) Average() T {
	if s.stream == nil {
		panic(PanicMissingChannel)
	}

	sum, ok := <-s.stream
	if !ok {
		panic(PanicNoSuchElement)
	}

	var cnt T = 1

	for val := range s.stream {
		sum += val
		cnt++
	}

	return sum / cnt
}
