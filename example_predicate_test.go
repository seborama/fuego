package fuego_test

import (
	"fmt"

	ƒ "github.com/seborama/fuego" // using ƒ as a short alias of fuego
)

// ExamplePredicates shows how to create a Map, add entries
// to it and stream filter the even-numbered keys out of it
// using a Predicate.
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
