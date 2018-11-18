package fuego

import (
	"reflect"
	"testing"

	"github.com/seborama/fuego"
	"github.com/stretchr/testify/require"
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
			}
			if got := si.Forward(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceIterator.Forward() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetIterator_Value(t *testing.T) {
	type fields struct {
		set Set
	}
	tests := []struct {
		name    string
		fields  fields
		want    interface{}
		wantErr string
	}{
		{
			name: "Should produce errNoValue for empty Set",
			fields: fields{
				set: NewSet(),
			},
			wantErr: fuego.PanicNoSuchElement,
		},
		{
			name: "Should return '2' for Set{7, 2, 3}",
			fields: fields{
				set: NewSet().
					Insert(EntryInt(7)).
					Insert(EntryInt(2)).
					Insert(EntryInt(3)),
			},
			want: EntryInt(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := SetIterator{
				set: tt.fields.set,
			}
			if tt.wantErr != "" {
				require.PanicsWithValue(t, tt.wantErr, func() { si.Value() })
				return
			}
			got := si.Value()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetIterator.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceIterator_Value(t *testing.T) {
	type fields struct {
		slice []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    interface{}
		wantErr string
	}{
		{
			name: "Should produce errNoValue for empty slice",
			fields: fields{
				slice: []interface{}{},
			},
			wantErr: fuego.PanicNoSuchElement,
		},
		{
			name: "Should return 3 for Set{7, 2, 3}",
			fields: fields{
				slice: []interface{}{7, 2, 3},
			},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := SliceIterator{
				slice: tt.fields.slice,
			}
			if tt.wantErr != "" {
				require.PanicsWithValue(t, tt.wantErr, func() { si.Value() })
				return
			}
			got := si.Value()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceIterator.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetIterator_Size(t *testing.T) {
	type fields struct {
		set Set
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Should return 0 for empty Set",
			fields: fields{
				set: NewSet(),
			},
			want: 0,
		},
		{
			name: "Should return 3 for Set{1, 2, 3}",
			fields: fields{
				set: NewSet().
					Insert(EntryInt(1)).
					Insert(EntryInt(2)).
					Insert(EntryInt(3)),
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := SetIterator{
				set: tt.fields.set,
			}
			if got := si.Size(); got != tt.want {
				t.Errorf("SetIterator.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceIterator_Size(t *testing.T) {
	type fields struct {
		slice []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Should return 0 for empty Set",
			fields: fields{
				slice: []interface{}{},
			},
			want: 0,
		},
		{
			name: "Should return 3 for Set{1, 2, 3}",
			fields: fields{
				slice: []interface{}{1, 2, 3},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := SliceIterator{
				slice: tt.fields.slice,
			}
			if got := si.Size(); got != tt.want {
				t.Errorf("SliceIterator.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
