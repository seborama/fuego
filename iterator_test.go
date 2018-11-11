package fuego

import (
	"reflect"
	"testing"
)

func TestSetIterator_Forward(t *testing.T) {
	type fields struct {
		set Set
	}
	tests := []struct {
		name   string
		fields fields
		want   Iterator
	}{
		{
			name: "Should provide nil when no data exist",
			fields: fields{
				set: NewSet(),
			},
			want: nil,
		},
		{
			name: "Should provide nil when no more data exists",
			fields: fields{
				set: NewSet().
					Insert(EntryInt(1)),
			},
			want: nil,
		},
		{
			name: "Should provide iterator when more data exists",
			fields: fields{
				set: NewSet().
					Insert(EntryInt(1)).
					Insert(EntryInt(2)),
			},
			want: NewSetIterator(NewSet().
				Insert(EntryInt(2))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := SetIterator{
				set: tt.fields.set,
			}
			if got := si.Forward(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetIterator.Forward() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceIterator_Forward(t *testing.T) {
	type fields struct {
		slice []interface{}
		size  int
	}
	tests := []struct {
		name   string
		fields fields
		want   Iterator
	}{
		{
			name: "Should provide nil when no data exist",
			fields: fields{
				slice: []interface{}{},
				size:  0,
			},
			want: nil,
		},
		{
			name: "Should provide nil when no more data exists",
			fields: fields{
				slice: []interface{}{1},
				size:  1,
			},
			want: nil,
		},
		{
			name: "Should provide iterator when more data exists",
			fields: fields{
				slice: []interface{}{1, 2, 3},
				size:  3,
			},
			want: NewSliceIterator([]interface{}{2, 3}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := SliceIterator{
				slice: tt.fields.slice,
				size:  tt.fields.size,
			}
			if got := si.Forward(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceIterator.Forward() = %v, want %v", got, tt.want)
			}
		})
	}
}
