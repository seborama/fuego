package fuego

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTuple3_Hash(t *testing.T) {
	type fields struct {
		E1 Entry
		E2 Entry
		E3 Entry
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "Should return hash for Tuple3(nil,nil,nil)",
			fields: fields{
				E1: nil,
				E2: nil,
				E3: nil,
			},
			want: 29791,
		},
		{
			name: `Should return hash for Tuple3(nil,hello,goodbye)`,
			fields: fields{
				E1: nil,
				E2: EntryString("hello"),
				E3: EntryString("goodbye"),
			},
			want: 2576643853,
		},
		{
			name: `Should return hash for Tuple3(hello,nil,nil)`,
			fields: fields{
				E1: EntryString("hello"),
				E2: nil,
				E3: nil,
			},
			want: 4102132069,
		},
		{
			name: `Should return hash for Tuple3(hello,bye,the end)`,
			fields: fields{
				E1: EntryString("hello"),
				E2: EntryString("bye"),
				E3: EntryString("the end"),
			},
			want: 3555541634,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t2 := Tuple3{
				E1: tt.fields.E1,
				E2: tt.fields.E2,
				E3: tt.fields.E3,
			}
			if got := t2.Hash(); got != tt.want {
				t.Errorf("Tuple3.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuple3_EqualVariations(t *testing.T) {
	for i := 0; i < 3; i++ {
		t.Run(fmt.Sprintf("Should not equal when element %d differs", i), func(t *testing.T) {
			tupleA := Tuple3{
				E1: EntryString("hi"),
				E2: EntryString("bye"),
				E3: EntryString("the end"),
			}

			values := []EntryString{
				EntryString("hi"),
				EntryString("bye"),
				EntryString("the end"),
			}
			values[i] += "-kludge"

			tupleB := Tuple3{
				E1: values[0],
				E2: values[1],
				E3: values[2],
			}

			assert.False(t, tupleA.Equal(tupleB))
		})
	}
}

func TestTuple3_Equal(t *testing.T) {
	type fields struct {
		E1 Entry
		E2 Entry
		E3 Entry
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
				E3: EntryString("the end"),
			},
			args: args{
				o: Tuple3{
					E1: EntryString("hi"),
					E2: EntryString("bye"),
					E3: EntryString("the end"),
				},
			},
			want: true,
		},
		{
			name: "Should not equal Tuple3 and Tuple1",
			fields: fields{
				E1: EntryString("hi"),
				E2: EntryString("bye"),
				E3: EntryString("the end"),
			},
			args: args{
				o: Tuple1{
					E1: EntryString("hi"),
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t2 := Tuple3{
				E1: tt.fields.E1,
				E2: tt.fields.E2,
				E3: tt.fields.E3,
			}
			if got := t2.Equal(tt.args.o); got != tt.want {
				t.Errorf("Tuple3.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuple3_Arity(t *testing.T) {
	type fields struct {
		E1 Entry
		E2 Entry
		E3 Entry
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "Should return 3 for Tuple3",
			fields: fields{},
			want:   3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t2 := Tuple3{
				E1: tt.fields.E1,
				E2: tt.fields.E2,
				E3: tt.fields.E3,
			}
			if got := t2.Arity(); got != tt.want {
				t.Errorf("Tuple3.Arity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuple3_ToSlice(t *testing.T) {
	type fields struct {
		E1 Entry
		E2 Entry
		E3 Entry
	}
	tests := []struct {
		name   string
		fields fields
		want   EntrySlice
	}{
		{
			name: "Should return 3-set with value",
			fields: fields{
				E1: EntryString("hi"),
				E2: EntryString("bye"),
				E3: EntryString("the end"),
			},
			want: EntrySlice{
				EntryString("hi"),
				EntryString("bye"),
				EntryString("the end"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t2 := Tuple3{
				E1: tt.fields.E1,
				E2: tt.fields.E2,
				E3: tt.fields.E3,
			}
			if got := t2.ToSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tuple3.ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuple3_Map(t *testing.T) {
	unit := Tuple3{
		E1: EntryInt(3),
		E2: EntryInt(-7),
		E3: EntryInt(10),
	}

	expected := Tuple3{
		E1: EntryInt(9),
		E2: EntryInt(-21),
		E3: EntryInt(30),
	}

	got := unit.Map(timesN(3))

	assert.EqualValues(t, expected, got)
}

func TestTuple3_MapMulti(t *testing.T) {
	unit := Tuple3{
		E1: EntryInt(3),
		E2: EntryInt(-7),
		E3: EntryInt(10),
	}

	expected := Tuple3{
		E1: EntryInt(9),
		E2: EntryInt(-35),
		E3: EntryInt(70),
	}

	got := unit.MapMulti(
		timesN(3),
		timesN(5),
		timesN(7),
	)

	assert.EqualValues(t, expected, got)
}
