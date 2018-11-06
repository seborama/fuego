package fuego

import (
	"testing"

	"github.com/raviqqe/hamt"
	"github.com/stretchr/testify/assert"
)

func TestSeq_Head_PanicsIfEmpty(t *testing.T) {
	unit := NewSeq()
	assert.Panics(t, func() { unit.Head() })
}

func TestSeq_Head(t *testing.T) {
	unit := NewSeq().
		Append(EntryInt(1)).
		Append(EntryInt(2)).
		Append(EntryInt(3))
	assert.Equal(t, unit.Head(), EntryInt(1))
}

func TestSeq_Get_PanicsIfEmpty(t *testing.T) {
	unit := NewSeq()
	assert.Panics(t, func() { unit.Get() })
}

func TestSeq_Get(t *testing.T) {
	unit := NewSeq().
		Append(EntryInt(1)).
		Append(EntryInt(2)).
		Append(EntryInt(3))
	assert.Equal(t, unit.Get(), EntryInt(1))
}

func TestSeq_Last_PanicsIfEmpty(t *testing.T) {
	unit := NewSeq()
	assert.Panics(t, func() { unit.Last() })
}

func TestSeq_Last(t *testing.T) {
	unit := NewSeq().
		Append(EntryInt(1)).
		Append(EntryInt(2)).
		Append(EntryInt(3))
	assert.Equal(t, unit.Last(), EntryInt(3))
}

func TestSeq_Tail_PanicsIfEmpty(t *testing.T) {
	unit := NewSeq()
	assert.Panics(t, func() { unit.Tail() })
}

func TestSeq_Tail(t *testing.T) {
	unit := NewSeq().
		Append(EntryInt(1)).
		Append(EntryInt(2)).
		Append(EntryInt(3))

	assert.EqualValues(t, unit.Tail(),
		NewSeq().
			Append(EntryInt(2)).
			Append(EntryInt(3)),
	)
}

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

func TestSeq_AsGo_Ordering(t *testing.T) {
	tests := []struct {
		name string
		seq  Seq
		want []hamt.Entry
	}{
		{
			name: "Should preserve the order of items - case 1,2,3",
			seq: NewSeq().
				Append(EntryInt(1)).
				Append(EntryInt(2)).
				Append(EntryInt(3)),
			want: []hamt.Entry{EntryInt(1), EntryInt(2), EntryInt(3)},
		},
		{
			name: "Should preserve the order of items - case 3,2,7 with valid Remove",
			seq: NewSeq().
				Append(EntryInt(3)).
				Append(EntryInt(1)).
				Append(EntryInt(2)).
				Remove(EntryInt(1)).
				Append(EntryInt(7)),
			want: []hamt.Entry{EntryInt(3), EntryInt(2), EntryInt(7)},
		},
		{
			name: "Should preserve the order of items - case 3,1,2,7 with invalid Remove",
			seq: NewSeq().
				Append(EntryInt(3)).
				Append(EntryInt(1)).
				Append(EntryInt(2)).
				Remove(EntryInt(123456)). //  remove an unknown element (should have no effect)
				Append(EntryInt(7)),
			want: []hamt.Entry{EntryInt(3), EntryInt(1), EntryInt(2), EntryInt(7)},
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

func TestSeq_IsEmpty(t *testing.T) {
	type fields struct {
		seq []hamt.Entry
	}
	tests := []struct {
		name string
		seq  Seq
		want bool
	}{
		{
			name: "Should be empty",
			seq:  NewSeq(),
			want: true,
		},
		{
			name: "Should not be empty",
			seq: NewSeq().
				Append(EntryInt(1)).
				Append(EntryInt(2)).
				Append(EntryInt(3)),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.seq.IsEmpty(); got != tt.want {
				t.Errorf("Seq.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_NonEmpty(t *testing.T) {
	type fields struct {
		seq []hamt.Entry
	}
	tests := []struct {
		name string
		seq  Seq
		want bool
	}{
		{
			name: "Should not be non-empty (i.e. is empty)",
			seq:  NewSeq(),
			want: false,
		},
		{
			name: "Should not be non-empty",
			seq: NewSeq().
				Append(EntryInt(1)).
				Append(EntryInt(2)).
				Append(EntryInt(3)),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.seq.NonEmpty(); got != tt.want {
				t.Errorf("Seq.NonEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_HashCode(t *testing.T) {
	type fields struct {
		seq []hamt.Entry
	}
	tests := []struct {
		name string
		seq  Seq
		want uint32
	}{
		{
			name: "Should be 2166136261 for an empty Seq",
			seq:  NewSeq(),
			want: 2166136261,
		},
		{
			name: "Should be 2034659765 value for Seq of EntryInts {1,2,3}",
			seq: NewSeq().
				Append(EntryInt(1)).
				Append(EntryInt(2)).
				Append(EntryInt(3)),
			want: 2034659765,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.seq.HashCode(); got != tt.want {
				t.Errorf("Seq.NonEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq_ImplementsTraversable(t *testing.T) {
	unit := NewSeq()
	assert.Implements(t, (*Traversable)(nil), unit)
}
