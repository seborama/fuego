package fuego_test

import (
	"fmt"

	"github.com/seborama/fuego"
)

// ExampleFunction shows how to use Function's.
// There are more interesting examples through the code.
// Search for `Function` or the Function signature.
func ExampleFunction() {
	timesTwoFunction := timesTwo()
	res := timesTwoFunction(7)
	fmt.Printf("res = %+v\n", res)

	// Output:
	// res = 14
}

// ExampleBiFunction shows how to use BiFunction's.
// There are more interesting examples through the code.
// Search for `BiFunction` or the BiFunction signature.
func ExampleBiFunction() {
	res := fuego.NewHamtSet().
		Insert(EntryString("four")).
		Insert(EntryString("twelve")).
		Insert(EntryString("one")).
		Insert(EntryString("six")).
		Insert(EntryString("three")).
		Stream().
		Reduce(concatenateStringsBiFunc)

	fmt.Printf("res = %+v\n", res)

	// Output:
	// res = one-three-twelve-six-four
}
