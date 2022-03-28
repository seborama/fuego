package fuego

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollector_GroupingBy_Mapping_Filtering_ToEntrySlice(t *testing.T) {
	stringLength := func(el string) int {
		return len(el)
	}

	stringToUpper := func(el string) string {
		return strings.ToUpper(el)
	}

	stringLengthGreaterThan := func(length int) Predicate[string] {
		return func(el string) bool {
			return len(el) > length
		}
	}

	strs := []string{
		"a",
		"bb",
		"cc",
		"ddd",
	}

	got :=
		Collect(
			NewStreamFromSlice(strs, 0),
			GroupingBy(
				stringLength,
				Mapping(
					stringToUpper,
					Filtering(
						stringLengthGreaterThan(1),
						ToSlice[string](),
					),
				),
			),
		)

	expected := map[int][]string{
		1: {},
		2: {
			"BB",
			"CC"},
		3: {
			"DDD"},
	}

	assert.EqualValues(t, expected, got)
}
