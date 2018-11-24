package fuego

// Function that accepts one argument and produces a result.
type Function func(e Entry) Entry

// BiFunction that accepts two arguments and produces a result.
type BiFunction func(e1, e2 Entry) Entry
