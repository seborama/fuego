package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollector_GroupingBy_Mapping_Filtering_ToEntrySlice(t *testing.T) {
	strs := EntrySlice{
		EntryString("a"),
		EntryString("bb"),
		EntryString("cc"),
		EntryString("ddd"),
	}

	stringLength := func(e Entry) Entry {
		return e.(EntryString).Len()
	}

	stringToUpper := func(e Entry) Entry {
		return e.(EntryString).ToUpper()
	}

	stringLengthGreaterThan := func(i int) Predicate {
		return func(e Entry) bool {
			return int(e.(EntryString).Len()) > i
		}
	}

	got := NewStreamFromSlice(strs, 1e3).
		Collect(
			GroupingBy(
				stringLength,
				Mapping(
					stringToUpper,
					Filtering(
						stringLengthGreaterThan(1),
						ToEntrySlice()))))

	expected := EntryMap{
		EntryInt(1): EntrySlice{},
		EntryInt(2): EntrySlice{
			EntryString("BB"),
			EntryString("CC")},
		EntryInt(3): EntrySlice{
			EntryString("DDD")},
	}

	assert.EqualValues(t, expected, got)
}
