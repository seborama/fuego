package fuego

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIntStream(t *testing.T) {
	emptyChannel := make(chan EntryInt)
	populatedChannel := func() chan EntryInt {
		c := make(chan EntryInt, 1e3)
		defer close(c)
		c <- EntryInt(1)
		c <- EntryInt(2)
		c <- EntryInt(3)
		c <- EntryInt(4)
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
			name: "Should create a IntStream with a nil channel",
			args: args{s: nil},
			want: IntStream{stream: nil},
		},
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
		s []EntryInt
	}
	tests := []struct {
		name string
		args args
		want []EntryInt
	}{
		{
			name: "Should create a IntStream with an empty channel",
			args: args{s: nil},
			want: []EntryInt{},
		},
		{
			name: "Should create an empty IntStream with an empty channel",
			args: args{s: []EntryInt{}},
			want: []EntryInt{},
		},
		{
			name: "Should create a IntStream with a populated channel",
			args: args{
				s: []EntryInt{
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
			if gotIntStream := NewIntStreamFromSlice(tt.args.s).stream; gotIntStream != nil {
				got = []EntryInt{}
				for val := range gotIntStream {
					got = append(got, val)
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIntStreamFromSlice() = %+v, want %+v", got, tt.want)
			}
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
		panic  bool
	}{
		{
			name:   "Should panic when closing nil channel",
			fields: fields{stream: nil},
			panic:  true,
		},
		{
			name:   "Should close an open channel",
			fields: fields{stream: make(chan EntryInt)},
			panic:  false,
		},
		{
			name: "Should panic when closing a closed channel",
			fields: fields{stream: func() chan EntryInt {
				c := make(chan EntryInt)
				defer close(c)
				return c
			}()},
			panic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := IntStream{
				stream: tt.fields.stream,
			}
			assert.Equal(t, tt.panic, assert.Panics(&testing.T{}, func() { s.Close() }))
		})
	}
}
