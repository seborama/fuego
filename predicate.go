package fuego

// Predicate represents a predicate (boolean-valued function) of one argument.
type Predicate[T any] func(t T) bool // TODO return EntryBool instead of bool??

// And is a composed predicate that represents a short-circuiting logical
// AND of this predicate and another.
func (p Predicate[T]) And(other Predicate[T]) Predicate[T] {
	return func(t T) bool {
		if p == nil || other == nil {
			return false
		}
		return p(t) && other(t)
	}
}

// Or is a composed predicate that represents a short-circuiting logical
// OR of two predicates.
func (p Predicate[T]) Or(other Predicate[T]) Predicate[T] {
	return func(t T) bool {
		if p == nil {
			p = False[T]
		}
		if other == nil {
			return p(t)
		}
		return p(t) || other(t)
	}
}

// Xor is a composed predicate that represents a short-circuiting logical
// XOR of two predicates.
func (p Predicate[T]) Xor(other Predicate[T]) Predicate[T] {
	return func(t T) bool {
		return p.Or(other).And(p.And(other).Negate())(t)
	}
}

// Negate is an alias for Not().
func (p Predicate[T]) Negate() Predicate[T] {
	return p.Not()
}

// Not is the logical negation of a predicate.
func (p Predicate[T]) Not() Predicate[T] {
	return func(t T) bool {
		return p == nil || !p(t)
	}
}

// FunctionPredicate creates a Predicate from a Function.
func FunctionPredicate[T any](f Function[T, bool]) Predicate[T] {
	return func(t T) bool {
		return f(t)
	}
}

// False is a predicate that returns always false.
func False[T any](t T) bool {
	return false
}

// True is a predicate that returns always true.
func True[T any](t T) bool {
	// return Predicate[T](False[T]).Negate()(t)
	return true
}
