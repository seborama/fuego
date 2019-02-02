package fuego

// A Tuple is a container of value(s).
// A special case is Tuple0 which does not hold any value.
type Tuple interface {
	Hash() uint32
	Equal(o Tuple) bool
}
