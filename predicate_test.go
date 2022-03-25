package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// intGreaterThanPredicate is a Predicate for numbers greater
// than 'rhs'.
func intGreaterThanPredicate(rhs int) Predicate[int] {
	return func(lhs int) bool {
		return lhs > rhs
	}
}

func TestNotPredicate(t *testing.T) {
	type args struct {
		p Predicate[int]
		t int
	}
	tt := map[string]struct {
		args args
		want bool
	}{
		"Should negate the predicate": {
			args: args{
				p: intGreaterThanPredicate(5),
				t: 7,
			},
			want: false,
		},
		"Should confirm the predicate": {
			args: args{
				p: intGreaterThanPredicate(10),
				t: 7,
			},
			want: true,
		},
		"Should return true when nil predicate": /* TODO: is that correct? */ {
			args: args{
				p: nil,
				t: 2,
			},
			want: true,
		},
	}

	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := tc.args.p.Negate()(tc.args.t)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFalsePredicate(t *testing.T) {
	type args struct {
		t any
	}
	tt := map[string]struct {
		args args
		want bool
	}{
		"Should return false when '123'": {
			args: args{
				t: 123,
			},
			want: false,
		},
		"Should return false when 'Hello World'": {
			args: args{
				t: "Hello World",
			},
			want: false,
		},
		"Should return false when 'true'": {
			args: args{
				t: true,
			},
			want: false,
		},
		"Should return false when 'false'": {
			args: args{
				t: false,
			},
			want: false,
		},
		"Should return false when 'nil'": {
			args: args{
				t: nil,
			},
			want: false,
		},
	}

	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := False[any]()(tc.args.t)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTruePredicate(t *testing.T) {
	type args struct {
		t any
	}
	tt := map[string]struct {
		args args
		want bool
	}{
		"Should return true when '123'": {
			args: args{
				t: 123,
			},
			want: true,
		},
		"Should return true when 'Hello World'": {
			args: args{
				t: "Hello World",
			},
			want: true,
		},
		"Should return true when 'true'": {
			args: args{
				t: true,
			},
			want: true,
		},
		"Should return true when 'false'": {
			args: args{
				t: false,
			},
			want: true,
		},
		"Should return true when 'nil'": {
			args: args{
				t: nil,
			},
			want: true,
		},
	}

	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := True[any]()(tc.args.t)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAndPredicate(t *testing.T) {
	tt := map[string]struct {
		p1   Predicate[int]
		p2   Predicate[int]
		want bool
	}{
		"Should return true for: true AND true": {
			p1:   True[int](),
			p2:   True[int](),
			want: true,
		},
		"Should return false for: true AND false": {
			p1:   True[int](),
			p2:   False[int](),
			want: false,
		},
		"Should return false for: false AND true": {
			p1:   False[int](),
			p2:   True[int](),
			want: false,
		},
		"Should return false for: false AND false": {
			p1:   False[int](),
			p2:   False[int](),
			want: false,
		},
		"Should return false for: nil AND true": {
			p1:   nil,
			p2:   True[int](),
			want: false,
		},
		"Should return false for: true AND nil": {
			p1:   True[int](),
			p2:   nil,
			want: false,
		},
		"Should return false for: nil AND false": {
			p1:   nil,
			p2:   False[int](),
			want: false,
		},
		"Should return false for: false AND nil": {
			p1:   False[int](),
			p2:   nil,
			want: false,
		},
		"Should return false for: nil AND nil": {
			p1:   nil,
			p2:   nil,
			want: false,
		},
	}
	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := tc.p1.And(tc.p2)(0)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestOrPredicate(t *testing.T) {
	tt := map[string]struct {
		p1   Predicate[int]
		p2   Predicate[int]
		want bool
	}{
		"Should return true for: true OR true": {
			p1:   True[int](),
			p2:   True[int](),
			want: true,
		},
		"Should return true for: true OR false": {
			p1:   True[int](),
			p2:   False[int](),
			want: true,
		},
		"Should return true for: false OR true": {
			p1:   False[int](),
			p2:   True[int](),
			want: true,
		},
		"Should return false for: false OR false": {
			p1:   False[int](),
			p2:   False[int](),
			want: false,
		},
		"Should return true for: nil OR true": {
			p1:   nil,
			p2:   True[int](),
			want: true,
		},
		"Should return true for: true OR nil": {
			p1:   True[int](),
			p2:   nil,
			want: true,
		},
		"Should return false for: nil OR false": {
			p1:   nil,
			p2:   False[int](),
			want: false,
		},
		"Should return false for: false OR nil": {
			p1:   False[int](),
			p2:   nil,
			want: false,
		},
		"Should return false for: nil OR nil": {
			p1:   nil,
			p2:   nil,
			want: false,
		},
	}

	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := tc.p1.Or(tc.p2)(0)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestXorPredicate(t *testing.T) {
	tt := map[string]struct {
		p1   Predicate[int]
		p2   Predicate[int]
		want bool
	}{
		"Should return false for: true XOR true": {
			p1:   True[int](),
			p2:   True[int](),
			want: false,
		},
		"Should return true for: true XOR false": {
			p1:   True[int](),
			p2:   False[int](),
			want: true,
		},
		"Should return true for: false XOR true": {
			p1:   False[int](),
			p2:   True[int](),
			want: true,
		},
		"Should return false for: false XOR false": {
			p1:   False[int](),
			p2:   False[int](),
			want: false,
		},
		"Should return true for: nil XOR true": {
			p1:   nil,
			p2:   True[int](),
			want: true,
		},
		"Should return true for: true XOR nil": {
			p1:   True[int](),
			p2:   nil,
			want: true,
		},
		"Should return false for: nil XOR false": {
			p1:   nil,
			p2:   False[int](),
			want: false,
		},
		"Should return false for: false XOR nil": {
			p1:   False[int](),
			p2:   nil,
			want: false,
		},
		"Should return false for: nil XOR nil": {
			p1:   nil,
			p2:   nil,
			want: false,
		},
	}

	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := tc.p1.Xor(tc.p2)(0)
			assert.Equal(t, tc.want, got)
		})
	}
}
