package fuego

// Function that accepts one argument and produces a result.
type Function func(e Entry) Entry

// BiFunction that accepts two arguments and produces a result.
type BiFunction func(e1, e2 Entry) Entry

// TODO: PartialFunction? See https://www.scala-lang.org/api/current/scala/PartialFunction.html
// A partial function applies a function to a type if the value is within the domain of the partial function.
// For instance, a partial function could return the doubles of the odd numbers in a series of mixed odd and even numbers.

// ToIntFunction that accepts one argument and produces an EntryInt result.
type ToIntFunction func(e Entry) EntryInt
