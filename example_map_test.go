package fuego_test

import (
	"fmt"

	"github.com/seborama/fuego"
)

// ExampleMap_EntrySet_Stream shows how to create a Map, add entries
// to it and stream filter the even-numbered keys out of it
// using a Predicate.
func ExampleMap_EntrySet_stream() {
	m := fuego.NewMap().
		Insert(EntryInt(1), "one").
		Insert(EntryInt(2), "two").
		Insert(EntryInt(3), "three").
		Insert(EntryInt(4), "four").
		Insert(EntryInt(5), "five").
		Insert(EntryInt(6), "six").
		Insert(EntryInt(7), "seven").
		Insert(EntryInt(8), "eight").
		Insert(EntryInt(9), "nine")

	res := m.EntrySet().
		Stream().
		Filter(isEvenNumber)
	fmt.Printf("res=%+v\n", res)

	// Output:
	// res={iterator:{slice:[{K:2 V:two} {K:4 V:four} {K:6 V:six} {K:8 V:eight}]}}
}

// ExampleMap_KeySet_Stream shows how to create a Map,
// add entries to it and stream filter the
// even-numbered keys out of it using a Predicate.
func ExampleMap_KeySet_stream() {
	m := fuego.NewMap().
		Insert(EntryInt(1), "one").
		Insert(EntryInt(2), "two").
		Insert(EntryInt(3), "three").
		Insert(EntryInt(4), "four").
		Insert(EntryInt(5), "five").
		Insert(EntryInt(6), "six").
		Insert(EntryInt(7), "seven").
		Insert(EntryInt(8), "eight").
		Insert(EntryInt(9), "nine")

	printEntry := func(value interface{}) {
		fmt.Printf("Number: %d\n", value.(EntryInt).Value())
	}

	m.KeySet().Stream().
		Filter(isOddNumber).
		ForEach(printEntry)

	// Output:
	// Number: 2
	// Number: 4
	// Number: 6
	// Number: 8
}
