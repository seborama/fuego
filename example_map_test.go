package fuego_test

import (
	"fmt"

	"github.com/raviqqe/hamt"
	"github.com/seborama/fuego"
)

// ExampleMap_EntrySet_Stream shows how to create a Map, add entries
// to it and stream filter the even-numbered keys out of it
// using a Predicate.
func ExampleMap_EntrySet_Stream() {
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
	// res={iterator:{slice:[{K:2 V:two} {K:4 V:four} {K:6 V:six} {K:8 V:eight}] size:4}}
}

// ExampleMap_KeySet_Stream shows how to create a Map, add entries to it and stream filter the
// even-numbered keys out of it using a Predicate.
func ExampleMap_KeySet_Stream() {
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

type EntryInt int

func (i EntryInt) Hash() uint32 {
	return uint32(i)
}

// TODO Call this FlatMap instead?
func (i EntryInt) Value() EntryInt {
	return i
}

func (i EntryInt) Equal(e hamt.Entry) bool {
	j, ok := e.(EntryInt)

	if !ok {
		return false
	}

	return i == j
}

func isEvenNumber(t interface{}) bool {
	k := (t.(fuego.MapEntry)).K.(EntryInt)
	return k.Value()&1 == 0
}

func isOddNumber(t interface{}) bool {
	v := t.(EntryInt)
	return v&1 == 0
}
