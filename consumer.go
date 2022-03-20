package fuego

// Consumer that accepts one argument and does not
// return any value.
type Consumer[E Entry[E]] func(i E)

// BiConsumer that accepts two arguments and does not
// return any value.
type BiConsumer[E Entry[E]] func(i, j E)
