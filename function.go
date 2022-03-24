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
