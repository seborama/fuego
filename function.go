package fuego

// Function that accepts one argument and produces a result.
type Function func(e Entry) Entry

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
