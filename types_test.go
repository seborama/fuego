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
