package fuego

import (
	"testing"

	"github.com/raviqqe/hamt"
	"github.com/stretchr/testify/assert"
)

func timesTwo() func(hamt.Entry) hamt.Entry {
	return func(e hamt.Entry) hamt.Entry {
		e1, ok := e.(EntryInt)
		if !ok {
			return EntryInt(0)
		}
		return hamt.Entry(2 * EntryInt(e1))
	}
}

func xTestMapFunctionOverCollection(t *testing.T) {
	unit := NewSet().
		Insert(EntryInt(1)).
		Insert(EntryInt(2)).
		Insert(EntryInt(3)).
		Map(timesTwo())

	expected := NewSeq().
		Append(EntryInt(2)).
		Append(EntryInt(4)).
		Append(EntryInt(6))

	assert.EqualValues(t, expected, unit.Values())
}

func TestSet_Map(t *testing.T) {
	type args struct {
		f func(hamt.Entry) hamt.Entry
	}
	tests := []struct {
		name string
		set  Set
		args args
		want Seq
	}{
		{
			name: "Should return a new map with value doubles",
			set: NewSet().
				Insert(EntryInt(1)).
				Insert(EntryInt(2)).
				Insert(EntryInt(3)),
			args: args{
				f: timesTwo(),
			},
			want: NewSeq().
				Append(EntryInt(2)).
				Append(EntryInt(4)).
				Append(EntryInt(6)),
		},
		{
			name: "Should return an empty Seq when Set is empty",
			set:  NewSet(),
			args: args{
				f: timesTwo(),
			},
			want: NewSeq(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.set.Map(tt.args.f).Values()
			assert.EqualValues(t, got, tt.want)
		})
	}
}
