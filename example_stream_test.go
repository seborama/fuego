package fuego_test

import (
	"fmt"

	. "github.com/seborama/fuego/v10"
)

// ExampleStream_GroupBy shows a use of Stream's with GroupBy.
func ExampleStream_GroupBy() {
	data := EntrySlice{
		Tuple2{E1: EntryInt(1), E2: EntryString("one")},
		Tuple2{E1: EntryInt(2), E2: EntryString("two")},
		Tuple2{E1: EntryInt(3), E2: EntryString("three")},
		Tuple2{E1: EntryInt(4), E2: EntryString("four")},
		Tuple2{E1: EntryInt(5), E2: EntryString("five")},
		Tuple2{E1: EntryInt(6), E2: EntryString("six")},
		Tuple2{E1: EntryInt(7), E2: EntryString("seven")},
		Tuple2{E1: EntryInt(8), E2: EntryString("eight")},
		Tuple2{E1: EntryInt(9), E2: EntryString("nine")}}

	resMap := map[Entry]interface{}{}
	NewStreamFromSlice(data, 0).
		GroupBy(func(i Entry) Entry {
			return i.(Tuple2).E1.(EntryInt) & 1
		}).
		Stream(0).
		ForEach(func(e Entry) {
			resMap[e.(Tuple2).E1] = e.(Tuple2).E2
		})

	for i := 0; i < len(resMap); i++ {
		fmt.Printf("%d => %v\n", i, resMap[EntryInt(i)])
	}

	// Output:
	// 0 => [{2 two} {4 four} {6 six} {8 eight}]
	// 1 => [{1 one} {3 three} {5 five} {7 seven} {9 nine}]
}

func ExampleStream_FlatMap() {
	a := EntrySlice{EntryInt(1), EntryInt(2), EntryInt(3)}
	b := EntrySlice{EntryInt(4), EntryInt(5)}
	c := EntrySlice{EntryInt(6), EntryInt(7), EntryInt(8)}

	sliceOfEntrySlicesOfEntryInts := EntrySlice{a, b, c}

	fmt.Printf("Before flattening: %+v\n", sliceOfEntrySlicesOfEntryInts)

	sliceOfEntryInts := NewStreamFromSlice(sliceOfEntrySlicesOfEntryInts, 0).
		FlatMap(FlattenEntrySliceToEntry(0)).
		Collect(ToEntrySlice())

	fmt.Printf("After flattening: %+v\n", sliceOfEntryInts)

	// Output:
	// Before flattening: [[1 2 3] [4 5] [6 7 8]]
	// After flattening: [1 2 3 4 5 6 7 8]
}
