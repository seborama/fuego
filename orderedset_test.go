package fuego

import (
	"reflect"
	"testing"
)

func TestOrderedSet_Stream(t *testing.T) {
	tests := []struct {
		name string
		set  OrderedSet
		want Stream
	}{
		{
			name: "Should return an empty stream when empty set",
			set:  NewOrderedSet(),
			want: NewStream(
				NewSliceIterator(
					[]Entry{})),
		},
		{
			name: "Should return value when one value set",
			set: NewOrderedSet().
				Insert(EntryInt(1)).(OrderedSet),
			want: NewStream(
				NewSliceIterator(
					[]Entry{EntryInt(1)})),
		},
		{
			name: "Should return values present in the Set, and in order",
			set: NewOrderedSet().
				Insert(EntryInt(1)).
				Insert(EntryInt(2)).
				Delete(EntryInt(1)).
				Insert(EntryInt(3)).
				Insert(EntryInt(1)).
				Insert(EntryInt(2)).(OrderedSet),
			want: NewStream(
				NewSliceIterator(
					[]Entry{
						EntryInt(2),
						EntryInt(3),
						EntryInt(1)})),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Stream(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderedSet.Stream() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderedSet_Merge(t *testing.T) {
	type fields struct {
		set OrderedSet
	}
	type args struct {
		t Set
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Set
	}{
		{
			name: "Should merge two excluding sets",
			fields: fields{
				set: NewOrderedSet().
					Insert(EntryInt(7)).
					Insert(EntryInt(2)).(OrderedSet),
			},
			args: args{
				t: NewOrderedSet().
					Insert(EntryInt(3)).
					Insert(EntryInt(9)),
			},
			want: NewOrderedSet().
				Insert(EntryInt(7)).
				Insert(EntryInt(2)).
				Insert(EntryInt(3)).
				Insert(EntryInt(9)),
		},
		{
			name: "Should merge two overlapping sets",
			fields: fields{
				set: NewOrderedSet().
					Insert(EntryInt(3)).
					Insert(EntryInt(1)).(OrderedSet),
			},
			args: args{
				t: NewOrderedSet().
					Insert(EntryInt(3)).
					Insert(EntryInt(1)),
			},
			want: NewOrderedSet().
				Insert(EntryInt(3)).
				Insert(EntryInt(1)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.set.Merge(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderedSet.Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderedSet_FirstRest(t *testing.T) {
	type fields struct {
		set Set
	}
	tests := []struct {
		name   string
		fields fields
		want   Entry
		want1  Set
	}{
		{
			name: "Should pop first and return rest",
			fields: fields{
				set: NewOrderedSet().
					Insert(EntryInt(3)).
					Insert(EntryInt(2)).
					Insert(EntryInt(7)),
			},
			want: EntryInt(3),
			want1: NewOrderedSet().
				Insert(EntryInt(2)).
				Insert(EntryInt(7)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.fields.set.FirstRest()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderedSet.FirstRest() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("OrderedSet.FirstRest() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestOrderedSet_Size(t *testing.T) {
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
				set: NewOrderedSet(),
			},
			want: 0,
		},
		{
			name: "Should return size",
			fields: fields{
				set: NewOrderedSet().
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
				t.Errorf("NewOrderedSet.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderedSet_Delete(t *testing.T) {
	type args struct {
		e Entry
	}
	tests := []struct {
		name string
		set  Set
		args args
		want Set
	}{
		{
			name: "Should return empty set when deleting from empty set",
			set:  NewOrderedSet(),
			args: args{},
			want: NewOrderedSet(),
		},
		{
			name: "Should return empty set when deleting unique entry from set",
			set: NewOrderedSet().
				Insert(EntryInt(1)),
			args: args{
				e: EntryInt(1),
			},
			want: NewOrderedSet(),
		},
		{
			name: "Should return original set when deleting non-existent entry from set",
			set: NewOrderedSet().
				Insert(EntryInt(3)).
				Insert(EntryInt(8)).
				Insert(EntryInt(1)).
				Insert(EntryInt(7)),
			args: args{
				e: EntryInt(-999),
			},
			want: NewOrderedSet().
				Insert(EntryInt(3)).
				Insert(EntryInt(8)).
				Insert(EntryInt(1)).
				Insert(EntryInt(7)),
		},
		{
			name: "Should return reduced set when deleting first entry from set",
			set: NewOrderedSet().
				Insert(EntryInt(3)).
				Insert(EntryInt(8)).
				Insert(EntryInt(1)).
				Insert(EntryInt(7)),
			args: args{
				e: EntryInt(3),
			},
			want: NewOrderedSet().
				Insert(EntryInt(8)).
				Insert(EntryInt(1)).
				Insert(EntryInt(7)),
		},
		{
			name: "Should return reduced set when deleting last entry from set",
			set: NewOrderedSet().
				Insert(EntryInt(3)).
				Insert(EntryInt(8)).
				Insert(EntryInt(1)).
				Insert(EntryInt(7)),
			args: args{
				e: EntryInt(7),
			},
			want: NewOrderedSet().
				Insert(EntryInt(3)).
				Insert(EntryInt(8)).
				Insert(EntryInt(1)),
		},
		{
			name: "Should return reduced set when deleting middle entry from set",
			set: NewOrderedSet().
				Insert(EntryInt(3)).
				Insert(EntryInt(8)).
				Insert(EntryInt(1)).
				Insert(EntryInt(7)),
			args: args{
				e: EntryInt(8),
			},
			want: NewOrderedSet().
				Insert(EntryInt(3)).
				Insert(EntryInt(1)).
				Insert(EntryInt(7)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Delete(tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderedSet.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}
