package fuego

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatStream_NewFloatStreamFromNilChannelPanics(t *testing.T) {
	assert.PanicsWithValue(t, PanicMissingChannel, func() { NewFloatStream(nil) })
}

func TestNewFloatStream(t *testing.T) {
	emptyChannel := func() chan EntryFloat {
		c := make(chan EntryFloat)
		go func() {
			defer close(c)
		}()
		return c
	}
	populatedChannel := func() chan EntryFloat {
		c := make(chan EntryFloat)
		go func() {
			defer close(c)
			c <- EntryFloat(1.23)
			c <- EntryFloat(2.789)
			c <- EntryFloat(-3.14)
			c <- EntryFloat(4.87653)
		}()
		return c
	}

	type args struct {
		s chan EntryFloat
	}
	tests := []struct {
		name string
		args args
		want FloatStream
	}{
		{
			name: "Should create an empty FloatStream with an empty channel",
			args: args{s: emptyChannel()},
			want: NewFloatStream(emptyChannel()),
		},
		{
			name: "Should create a FloatStream with a populated channel",
			args: args{
				s: populatedChannel(),
			},
			want: NewFloatStream(populatedChannel()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValues := []EntryFloat{}
			gotStream := NewFloatStream(tt.args.s)
			gotStream.ForEach(func(e Entry) {
				gotValues = append(gotValues, e.(EntryFloat))
			})

			wantValues := []EntryFloat{}
			tt.want.ForEach(func(e Entry) {
				wantValues = append(wantValues, e.(EntryFloat))
			})

			assert.IsType(t, tt.want, gotStream)
			assert.EqualValues(t, wantValues, gotValues)
		})
	}
}

// closeFloatStream is a test utility function that closes a channel
// associated with a Stream. This is only for testing purpose.
// A formal FloatStream.Close() method is an anti-pattern because
// closing a channel is a producer concern.
func closeFloatStream(s FloatStream) bool {
	closed := false
	func() {
		defer func() { _ = recover() }()
		close(s.stream)
		closed = true
	}()
	return closed
}

func TestNewFloatStreamFromSlice(t *testing.T) {
	type args struct {
		slice []EntryFloat
	}
	tests := []struct {
		name string
		args args
		want []EntryFloat
	}{
		{
			name: "Should create a FloatStream with an empty channel",
			args: args{slice: nil},
			want: []EntryFloat{},
		},
		{
			name: "Should create an empty FloatStream with an empty channel",
			args: args{slice: []EntryFloat{}},
			want: []EntryFloat{},
		},
		{
			name: "Should create a FloatStream with a populated channel",
			args: args{
				slice: []EntryFloat{
					EntryFloat(1.23),
					EntryFloat(2.789),
					EntryFloat(-3.14),
					EntryFloat(4.87653),
				},
			},
			want: []EntryFloat{
				EntryFloat(1.23),
				EntryFloat(2.789),
				EntryFloat(-3.14),
				EntryFloat(4.87653),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []EntryFloat
			gotFloatStream := NewFloatStreamFromSlice(tt.args.slice, 0)

			if streamChannel := gotFloatStream.stream; streamChannel != nil {
				got = []EntryFloat{}
				for val := range streamChannel {
					got = append(got, val.(EntryFloat))
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFloatStreamFromSlice() = %+v, want %+v", got, tt.want)
			}

			assert.False(t, closeFloatStream(gotFloatStream), "the stream was not closed but should have been")
		})
	}
}

func TestFloatStream_AggregateFunctionPanicsWhenNilInstream(t *testing.T) {
	tests := []struct {
		name     string
		function func()
	}{
		{
			name:     "Should panic with Max(nil)",
			function: func() { FloatStream{Stream{stream: nil}}.Max() },
		},
		{
			name:     "Should panic with Min(nil)",
			function: func() { FloatStream{Stream{stream: nil}}.Min() },
		},
		{
			name:     "Should panic with Sum(nil)",
			function: func() { FloatStream{Stream{stream: nil}}.Sum() },
		},
		{
			name:     "Should panic with Average(nil)",
			function: func() { FloatStream{Stream{stream: nil}}.Average() },
		},
	}

	for _, tt := range tests {
		assert.PanicsWithValue(t, PanicMissingChannel, tt.function)
	}
}

func TestFloatStream_MaxPanicsWhenEmptyInstream(t *testing.T) {
	emptyStream := func() chan EntryFloat {
		c := make(chan EntryFloat)
		go func() {
			defer close(c)
		}()
		return c
	}

	tests := []struct {
		name     string
		function func()
	}{
		{
			name:     "Should panic with Max(<empty>)",
			function: func() { NewFloatStream(emptyStream()).Max() },
		},
		{
			name:     "Should panic with Min(<empty>)",
			function: func() { NewFloatStream(emptyStream()).Min() },
		},
		{
			name:     "Should panic with Sum(<empty>)",
			function: func() { NewFloatStream(emptyStream()).Sum() },
		},
		{
			name:     "Should panic with Average(<empty>)",
			function: func() { NewFloatStream(emptyStream()).Average() },
		},
	}

	for _, tt := range tests {
		assert.PanicsWithValue(t, PanicNoSuchElement, tt.function)
	}
}

func TestFloatStream_AggregateFunction(t *testing.T) {
	singleEntryStream := func() chan EntryFloat {
		c := make(chan EntryFloat)
		go func() {
			defer close(c)
			c <- EntryFloat(1.618)
		}()
		return c
	}
	multipleEntryStream := func() chan EntryFloat {
		c := make(chan EntryFloat)
		go func() {
			defer close(c)
			c <- EntryFloat(2.1)
			c <- EntryFloat(-2.14)
			c <- EntryFloat(31.415)
			c <- EntryFloat(10)
		}()
		return c
	}

	tests := []struct {
		name     string
		function func() EntryFloat
		want     EntryFloat
	}{
		{
			name:     "Max: should return the sole entry in the stream",
			function: func() EntryFloat { return NewFloatStream(singleEntryStream()).Max() },
			want:     EntryFloat(1.618),
		},
		{
			name:     "Max: Should return 31.415",
			function: func() EntryFloat { return NewFloatStream(multipleEntryStream()).Max() },
			want:     EntryFloat(31.415),
		},
		{
			name:     "Min: should return the sole entry in the stream",
			function: func() EntryFloat { return NewFloatStream(singleEntryStream()).Min() },
			want:     EntryFloat(1.618),
		},
		{
			name:     "Min: Should return -2.14",
			function: func() EntryFloat { return NewFloatStream(multipleEntryStream()).Min() },
			want:     EntryFloat(-2.14),
		},
		{
			name:     "Sum: should return the sole entry in the stream",
			function: func() EntryFloat { return NewFloatStream(singleEntryStream()).Sum() },
			want:     EntryFloat(1.618),
		},
		{
			name:     "Sum: Should return 41.375",
			function: func() EntryFloat { return NewFloatStream(multipleEntryStream()).Sum() },
			want:     EntryFloat(41.375),
		},
		{
			name:     "Average: should return the sole entry in the stream",
			function: func() EntryFloat { return NewFloatStream(singleEntryStream()).Average() },
			want:     EntryFloat(1.618),
		},
		{
			name:     "Average: Should return 10.34375",
			function: func() EntryFloat { return NewFloatStream(multipleEntryStream()).Average() },
			want:     EntryFloat(10.34375),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.function(); got != tt.want {
				t.Errorf("FloatStream aggregate function = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloatStream_StreamInheritance(t *testing.T) {
	s := NewFloatStream(func() chan EntryFloat {
		c := make(chan EntryFloat)
		go func() {
			defer close(c)
			c <- EntryFloat(17)
			c <- EntryFloat(13)
			c <- EntryFloat(23)
		}()
		return c
	}())

	got := s.Map(func(ei Entry) Entry {
		return EntryString(strconv.Itoa(int(ei.(EntryFloat))))
	}).
		Intersperse(EntryString(" - ")).
		Reduce(func(i, j Entry) Entry {
			iStr := i.(EntryString)
			jStr := j.(EntryString)
			return iStr + jStr
		})
	expected := EntryString("17 - 13 - 23")
	assert.Equal(t, expected, got)
}
