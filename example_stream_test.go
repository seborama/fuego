package fuego_test

import (
	"fmt"
	ƒ "github.com/seborama/fuego"
)

// ExampleStream_GroupBy shows a use of Stream's with GroupBy.
func ExampleStream_GroupBy() {
	data := []ƒ.Entry{
		ƒ.Tuple2{E1: EntryInt(1), E2: EntryString("one")},
		ƒ.Tuple2{E1: EntryInt(2), E2: EntryString("two")},
		ƒ.Tuple2{E1: EntryInt(3), E2: EntryString("three")},
		ƒ.Tuple2{E1: EntryInt(4), E2: EntryString("four")},
		ƒ.Tuple2{E1: EntryInt(5), E2: EntryString("five")},
		ƒ.Tuple2{E1: EntryInt(6), E2: EntryString("six")},
		ƒ.Tuple2{E1: EntryInt(7), E2: EntryString("seven")},
		ƒ.Tuple2{E1: EntryInt(8), E2: EntryString("eight")},
		ƒ.Tuple2{E1: EntryInt(9), E2: EntryString("nine")}}

	resMap := map[ƒ.Entry]interface{}{}
	ƒ.NewStreamFromSlice(data).
		GroupBy(func(i ƒ.Entry) ƒ.Entry {
			return i.(ƒ.Tuple2).E1.(EntryInt) & 1
		}).
		Stream().
		ForEach(func(e ƒ.Entry) {
			resMap[e.(ƒ.Tuple2).E1] = e.(ƒ.Tuple2).E2
		})

	for i := 0; i < len(resMap); i++ {
		fmt.Printf("%d => %v\n", i, resMap[EntryInt(i)])
	}

	// Output:
	// 0 => [{2 two} {4 four} {6 six} {8 eight}]
	// 1 => [{1 one} {3 three} {5 five} {7 seven} {9 nine}]
}
