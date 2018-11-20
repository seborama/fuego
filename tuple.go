package fuego

type Tuple interface {
	Hash() uint32
	Equal(o Tuple) bool
	Arity() int
	ToSet() Set
	// Append, Concat, Swap, Update, UpdateN (Update1, Update2, ...)
}
