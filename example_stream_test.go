package fuego_test

import (
	"fmt"
	ƒ "github.com/seborama/fuego"
)

// ExampleStream_GroupBy shows a use of Stream's with GroupBy.
func ExampleStream_GroupBy() {
	m := ƒ.NewOrderedMap().
		Insert(EntryInt(1), "one").
		Insert(EntryInt(2), "two").
		Insert(EntryInt(3), "three").
		Insert(EntryInt(4), "four").
		Insert(EntryInt(5), "five").
		Insert(EntryInt(6), "six").
		Insert(EntryInt(7), "seven").
		Insert(EntryInt(8), "eight").
		Insert(EntryInt(9), "nine")

	resMap := map[ƒ.Entry]interface{}{}
	m.EntrySet().
		Stream().
		GroupBy(func(i ƒ.Entry) ƒ.Entry {
			return i.(ƒ.MapEntry).K.(EntryInt) & 1
		}).
		EntrySet().
		Stream().
		ForEach(func(e ƒ.Entry) {
			resMap[e.(ƒ.MapEntry).K] = e.(ƒ.MapEntry).V
		})

	for i := 0; i < len(resMap); i++ {
		fmt.Printf("%d => %v\n", i, resMap[EntryInt(i)])
	}

	// Output:
	// 0 => {[{2 two} {4 four} {6 six} {8 eight}]}
	// 1 => {[{1 one} {3 three} {5 five} {7 seven} {9 nine}]}
}
