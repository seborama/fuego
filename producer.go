package fuego

// Supplier accepts no argument and returns an Entry.
type Supplier[E any] func() E
