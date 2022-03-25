package fuego

type ComparableStream[T Comparable] struct {
	Stream[T]
}

func (s ComparableStream[T]) Max() interface{} {
	if s.stream == nil {
		panic(PanicMissingChannel)
	}

	val, ok := <-s.stream
	if !ok {
		panic(PanicNoSuchElement)
	}

	max := val

	for val = range s.stream {
		max = Max(val, max)
	}

	return max
}

func (s ComparableStream[T]) Min() interface{} {
	if s.stream == nil {
		panic(PanicMissingChannel)
	}

	val, ok := <-s.stream
	if !ok {
		panic(PanicNoSuchElement)
	}

	min := val

	for val = range s.stream {
		min = Min(val, min)
	}

	return min
}
