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

	result := NewStreamFromSlice(strs, 1e3).
		Collect(GroupingBy(func(e Entry) Entry {
			return Tuple2{
				E1: e.(EntryString).Len(),
				E2: e.(EntryString),
			}
		}))
	log.Printf("DEBUG - result: %+v\n", result) // {1=[], 2=[BB, CC], 3=[DDD]}

}

func GroupingBy(classifier Function) Collector {
	supplier := func() Entry { // use chan Entry instead?
		return EntryMap{}
	}
	accumulator := func(supplier Entry, entry Entry) Entry {
		t2 := classifier(entry).(Tuple2)
		return supplier.(EntryMap).Append(t2)
	}
	finisher := func(e Entry) Entry {
		return e
	}
	return NewCollector(supplier, accumulator, finisher)
}
