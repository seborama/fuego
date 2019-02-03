package fuego

import (
	"reflect"
	"testing"
)

// intGreaterThanPredicate is a Predicate for numbers greater
// than 'rhs'.
func intGreaterThanPredicate(rhs int) Predicate {
	return func(lhs Entry) bool {
		return int(lhs.(EntryInt)) > rhs
	}
}

func TestNotPredicate(t *testing.T) {
	type args struct {
		p Predicate
		t Entry
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should negate the predicate",
			args: args{
				p: intGreaterThanPredicate(5),
				t: EntryInt(7),
			},
			want: false,
		},
		{
			name: "Should negate the predicate",
			args: args{
				p: intGreaterThanPredicate(5),
				t: EntryInt(2),
			},
			want: true,
		},
		{
			name: "Should return true when nil predicate",
			args: args{
				p: nil,
				t: EntryInt(2),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.p.Not()(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Not() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFalsePredicate(t *testing.T) {
	type args struct {
		t Entry
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should return false when '123'",
			args: args{
				t: EntryInt(123),
			},
			want: false,
		},
		{
			name: "Should return false when 'Hello World'",
			args: args{
				t: EntryString("Hello World"),
			},
			want: false,
		},
		{
			name: "Should return false when 'true'",
			args: args{
				t: EntryBool(true),
			},
			want: false,
		},
		{
			name: "Should return false when 'false'",
			args: args{
				t: EntryBool(false),
			},
			want: false,
		},
		{
			name: "Should return false when 'nil'",
			args: args{
				t: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := False(tt.args.t); !got == tt.want {
				t.Errorf("False = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTruePredicate(t *testing.T) {
	type args struct {
		t Entry
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should return true when '123'",
			args: args{
				t: EntryInt(123),
			},
			want: true,
		},
		{
			name: "Should return true when 'Hello World'",
			args: args{
				t: EntryString("Hello World"),
			},
			want: true,
		},
		{
			name: "Should return true when 'true'",
			args: args{
				t: EntryBool(true),
			},
			want: true,
		},
		{
			name: "Should return true when 'false'",
			args: args{
				t: EntryBool(false),
			},
			want: true,
		},
		{
			name: "Should return true when 'nil'",
			args: args{
				t: nil,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := True(tt.args.t); !got == tt.want {
				t.Errorf("False.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAndPredicate(t *testing.T) {
	type args struct {
		p1 Predicate
		p2 Predicate
		t  Entry
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should return true for: true AND true",
			args: args{
				p1: True,
				p2: True,
				t:  EntryInt(1),
			},
			want: true,
		},
		{
			name: "Should return false for: true AND false",
			args: args{
				p1: True,
				p2: False,
				t:  EntryInt(1),
			},
			want: false,
		},
		{
			name: "Should return false for: false AND true",
			args: args{
				p1: False,
				p2: True,
				t:  EntryInt(1),
			},
			want: false,
		},
		{
			name: "Should return false for: false AND false",
			args: args{
				p1: False,
				p2: False,
				t:  EntryInt(1),
			},
			want: false,
		},
		{
			name: "Should return false for: nil AND true",
			args: args{
				p1: nil,
				p2: True,
				t:  EntryInt(1),
			},
			want: false,
		},
		{
			name: "Should return false for: true AND nil",
			args: args{
				p1: True,
				p2: nil,
				t:  EntryInt(1),
			},
			want: false,
		},
		{
			name: "Should return false for: nil AND false",
			args: args{
				p1: nil,
				p2: False,
				t:  EntryInt(1),
			},
			want: false,
		},
		{
			name: "Should return false for: false AND nil",
			args: args{
				p1: False,
				p2: nil,
				t:  EntryInt(1),
			},
			want: false,
		},
		{
			name: "Should return false for: nil AND nil",
			args: args{
				p1: nil,
				p2: nil,
				t:  EntryInt(1),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.p1.And(tt.args.p2)(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("And() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrPredicate(t *testing.T) {
	type args struct {
		p1 Predicate
		p2 Predicate
		t  Entry
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should return true for: true OR true",
			args: args{
				p1: True,
				p2: True,
				t:  EntryInt(1),
			},
			want: true,
		},
		{
			name: "Should return true for: true OR false",
			args: args{
				p1: True,
				p2: False,
				t:  EntryInt(1),
			},
			want: true,
		},
		{
			name: "Should return true for: false OR true",
			args: args{
				p1: False,
				p2: True,
				t:  EntryInt(1),
			},
			want: true,
		},
		{
			name: "Should return false for: false OR false",
			args: args{
				p1: False,
				p2: False,
				t:  EntryInt(1),
			},
			want: false,
		},
		{
			name: "Should return true for: nil OR true",
			args: args{
				p1: nil,
				p2: True,
				t:  EntryInt(1),
			},
			want: true,
		},
		{
			name: "Should return true for: true OR nil",
			args: args{
				p1: True,
				p2: nil,
				t:  EntryInt(1),
			},
			want: true,
		},
		{
			name: "Should return false for: nil OR false",
			args: args{
				p1: nil,
				p2: False,
				t:  EntryInt(1),
			},
			want: false,
		},
		{
			name: "Should return false for: false OR nil",
			args: args{
				p1: False,
				p2: nil,
				t:  EntryInt(1),
			},
			want: false,
		},
		{
			name: "Should return false for: nil OR nil",
			args: args{
				p1: nil,
				p2: nil,
				t:  EntryInt(1),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.p1.Or(tt.args.p2)(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionPredicate(t *testing.T) {
	type args struct {
		f Function
	}
	tests := []struct {
		name string
		args args
		want Predicate
	}{
		{
			name: "Should apply the predicate 'equals to 25' and return false with 100",
			args: args{entryIntEqualsTo(25)},
			want: False,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FunctionPredicate(tt.args.f); !reflect.DeepEqual(got(EntryInt(100)), tt.want(nil)) {
				t.Errorf("FunctionPredicate() = %v, want %v", got, tt.want)
			}
		})
	}
}
