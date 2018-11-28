package fuego

import (
	"reflect"
	"testing"
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
			name: "Should panic for nil value Tuple2(nil,nil)",
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
		o Tuple
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Should equal with E1 nil and args nil",
			fields: fields{
				E1: nil,
				E2: EntryString("bye"),
			},
			args: args{o: nil},
			want: false,
		},
		{
			name: "Should NOT equal with E1 nil",
			fields: fields{
				E1: nil,
				E2: EntryString("hi"),
			},
			args: args{
				o: Tuple2{
					E1: nil,
					E2: EntryString("bye")}},
			want: false,
		},
		{
			name: "Should equal with E1 nil",
			fields: fields{
				E1: nil,
				E2: EntryString("bye"),
			},
			args: args{
				o: Tuple2{
					E1: nil,
					E2: EntryString("bye")}},
			want: true,
		},
		{
			name: "Should NOT equal with E2 nil",
			fields: fields{
				E1: EntryString("hi"),
				E2: nil,
			},
			args: args{
				o: Tuple2{
					E1: EntryString("hi"),
					E2: EntryString("bye")}},
			want: false,
		},
		{
			name: "Should NOT equal with E2 nil",
			fields: fields{
				E1: EntryString("hi"),
				E2: nil,
			},
			args: args{
				o: Tuple2{
					E1: EntryString("bye"),
					E2: nil}},
			want: false,
		},
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
			name: "Should not equal on E1 difference",
			fields: fields{
				E1: EntryString("hi"),
				E2: EntryString("bye"),
			},
			args: args{
				o: Tuple2{
					E1: EntryString("different"),
					E2: EntryString("bye")}},
			want: false,
		},
		{
			name: "Should not equal on E2 difference",
			fields: fields{
				E1: EntryString("hi"),
				E2: EntryString("bye"),
			},
			args: args{
				o: Tuple2{
					E1: EntryString("hi"),
					E2: EntryString("different")}},
			want: false,
		},
		{
			name: "Should not equal on E1 and E2 difference",
			fields: fields{
				E1: EntryString("hi"),
				E2: EntryString("bye"),
			},
			args: args{
				o: Tuple2{
					E1: EntryString("different 1"),
					E2: EntryString("different 2")}},
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

func TestTuple2_ToSet(t *testing.T) {
	type fields struct {
		E1 Entry
		E2 Entry
	}
	tests := []struct {
		name   string
		fields fields
		want   Set
	}{
		{
			name: "Should return 2-set with value",
			fields: fields{
				E1: EntryString("hi"),
				E2: EntryString("bye"),
			},
			want: NewOrderedSet().
				Insert(EntryString("hi")).
				Insert(EntryString("bye")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t2 := Tuple2{
				E1: tt.fields.E1,
				E2: tt.fields.E2,
			}
			if got := t2.ToSet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tuple2.ToSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
