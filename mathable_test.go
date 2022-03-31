package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMathableStream_Sum(t *testing.T) {
	tt := map[string]struct {
		got  func() interface{}
		want interface{}
	}{
		"int": {
			got: func() interface{} {
				return MC(NewStreamFromSlice([]int{1, -2, 3, 21, 5, 8, 13}, 0).
					Map(ToAny[int]), Int).
					Sum()
			},
			want: 49,
		},
		"uint": {
			got: func() interface{} {
				return MC(NewStreamFromSlice([]uint{2, 1, 3, 21, 5, 8, 13}, 0).
					Map(ToAny[uint]), Uint).
					Sum()
			},
			want: uint(53),
		},
		"float32": {
			got: func() interface{} {
				return MC(NewStreamFromSlice([]float32{1.1, -2.2, 3.3, 21.21, 5.5, 8.8, 13.13}, 0).
					Map(ToAny[float32]), Float32).
					Sum()
			},
			want: float32(50.84),
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.got())
		})
	}
}

func TestMathableStream_Average(t *testing.T) {
	tt := map[string]struct {
		got  func() interface{}
		want interface{}
	}{
		"int": {
			got: func() interface{} {
				return MC(NewStreamFromSlice([]int{1, -2, 3, 21, 5, 8, 13}, 0).
					Map(ToAny[int]), Int).
					Average()
			},
			want: 7,
		},
		"uint": {
			got: func() interface{} {
				return MC(NewStreamFromSlice([]uint{2, 1, 3, 21, 5, 8, 13}, 0).
					Map(ToAny[uint]), Uint).
					Average()
			},
			want: uint(7),
		},
		"float32": {
			got: func() interface{} {
				return MC(NewStreamFromSlice([]float32{1.1, -2.2, 3.3, 21.21, 5.5, 8.8, 13.13}, 0).
					Map(ToAny[float32]), Float32).
					Average()
			},
			want: float32(7.2628571429),
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.got())
		})
	}
}
