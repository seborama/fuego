package fuego

// SC is a typed Stream cast function from a non-parameterised Stream[Any] to a parameterised Stream[U].
// SC receives a typed Stream[U].
//
// SC exists to address the current lack of support in Go for parameterised methods.
// See doc.go for more details.
//
// See C for A typed cast.
func SC[U any](from Stream[Any], to Stream[U]) Stream[U] {
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

// C is a typed cast function from a non-parameterised Stream[Any] to a parameterised type Stream[U].
// C receives a type U and creates a Stream[U].
//
// C exists to address the current lack of support in Go for parameterised methods.
// See doc.go for more details.
//
// See SC for A Stream cast.
func C[U any](from Stream[Any], to U) Stream[U] {
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

// CC is a typed cast function from a non-parameterised Stream[Any] to a parameterised type ComparableStream[U].
// CC receives a type U and creates a ComparableStream[U].
//
// CC exists to address the current lack of support in Go for parameterised methods and a performance issue with Go 1.18.
// See doc.go for more details.
func CC[U Comparable](from Stream[Any], to U) ComparableStream[U] {
	toCh := make(chan U, from.concurrency)

	toStream := NewConcurrentStream(toCh, from.concurrency)

	go func() {
		defer close(toStream.stream)

		for f := range from.stream {
			toStream.stream <- interface{}(f).(U)
		}
	}()

	return ComparableStream[U]{toStream}
}

// MC is a typed cast function from a non-parameterised Stream[Any] to a parameterised type MathableStream[U].
// MC receives a type U and creates a MathableStream[U].
//
// MC exists to address the current lack of support in Go for parameterised methods and a performance issue with Go 1.18.
// See doc.go for more details.
func MC[U Mathable](from Stream[Any], to U) MathableStream[U] {
	toCh := make(chan U, from.concurrency)

	toStream := NewConcurrentStream(toCh, from.concurrency)

	go func() {
		defer close(toStream.stream)

		for f := range from.stream {
			toStream.stream <- interface{}(f).(U)
		}
	}()

	return MathableStream[U]{toStream}
}
