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
				return NewStreamFromSlice([]int{1, -2, 3, 21, 5, 8, 13}, 0).
					MapToInt(Identity[int]).
					Max()
			},
			want: 21,
		},
		"uint": {
			got: func() interface{} {
				return NewStreamFromSlice([]uint{1, 2, 3, 21, 5, 8, 13}, 0).
					MapToUint(Identity[uint]).
					Max()
			},
			want: uint(21),
		},
		"float32": {
			got: func() interface{} {
				return NewStreamFromSlice([]float32{1.1, -2.2, 3.3, 21.21, 5.5, 8.8, 13.12}, 0).
					MapToFloat32(Identity[float32]).
					Max()
			},
			want: float32(21.21),
		},
		"string": {
			got: func() interface{} {
				return NewStreamFromSlice([]string{"D", "K", "A", "Y", "M", "O"}, 0).
					MapToString(Identity[string]).
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
				return NewStreamFromSlice([]int{1, -2, 3, 21, 5, 8, 13}, 0).
					MapToInt(Identity[int]).
					Min()
			},
			want: -2,
		},
		"uint": {
			got: func() interface{} {
				return NewStreamFromSlice([]uint{2, 1, 3, 21, 5, 8, 13}, 0).
					MapToUint(Identity[uint]).
					Min()
			},
			want: uint(1),
		},
		"float32": {
			got: func() interface{} {
				return NewStreamFromSlice([]float32{1.1, -2.2, 3.3, 21.21, 5.5, 8.8, 13.13}, 0).
					MapToFloat32(Identity[float32]).
					Min()
			},
			want: float32(-2.2),
		},
		"string": {
			got: func() interface{} {
				return NewStreamFromSlice([]string{"D", "K", "A", "Y", "M", "O"}, 0).
					MapToString(Identity[string]).
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

func TestComparableStream_Sum(t *testing.T) {
	tt := map[string]struct {
		got  func() interface{}
		want interface{}
	}{
		"int": {
			got: func() interface{} {
				return NewStreamFromSlice([]int{1, -2, 3, 21, 5, 8, 13}, 0).
					MapToInt(Identity[int]).
					Sum()
			},
			want: 49,
		},
		"uint": {
			got: func() interface{} {
				return NewStreamFromSlice([]uint{2, 1, 3, 21, 5, 8, 13}, 0).
					MapToUint(Identity[uint]).
					Sum()
			},
			want: uint(53),
		},
		"float32": {
			got: func() interface{} {
				return NewStreamFromSlice([]float32{1.1, -2.2, 3.3, 21.21, 5.5, 8.8, 13.13}, 0).
					MapToFloat32(Identity[float32]).
					Sum()
			},
			want: float32(50.84),
		},
		"string": {
			got: func() interface{} {
				return NewStreamFromSlice([]string{"D", "K", "A", "Y", "M", "O"}, 0).
					MapToString(Identity[string]).
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
