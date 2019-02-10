package fuego

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntStream_NewIntStreamFromNilChannelPanics(t *testing.T) {
	assert.PanicsWithValue(t, PanicMissingChannel, func() { NewIntStream(nil) })
}

func TestNewIntStream(t *testing.T) {
	t.SkipNow() // TODO: IMPROVE test to permit this scenario!
	emptyChannel := make(chan EntryInt)
	populatedChannel := func() chan EntryInt {
		c := make(chan EntryInt)
		go func() {
			defer close(c)
			c <- EntryInt(1)
			c <- EntryInt(2)
			c <- EntryInt(3)
			c <- EntryInt(4)
		}()
		return c
	}()

	type args struct {
		s chan EntryInt
	}
	tests := []struct {
		name string
		args args
		want IntStream
	}{
		{
			name: "Should create an empty IntStream with an empty channel",
			args: args{s: emptyChannel},
			want: NewIntStream(emptyChannel),
		},
		{
			name: "Should create a IntStream with a populated channel",
			args: args{
				s: populatedChannel,
			},
			want: NewIntStream(populatedChannel),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIntStream(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIntStream() = %v, want %v", got, tt.want)
			}
		})
	}
}

// closeIntStream is a test utility function that closes a channel
// associated with a Stream. This is only for testing purpose.
// A formal IntStream.Close() method is an anti-pattern because
// closing a channel is a producer concern.
func closeIntStream(s IntStream) bool {
	closed := false
	func() {
		defer func() { _ = recover() }()
		close(s.stream)
		closed = true
	}()
	return closed
}

func TestNewIntStreamFromSlice(t *testing.T) {
	type args struct {
		slice []EntryInt
	}
	tests := []struct {
		name string
		args args
		want []EntryInt
	}{
		{
			name: "Should create a IntStream with an empty channel",
			args: args{slice: nil},
			want: []EntryInt{},
		},
		{
			name: "Should create an empty IntStream with an empty channel",
			args: args{slice: []EntryInt{}},
			want: []EntryInt{},
		},
		{
			name: "Should create a IntStream with a populated channel",
			args: args{
				slice: []EntryInt{
					EntryInt(1),
					EntryInt(2),
					EntryInt(3),
					EntryInt(4),
				},
			},
			want: []EntryInt{
				EntryInt(1),
				EntryInt(2),
				EntryInt(3),
				EntryInt(4),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []EntryInt
			gotIntStream := NewIntStreamFromSlice(tt.args.slice, 0)

			if streamChannel := gotIntStream.stream; streamChannel != nil {
				got = []EntryInt{}
				for val := range streamChannel {
					got = append(got, val.(EntryInt))
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIntStreamFromSlice() = %+v, want %+v", got, tt.want)
			}

			assert.False(t, closeIntStream(gotIntStream), "the stream was not closed but should have been")
		})
	}
}

func TestIntStream_AggregateFunctionPanicsWhenNilInstream(t *testing.T) {
	tests := []struct {
		name     string
		function func()
	}{
		{
			name:     "Should panic with Max(nil)",
			function: func() { IntStream{Stream{stream: nil}}.Max() },
		},
		{
			name:     "Should panic with Min(nil)",
			function: func() { IntStream{Stream{stream: nil}}.Min() },
		},
		{
			name:     "Should panic with Sum(nil)",
			function: func() { IntStream{Stream{stream: nil}}.Sum() },
		},
		{
			name:     "Should panic with Average(nil)",
			function: func() { IntStream{Stream{stream: nil}}.Average() },
		},
	}

	for _, tt := range tests {
		assert.PanicsWithValue(t, PanicMissingChannel, tt.function)
	}
}

func TestIntStream_MaxPanicsWhenEmptyInstream(t *testing.T) {
	emptyStream := func() chan EntryInt {
		c := make(chan EntryInt)
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
			function: func() { NewIntStream(emptyStream()).Max() },
		},
		{
			name:     "Should panic with Min(<empty>)",
			function: func() { NewIntStream(emptyStream()).Min() },
		},
		{
			name:     "Should panic with Sum(<empty>)",
			function: func() { NewIntStream(emptyStream()).Sum() },
		},
		{
			name:     "Should panic with Average(<empty>)",
			function: func() { NewIntStream(emptyStream()).Average() },
		},
	}

	for _, tt := range tests {
		assert.PanicsWithValue(t, PanicNoSuchElement, tt.function)
	}
}

func TestIntStream_AggregateFunction(t *testing.T) {
	singleEntryStream := func() chan EntryInt {
		c := make(chan EntryInt)
		go func() {
			defer close(c)
			c <- EntryInt(1618)
		}()
		return c
	}
	multipleEntryStream := func() chan EntryInt {
		c := make(chan EntryInt)
		go func() {
			defer close(c)
			c <- EntryInt(21)
			c <- EntryInt(-210)
			c <- EntryInt(31415)
			c <- EntryInt(10)
		}()
		return c
	}

	tests := []struct {
		name     string
		function func() EntryInt
		want     EntryInt
	}{
		{
			name:     "Max: should return the sole entry in the stream",
			function: func() EntryInt { return NewIntStream(singleEntryStream()).Max() },
			want:     EntryInt(1618),
		},
		{
			name:     "Max: Should return 31415",
			function: func() EntryInt { return NewIntStream(multipleEntryStream()).Max() },
			want:     EntryInt(31415),
		},
		{
			name:     "Min: should return the sole entry in the stream",
			function: func() EntryInt { return NewIntStream(singleEntryStream()).Min() },
			want:     EntryInt(1618),
		},
		{
			name:     "Min: Should return -210",
			function: func() EntryInt { return NewIntStream(multipleEntryStream()).Min() },
			want:     EntryInt(-210),
		},
		{
			name:     "Sum: should return the sole entry in the stream",
			function: func() EntryInt { return NewIntStream(singleEntryStream()).Sum() },
			want:     EntryInt(1618),
		},
		{
			name:     "Sum: Should return 31236",
			function: func() EntryInt { return NewIntStream(multipleEntryStream()).Sum() },
			want:     EntryInt(31236),
		},
		{
			name:     "Average: should return the sole entry in the stream",
			function: func() EntryInt { return NewIntStream(singleEntryStream()).Average() },
			want:     EntryInt(1618),
		},
		{
			name:     "Average: Should return 7809",
			function: func() EntryInt { return NewIntStream(multipleEntryStream()).Average() },
			want:     EntryInt(7809),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.function(); got != tt.want {
				t.Errorf("IntStream aggregate function = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntStream_StreamInheritance(t *testing.T) {
	s := NewIntStream(func() chan EntryInt {
		c := make(chan EntryInt)
		go func() {
			defer close(c)
			c <- EntryInt(17)
			c <- EntryInt(13)
			c <- EntryInt(23)
		}()
		return c
	}())

	got := s.Map(func(ei Entry) Entry {
		return EntryString(strconv.Itoa(int(ei.(EntryInt))))
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
