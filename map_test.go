package fuego

import (
	"reflect"
	"testing"

	"github.com/raviqqe/hamt"
	"github.com/stretchr/testify/assert"
)

func TestMap_Insert(t *testing.T) {
	type fields struct {
		myMap Map
	}
	type args struct {
		k hamt.Entry
		v interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Map
	}{
		{
			name: "Should Insert entries into Map",
			fields: fields{
				myMap: NewMap().
					Insert(EntryInt(1), "one"),
			},
			args: args{
				k: EntryInt(5),
				v: "five",
			},
			want: NewMap().
				Insert(EntryInt(1), "one").
				Insert(EntryInt(5), "five"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.myMap.Insert(tt.args.k, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap_Delete(t *testing.T) {
	type fields struct {
		myMap Map
	}
	type args struct {
		k hamt.Entry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Map
	}{
		{
			name: "Should Insert entries into Map",
			fields: fields{
				myMap: NewMap().
					Insert(EntryInt(1), "one").
					Insert(EntryInt(5), "five").
					Insert(EntryInt(2), "two"),
			},
			args: args{
				k: EntryInt(5),
			},
			want: NewMap().
				Insert(EntryInt(1), "one").
				Insert(EntryInt(2), "two"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.myMap.Delete(tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap_Size(t *testing.T) {
	m := NewMap().
		Insert(EntryInt(1), "one").
		Insert(EntryInt(5), "five").
		Insert(EntryInt(2), "two")

	assert.Equal(t, 3, m.Size())
}
