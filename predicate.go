package fuego

// Predicate represents a predicate (boolean-valued function) of one argument.
type Predicate func(t Entry) bool // TODO return EntryBool instead of bool??

// And is a composed predicate that represents a short-circuiting logical
// AND of this predicate and another.
func (p Predicate) And(other Predicate) Predicate {
	return func(t Entry) bool {
		if p == nil || other == nil {
			return false
		}
		return p(t) && other(t)
	}
}

// Or is a composed predicate that represents a short-circuiting logical
// OR of two predicates.
func (p Predicate) Or(other Predicate) Predicate {
	return func(t Entry) bool {
		if p == nil {
			p = False
		}
		if other == nil {
			return p(t)
		}
		return p(t) || other(t)
	}
}

// Xor is a composed predicate that represents a short-circuiting logical
// XOR of two predicates.
func (p Predicate) Xor(other Predicate) Predicate {
	return func(t Entry) bool {
		if p == nil {
			p = False
		}
		if other == nil {
			other = False
		}
		return p.Or(other).And(p.And(other).Negate())(t)
	}
}

// Negate is an alias for Not().
func (p Predicate) Negate() Predicate {
	return p.Not()
}

// Not is the logical negation of a predicate.
func (p Predicate) Not() Predicate {
	return func(t Entry) bool {
		return p == nil || !p(t)
	}
}

// FunctionPredicate creates a Predicate from a Function.
func FunctionPredicate(f Function) Predicate {
	return func(t Entry) bool {
		return bool(f(t).(EntryBool))
	}
}

// False is a predicate that returns always false.
func False(t Entry) bool {
	return false
}

// True is a predicate that returns always true.
func True(t Entry) bool {
	return Predicate(False).Negate()(nil)
}
