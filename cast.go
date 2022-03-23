package fuego

// SC is a typed Stream cast function from a non-parameterised Stream[R] to a parameterised Stream[U].
// SC receives a typed Stream[U].
//
// SC exists to address the current lack of support in Go for parameterised methods.
// See doc.go for more details.
//
// See C for A typed cast.
func SC[U any](from Stream[R], to Stream[U]) Stream[U] {
	toCh := make(chan U, from.concurrency)
	to.stream = toCh

	go func() {
		defer close(to.stream)

		for f := range from.stream {
			to.stream <- interface{}(f).(U)
		}
	}()

	return to
}

// C is a typed cast function from a non-parameterised Stream[R] to a parameterised type Stream[U].
// C receives a type U and creates a Stream[U].
//
// C exists to address the current lack of support in Go for parameterised methods.
// See doc.go for more details.
//
// See SC for A Stream cast.
func C[U any](from Stream[R], to U) Stream[U] {
	toCh := make(chan U, from.concurrency)

	toStream := NewConcurrentStream(toCh, from.concurrency)

	go func() {
		defer close(toStream.stream)

		for f := range from.stream {
			toStream.stream <- interface{}(f).(U)
		}
	}()

	return toStream
}
