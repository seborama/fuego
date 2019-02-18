package fuego

// Consumer that accepts one argument and does not
// return any value.
type Consumer func(i Entry)

// BiConsumer that accepts two arguments and does not
// return any value.
type BiConsumer func(i, j Entry)
