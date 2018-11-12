package fuego

import (
	"reflect"
	"testing"
)

func TestSet_Stream(t *testing.T) {
	tests := []struct {
		name string
		set  Set
		want Stream
	}{
		{
			name: "Should return an empty stream when empty set",
			set:  NewSet(),
			want: NewStream(
				NewSetIterator(
					NewSet())),
		},
		{
			name: "Should return value when one value set",
			set:  NewSet().Insert(EntryInt(1)),
			want: NewStream(
				NewSetIterator(
					NewSet().
						Insert(EntryInt(1)))),
		},
		{
			name: "Should return values in same order when three value set with operations",
			set: NewSet().
				Insert(EntryInt(1)).
				Insert(EntryInt(2)).
				Delete(EntryInt(1)).
				Insert(EntryInt(3)).
				Insert(EntryInt(1)).
				Insert(EntryInt(2)),
			want: NewStream(
				NewSetIterator(
					NewSet().
						Insert(EntryInt(2)).
						Insert(EntryInt(3)).
						Insert(EntryInt(1)))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Stream(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set.Stream() = %v, want %v", got, tt.want)
			}
		})
	}
}
