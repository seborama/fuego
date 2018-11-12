package fuego

// Function that accepts one argument and produces a result.
type Function interface {
	Apply(i interface{}) interface{}
}
