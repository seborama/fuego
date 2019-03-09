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
					return EntryString(byte(e.(EntryByte)))
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
