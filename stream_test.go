package fuego

import (
	"reflect"
	"testing"

	"github.com/raviqqe/hamt"
	"github.com/stretchr/testify/assert"
)

func functionTimesTwo() Function {
	return func(i interface{}) interface{} {
		num := i.(EntryInt).Value()
		return interface{}(2 * num)
	}
}

func TestReferenceStream_Map(t *testing.T) {
	type fields struct {
		iterator Iterator
	}
	type args struct {
		mapper Function
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Stream
	}{
		{
			name:   "Should return a Stream of nil",
			fields: fields{iterator: nil},
			args: args{
				mapper: functionTimesTwo(),
			},
			want: NewStream(
				NewEntrySliceIterator([]hamt.Entry{})),
		},
		{
			name: "Should return a Stream of one double",
			fields: fields{
				iterator: NewSetIterator(NewHamtSet().
					Insert(EntryInt(4)))},
			args: args{
				mapper: functionTimesTwo(),
			},
			want: NewStream(
				NewEntrySliceIterator([]hamt.Entry{EntryInt(8)})),
		},
		{
			name: "Should return a Stream of 3 doubles",
			fields: fields{
				iterator: NewSetIterator(NewHamtSet().
					Insert(EntryInt(1)).
					Insert(EntryInt(2)).
					Insert(EntryInt(3)))},
			args: args{
				mapper: functionTimesTwo(),
			},
			want: NewStream(
				NewEntrySliceIterator([]hamt.Entry{
					EntryInt(2),
					EntryInt(4),
					EntryInt(6)})),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := ReferenceStream{
				iterator: tt.fields.iterator,
			}
			_ = rp.iterator
			if got := rp.Map(tt.args.mapper); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReferenceStream.Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func entryIntEqualsTo(number EntryInt) Function {
	return func(subject interface{}) interface{} {
		subjectEntryInt, ok := subject.(EntryInt)
		if !ok {
			return false
		}

		return number.Equal(subjectEntryInt)
	}
}

func TestReferenceStream_Filter(t *testing.T) {
	type fields struct {
		iterator Iterator
	}
	type args struct {
		predicate Predicate
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Stream
	}{
		{
			name: "Should give {1,3} from {4,1,3} when Predicate wants 1 or 3",
			fields: fields{
				iterator: NewSetIterator(NewHamtSet().
					Insert(EntryInt(4)).
					Insert(EntryInt(17)).
					Insert(EntryInt(3)))},
			args: args{
				predicate: FunctionPredicate(entryIntEqualsTo(EntryInt(17))).
					Or(FunctionPredicate(entryIntEqualsTo(EntryInt(3)))),
			},
			want: NewStream(
				NewEntrySliceIterator([]hamt.Entry{
					EntryInt(3),
					EntryInt(17)})),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := ReferenceStream{
				iterator: tt.fields.iterator,
			}
			if got := rp.Filter(tt.args.predicate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReferenceStream.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReferenceStream_ForEach(t *testing.T) {
	total := 0
	computeSumTotal := func(value interface{}) {
		total += int(value.(EntryInt).Value())
	}

	iterator := NewSetIterator(NewHamtSet().
		Insert(EntryInt(4)).
		Insert(EntryInt(1)).
		Insert(EntryInt(3)))

	NewStream(iterator).ForEach(computeSumTotal)
	assert.Equal(t, 8, total)
}

func TestReferenceStream_LeftReduce(t *testing.T) {
	concatenateStringsBiFunc := func(i, j interface{}) interface{} {
		iStr := i.(EntryString)
		jStr := j.(EntryString)
		return EntryString(iStr + "-" + jStr)
	}

	type fields struct {
		iterator Iterator
	}
	type args struct {
		f2 BiFunction
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "Should return nil for an empty Stream",
			fields: fields{
				iterator: NewSetIterator(NewHamtSet()),
			},
			args: args{f2: concatenateStringsBiFunc},
			want: nil,
		},
		{
			name: "Should return reduction of Set of single element",
			fields: fields{
				iterator: NewSetIterator(NewHamtSet().
					Insert(EntryString("three"))),
			},
			args: args{f2: concatenateStringsBiFunc},
			want: EntryString("three"),
		},
		{
			name: "Should return reduction of Set",
			fields: fields{
				iterator: NewSetIterator(NewHamtSet().
					Insert(EntryString("four")).
					Insert(EntryString("twelve")).
					Insert(EntryString("one")).
					Insert(EntryString("six")).
					Insert(EntryString("three"))),
			},
			args: args{f2: concatenateStringsBiFunc},
			want: EntryString("one-three-twelve-six-four"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := ReferenceStream{
				iterator: tt.fields.iterator,
			}
			if gotReduce := rp.Reduce(tt.args.f2); !assert.Equal(t, tt.want, gotReduce) {
				return
			}

			if gotLeftReduce := rp.LeftReduce(tt.args.f2); !assert.Equal(t, tt.want, gotLeftReduce) {
				return
			}
		})
	}
}

func TestReferenceStream_RightReduce(t *testing.T) {
	concatenateStringsBiFunc := func(i, j interface{}) interface{} {
		iStr := i.(EntryString)
		jStr := j.(EntryString)
		return EntryString(iStr + "-" + jStr)
	}

	type fields struct {
		iterator Iterator
	}
	type args struct {
		f2 BiFunction
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "Should return nil for an empty Stream",
			fields: fields{
				iterator: NewSetIterator(NewHamtSet()),
			},
			args: args{f2: concatenateStringsBiFunc},
			want: nil,
		},
		{
			name: "Should return reduction of Set of single element",
			fields: fields{
				iterator: NewSetIterator(NewHamtSet().
					Insert(EntryString("three"))),
			},
			args: args{f2: concatenateStringsBiFunc},
			want: EntryString("three"),
		},
		{
			name: "Should return reduction of Set",
			fields: fields{
				iterator: NewSetIterator(NewHamtSet().
					Insert(EntryString("four")).
					Insert(EntryString("twelve")).
					Insert(EntryString("one")).
					Insert(EntryString("six")).
					Insert(EntryString("three"))),
			},
			args: args{f2: concatenateStringsBiFunc},
			want: EntryString("four-six-twelve-three-one"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := ReferenceStream{
				iterator: tt.fields.iterator,
			}
			got := rp.RightReduce(tt.args.f2)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReferenceStream_Intersperse(t *testing.T) {
	type fields struct {
		iterator Iterator
	}
	type args struct {
		e hamt.Entry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Stream
	}{
		{
			name:   "Should return a Stream of nil",
			fields: fields{iterator: nil},
			args: args{
				e: EntryString(" - "),
			},
			want: NewStream(
				NewEntrySliceIterator([]hamt.Entry{})),
		},
		{
			name: "Should return the original Set when it has a single value",
			fields: fields{
				iterator: NewSetIterator(NewHamtSet().
					Insert(EntryString("four")))},
			args: args{
				e: EntryString(" - "),
			},
			want: NewStream(
				NewEntrySliceIterator([]hamt.Entry{
					EntryString("four")})),
		},
		{
			name: "Should return the Set with given value interspersed",
			fields: fields{
				iterator: NewSetIterator(NewOrderedSet().
					Insert(EntryString("four")).
					Insert(EntryString("twelve")).
					Insert(EntryString("one")).
					Insert(EntryString("six")).
					Insert(EntryString("three"))),
			},
			args: args{
				e: EntryString(" - "),
			},
			want: NewStream(
				NewEntrySliceIterator([]hamt.Entry{
					EntryString("four"),
					EntryString(" - "),
					EntryString("twelve"),
					EntryString(" - "),
					EntryString("one"),
					EntryString(" - "),
					EntryString("six"),
					EntryString(" - "),
					EntryString("three")})),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := ReferenceStream{
				iterator: tt.fields.iterator,
			}
			if got := rp.Intersperse(tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReferenceStream.Intersperse() = %v, want %v", got, tt.want)
			}
		})
	}
}
