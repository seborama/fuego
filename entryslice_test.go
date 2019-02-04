package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntrySlice(t *testing.T) {
	t.Log(EntrySlice([]Entry{}))
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
