package fuego

import (
	"reflect"
	"testing"

	"github.com/seborama/fuego"
	"github.com/stretchr/testify/assert"
)

func TestSliceIterator_Forward(t *testing.T) {
	type fields struct {
		slice []interface{}
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
			},
			want: nil,
		},
		{
			name: "Should provide nil when no more data exists",
			fields: fields{
				slice: []interface{}{1},
			},
			want: nil,
		},
		{
			name: "Should provide iterator when more data exists",
			fields: fields{
				slice: []interface{}{7, 6, 1, 2},
			},
			want: NewSliceIterator([]interface{}{6, 1, 2}),
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
			name: "Should produce PanicNoSuchElement for empty slice",
			fields: fields{
				slice: []interface{}{},
			},
			wantErr: fuego.PanicNoSuchElement,
		},
		{
			name: "Should return the current value",
			fields: fields{
				slice: []interface{}{7, 6, 1, 2},
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
				assert.PanicsWithValue(t, tt.wantErr, func() { si.Value() })
				return
			}
			got := si.Value()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceIterator.Value() = %v, want %v", got, tt.want)
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
			name: "Should return size",
			fields: fields{
				slice: []interface{}{7, 6, 1, 2},
			},
			want: 4,
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

func TestSliceIterator_Reverse(t *testing.T) {
	panic("not yet written")
	type fields struct {
		slice []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   Iterator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := SliceIterator{
				slice: tt.fields.slice,
			}
			if got := si.Reverse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceIterator.Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
