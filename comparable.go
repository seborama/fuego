package fuego

// ComparableStream is a stream of Comparable type.
type ComparableStream[T Comparable] struct {
	Stream[T]
}

func (s ComparableStream[T]) Max() T {
	if s.stream == nil {
		panic(PanicMissingChannel)
	}

	val, ok := <-s.stream
	if !ok {
		panic(PanicNoSuchElement)
	}

	max := val

	for val = range s.stream {
		max = Max(max, val)
	}

	return max
}

func (s ComparableStream[T]) Min() T {
	if s.stream == nil {
		panic(PanicMissingChannel)
	}

	val, ok := <-s.stream
	if !ok {
		panic(PanicNoSuchElement)
	}

	min := val

	for val = range s.stream {
		min = Min(min, val)
	}

	return min
}
