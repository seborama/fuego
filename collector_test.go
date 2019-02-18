package fuego

import (
	"log"
	"testing"
)

func TestNewCollector(t *testing.T) {
	strs := []Entry{
		EntryString("a"),
		EntryString("bb"),
		EntryString("cc"),
		EntryString("ddd"),
	}

	NewStreamFromSlice(strs, 1e3).
		Map(func(e Entry) Entry {
			return e.(EntryString).ToUpper()
		}).
		Filter(func(e Entry) bool {
			return e.(EntryString).Len() > 1
		}).
		GroupBy(
			func(e Entry) Entry {
				return e.(EntryString).Len()
			},
		).
		Stream(1e3).
		ForEach(func(e Entry) {
			log.Printf("DEBUG - ForEach: %+v\n", e) // {1=[], 2=[BB, CC], 3=[DDD]}
		})

	stringLength := func(e Entry) Entry {
		t2 := Tuple2{
			E1: e.(EntryString).Len(),
			E2: e,
		}
		return t2
	}

	stringToUpper := func(e Entry) Entry {
		return Tuple2{
			E1: e.(Tuple2).E1,
			E2: e.(Tuple2).E2.(EntryString).ToUpper(),
		}
	}

	result := NewStreamFromSlice(strs, 1e3).
		Collect(
			GroupingBy(
				stringLength,
				Mapping(
					stringToUpper,
					ToEntryMap())))
	log.Printf("DEBUG - result: %+v\n", result) // {1=[], 2=[BB, CC], 3=[DDD]}

}
