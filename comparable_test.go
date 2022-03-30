package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparableStream_Max(t *testing.T) {
	tt := map[string]struct {
		got  func() interface{}
		want interface{}
	}{
		"int": {
			got: func() interface{} {
				return CC(NewStreamFromSlice([]int{1, -2, 3, 21, 5, 8, 13}, 0).
					Map(ToR[int]), Int).
					Max()
			},
			want: 21,
		},
		"uint": {
			got: func() interface{} {
				return CC(NewStreamFromSlice([]uint{1, 2, 3, 21, 5, 8, 13}, 0).
					Map(ToR[uint]), Uint).
					Max()
			},
			want: uint(21),
		},
		"float32": {
			got: func() interface{} {
				return CC(NewStreamFromSlice([]float32{1.1, -2.2, 3.3, 21.21, 5.5, 8.8, 13.12}, 0).
					Map(ToR[float32]), Float32).
					Max()
			},
			want: float32(21.21),
		},
		"string": {
			got: func() interface{} {
				return CC(NewStreamFromSlice([]string{"D", "K", "A", "Y", "M", "O"}, 0).
					Map(ToR[string]), String).
					Max()
			},
			want: "Y",
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.got())
		})
	}
}

func TestComparableStream_Min(t *testing.T) {
	tt := map[string]struct {
		got  func() interface{}
		want interface{}
	}{
		"int": {
			got: func() interface{} {
				return CC(NewStreamFromSlice([]int{1, -2, 3, 21, 5, 8, 13}, 0).
					Map(ToR[int]), Int).
					Min()
			},
			want: -2,
		},
		"uint": {
			got: func() interface{} {
				return CC(NewStreamFromSlice([]uint{2, 1, 3, 21, 5, 8, 13}, 0).
					Map(ToR[uint]), Uint).
					Min()
			},
			want: uint(1),
		},
		"float32": {
			got: func() interface{} {
				return CC(NewStreamFromSlice([]float32{1.1, -2.2, 3.3, 21.21, 5.5, 8.8, 13.13}, 0).
					Map(ToR[float32]), Float32).
					Min()
			},
			want: float32(-2.2),
		},
		"string": {
			got: func() interface{} {
				return CC(NewStreamFromSlice([]string{"D", "K", "A", "Y", "M", "O"}, 0).
					Map(ToR[string]), String).
					Min()
			},
			want: "A",
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.got())
		})
	}
}
