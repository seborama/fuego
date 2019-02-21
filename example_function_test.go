package fuego_test

import (
	"fmt"

	. "github.com/seborama/fuego"
)

// ExampleFunction shows how to use Function's.
// There are more interesting examples through the code.
// Search for `Function` or the Function signature.
func ExampleFunction() {
	timesTwoFunction := timesTwo()
	res := timesTwoFunction(EntryInt(7))
	fmt.Printf("res = %+v\n", res)

	// Output:
	// res = 14
}

// ExampleBiFunction shows how to use BiFunction's.
// There are more interesting examples through the code.
// Search for `BiFunction` or the BiFunction signature.
func ExampleBiFunction() {
	data := []Entry{
		EntryString("four"),
		EntryString("twelve"),
		EntryString("one"),
		EntryString("six"),
		EntryString("three")}

	res := NewStreamFromSlice(data, 0).
		Reduce(concatenateStringsBiFunc)

	fmt.Printf("res = %+v\n", res)

	// Output:
	// res = four-twelve-one-six-three
}
