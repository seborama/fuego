package fuego

// Predicate represents a predicate (boolean-valued function) of one argument.
type Predicate[E Entry] func(t E) EntryBool[E] // TODO return EntryBool instead of bool??

// And is a composed predicate that represents a short-circuiting logical
// AND of this predicate and another.
func (p Predicate[E]) And(other Predicate[E]) Predicate[E] {
	return func(t E) EntryBool[E] {
		if p == nil || other == nil {
			return false
		}
		return p(t) && other(t)
	}
}

// Or is a composed predicate that represents a short-circuiting logical
// OR of two predicates.
func (p Predicate[E]) Or(other Predicate[E]) Predicate[E] {
	return func(t E) EntryBool[E] {
		if p == nil {
			p = False[E]
		}
		if other == nil {
			return p(t)
		}
		return p(t) || other(t)
	}
}

// Xor is a composed predicate that represents a short-circuiting logical
// XOR of two predicates.
func (p Predicate[E]) Xor(other Predicate[E]) Predicate[E] {
	return func(t E) EntryBool[E] {
		return p.Or(other).And(p.And(other).Negate())(t)
	}
}

// Negate is an alias for Not().
func (p Predicate[E]) Negate() Predicate[E] {
	return p.Not()
}

// Not is the logical negation of a predicate.
func (p Predicate[E]) Not() Predicate[E] {
	return func(t E) EntryBool[E] {
		return p == nil || !p(t)
	}
}

// FunctionPredicate creates a Predicate from a Function.
func FunctionPredicate[E EntryBool](f Function[E]) Predicate[E] {
	return func(t E) EntryBool[E] {
		return f(t)
	}
}

// False is a predicate that returns always false.
func False[E Entry](t E) EntryBool[E] {
	return false
}

// True is a predicate that returns always true.
func True[E Entry](t E) EntryBool[E] {
	return Predicate[E](False[E]).Negate()(t)
}
