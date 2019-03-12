package fuego

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTuple2_Hash(t *testing.T) {
	type fields struct {
		E1 Entry
		E2 Entry
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "Should return hash for Tuple2(nil,nil)",
			fields: fields{
				E1: nil,
				E2: nil,
			},
			want: 961,
		},
		{
			name: `Should return hash for Tuple2(nil,hello)`,
			fields: fields{
				E1: nil,
				E2: EntryString("hello"),
			},
			want: 907061831,
		},
		{
			name: `Should return hash for Tuple2(hello,nil)`,
			fields: fields{
				E1: EntryString("hello"),
				E2: nil,
			},
			want: 2349084155,
		},
		{
			name: `Should return hash for Tuple2(hello,nil)`,
			fields: fields{
				E1: EntryString("hello"),
				E2: EntryString("bye"),
			},
			want: 54247215,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t2 := Tuple2{
				E1: tt.fields.E1,
				E2: tt.fields.E2,
			}
			if got := t2.Hash(); got != tt.want {
				t.Errorf("Tuple2.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuple2_Equal(t *testing.T) {
	type fields struct {
		E1 Entry
		E2 Entry
	}
	type args struct {
		o Entry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Should equal",
			fields: fields{
				E1: EntryString("hi"),
				E2: EntryString("bye"),
			},
			args: args{
				o: Tuple2{
					E1: EntryString("hi"),
					E2: EntryString("bye")}},
			want: true,
		},
		{
			name: "Should not equal on E1",
			fields: fields{
				E1: EntryString("hi"),
				E2: EntryString("bye"),
			},
			args: args{
				o: Tuple2{
					E1: EntryString("hi2"),
					E2: EntryString("bye")}},
			want: false,
		},
		{
			name: "Should not equal on E2",
			fields: fields{
				E1: EntryString("hi"),
				E2: EntryString("bye"),
			},
			args: args{
				o: Tuple2{
					E1: EntryString("hi"),
					E2: EntryString("bye2")}},
			want: false,
		},
		{
			name: "Should not equal Tuple2 and Tuple1",
			fields: fields{
				E1: EntryString("hi"),
				E2: EntryString("bye"),
			},
			args: args{
				o: Tuple1{
					E1: EntryString("hi"),
				}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t2 := Tuple2{
				E1: tt.fields.E1,
				E2: tt.fields.E2,
			}
			if got := t2.Equal(tt.args.o); got != tt.want {
				t.Errorf("Tuple2.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuple2_Arity(t *testing.T) {
	type fields struct {
		E1 Entry
		E2 Entry
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "Should return 2 for Tuple2",
			fields: fields{},
			want:   2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t2 := Tuple2{
				E1: tt.fields.E1,
				E2: tt.fields.E2,
			}
			if got := t2.Arity(); got != tt.want {
				t.Errorf("Tuple2.Arity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuple2_ToSlice(t *testing.T) {
	type fields struct {
		E1 Entry
		E2 Entry
	}
	tests := []struct {
		name   string
		fields fields
		want   EntrySlice
	}{
		{
			name: "Should return 2-set with value",
			fields: fields{
				E1: EntryString("hi"),
				E2: EntryString("bye"),
			},
			want: EntrySlice{
				EntryString("hi"),
				EntryString("bye")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t2 := Tuple2{
				E1: tt.fields.E1,
				E2: tt.fields.E2,
			}
			if got := t2.ToSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tuple2.ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuple2_Map(t *testing.T) {
	unit := Tuple2{
		E1: EntryInt(3),
		E2: EntryInt(-7),
	}

	expected := Tuple2{
		E1: EntryInt(9),
		E2: EntryInt(49),
	}

	got := unit.Map(func(e Entry) Entry {
		return e.(EntryInt) * e.(EntryInt)
	})

	assert.EqualValues(t, expected, got)
}
