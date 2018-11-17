package fuego_test

import (
	"fmt"

	"github.com/seborama/fuego"
)

// ExampleFunction shows how to use Function's.
// There are more interesting examples through the code.
// Search for `Function` or the Function signature such as
// `func([^)]*interface{}) interface{}` (as a regexp).
func ExampleFunction() {
	timesTwoFunction := timesTwo()
	res := timesTwoFunction(7)
	fmt.Printf("res = %+v\n", res)

	// Output:
	// res = 14
}

// timesTwo returns a fuego.Function
func timesTwo() fuego.Function {
	return func(i interface{}) interface{} {
		return (EntryInt(2 * i.(int)))
	}
}
