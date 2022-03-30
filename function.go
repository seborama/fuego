package fuego

// Function that accepts one argument and produces a result.
type Function[T, R any] func(T) R

// BiFunction that accepts two arguments and produces a result.
type BiFunction[T, U, R any] func(T, U) R

// BinaryOperator that accepts two arguments of the same type and produces a result of the same type.
// This is a special case of BiFunction.
type BinaryOperator[T any] func(T, T) T

// StreamFunction that accepts one argument and produces a Stream[R].
//
// It is used with when "flat mapping" a `Stream`.
// In effect, this is a one to many operation, such as exploding the individual values of a slice
// into a Stream[R].
type StreamFunction[T, R any] func(T) Stream[R]

// FlattenSlice is a StreamFunction that flattens a []T slice to a Stream[Any] of its elements.
func FlattenSlice[T any](bufsize int) StreamFunction[[]T, Any] {
	return func(el []T) Stream[Any] {
		return NewStreamFromSlice(el, bufsize).StreamAny()
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
func Sum[T Mathable](a, b T) T {
	return a + b
}

// Concatenate is BiFunction that returns the concatenation of the two strings passed to it.
func Concatenate[T string](a, b T) T {
	return a + b
}

// Identity is a basic Function that returns the original value passed to it, unmodified.
func Identity[T any](v T) T {
	return v
}

// ToAny is a basic Function that returns the original value passed to it, cast to an 'Any' type.
func ToAny[T any](v T) Any {
	return v
}
