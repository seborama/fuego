package fuego

import (
	"reflect"
	"testing"
)

type FunctionTimesTwo int

func NewFunctionTimesTwo() FunctionTimesTwo {
	return *(new(FunctionTimesTwo))
}

func (f FunctionTimesTwo) Apply(i interface{}) interface{} {
	num := i.(EntryInt).Value()
	result := 2 * num
	return interface{}(result)
}

func TestReferenceStream_Map(t *testing.T) {
	fi := NewFunctionTimesTwo()

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
			want: NewStream(
				NewBaseIterable(
					NewSliceIterator([]interface{}{}),
				)),
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
				NewBaseIterable(
					NewSliceIterator([]interface{}{
						EntryInt(8)}))),
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
				NewBaseIterable(
					NewSliceIterator([]interface{}{
						EntryInt(2),
						EntryInt(4),
						EntryInt(6)}))),
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
