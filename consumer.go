package fuego

// Consumer that accepts one argument and does not
// return any value.
type Consumer[T any] func(T)

// BiConsumer that accepts two arguments and does not
// return any value.
type BiConsumer[T, U any] func(T, U)
