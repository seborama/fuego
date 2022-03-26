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
		max = Max(max, val)
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
		min = Min(min, val)
	}

	return min
}

func (s ComparableStream[T]) Sum() interface{} {
	if s.stream == nil {
		panic(PanicMissingChannel)
	}

	val, ok := <-s.stream
	if !ok {
		panic(PanicNoSuchElement)
	}

	sum := val

	for val = range s.stream {
		sum = Sum(sum, val)
	}

	return sum
}
