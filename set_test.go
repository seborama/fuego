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

func timesTwoI() func(hamt.Entry) interface{} {
	return func(e hamt.Entry) interface{} {
		e1, ok := e.(EntryInt)
		if !ok {
			return EntryInt(0)
		}
		return hamt.Entry(2 * EntryInt(e1))
	}
}

func TestSet_Map(t *testing.T) {
	type args struct {
		f func(hamt.Entry) hamt.Entry
	}
	tests := []struct {
		name string
		set  Set
		args args
		want []hamt.Entry
	}{
		{
			name: "Should return a new Set with value doubles",
			set: NewSet().
				Insert(EntryInt(1)).
				Insert(EntryInt(2)).
				Insert(EntryInt(3)),
			args: args{
				f: timesTwo(),
			},
			want: []hamt.Entry{
				EntryInt(2),
				EntryInt(4),
				EntryInt(6),
			},
		},
		{
			name: "Should return an empty Set",
			set:  NewSet(),
			args: args{
				f: timesTwo(),
			},
			want: []hamt.Entry{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.set.Map(tt.args.f).Values()
			assert.EqualValues(t, got, tt.want)
		})
	}
}

func TestSet_MapC(t *testing.T) {
	type args struct {
		f func(hamt.Entry) interface{}
	}
	tests := []struct {
		name string
		set  Set
		args args
		want []hamt.Entry
	}{
		{
			name: "Should return a channel with value doubles",
			set: NewSet().
				Insert(EntryInt(1)).
				Insert(EntryInt(2)).
				Insert(EntryInt(3)),
			args: args{
				f: timesTwoI(),
			},
			want: []hamt.Entry{
				EntryInt(2),
				EntryInt(4),
				EntryInt(6),
			},
		},
		{
			name: "Should return a channel without data",
			set:  NewSet(),
			args: args{
				f: timesTwoI(),
			},
			want: []hamt.Entry{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := []hamt.Entry{}
			streams := tt.set.MapC(tt.args.f)
			stream := streams[0]
			for v := range stream {
				got = append(got, v.(hamt.Entry))
			}
			assert.EqualValues(t, got, tt.want)
		})
	}
}
