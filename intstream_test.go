package fuego

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntStream_NewIntStreamFromNilChannelPanics(t *testing.T) {
	assert.PanicsWithValue(t, PanicMissingChannel, func() { NewIntStream(nil) })
}

func TestNewIntStream(t *testing.T) {
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
			want: IntStream{stream: emptyChannel},
		},
		{
			name: "Should create a IntStream with a populated channel",
			args: args{
				s: populatedChannel,
			},
			want: IntStream{
				stream: populatedChannel,
			},
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
			gotIntStream := NewIntStreamFromSlice(tt.args.slice)

			if streamChannel := gotIntStream.stream; streamChannel != nil {
				got = []EntryInt{}
				for val := range streamChannel {
					got = append(got, val)
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIntStreamFromSlice() = %+v, want %+v", got, tt.want)
			}

			assert.False(t, gotIntStream.Close(), "the stream was not closed but should have been")
		})
	}
}

func TestIntStream_Close(t *testing.T) {
	type fields struct {
		stream chan EntryInt
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "Should return false when closing nil channel",
			fields: fields{stream: nil},
			want:   false,
		},
		{
			name:   "Should return true when closing an open channel",
			fields: fields{stream: make(chan EntryInt)},
			want:   true,
		},
		{
			name: "Should not panic when closing a closed channel",
			fields: fields{
				stream: func() chan EntryInt {
					c := make(chan EntryInt)
					defer close(c)
					return c
				}(),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := IntStream{
				stream: tt.fields.stream,
			}
			assert.Equal(t, tt.want, s.Close())
		})
	}
}

func TestIntStream_MaxPanicsWhenNilInstream(t *testing.T) {
	assert.PanicsWithValue(t, PanicMissingChannel, func() { IntStream{stream: nil}.Max() })
}

func TestIntStream_MaxPanicsWhenEmptyInstream(t *testing.T) {
	assert.PanicsWithValue(t, PanicNoSuchElement, func() {
		NewIntStream(func() chan EntryInt {
			c := make(chan EntryInt)
			go func() {
				defer close(c)
			}()
			return c
		}()).Max()
	})
}

func TestIntStream_Max(t *testing.T) {
	type fields struct {
		stream chan EntryInt
	}
	tests := []struct {
		name   string
		fields fields
		want   EntryInt
	}{
		{
			name: "Should return the sole entry in the stream",
			fields: fields{
				stream: func() chan EntryInt {
					c := make(chan EntryInt)
					go func() {
						defer close(c)
						c <- EntryInt(1618)
					}()
					return c
				}(),
			},
			want: EntryInt(1618),
		},
		{
			name: "Should return 31415",
			fields: fields{
				stream: func() chan EntryInt {
					c := make(chan EntryInt)
					go func() {
						defer close(c)
						c <- EntryInt(-210)
						c <- EntryInt(21)
						c <- EntryInt(31415)
						c <- EntryInt(10)
					}()
					return c
				}(),
			},
			want: EntryInt(31415),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := IntStream{
				stream: tt.fields.stream,
			}
			if got := is.Max(); got != tt.want {
				t.Errorf("IntStream.Max() = %v, want %v", got, tt.want)
			}
		})
	}
}
