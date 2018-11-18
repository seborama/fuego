package fuego

import (
	"reflect"
	"testing"

	"github.com/raviqqe/hamt"
	"github.com/seborama/fuego"
	"github.com/stretchr/testify/assert"
)

func TestEntryEntry_Forward(t *testing.T) {
	type fields struct {
		slice []hamt.Entry
	}
	tests := []struct {
		name   string
		fields fields
		want   Iterator
	}{
		{
			name: "Should provide nil when no data exist",
			fields: fields{
				slice: []hamt.Entry{},
			},
			want: nil,
		},
		{
			name: "Should provide nil when no more data exists",
			fields: fields{
				slice: []hamt.Entry{EntryInt(1)},
			},
			want: nil,
		},
		{
			name: "Should provide iterator when more data exists",
			fields: fields{
				slice: []hamt.Entry{
					EntryInt(7),
					EntryInt(6),
					EntryInt(1),
					EntryInt(2)},
			},
			want: NewEntrySliceIterator(
				[]hamt.Entry{
					EntryInt(6),
					EntryInt(1),
					EntryInt(2)}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := EntrySliceIterator{
				slice: tt.fields.slice,
			}
			if got := si.Forward(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceIterator.Forward() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntrySliceIterator_Value(t *testing.T) {
	type fields struct {
		slice []hamt.Entry
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
				slice: []hamt.Entry{},
			},
			wantErr: fuego.PanicNoSuchElement,
		},
		{
			name: "Should return the current value",
			fields: fields{
				slice: []hamt.Entry{
					EntryInt(7),
					EntryInt(6),
					EntryInt(1),
					EntryInt(2)},
			},
			want: EntryInt(7),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := EntrySliceIterator{
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

func TestEntrySliceIterator_Size(t *testing.T) {
	type fields struct {
		slice []hamt.Entry
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Should return 0 for empty Set",
			fields: fields{
				slice: []hamt.Entry{},
			},
			want: 0,
		},
		{
			name: "Should return size",
			fields: fields{
				slice: []hamt.Entry{
					EntryInt(7),
					EntryInt(6),
					EntryInt(1),
					EntryInt(2)},
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := EntrySliceIterator{
				slice: tt.fields.slice,
			}
			if got := si.Size(); got != tt.want {
				t.Errorf("Entry.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
