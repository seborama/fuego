package fuego

import (
	"reflect"
	"testing"
)

func TestTuple0_Hash(t *testing.T) {
	tests := []struct {
		name string
		t    Tuple0
		want uint32
	}{
		{
			name: "Should return 1 for Tuple0",
			t:    Tuple0{},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Hash(); got != tt.want {
				t.Errorf("Tuple0.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuple0_Equal(t *testing.T) {
	refTuple := Tuple0{}
	type args struct {
		o Entry
	}
	tests := []struct {
		name string
		t    Tuple0
		args args
		want bool
	}{
		{
			name: "Should not be equal",
			t:    refTuple,
			args: args{o: refTuple}, // passing by value, not by reference
			want: false,
		},
		{
			name: "Should not be equal",
			t:    Tuple0{},
			args: args{o: Tuple0{}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Equal(tt.args.o); got != tt.want {
				t.Errorf("Tuple0.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuple0_Arity(t *testing.T) {
	tests := []struct {
		name string
		t    Tuple0
		want int
	}{
		{
			name: "Should return 0 for Tuple0",
			t:    Tuple0{},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Arity(); got != tt.want {
				t.Errorf("Tuple0.Arity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTuple0_ToSlice(t *testing.T) {
	tests := []struct {
		name string
		t    Tuple0
		want EntrySlice
	}{
		{
			name: "Should return empty set for Tuple0",
			t:    Tuple0{},
			want: EntrySlice{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.ToSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tuple0.ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
