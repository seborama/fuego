package fuego

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestC(t *testing.T) {
	s := Stream[float32]{
		stream: func() chan float32 {
			c := make(chan float32)
			go func() { defer close(c); c <- 1.1; c <- 2.2; c <- 7.7; c <- 4.4 }()
			return c
		}(),
		concurrency: 2,
	}

	result := []int{}

	C(C(C(s.
		Map(float2int), Int).
		Map(int2string), String).
		Map(string2int), Int).
		ForEach(func(i int) { result = append(result, i) })

	expected := []int{1, 2, 7, 4}
	if !cmp.Equal(expected, result) {
		t.Error(cmp.Diff(expected, result))
	}
}

func TestSC(t *testing.T) {
	s := Stream[float32]{
		stream: func() chan float32 {
			c := make(chan float32)
			go func() { defer close(c); c <- 1.1; c <- 2.2; c <- 7.7; c <- 4.4 }()
			return c
		}(),
		concurrency: 2,
	}

	result := []int{}

	SC(SC(SC(s.
		Map(float2int), Stream[int]{}).
		Map(int2string), Stream[string]{}).
		Map(string2int), Stream[int]{}).
		ForEach(func(i int) { result = append(result, i) })

	expected := []int{1, 2, 7, 4}
	if !cmp.Equal(expected, result) {
		t.Error(cmp.Diff(expected, result))
	}
}

func TestStream_Map(t *testing.T) {
	tt := map[string]struct {
		stream Stream[int]
		mapper Function[int, R]
		want   []int
	}{
		"Should return an empty Stream": {
			stream: Stream[int]{stream: nil},
			mapper: functionTimesTwo,
			want:   nil,
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
			mapper: functionTimesTwo,
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
			var resultStream Stream[int] = C(tc.stream.Map(tc.mapper), Int)
			if gotStream := resultStream.stream; gotStream != nil {
				for val := range gotStream {
					got = append(got, val)
				}
			}

			if !cmp.Equal(tc.want, got) {
				t.Error(cmp.Diff(tc.want, got))
			}
		})
	}
}

func TestStream_Map_Concurrent(t *testing.T) {
	const numEntries = 300
	const concurrencyLevel = numEntries / 10

	sourceStream := func() chan int {
		c := make(chan int, 10)
		go func() {
			defer close(c)
			for i := 0; i < numEntries; i++ {
				c <- i
			}
		}()
		return c
	}()

	want := func() []int {
		ints := []int{}
		for i := 0; i < numEntries; i++ {
			ints = append(ints, functionTimesTwo(i).(int))
		}
		return ints
	}()

	unitStream := Stream[int]{
		stream: sourceStream,
	}

	result := []int{}

	start := time.Now()

	// functionSlowTimesTwo: use slow function to illustrate the performance improvement
	C(unitStream.
		Concurrent(concurrencyLevel).
		Map(functionSlowTimesTwo), Int).
		ForEach(func(i int) { result = append(result, i) })

	end := time.Now()

	if !cmp.Equal(want, result) {
		t.Error(cmp.Diff(want, result))
	}

	// if concurrency is not effective, the test will take 15 seconds for concurrency of 0 or
	// 7.5 seconds for a concurrency of 1.
	assert.WithinDuration(t, end, start, 3*time.Second) // 3 seconds should be plenty enough...
}

func TestStream_FlatMap(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5}
	c := []int{6, 7, 8}

	sliceOfSlicesOfInts := [][]int{a, b, c}

	result := []int{}

	C(NewStreamFromSlice(sliceOfSlicesOfInts, 0).
		FlatMap(FlattenSlice[int](0)), Int).
		ForEach(func(i int) { result = append(result, i) })

	// Before flattening: [[1 2 3] [4 5] [6 7 8]]
	// After flattening: [1 2 3 4 5 6 7 8]
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8}
	if !cmp.Equal(expected, result) {
		t.Error(cmp.Diff(expected, result))
	}
}

func TestStream_FlatMap_Concurrent(t *testing.T) {
	const numEntries = 300
	const concurrencyLevel = numEntries / 10

	sourceStream := func() chan []int {
		c := make(chan []int, 10)
		go func() {
			defer close(c)
			for i := 0; i < numEntries; i++ {
				p1 := []int{1, 2, 3}
				p2 := []int{4, 5}
				p3 := []int{6, 7, 8}

				c <- p1
				c <- p2
				c <- p3
			}
		}()
		return c
	}()

	want := func() []int {
		ints := []int{}
		for i := 0; i < numEntries; i++ {
			ints = append(ints, 1, 2, 3, 4, 5, 6, 7, 8)
		}
		return ints
	}()

	unitStream := Stream[[]int]{
		stream: sourceStream,
	}

	result := []int{}

	start := time.Now()

	/*slowFlattening :=*/
	func() StreamFunction[[]int, int] {
		// slow down the execution to illustrate the performance improvement of the concurrent stream
		time.Sleep(50 * time.Millisecond)
		return FlattenTypedSlice[int](0)
	}()

	slowFlattening := func() StreamFunction[[]int, R] {
		// slow down the execution to illustrate the performance improvement of the concurrent stream
		time.Sleep(50 * time.Millisecond)
		return FlattenSlice[int](0)
	}()

	unitStream.
		Concurrent(concurrencyLevel).
		FlatMap(slowFlattening).
		ForEach(func(i R) { result = append(result, i.(int)) })

	end := time.Now()

	if !cmp.Equal(want, result) {
		t.Error(cmp.Diff(want, result))
	}

	// if concurrency is not effective, the test will take 15 seconds for concurrency of 0 or
	// 7.5 seconds for a concurrency of 1.
	assert.WithinDuration(t, end, start, 3*time.Second) // 3 seconds should be plenty enough...
}

func TestStream_Filter(t *testing.T) {
	tests := map[string]struct {
		stream    chan int
		predicate Predicate[int]
		want      []int
	}{
		"Should return nil for an empty Stream": {
			stream:    nil,
			predicate: intGreaterThanPredicate(5),
			want:      nil,
		},
		"Should give produce filtered values as per predicate": {
			stream: func() chan int {
				c := make(chan int, 1)
				go func() {
					defer close(c)
					c <- 17
					c <- 8
					c <- 2
				}()
				return c
			}(),
			predicate: intGreaterThanPredicate(5),
			want:      []int{17, 8},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := Stream[int]{
				stream: tt.stream,
			}

			var got []int
			if gotStream := s.Filter(tt.predicate).stream; gotStream != nil {
				for val := range gotStream {
					got = append(got, val)
				}
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStream_LeftReduce(t *testing.T) {
	tt := map[string]struct {
		stream chan string
		want   string
	}{
		"Should return nil for a Stream of nil": {
			stream: nil,
			want:   "",
		},
		"Should return reduction of set of single element": {
			stream: func() chan string {
				c := make(chan string)
				go func() {
					defer close(c)
					c <- "three"
				}()
				return c
			}(),
			want: "three",
		},
		"Should return reduction of set of multiple elements": {
			stream: func() chan string {
				c := make(chan string)
				go func() {
					defer close(c)
					c <- "four-"
					c <- "twelve-"
					c <- "one-"
					c <- "six-"
					c <- "three"
				}()
				return c
			}(),
			want: "four-twelve-one-six-three",
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[string]{
				stream: tc.stream,
			}
			got := s.Reduce(Sum[string])
			assert.Exactly(t, tc.want, got)
		})
	}
}

func TestStream_ForEach(t *testing.T) {
	computeSumTotal := func(callCount, total *int) Consumer[int] {
		return func(value int) {
			*callCount++
			*total += value
		}
	}

	type want struct {
		total, count int
	}
	tt := map[string]struct {
		stream   chan int
		consumer func(callCount, total *int) Consumer[int]
		want     want
	}{
		"Should not call consumer for a Stream of nil": {
			stream:   nil,
			consumer: computeSumTotal,
			want: want{
				count: 0,
				total: 0,
			},
		},
		"Should give produce filtered values as per predicate": {
			stream: func() chan int {
				c := make(chan int)
				go func() {
					defer close(c)
					c <- 4
					c <- 1
					c <- 3
				}()
				return c
			}(),
			consumer: computeSumTotal,
			want: want{
				count: 3,
				total: 8,
			},
		},
	}

	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			callCount, total := 0, 0
			s := Stream[int]{
				stream: tc.stream,
			}

			s.ForEach(tc.consumer(&callCount, &total))
			assert.Equal(t, tc.want.count, callCount)
			assert.Equal(t, tc.want.total, total)
		})
	}
}

var float2int = func() Function[float32, R] {
	return func(f float32) R {
		return int(f)
	}
}()

var int2string = func() Function[int, R] {
	return func(i int) R {
		return fmt.Sprintf("%d", i)
	}
}()

var string2int = func() Function[string, R] {
	return func(s string) R {
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		return i
	}
}()

var functionTimesTwo = func() Function[int, R] {
	return func(i int) R {
		return 2 * i
	}
}()

var functionSlowTimesTwo = func() Function[int, R] {
	return func(i int) R {
		time.Sleep(50 * time.Millisecond)
		return 2 * i
	}
}()

var intFunctionSlowTimesTwo = func() Function[int, int] {
	return func(i int) int {
		time.Sleep(50 * time.Millisecond)
		return 2 * i
	}
}()
