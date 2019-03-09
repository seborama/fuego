package fuego

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPStream(t *testing.T) {
	type args struct {
		channels []chan Entry
	}
	tests := []struct {
		name string
		args args
		want []chan Entry
	}{
		{
			name: "Should add empty set",
			args: args{
				channels: []chan Entry{},
			},
			want: []chan Entry{},
		},
		{
			name: "Should added a set of 1",
			args: args{
				channels: channels1(),
			},
			want: channels1(),
		},
		{
			name: "Should added a set of 4",
			args: args{
				channels: channels4(),
			},
			want: channels4(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCStream(tt.args.channels)
			assert.Equal(t, len(tt.want), len(got.streams))

			wantSliceOfEntrySlices := consumeChannels(tt.want)
			for idx, gotStream := range got.streams {
				assert.EqualValues(t, wantSliceOfEntrySlices[idx], gotStream.ToSlice())
			}
		})
	}
}

func TestPStream_AddStreamFromChannels(t *testing.T) {
	type fields struct {
		streams []Stream
	}
	type args struct {
		channels []chan Entry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []chan Entry
	}{
		{
			name: "Should add empty set to CStream",
			fields: fields{
				streams: NewCStream(channels1()).streams,
			},
			args: args{
				channels: []chan Entry{},
			},
			want: channels1(),
		},
		{
			name: "Should add set of 1 to CStream",
			fields: fields{
				streams: NewCStream(channels1()).streams,
			},
			args: args{
				channels: channels1(),
			},
			want: append(channels1(), channels1()...),
		},
		{
			name: "Should add set of 4 to CStream",
			fields: fields{
				streams: NewCStream(channels1()).streams,
			},
			args: args{
				channels: channels4(),
			},
			want: append(channels1(), channels4()...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := CStream{
				streams: tt.fields.streams,
			}

			got := cs.AddStreamFromChannels(tt.args.channels)
			assert.Equal(t, len(tt.want), len(got.streams))

			wantSliceOfEntrySlices := consumeChannels(tt.want)
			for idx, gotStream := range got.streams {
				assert.EqualValues(t, wantSliceOfEntrySlices[idx], gotStream.ToSlice())
			}
		})
	}
}

func TestPStream_AddStreamFromSlices(t *testing.T) {
	type fields struct {
		streams []Stream
	}
	type args struct {
		slices []EntrySlice
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []EntrySlice
	}{
		{
			name: "Should add empty set to CStream",
			fields: fields{
				streams: NewCStream(channels1()).streams,
			},
			args: args{
				slices: []EntrySlice{},
			},
			want: []EntrySlice{slice1},
		},
		{
			name: "Should add set of 1 to CStream",
			fields: fields{
				streams: NewCStream(channels1()).streams,
			},
			args: args{
				slices: []EntrySlice{slice2},
			},
			want: []EntrySlice{slice1, slice2},
		},
		{
			name: "Should add set of 4 to CStream",
			fields: fields{
				streams: NewCStream(channels1()).streams,
			},
			args: args{
				slices: []EntrySlice{slice1, slice2, slice3, slice4},
			},
			want: []EntrySlice{slice1, slice1, slice2, slice3, slice4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := CStream{
				streams: tt.fields.streams,
			}

			got := cs.AddStreamFromSlices(tt.args.slices, 0)
			assert.Equal(t, len(tt.want), len(got.streams))
			for idx, gotStream := range got.streams {
				assert.EqualValues(t, tt.want[idx], gotStream.ToSlice())
			}
		})
	}
}

func TestCStream_AddStreams(t *testing.T) {
	type fields struct {
		streams []Stream
	}
	type args struct {
		streams []Stream
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   CStream
	}{
		{
			name: "Should add empty set to CStream",
			fields: fields{
				streams: NewCStream(channels1()).streams,
			},
			args: args{
				streams: []Stream{},
			},
			want: CStream{[]Stream{
				NewStream(channels1()[0]),
			}},
		},
		{
			name: "Should add set of 1 to CStream",
			fields: fields{
				streams: NewCStream(channels1()).streams,
			},
			args: args{
				streams: []Stream{NewStreamFromSlice(slice2, 0)},
			},
			want: CStream{[]Stream{
				NewStream(channels1()[0]),
				NewStreamFromSlice(slice2, 0),
			}},
		},
		{
			name: "Should add set of 4 to CStream",
			fields: fields{
				streams: NewCStream(channels1()).streams,
			},
			args: args{
				streams: []Stream{
					NewStreamFromSlice(slice1, 0),
					NewStreamFromSlice(slice2, 0),
					NewStreamFromSlice(slice3, 0),
					NewStreamFromSlice(slice4, 0),
				},
			},
			want: CStream{[]Stream{
				NewStream(channels1()[0]),
				NewStreamFromSlice(slice1, 0),
				NewStreamFromSlice(slice2, 0),
				NewStreamFromSlice(slice3, 0),
				NewStreamFromSlice(slice4, 0),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := CStream{
				streams: tt.fields.streams,
			}

			got := cs.AddStreams(tt.args.streams)
			assert.Equal(t, len(tt.want.streams), len(got.streams))
			for idx, gotStream := range got.streams {
				assert.EqualValues(t, tt.want.streams[idx].ToSlice(), gotStream.ToSlice())
			}
		})
	}
}

func TestCStream_ForEach(t *testing.T) {
	channel := make(chan struct{})

	consumer := func(e Entry) {
		channel <- struct{}{}
	}

	count := 0

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for range channel {
			count++
		}
	}()

	cs := CStream{
		streams: []Stream{
			NewStreamFromSlice(EntrySlice{
				EntryString("a"),
				EntryString("b"),
				EntryString("c"),
			}, 0),
			NewStreamFromSlice(EntrySlice{
				EntryInt(1),
				EntryInt(2),
				EntryInt(3),
			}, 0),
			NewStreamFromSlice(EntrySlice{
				EntryBool(true),
				EntryBool(false),
			}, 0),
			NewStreamFromSlice(EntrySlice{
				EntryFloat(1.2),
				EntryFloat(3.4),
			}, 0),
		},
	}
	cs.ForEach(consumer)
	close(channel)

	wg.Wait()

	assert.Equal(t, 10, count)
}

func TestCStream_Filter(t *testing.T) {
	predicate := func(e Entry) bool {
		if _, ok := e.(EntryString); ok {
			return true
		}
		if _, ok := e.(EntryBool); ok {
			return true
		}
		return false
	}

	cs := CStream{
		streams: []Stream{
			NewStreamFromSlice(EntrySlice{
				EntryString("a"),
				EntryString("b"),
				EntryString("c"),
			}, 0),
			NewStreamFromSlice(EntrySlice{
				EntryInt(1),
				EntryInt(2),
				EntryInt(3),
			}, 0),
			NewStreamFromSlice(EntrySlice{
				EntryBool(true),
				EntryBool(false),
			}, 0),
			NewStreamFromSlice(EntrySlice{
				EntryFloat(1.2),
				EntryFloat(3.4),
			}, 0),
		},
	}
	cstream := cs.Filter(predicate)

	expected := []Stream{
		NewStreamFromSlice(EntrySlice{
			EntryString("a"),
			EntryString("b"),
			EntryString("c"),
		}, 0),
		{},
		NewStreamFromSlice(EntrySlice{
			EntryBool(true),
			EntryBool(false),
		}, 0),
		{},
	}

	require.Equal(t, len(expected), len(cstream.streams))

	for idx, gotCStream := range cstream.streams {
		assert.EqualValues(t, expected[idx].ToSlice(), gotCStream.ToSlice())
	}
}

func consumeChannels(channels []chan Entry) []EntrySlice {
	es := make([]EntrySlice, len(channels))

	for idx, channel := range channels {
		for entry := range channel {
			es[idx] = append(es[idx], entry)
		}
	}

	return es
}

func makeChanEntry(es EntrySlice) chan Entry {
	c := make(chan Entry)
	go func() {
		defer close(c)
		for _, val := range es {
			c <- val
		}
	}()
	return c
}

func channels1() []chan Entry {
	return []chan Entry{
		makeChanEntry(slice1),
	}
}

func channels4() []chan Entry {
	return []chan Entry{
		makeChanEntry(slice1),
		makeChanEntry(slice2),
		makeChanEntry(slice3),
		makeChanEntry(slice4),
	}
}

func makeEntrySlice() EntrySlice {
	es := EntrySlice{}
	for i := 1; i <= rand.Intn(10); i++ {
		es = append(es, EntryInt(rand.Intn(1000)))
	}
	return es
}

var slice1 = EntrySlice{
	makeEntrySlice(),
}

var slice2 = EntrySlice{
	makeEntrySlice(),
}

var slice3 = EntrySlice{
	makeEntrySlice(),
}

var slice4 = EntrySlice{
	makeEntrySlice(),
}
