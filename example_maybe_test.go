package fuego_test

import (
	"fmt"
	ƒ "github.com/seborama/fuego"
	"github.com/stretchr/testify/assert"
)

// ExampleMaybe shows ways to use a Maybe.
func ExampleMaybe() {
	m1 := ƒ.MaybeOf(EntryString("Hello World"))
	fmt.Printf("m1.Get=%v\n", m1.Get())
	fmt.Printf("m1.GetOrElse=%v\n", m1.GetOrElse(EntryString("Bonjour le monde")))

	m2 := ƒ.MaybeOf(nil)
	if assert.PanicsWithValue(nil, ƒ.PanicNoSuchElement, func() {
		fmt.Printf("m2.Get=%v\n", m2.Get())
	}) {
		fmt.Println("m2.Get() panics with fuego.PanicNoSuchElement")
	}

	fmt.Printf("m2.GetOrElse=%v\n", m2.GetOrElse(EntryString("Bonjour le monde")))

	// Output:
	// m1.Get=Hello World
	// m1.GetOrElse=Hello World
	// m2.Get() panics with fuego.PanicNoSuchElement
	// m2.GetOrElse=Bonjour le monde
}
