package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlattenEntrySliceToEntry(t *testing.T) {
	tests := []struct {
		name  string
		input Entry
		want  EntrySlice
	}{
		{
			name: "Should flatten EntrySlice to its elements",
			input: EntrySlice{
				EntryInt(1),
				EntryInt(3),
				EntryInt(5),
			},
			want: EntrySlice{
				EntryInt(1),
				EntryInt(3),
				EntryInt(5),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FlattenEntrySliceToEntry(0)(tt.input).ToSlice()
			assert.EqualValues(t, tt.want, got)
		})
	}
}
