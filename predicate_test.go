package fuego

import (
	"reflect"
	"testing"
)

type intGreaterThanPredicate struct {
	number int
}

func newIntGreaterThanPredicate(number int) intGreaterThanPredicate {
	return intGreaterThanPredicate{
		number: number,
	}
}

func (p intGreaterThanPredicate) Test(t interface{}) bool {
	return t.(int) > p.number
}

type intLessThanPredicate struct {
	number int
}

func newIntLessThanPredicate(number int) intLessThanPredicate {
	return intLessThanPredicate{
		number: number,
	}
}

func (p intLessThanPredicate) Test(t interface{}) bool {
	return t.(int) < p.number
}

func TestNotPredicate(t *testing.T) {
	type args struct {
		p Predicate
		t int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should negate the predicate",
			args: args{
				p: newIntGreaterThanPredicate(5),
				t: 7,
			},
			want: false,
		},
		{
			name: "Should negate the predicate",
			args: args{
				p: newIntGreaterThanPredicate(5),
				t: 2,
			},
			want: true,
		},
		{
			name: "Should return true when nil predicate",
			args: args{
				p: nil,
				t: 2,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Not(tt.args.p).Test(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Not() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFalsePredicate(t *testing.T) {
	type args struct {
		t interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should return false when '123'",
			args: args{
				t: 123,
			},
			want: false,
		},
		{
			name: "Should return false when 'Hello World'",
			args: args{
				t: "Hello World",
			},
			want: false,
		},
		{
			name: "Should return false when 'true'",
			args: args{
				t: true,
			},
			want: false,
		},
		{
			name: "Should return false when 'false'",
			args: args{
				t: false,
			},
			want: false,
		},
		{
			name: "Should return false when '[]int{1,2,3}'",
			args: args{
				t: []int{1, 2, 3},
			},
			want: false,
		},
		{
			name: "Should return false when 'struct {}'",
			args: args{
				t: struct{}{},
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
			if got := False().Test(tt.args.t); !got == tt.want {
				t.Errorf("False.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTruePredicate(t *testing.T) {
	type args struct {
		t interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should return true when '123'",
			args: args{
				t: 123,
			},
			want: true,
		},
		{
			name: "Should return true when 'Hello World'",
			args: args{
				t: "Hello World",
			},
			want: true,
		},
		{
			name: "Should return true when 'true'",
			args: args{
				t: true,
			},
			want: true,
		},
		{
			name: "Should return true when 'false'",
			args: args{
				t: false,
			},
			want: true,
		},
		{
			name: "Should return true when '[]int{1,2,3}'",
			args: args{
				t: []int{1, 2, 3},
			},
			want: true,
		},
		{
			name: "Should return true when 'struct {}'",
			args: args{
				t: struct{}{},
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
			if got := True().Test(tt.args.t); !got == tt.want {
				t.Errorf("False.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAndPredicate(t *testing.T) {
	type args struct {
		p1 Predicate
		p2 Predicate
		t  int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should return true for: true AND true",
			args: args{
				p1: True(),
				p2: True(),
				t:  1,
			},
			want: true,
		},
		{
			name: "Should return false for: true AND false",
			args: args{
				p1: True(),
				p2: False(),
				t:  1,
			},
			want: false,
		},
		{
			name: "Should return false for: false AND true",
			args: args{
				p1: False(),
				p2: True(),
				t:  1,
			},
			want: false,
		},
		{
			name: "Should return false for: false AND false",
			args: args{
				p1: False(),
				p2: False(),
				t:  1,
			},
			want: false,
		},
		{
			name: "Should return false for: nil AND true",
			args: args{
				p1: nil,
				p2: True(),
				t:  1,
			},
			want: false,
		},
		{
			name: "Should return false for: true AND nil",
			args: args{
				p1: True(),
				p2: nil,
				t:  1,
			},
			want: false,
		},
		{
			name: "Should return false for: nil AND false",
			args: args{
				p1: nil,
				p2: False(),
				t:  1,
			},
			want: false,
		},
		{
			name: "Should return false for: false AND nil",
			args: args{
				p1: False(),
				p2: nil,
				t:  1,
			},
			want: false,
		},
		{
			name: "Should return false for: nil AND nil",
			args: args{
				p1: nil,
				p2: nil,
				t:  1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := And(tt.args.p1, tt.args.p2).Test(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("And() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrPredicate(t *testing.T) {
	type args struct {
		p1 Predicate
		p2 Predicate
		t  int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should return true for: true OR true",
			args: args{
				p1: True(),
				p2: True(),
				t:  1,
			},
			want: true,
		},
		{
			name: "Should return true for: true OR false",
			args: args{
				p1: True(),
				p2: False(),
				t:  1,
			},
			want: true,
		},
		{
			name: "Should return true for: false OR true",
			args: args{
				p1: False(),
				p2: True(),
				t:  1,
			},
			want: true,
		},
		{
			name: "Should return false for: false OR false",
			args: args{
				p1: False(),
				p2: False(),
				t:  1,
			},
			want: false,
		},
		{
			name: "Should return true for: nil OR true",
			args: args{
				p1: nil,
				p2: True(),
				t:  1,
			},
			want: true,
		},
		{
			name: "Should return true for: true OR nil",
			args: args{
				p1: True(),
				p2: nil,
				t:  1,
			},
			want: true,
		},
		{
			name: "Should return false for: nil OR false",
			args: args{
				p1: nil,
				p2: False(),
				t:  1,
			},
			want: false,
		},
		{
			name: "Should return false for: false OR nil",
			args: args{
				p1: False(),
				p2: nil,
				t:  1,
			},
			want: false,
		},
		{
			name: "Should return false for: nil OR nil",
			args: args{
				p1: nil,
				p2: nil,
				t:  1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Or(tt.args.p1, tt.args.p2).Test(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Or() = %v, want %v", got, tt.want)
			}
		})
	}
}
