package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollector_Collect_GroupingBy(t *testing.T) {
	stringLength := func(e Entry) Entry {
		return e.(EntryString).Len()
	}
	collectorWithNilFinisher := ToEntrySlice()
	collectorWithNilFinisher.finisher = nil

	strs := EntrySlice{}
	got := NewStreamFromSlice(strs, 0).
		Collect(
			GroupingBy(
				stringLength,
				collectorWithNilFinisher))

	assert.Equal(t, EntryMap{}, got)
}

func TestCollector_GroupingBy_Mapping_Filtering_ToEntrySlice(t *testing.T) {
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

	strs := EntrySlice{
		EntryString("a"),
		EntryString("bb"),
		EntryString("cc"),
		EntryString("ddd"),
	}
	got := NewStreamFromSlice(strs, 0).
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

func TestCollector_GroupingBy_Mapping_FlatMapping_Filtering_Mapping_Reducing(t *testing.T) {
	stringLength :=
		func(e Entry) Entry {
			return e.(EntryString).Len()
		}

	toStringList :=
		func(e Entry) Entry {
			return e.(EntryString).MapToEntryBytes(0).
				Map(func(e Entry) Entry {
					return EntryString(string(e.(EntryByte)))
				}).
				Collect(ToEntrySlice()).(EntrySlice)
		}

	flattenStringListToDistinct :=
		func(e Entry) Stream {
			return NewStreamFromSlice(e.(EntrySlice), 0).Distinct()
		}

	stringToUpper :=
		func(e Entry) Entry {
			return e.(EntryString).ToUpper()
		}

	concatenateStringsBiFunc := func(i, j Entry) Entry {
		iStr := i.(EntryString)
		jStr := j.(EntryString)
		return iStr + jStr
	}

	strs := EntrySlice{
		EntryString("a"),
		EntryString("bb"),
		EntryString("cc"),
		EntryString("ee"),
		EntryString("ddd"),
	}

	got := NewStreamFromSlice(strs, 0).
		Collect(
			GroupingBy(
				stringLength,
				Mapping(
					toStringList,
					FlatMapping(flattenStringListToDistinct,
						Mapping(stringToUpper,
							Reducing(concatenateStringsBiFunc))))))

	expected := EntryMap{
		EntryInt(1): EntryString("A"),
		EntryInt(2): EntryString("BCE"),
		EntryInt(3): EntryString("D"),
	}

	assert.EqualValues(t, expected, got)
}

func TestCollector_Collect_Reducing(t *testing.T) {
	s := NewIntStreamFromSlice([]EntryInt{5, 10, 20, 50}, 0)

	got := s.Collect(
		Reducing(
			func(integer, integer2 Entry) Entry {
				return integer2.(EntryInt) - integer.(EntryInt)
			}))

	assert.Equal(t, EntryInt(35), got)
}

func TestCollector_Collect_ToEntryMap(t *testing.T) {
	tt := map[string]struct {
		inputData     EntrySlice
		expected      EntryMap
		expectedPanic string
	}{
		"panics when key exists": {
			inputData: EntrySlice{
				employee{
					id:   1,
					name: "One",
				},
				employee{
					id:   1000,
					name: "One",
				},
			},
			expectedPanic: PanicDuplicateKey + ": 'One'",
		},
		"returns a map of employee (name, id)": {
			inputData: getEmployeesSample(),
			expected: EntryMap{
				EntryString("One"):   EntryInt(1),
				EntryString("Two"):   EntryInt(2),
				EntryString("Three"): EntryInt(3),
				EntryString("Four"):  EntryInt(4),
				EntryString("Five"):  EntryInt(5),
			},
		},
	}

	for name, tc := range tt {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			employeeNameByID := func() EntryMap {
				return NewStreamFromSlice(tc.inputData, 0).
					Collect(
						ToEntryMap(employeeName, employeeID)).(EntryMap)
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

	overwriteKeyMergeFn := func(v1, v2 Entry) Entry {
		return v2
	}

	employeeNameByID :=
		NewStreamFromSlice(employees, 0).
			Collect(
				ToEntryMapWithKeyMerge(employeeDepartment, employeeID, overwriteKeyMergeFn))

	expected := EntryMap{
		EntryString("HR"):        EntryInt(5),
		EntryString("IT"):        EntryInt(3),
		EntryString("Marketing"): EntryInt(1),
	}

	assert.EqualValues(t, expected, employeeNameByID)
}

func TestIdentityFinisher(t *testing.T) {
	type args struct {
		e Entry
	}
	tests := []struct {
		name string
		args args
		want Entry
	}{
		{
			name: "Should return identity for nil",
			args: args{e: nil},
			want: nil,
		},
		{
			name: "Should return identity for a given simple Entry",
			args: args{e: EntryString("~å∫√çß∆")},
			want: EntryString("~å∫√çß∆"),
		},
		{
			name: "Should return identity for a given complex Entry",
			args: args{e: EntryMap{
				EntryString("1"): EntrySlice{EntryBool(true)},
				EntryBool(true):  EntrySlice{EntryBool(true), EntryString("abc")},
			}},
			want: EntryMap{
				EntryString("1"): EntrySlice{EntryBool(true)},
				EntryBool(true):  EntrySlice{EntryBool(true), EntryString("abc")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IdentityFinisher(tt.args.e)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCollector_Filtering(t *testing.T) {
	employees := getEmployeesSample()

	highestPaidEmployeesByDepartment :=
		NewStreamFromSlice(employees, 0).Collect(
			GroupingBy(employeeDepartment,
				Filtering(func(e Entry) bool {
					return e.(employee).Salary() > 2000
				},
					ToEntrySlice())))

	expected := EntryMap{
		EntryString("HR"): EntrySlice{
			employee{
				id:         5,
				name:       "Five",
				department: "HR",
				salary:     2300,
			}},
		EntryString("IT"): EntrySlice{
			employee{
				id:         2,
				name:       "Two",
				department: "IT",
				salary:     2500,
			},
			employee{
				id:         3,
				name:       "Three",
				department: "IT",
				salary:     2200,
			}},
		EntryString("Marketing"): EntrySlice{},
	}

	assert.EqualValues(t, expected, highestPaidEmployeesByDepartment)
}

var (
	salary1 float32 = 10.0
	salary2 float32 = 20.0
	salary3 float32 = 30.0
	salary4 float32 = 40.0
	salary5 float32 = 50.0
)

func getEmployeesSample() EntrySlice {
	return EntrySlice{
		employee{
			id:         1,
			name:       "One",
			department: "Marketing",
			salary:     1500,
		},
		employee{
			id:         2,
			name:       "Two",
			department: "IT",
			salary:     2500,
		},
		employee{
			id:         3,
			name:       "Three",
			department: "IT",
			salary:     2200,
		},
		employee{
			id:         4,
			name:       "Four",
			department: "HR",
			salary:     1800,
		},
		employee{
			id:         5,
			name:       "Five",
			department: "HR",
			salary:     2300,
		},
	}
}

type employee struct {
	id         uint32
	name       string
	department string
	salary     float32
}

func (e employee) Hash() uint32 {
	return e.id
}

func (e employee) Equal(other Entry) bool {
	if other == nil {
		return false
	}
	_, ok := other.(employee)
	return ok && (e.Hash() == other.Hash())
}

func (e employee) ID() EntryInt {
	return EntryInt(e.id)
}

func (e employee) Name() EntryString {
	return EntryString(e.name)
}

func (e employee) Department() EntryString {
	return EntryString(e.department)
}

func (e employee) Salary() EntryFloat {
	return EntryFloat(e.salary)
}

var employeeID = func(e Entry) Entry {
	return e.(employee).ID()
}

var employeeName = func(e Entry) Entry {
	return e.(employee).Name()
}

var employeeDepartment = func(e Entry) Entry {
	return e.(employee).Department()
}

var employeeSalary = func(e Entry) Entry {
	return e.(employee).Salary()
}
