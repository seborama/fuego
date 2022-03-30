package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	tt := map[string]struct {
		max  func() interface{}
		want interface{}
	}{
		"Max(-1,3)": {
			max:  func() interface{} { return Max(-1, 3) },
			want: 3,
		},
		"Max(3,-1)": {
			max:  func() interface{} { return Max(3, -1) },
			want: 3,
		},
		"Max(-1.2,3.4)": {
			max:  func() interface{} { return Max(-1.2, 3.4) },
			want: 3.4,
		},
		"Max(3.4,-1.2)": {
			max:  func() interface{} { return Max(3.4, -1.2) },
			want: 3.4,
		},
		`Max("abc","zyx")`: {
			max:  func() interface{} { return Max("abc", "zyx") },
			want: "zyx",
		},
		`Max("zyx","abc")`: {
			max:  func() interface{} { return Max("zyx", "abc") },
			want: "zyx",
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			got := tc.max()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMin(t *testing.T) {
	tt := map[string]struct {
		min  func() interface{}
		want interface{}
	}{
		"Min(-1,3)": {
			min:  func() interface{} { return Min(-1, 3) },
			want: -1,
		},
		"Min(3,-1)": {
			min:  func() interface{} { return Min(3, -1) },
			want: -1,
		},
		"Min(-1.2,3.4)": {
			min:  func() interface{} { return Min(-1.2, 3.4) },
			want: -1.2,
		},
		"Min(3.4,-1.2)": {
			min:  func() interface{} { return Min(3.4, -1.2) },
			want: -1.2,
		},
		`Min("abc","zyx")`: {
			min:  func() interface{} { return Min("abc", "zyx") },
			want: "abc",
		},
		`Min("zyx","abc")`: {
			min:  func() interface{} { return Min("zyx", "abc") },
			want: "abc",
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			got := tc.min()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSum(t *testing.T) {
	tt := map[string]struct {
		sum  func() interface{}
		want interface{}
	}{
		"Sum(-1,3)": {
			sum:  func() interface{} { return Sum(-1, 3) },
			want: 2,
		},
		"Sum(3,-1)": {
			sum:  func() interface{} { return Sum(3, -1) },
			want: 2,
		},
		"Sum(-1.2,3.4)": {
			sum:  func() interface{} { return Sum(-1.2, 3.4) },
			want: 2.2,
		},
		"Sum(3.4,-1.2)": {
			sum:  func() interface{} { return Sum(3.4, -1.2) },
			want: 2.2,
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			got := tc.sum()
			assert.Equal(t, tc.want, got)
		})
	}
}
