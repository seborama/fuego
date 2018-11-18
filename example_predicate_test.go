package fuego_test

import (
	"fmt"

	ƒ "github.com/seborama/fuego" // using ƒ as a short alias of fuego
)

// ExamplePredicate shows how to use and combine Predicates.
func ExamplePredicate() {
	res := ƒ.Predicate(ƒ.False).Not()(1)
	fmt.Printf("Not False == %+v\n", res)

	res = ƒ.Predicate(ƒ.True).And(ƒ.False)(1)
	fmt.Printf("True and False == %+v\n", res)

	res = ƒ.Predicate(ƒ.True).Or(ƒ.False)(1)
	fmt.Printf("True or False == %+v\n", res)

	// You can use associativity too - part 1 of 2:
	// False And False Or True == true
	res = ƒ.Predicate(ƒ.False).And(ƒ.False).Or(ƒ.True)(1)
	fmt.Printf("False And False Or True == %+v\n", res)

	// You can use associativity too - part 2 of 2:
	// False And (False Or True) == false
	res = ƒ.Predicate(ƒ.False).And(ƒ.Predicate(ƒ.False).Or(ƒ.True))(1)
	fmt.Printf("False And (False Or True) == %+v\n", res)

	// Output:
	// Not False == true
	// True and False == false
	// True or False == true
	// False And False Or True == true
	// False And (False Or True) == false
}

// ExamplePredicate_custom1 shows how to create a custom Predicate using
// the utility function fuego.FunctionPredicate().
func ExamplePredicate_functionPredicate() {
	isEvenNumberPredicate := ƒ.FunctionPredicate(isEvenNumberFunction)

	res := ƒ.Predicate(isEvenNumberPredicate).And(ƒ.True)(23)
	fmt.Printf("res = %v", res)

	// Output:
	// res = false
}

// ExamplePredicate_custom2 shows how to create a custom Predicate from
// scratch.
// Notice how we get all Predicate helpers (And, Or, Not, etc) for "free".
func ExamplePredicate_predicate() {
	res := ƒ.Predicate(intGreaterThanPredicate(50)).And(ƒ.True).Not()(23)
	fmt.Printf("res = %v", res)

	// Output:
	// res = true
}
