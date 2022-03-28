package fuego

// Supplier accepts no argument and returns a T.
type Supplier[T any] func() T
