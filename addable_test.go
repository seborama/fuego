package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparableStream_Sum(t *testing.T) {
	tt := map[string]struct {
		got  func() interface{}
		want interface{}
	}{
		"int": {
			got: func() interface{} {
				return AC(NewStreamFromSlice([]int{1, -2, 3, 21, 5, 8, 13}, 0).
					Map(ToR[int]), Int).
					Sum()
			},
			want: 49,
		},
		"uint": {
			got: func() interface{} {
				return AC(NewStreamFromSlice([]uint{2, 1, 3, 21, 5, 8, 13}, 0).
					Map(ToR[uint]), Uint).
					Sum()
			},
			want: uint(53),
		},
		"float32": {
			got: func() interface{} {
				return AC(NewStreamFromSlice([]float32{1.1, -2.2, 3.3, 21.21, 5.5, 8.8, 13.13}, 0).
					Map(ToR[float32]), Float32).
					Sum()
			},
			want: float32(50.84),
		},
		"string": {
			got: func() interface{} {
				return AC(NewStreamFromSlice([]string{"D", "K", "A", "Y", "M", "O"}, 0).
					Map(ToR[string]), String).
					Sum()
			},
			want: "DKAYMO",
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.got())
		})
	}
}
