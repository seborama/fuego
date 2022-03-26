package fuego

// Function that accepts one argument and produces a result.
type Function[T any, R any] func(T) R

// BiFunction that accepts two arguments and produces a result.
type BiFunction[T any, U any, R any] func(T, U) R

// StreamFunction that accepts one argument and produces a Stream[R].
//
// It is used with when "flat mapping" a `Stream`.
// In effect, this is a one to many operation, such as exploding the individual values of a slice
// into a Stream[R].
type StreamFunction[T any, R any] func(T) Stream[R]

// FlattenSlice is a StreamFunction that flattens a []T slice to a Stream[R] of its elements.
func FlattenSlice[T any](bufsize int) StreamFunction[[]T, R] {
	return func(el []T) Stream[R] {
		return NewStreamFromSlice(el, bufsize).StreamR()
	}
}

// FlattenTypedSlice is a StreamFunction that flattens a []T slice to a Stream[T] of its elements.
func FlattenTypedSlice[T any](bufsize int) StreamFunction[[]T, T] {
	return func(el []T) Stream[T] {
		return NewStreamFromSlice(el, bufsize)
	}
}

// Max is a BiFunction that returns the greatest of 2 values.
func Max[T Comparable](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min is a BiFunction that returns the smallest of 2 values.
func Min[T Comparable](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Sum is a BiFunction that returns the sum of 2 values.
func Sum[T Addable](a, b T) T {
	return a + b
}

// Identity is a basic Function that returns the original value passed to it, unmodified.
func Identity[T any](v T) T {
	return v
}

// ToR is a basic Function that returns the original value passed to it, cast to an 'R' type.
func ToR[T any](v T) R {
	return v
}
