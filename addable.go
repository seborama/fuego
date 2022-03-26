package fuego

type AddableStream[T Addable] struct {
	Stream[T]
}

func (s AddableStream[T]) Sum() T {
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
