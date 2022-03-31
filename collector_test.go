package fuego

import (
	"hash/crc32"
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

func TestCollector_Collect_ToEntryMap(t *testing.T) {
	tt := map[string]struct {
		inputData     []employee
		expected      map[string]int
		expectedPanic string
	}{
		"panics when key exists": {
			inputData: []employee{
				{
					id:   1,
					name: "One",
				},
				{
					id:   1000,
					name: "One",
				},
			},
			expectedPanic: PanicDuplicateKey + ": 'One'",
		},
		"returns a map of employee (name, id)": {
			inputData: getEmployeesSample(),
			expected: map[string]int{
				"One":   1,
				"Two":   2,
				"Three": 3,
				"Four":  4,
				"Five":  5,
			},
		},
	}

	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			employeeNameByID := func() map[string]int {
				return Collect(
					NewStreamFromSlice(tc.inputData, 0),
					ToMap(employee.Name, employee.ID),
				)
			}

			if tc.expectedPanic != "" {
				assert.PanicsWithValue(t, tc.expectedPanic, func() { _ = employeeNameByID() })
				return
			}

			assert.EqualValues(t, tc.expected, employeeNameByID())
		})
	}
}

func TestCollector_Collect_ToEntryMapWithKeyMerge(t *testing.T) {
	employees := getEmployeesSample()

	overwriteKeyMergeFn := func(v1, v2 int) int {
		return v2
	}

	employeeNameByID :=
		Collect(
			NewStreamFromSlice(employees, 0),
			ToMapWithMerge(employee.Department, employee.ID, overwriteKeyMergeFn))

	expected := map[string]int{
		"HR":        5,
		"IT":        3,
		"Marketing": 1,
	}

	assert.EqualValues(t, expected, employeeNameByID)
}

func TestIdentityFinisher(t *testing.T) {
	tests := []struct {
		name    string
		element any
		want    any
	}{
		{
			name:    "Should return identity for nil",
			element: nil,
			want:    nil,
		},
		{
			name:    "Should return identity for a given simple Entry",
			element: "~å∫√çß∆",
			want:    "~å∫√çß∆",
		},
		{
			name: "Should return identity for a given complex Entry",
			element: map[int][]any{
				1: {true},
				2: {true, "abc"},
			},
			want: map[int][]any{
				1: {true},
				2: {true, "abc"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IdentityFinisher(tt.element)
			assert.Equal(t, tt.want, got)
		})
	}
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

func TestCollector_GroupingBy_Mapping_FlatMapping_Filtering_Mapping_Reducing(t *testing.T) {
	stringLength :=
		func(e string) int {
			return len(e)
		}

	toStringList :=
		func(e string) []string {
			r := []string{}
			for _, c := range e {
				r = append(r, string(c))
			}
			return r
		}

	flattenStringListToDistinct :=
		func(e []string) Stream[string] {
			return NewStreamFromSlice(e, 0).
				Distinct(func(s string) uint32 { return crc32.ChecksumIEEE([]byte(s)) })
		}

	stringToUpper :=
		func(e string) string {
			return strings.ToUpper(e)
		}

	concatenateStringsBiFunc := func(i, j string) string {
		iStr := i
		jStr := j
		return iStr + jStr
	}

	strs := []string{
		"a",
		"bb",
		"cc",
		"ee",
		"ddd",
	}

	got :=
		Collect(
			NewStreamFromSlice(strs, 0),
			GroupingBy(
				stringLength,
				Mapping(
					toStringList,
					FlatMapping(flattenStringListToDistinct,
						Mapping(stringToUpper,
							Reducing(concatenateStringsBiFunc),
						),
					),
				),
			),
		)

	expected := map[int]string{
		1: "A",
		2: "BCE",
		3: "D",
	}

	assert.EqualValues(t, expected, got)
}

type employee struct {
	id         int
	name       string
	department string
	salary     float32
}

func (e employee) ID() int {
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
