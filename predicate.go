package fuego

// Represents a predicate (boolean-valued function) of one argument.
type Predicate interface {
	// Test evaluates this predicate on the given argument.
	Test(t interface{}) bool
}

// FunctionPredicate is a Predicate based on a Function.
type FunctionPredicate struct {
	function Function
}

// NewFunctionPredicate creates a FunctionPredicate.
func NewFunctionPredicate(f Function) Predicate {
	return FunctionPredicate{
		function: f,
	}
}

// Test evaluates this predicate on the given argument.
func (p FunctionPredicate) Test(t interface{}) bool {
	return (p.function.Apply(t)).(bool)
}

// A FalsePredicate always returns 'false'.
type FalsePredicate struct{}

// False creates a FalsePredicate.
func False() Predicate {
	return FalsePredicate{}
}

// Test evaluates this predicate on the given argument.
func (p FalsePredicate) Test(t interface{}) bool {
	return false
}

// True creates a TruePredicate.
func True() Predicate {
	return Not(False())
}

// AndPredicate is a composed predicate that represents a short-circuiting logical
// AND of this predicate and another.
type AndPredicate struct {
	p1 Predicate
	p2 Predicate
}

// And creates a AndPredicate.
func And(p1, p2 Predicate) Predicate {
	if p1 == nil || p2 == nil {
		return False()
	}
	return AndPredicate{
		p1: p1,
		p2: p2,
	}
}

// Test evaluates this predicate on the given argument.
func (p AndPredicate) Test(t interface{}) bool {
	return p.p1.Test(t) && p.p2.Test(t)
}

// NotPredicate is te logical negation of a predicate.
type NotPredicate struct {
	predicate Predicate
}

// Not creates a NotPredicate.
func Not(predicate Predicate) Predicate {
	if predicate == nil {
		return True() // TODO is that correct??
	}
	return NotPredicate{predicate: predicate}
}

// Test evaluates this predicate on the given argument.
func (p NotPredicate) Test(t interface{}) bool {
	return !p.predicate.Test(t)
}

// OrPredicate is a composed predicate that represents a short-circuiting logical
// OR of two predicates.
type OrPredicate struct {
	p1 Predicate
	p2 Predicate
}

// Or creates a OrPredicate.
func Or(p1, p2 Predicate) Predicate {
	if p1 == nil {
		p1 = False()
	}
	if p2 == nil {
		p2 = False()
	}
	return OrPredicate{
		p1: p1,
		p2: p2,
	}
}

// Test evaluates this predicate on the given argument.
func (p OrPredicate) Test(t interface{}) bool {
	return p.p1.Test(t) || p.p2.Test(t)
}
