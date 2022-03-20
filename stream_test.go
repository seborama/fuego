package fuego

import (
	"reflect"
	"testing"
	"time"
)

func functionTimesTwo() Function[int, int] {
	return func(i int) int {
		return 2 * i
	}
}

func functionSlowTimesTwo() Function[int, int] {
	return func(i int) int {
		time.Sleep(50 * time.Millisecond)
		return 2 * i
	}
}

func entryIntEqualsTo(number int) Function[int, bool] {
	return func(subject int) bool {
		return number == subject
	}
}

func TestStream_Map(t *testing.T) {
	tt := map[string]struct {
		stream Stream[int]
		mapper Function[int, int]
		want   []int
	}{
		"Should return an empty Stream": {

			stream: Stream[int]{stream: nil},
			mapper: functionTimesTwo(),
			want:   []int{},
		},
		"Should return a Stream of doubled integers": {
			stream: Stream[int]{
				stream: func() chan int {
					c := make(chan int, 1)
					go func() {
						defer close(c)
						c <- 1
						c <- 3
						c <- 2
					}()
					return c
				}()},
			mapper: functionTimesTwo(),
			want: []int{
				2,
				6,
				4,
			},
		},
	}

	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			var got []int
			var resultStream Stream[int] = tc.stream.Map(tc.mapper)
			if gotStream := resultStream.stream; gotStream != nil {
				for val := range gotStream {
					got = append(got, val)
				}
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Stream.Map() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStream_Map_Concurrent(t *testing.T) {
	const numEntries = 300
	const concurrencyLevel = numEntries / 10

	sourceStream := func() chan Entry {
		c := make(chan Entry, 10)
		go func() {
			defer close(c)
			for i := 0; i < numEntries; i++ {
				c <- EntryInt(i)
			}
		}()
		return c
	}()

	want := func() EntrySlice {
		es := EntrySlice{}
		for i := 0; i < numEntries; i++ {
			es = es.Append(functionTimesTwo()(EntryInt(i)))
		}
		return es
	}()

	unitStream := Stream{
		stream: sourceStream,
	}

	// functionSlowTimesTwo: use slow function to illustrate the performance improvement
	gotStream := unitStream.Concurrent(concurrencyLevel).Map(functionSlowTimesTwo()).ToSlice()

	if !reflect.DeepEqual(gotStream, want) {
		t.Errorf("Stream.Map() = %v, want %v", gotStream, want)
	}
}
