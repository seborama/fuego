package fuego

// Function that accepts one argument and produces a result.
type Function func(e Entry) Entry

// StreamFunction that accepts one argument and produces a stream.
//
// It is worth noting is `StreamFunction` which accepts one argument
//  and produces a stream. It is used with when "flat mapping" and
// `Stream`. This effectively is a one to many operation, such
// as exploding the individual values of an EntrySlice into a Stream.
//
// This effectively is a one to many operation, such as exploding
// the individual values of an EntrySlice into a Stream.
type StreamFunction func(e Entry) Stream

// BiFunction that accepts two arguments and produces a result.
type BiFunction func(e1, e2 Entry) Entry

// A PartialFunction applies a function to a type if the value
// is within the domain of the partial function.
// For instance, a partial function could return the doubles
// of the odd numbers in a series of mixed odd and even numbers.
// type PartialFunction func(e Entry, p Predicate) Entry

// ToIntFunction that accepts one argument and produces an EntryInt result.
type ToIntFunction func(e Entry) EntryInt

// ToFloatFunction that accepts one argument and produces an EntryFloat result.
type ToFloatFunction func(e Entry) EntryFloat

// FlattenEntrySliceToEntry is a StreamFunction that flattens
// an EntrySlice to a Stream of its elements.
func FlattenEntrySliceToEntry(bufsize int) StreamFunction {
	return func(e Entry) Stream {
		return NewStreamFromSlice(e.(EntrySlice), bufsize)
	}
}
