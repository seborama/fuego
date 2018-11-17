package fuego

// Predicate represents a predicate (boolean-valued function) of one argument.
type Predicate func(t interface{}) bool

// And is a composed predicate that represents a short-circuiting logical
// AND of this predicate and another.
func (p Predicate) And(other Predicate) Predicate {
	return func(t interface{}) bool {
		if p == nil || other == nil {
			return false
		}
		return p(t) && other(t)
	}
}

// Or is a composed predicate that represents a short-circuiting logical
// OR of two predicates.
func (p Predicate) Or(other Predicate) Predicate {
	return func(t interface{}) bool {
		if p == nil {
			p = False
		}
		if other == nil {
			return p(t)
		}
		return p(t) || other(t)
	}
}

// Not is the logical negation of a predicate.
func (p Predicate) Not() Predicate {
	return func(t interface{}) bool {
		return p == nil || !p(t)
	}
}

// FunctionPredicate creates a Predicate from a Function.
func FunctionPredicate(f Function) Predicate {
	return func(t interface{}) bool {
		return f(t).(bool)
	}
}

// False is a predicate that returns always false.
func False(t interface{}) bool {
	return false
}

// True is a predicate that returns always true.
func True(t interface{}) bool {
	return Predicate(False).Not()(nil)
}
