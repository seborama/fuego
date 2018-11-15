package fuego

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type functionTimesTwo int

func newFunctionTimesTwo() functionTimesTwo {
	return *(new(functionTimesTwo))
}

func (f functionTimesTwo) Apply(i interface{}) interface{} {
	num := i.(EntryInt).Value()
	return interface{}(2 * num)
}

func TestReferenceStream_Map(t *testing.T) {
	fi := newFunctionTimesTwo()

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
			name: "Should return a Stream of nil",
			fields: fields{
				iterator: nil},
			args: args{
				mapper: fi,
			},
			want: NewStream(NewSliceIterator([]interface{}{})),
		},
		{
			name: "Should return a Stream of one double",
			fields: fields{
				iterator: NewSetIterator(NewSet().
					Insert(EntryInt(4)))},
			args: args{
				mapper: fi,
			},
			want: NewStream(
				NewSliceIterator(
					[]interface{}{EntryInt(8)})),
		},
		{
			name: "Should return a Stream of 3 doubles",
			fields: fields{
				iterator: NewSetIterator(NewSet().
					Insert(EntryInt(1)).
					Insert(EntryInt(2)).
					Insert(EntryInt(3)))},
			args: args{
				mapper: fi,
			},
			want: NewStream(
				NewSliceIterator([]interface{}{
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

type entryIntEqualsTo struct {
	entryInt EntryInt
}

func newEntryIntEqualsTo(e EntryInt) Function {
	return entryIntEqualsTo{
		entryInt: e,
	}
}
func (e entryIntEqualsTo) Apply(subject interface{}) interface{} {
	subjectEntryInt, ok := subject.(EntryInt)
	if !ok {
		return false
	}

	return e.entryInt.Equal(subjectEntryInt)
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
				iterator: NewSetIterator(NewSet().
					Insert(EntryInt(4)).
					Insert(EntryInt(1)).
					Insert(EntryInt(3)))},
			args: args{
				predicate: Or(
					NewFunctionPredicate(newEntryIntEqualsTo(EntryInt(1))),
					NewFunctionPredicate(newEntryIntEqualsTo(EntryInt(3))),
				),
			},
			want: NewStream(
				NewSliceIterator([]interface{}{
					EntryInt(1),
					EntryInt(3)})),
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

	iterator := NewSetIterator(NewSet().
		Insert(EntryInt(4)).
		Insert(EntryInt(1)).
		Insert(EntryInt(3)))

	NewStream(iterator).ForEach(computeSumTotal)
	assert.Equal(t, 8, total)
}
