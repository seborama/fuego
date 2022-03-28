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

func TestCollector_Filtering(t *testing.T) {
	employees := getEmployeesSample()

	highestPaidEmployeesByDepartment :=
		Collect(
			NewStreamFromSlice(employees, 0),
			GroupingBy(employee.Department,
				Filtering(func(e employee) bool {
					return e.Salary() > 2000
				},
					ToSlice[employee]())))

	expected := map[string][]employee{
		"HR": {
			{
				id:         5,
				name:       "Five",
				department: "HR",
				salary:     2300,
			}},
		"IT": {
			{
				id:         2,
				name:       "Two",
				department: "IT",
				salary:     2500,
			},
			{
				id:         3,
				name:       "Three",
				department: "IT",
				salary:     2200,
			}},
		"Marketing": {},
	}

	assert.EqualValues(t, expected, highestPaidEmployeesByDepartment)
}

type employee struct {
	id         uint32
	name       string
	department string
	salary     float32
}

func (e employee) ID() uint32 {
	return e.id
}

func (e employee) Name() string {
	return e.name
}

func (e employee) Department() string {
	return e.department
}

func (e employee) Salary() float32 {
	return e.salary
}

func getEmployeesSample() []employee {
	return []employee{
		{
			id:         1,
			name:       "One",
			department: "Marketing",
			salary:     1500,
		},
		{
			id:         2,
			name:       "Two",
			department: "IT",
			salary:     2500,
		},
		{
			id:         3,
			name:       "Three",
			department: "IT",
			salary:     2200,
		},
		{
			id:         4,
			name:       "Four",
			department: "HR",
			salary:     1800,
		},
		{
			id:         5,
			name:       "Five",
			department: "HR",
			salary:     2300,
		},
	}
}
