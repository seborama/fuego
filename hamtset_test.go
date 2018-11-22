package fuego

import (
	"reflect"
	"testing"
)

func TestHamtSet_Stream(t *testing.T) {
	tests := []struct {
		name string
		set  HamtSet
		want Stream
	}{
		{
			name: "Should return an empty stream when empty set",
			set:  NewHamtSet(),
			want: NewStream(
				NewSetIterator(
					NewHamtSet())),
		},
		{
			name: "Should return value when one value set",
			set:  NewHamtSet().Insert(EntryInt(1)).(HamtSet),
			want: NewStream(
				NewSetIterator(
					NewHamtSet().
						Insert(EntryInt(1)))),
		},
		{
			name: "Should return values present in the Set",
			set: NewHamtSet().
				Insert(EntryInt(1)).
				Insert(EntryInt(2)).
				Delete(EntryInt(1)).
				Insert(EntryInt(3)).
				Insert(EntryInt(1)).
				Insert(EntryInt(2)).(HamtSet),
			want: NewStream(
				NewSetIterator(
					NewHamtSet().
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

func TestHamtSet_Merge(t *testing.T) {
	type fields struct {
		set HamtSet
	}
	type args struct {
		t HamtSet
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   HamtSet
	}{
		{
			name: "Should merge two excluding sets",
			fields: fields{
				set: NewHamtSet().
					Insert(EntryInt(7)).
					Insert(EntryInt(2)).(HamtSet),
			},
			args: args{
				t: NewHamtSet().
					Insert(EntryInt(3)).
					Insert(EntryInt(9)).(HamtSet),
			},
			want: NewHamtSet().
				Insert(EntryInt(7)).
				Insert(EntryInt(2)).
				Insert(EntryInt(3)).
				Insert(EntryInt(9)).(HamtSet),
		},
		{
			name: "Should merge two overlapping sets",
			fields: fields{
				set: NewHamtSet().
					Insert(EntryInt(3)).
					Insert(EntryInt(1)).(HamtSet),
			},
			args: args{
				t: NewHamtSet().
					Insert(EntryInt(3)).
					Insert(EntryInt(1)).(HamtSet),
			},
			want: NewHamtSet().
				Insert(EntryInt(3)).
				Insert(EntryInt(1)).(HamtSet),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.set.Merge(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set.Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHamtSet_FirstRest(t *testing.T) {
	type fields struct {
		set HamtSet
	}
	tests := []struct {
		name   string
		fields fields
		want   Entry
		want1  HamtSet
	}{
		{
			name: "Should pop first and return rest",
			fields: fields{
				set: NewHamtSet().
					Insert(EntryInt(3)).
					Insert(EntryInt(2)).
					Insert(EntryInt(7)).(HamtSet),
			},
			// note: hamt.Set entries are seemingly sorted based on their hash
			want: EntryInt(2),
			want1: NewHamtSet().
				Insert(EntryInt(3)).
				Insert(EntryInt(7)).(HamtSet),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.fields.set.FirstRest()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set.FirstRest() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Set.FirstRest() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHamtSet_Size(t *testing.T) {
	type fields struct {
		set Set
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Should return 0 for empty Set",
			fields: fields{
				set: NewHamtSet(),
			},
			want: 0,
		},
		{
			name: "Should return size",
			fields: fields{
				set: NewHamtSet().
					Insert(EntryInt(7)).
					Insert(EntryInt(6)).
					Insert(EntryInt(1)).
					Insert(EntryInt(2)),
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.set.Size(); got != tt.want {
				t.Errorf("HamtSet.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
