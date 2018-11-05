package fuego

import (
	"testing"

	"github.com/raviqqe/hamt"
	"github.com/stretchr/testify/assert"
)

func TestSeq_AsGo_Immutability(t *testing.T) {
	unit := NewSeq()
	unit = unit.Append(EntryInt(1))
	unit = unit.Append(EntryInt(2))
	unit = unit.Append(EntryInt(3))
	assert.Equal(t, 3, unit.Size())

	seq := unit.AsGo()
	assert.EqualValues(t, unit.seq, seq)

	seq[0] = EntryInt(1000)
	assert.NotEqual(t, unit.seq[0], seq[0])
}

func TestSeq_AsGo(t *testing.T) {
	tests := []struct {
		name string
		seq  Seq
		want []hamt.Entry
	}{
		{
			name: "Should preserve the order of items - case 1",
			seq: NewSeq().
				Append(EntryInt(1)).
				Append(EntryInt(2)).
				Append(EntryInt(3)),
			want: []hamt.Entry{EntryInt(1), EntryInt(2), EntryInt(3)},
		},
		{
			name: "Should preserve the order of items - case 2",
			seq: NewSeq().
				Append(EntryInt(3)).
				Append(EntryInt(1)).
				Append(EntryInt(2)).
				Remove(EntryInt(1)).
				Append(EntryInt(7)),
			want: []hamt.Entry{EntryInt(3), EntryInt(2), EntryInt(7)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.seq
			got := s.AsGo()
			assert.EqualValues(t, tt.want, got)
		})
	}
}
