package fuego

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetIterator_Forward(t *testing.T) {
	type fields struct {
		set Set
	}
	tests := []struct {
		name   string
		fields fields
		want   Iterator
	}{
		{
			name: "Should provide nil when no data exist",
			fields: fields{
				set: NewHamtSet(),
			},
			want: nil,
		},
		{
			name: "Should provide nil when no more data exists",
			fields: fields{
				set: NewHamtSet().
					Insert(EntryInt(1)),
			},
			want: nil,
		},
		{
			name: "Should provide iterator when more data exists",
			fields: fields{
				set: NewHamtSet().
					Insert(EntryInt(7)).
					Insert(EntryInt(6)).
					Insert(EntryInt(1)).
					Insert(EntryInt(2)),
			},
			want: NewSetIterator(NewHamtSet().
				Insert(EntryInt(2)).
				Insert(EntryInt(6)).
				Insert(EntryInt(7))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := SetIterator{
				set: tt.fields.set,
			}
			if got := si.Forward(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetIterator.Forward() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetIterator_Value(t *testing.T) {
	type fields struct {
		set Set
	}
	tests := []struct {
		name    string
		fields  fields
		want    interface{}
		wantErr string
	}{
		{
			name: "Should produce PanicNoSuchElement for nil Set",
			fields: fields{
				set: nil,
			},
			wantErr: PanicNoSuchElement,
		},
		{
			name: "Should produce PanicNoSuchElement for empty Set",
			fields: fields{
				set: NewHamtSet(),
			},
			wantErr: PanicNoSuchElement,
		},
		{
			name: "Should return the current value",
			fields: fields{
				set: NewHamtSet().
					Insert(EntryInt(7)).
					Insert(EntryInt(6)).
					Insert(EntryInt(1)).
					Insert(EntryInt(2)),
			},
			want: EntryInt(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := SetIterator{
				set: tt.fields.set,
			}
			if tt.wantErr != "" {
				assert.PanicsWithValue(t, tt.wantErr, func() { si.Value() })
				return
			}
			got := si.Value()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetIterator.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetIterator_Size(t *testing.T) {
	type fields struct {
		set Set
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Should return 0 for nil Set",
			fields: fields{
				set: nil,
			},
			want: 0,
		},
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
			si := SetIterator{
				set: tt.fields.set,
			}
			if got := si.Size(); got != tt.want {
				t.Errorf("SetIterator.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetIterator_Reverse(t *testing.T) {
	type fields struct {
		set Set
	}
	tests := []struct {
		name   string
		fields fields
		want   Iterator
	}{
		{
			name: "Should provide an empty set when no data exist",
			fields: fields{
				set: NewHamtSet(),
			},
			want: NewSetIterator(NewOrderedSet()),
		},
		{
			name: "Should provide a same set for a single entry set",
			fields: fields{
				set: NewHamtSet().
					Insert(EntryInt(1)),
			},
			want: NewSetIterator(NewOrderedSet().
				Insert(EntryInt(1))),
		},
		{
			name: "Should provide reverse set",
			fields: fields{
				// reminder: hamt.Set is unnaturally ordered: 1, 2, 6, 7
				set: NewHamtSet().
					Insert(EntryInt(7)).
					Insert(EntryInt(6)).
					Insert(EntryInt(1)).
					Insert(EntryInt(2)),
			},
			want: NewSetIterator(NewOrderedSet().
				Insert(EntryInt(7)).
				Insert(EntryInt(6)).
				Insert(EntryInt(2)).
				Insert(EntryInt(1))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := SetIterator{
				set: tt.fields.set,
			}
			if got := si.Reverse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetIterator.Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
