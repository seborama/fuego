package fuego

import (
	"reflect"
	"testing"
)

func TestTuple1_Hash(t *testing.T) {
	type fields struct {
		E1 Entry
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name:   "Should panic for nil value Tuple1",
			fields: fields{E1: nil},
			want:   0,
		},
		{
			name:   `Should return hash for Tuple1`,
			fields: fields{E1: EntryString("hello")},
			want:   907060870,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t1 := Tuple1{
				E1: tt.fields.E1,
			}
			if got := t1.Hash(); got != tt.want {
				t.Errorf("Tuple1.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuple1_Equal(t *testing.T) {
	type fields struct {
		E1 Entry
	}
	type args struct {
		o Tuple
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Should not equal: nil != hi",
			fields: fields{
				E1: nil,
			},
			args: args{
				o: Tuple1{
					E1: EntryString("hi")}},
			want: false,
		},
		{
			name: "Should not equal: hi != nil",
			fields: fields{
				E1: EntryString("hi"),
			},
			args: args{o: nil},
			want: false,
		},
		{
			name:   "Should equal: nil == nil",
			fields: fields{E1: nil},
			args:   args{o: nil},
			want:   false,
		},
		{
			name: "Should equal: hi == hi",
			fields: fields{
				E1: EntryString("hi"),
			},
			args: args{
				o: Tuple1{
					E1: EntryString("hi")}},
			want: true,
		},
		{
			name: "Should not equal: hi != bye",
			fields: fields{
				E1: EntryString("hi"),
			},
			args: args{
				o: Tuple1{E1: EntryString("bye")}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t1 := Tuple1{
				E1: tt.fields.E1,
			}
			if got := t1.Equal(tt.args.o); got != tt.want {
				t.Errorf("Tuple1.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuple1_Arity(t *testing.T) {
	type fields struct {
		E1 Entry
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "Should return 1 for Tuple1",
			fields: fields{},
			want:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t1 := Tuple1{
				E1: tt.fields.E1,
			}
			if got := t1.Arity(); got != tt.want {
				t.Errorf("Tuple1.Arity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuple1_ToSet(t *testing.T) {
	type fields struct {
		E1 Entry
	}
	tests := []struct {
		name   string
		fields fields
		want   Set
	}{
		{
			name:   "Should return 1-set with value",
			fields: fields{E1: EntryString("hi")},
			want: NewOrderedSet().
				Insert(EntryString("hi")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t1 := Tuple1{
				E1: tt.fields.E1,
			}
			if got := t1.ToSet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tuple1.ToSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
