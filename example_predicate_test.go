package fuego_test

import (
	"fmt"

	"github.com/seborama/fuego" // using ƒ as a short alias of fuego
)

// ExamplePredicate shows how to use and combine Predicates.
func ExamplePredicate() {
	res := fuego.Predicate(fuego.False).Not()(1)
	fmt.Printf("Not False == %+v\n", res)

	res = fuego.Predicate(fuego.True).And(fuego.False)(1)
	fmt.Printf("True and False == %+v\n", res)

	res = fuego.Predicate(fuego.True).Or(fuego.False)(1)
	fmt.Printf("True or False == %+v\n", res)

	// You can use associativity too - part 1 of 2:
	// False And False Or True == true
	res = fuego.Predicate(fuego.False).And(fuego.False).Or(fuego.True)(1)
	fmt.Printf("False And False Or True == %+v\n", res)

	// You can use associativity too - part 2 of 2:
	// False And (False Or True) == false
	res = fuego.Predicate(fuego.False).And(fuego.Predicate(fuego.False).Or(fuego.True))(1)
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
	isEvenNumberPredicate := fuego.FunctionPredicate(isEvenNumberFunction)

	res := fuego.Predicate(isEvenNumberPredicate).And(fuego.True)(23)
	fmt.Printf("res = %v", res)

	// Output:
	// res = false
}

// ExamplePredicate_custom2 shows how to create a custom Predicate from
// scratch.
// Notice how we get all Predicate helpers (And, Or, Not, etc) for "free".
func ExamplePredicate_predicate() {
	res := fuego.Predicate(intGreaterThanPredicate(50)).And(fuego.True).Not()(23)
	fmt.Printf("res = %v", res)

	// Output:
	// res = true
}
