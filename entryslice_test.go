package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntrySlice(t *testing.T) {
	t.Log(EntrySlice(EntrySlice{}))
}

func TestEntrySliceHash(t *testing.T) {
	type fields struct {
		slice EntrySlice
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name:   "Should return 0 for a nil slice",
			fields: fields{slice: nil},
			want:   0,
		},
		{
			name:   "Should return 0 for an empty slice",
			fields: fields{slice: EntrySlice{}},
			want:   0,
		},
		{
			name:   "Should return hash for a single-item slice",
			fields: fields{slice: EntrySlice{EntryInt(1)}},
			want:   0x20,
		},
		{
			name: "Should return hash for a multi-item slice",
			fields: fields{slice: EntrySlice{
				EntryInt(7),
				EntryInt(13),
				EntryInt(28)}},
			want: 0x9055,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.fields.slice.Hash())
		})
	}
}

func TestEntrySliceEqual(t *testing.T) {
	type fields struct {
		slice1 EntrySlice
		slice2 Entry
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Should return false for when comparee is not an EntrySlice",
			fields: fields{
				slice1: EntrySlice{},
				slice2: EntryMap{}},
			want: false,
		},
		{
			name: "Should return true for two empty slices",
			fields: fields{
				slice1: EntrySlice{},
				slice2: EntrySlice{}},
			want: true,
		},
		{
			name: "Should return false for one empty slice and one non-empty slice",
			fields: fields{
				slice1: EntrySlice{},
				slice2: EntrySlice{EntryInt(7)}},
			want: false,
		},
		{
			name: "Should return true for two identical multi-item slices",
			fields: fields{
				slice1: EntrySlice{
					EntryInt(7),
					EntryInt(13),
					EntryInt(28)},
				slice2: EntrySlice{
					EntryInt(7),
					EntryInt(13),
					EntryInt(28)}},
			want: true,
		},
		{
			name: "Should return false for two slices with same value items but in different order",
			fields: fields{
				slice1: EntrySlice{
					EntryInt(7),
					EntryInt(13),
					EntryInt(28)},
				slice2: EntrySlice{
					EntryInt(13),
					EntryInt(7),
					EntryInt(28)}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.fields.slice1.Equal(tt.fields.slice2))
		})
	}
}

func TestEntrySlice_Append(t *testing.T) {
	type args struct {
		e Entry
	}
	tests := []struct {
		name string
		es   EntrySlice
		args args
		want EntrySlice
	}{
		{
			name: "Should append entry when slice is nil",
			es:   nil,
			args: args{
				e: EntryInt(1),
			},
			want: EntrySlice{EntryInt(1)},
		},
		{
			name: "Should append Entry to slice when Entry does not exist and should not modify original slice",
			es: EntrySlice{
				EntryInt(2),
				EntryInt(3),
			},
			args: args{
				e: EntryInt(1),
			},
			want: EntrySlice{
				EntryInt(2),
				EntryInt(3),
				EntryInt(1),
			},
		},
		{
			name: "Should append Entry to slice when Entry exists and should not modify original slice",
			es: EntrySlice{
				EntryInt(2),
				EntryInt(3),
			},
			args: args{
				e: EntryInt(3),
			},
			want: EntrySlice{
				EntryInt(2),
				EntryInt(3),
				EntryInt(3),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := make(EntrySlice, tt.es.Len())
			if tt.es == nil {
				original = nil
			} else {
				copy(original, tt.es)
			}
			got := tt.es.Append(tt.args.e)
			assert.EqualValues(t, tt.want, got, "unexpected Append() behaviour")
			assert.EqualValues(t, original, tt.es, "original slice was not preserved")
		})
	}
}

func TestEntrySlice_Len(t *testing.T) {
	tests := []struct {
		name string
		es   EntrySlice
		want int
	}{
		{
			name: "Should return 0 for nil slice",
			es:   nil,
			want: 0,
		},
		{
			name: "Should return 0 for empty slice",
			es:   EntrySlice{},
			want: 0,
		},
		{
			name: "Should return 1 for slice of 1 Entry",
			es:   EntrySlice{EntryInt(123)},
			want: 1,
		},
		{
			name: "Should return 3 for slice of 3 Entries",
			es: EntrySlice{
				EntryInt(123),
				EntryInt(12),
				EntryInt(3),
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.es.Len(); got != tt.want {
				t.Errorf("EntrySlice.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}
