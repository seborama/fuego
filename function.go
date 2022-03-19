package fuego

// Function that accepts one argument and produces a result.
// TODO: is this correct?? Should it be Function[E2,R] where is defined as Function[E1,E2 any] func(e E1) E2
// TODO: in the same line of thoughts, all *Function's should be re-written to accommodate the specification of the result type.
type Function[T, R any] func(T) R

// StreamFunction that accepts one argument and produces a stream.
//
// It is worth noting is `StreamFunction` which accepts one argument
//  and produces a stream. It is used with when "flat mapping" and
// `Stream`. This effectively is a one to many operation, such
// as exploding the individual values of an EntrySlice into a Stream.
//
// This effectively is a one to many operation, such as exploding
// the individual values of an EntrySlice into a Stream.
type StreamFunction[E Entry[E]] func(e E) Stream[E]

// BiFunction that accepts two arguments and produces a result.
type BiFunction[T, U, R any] func(T, U) R

// A PartialFunction applies a function to a type if the value
// is within the domain of the partial function.
// For instance, a partial function could return the doubles
// of the odd numbers in a series of mixed odd and even numbers.
// type PartialFunction func(e Entry, p Predicate) Entry

// ToIntFunction that accepts one argument and produces an EntryInt result.
type ToIntFunction[E Entry] func(e E) EntryInt

// ToFloatFunction that accepts one argument and produces an EntryFloat result.
type ToFloatFunction[E Entry] func(e E) EntryFloat

// FlattenEntrySliceToEntry is a StreamFunction that flattens
// an EntrySlice to a Stream of its elements.
func FlattenEntrySliceToEntry[E Entry](bufsize int) StreamFunction[E] {
	return func(e E) Stream[E] {
		return NewStreamFromSlice(e.(EntrySlice), bufsize)
	}
}
