package fuego

import (
	"fmt"
	"hash/crc32"
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
		mapper Function[int, Any]
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

	slowFlattening := func() StreamFunction[[]int, Any] {
		// slow down the execution to illustrate the performance improvement of the concurrent stream
		time.Sleep(50 * time.Millisecond)
		return FlattenSlice[int](0)
	}()

	unitStream.
		Concurrent(concurrencyLevel).
		FlatMap(slowFlattening).
		ForEach(func(i Any) { result = append(result, i.(int)) })

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
			got := s.Reduce(Concatenate[string])
			assert.Exactly(t, tc.want, got)
		})
	}
}

func TestStream_Intersperse(t *testing.T) {
	tests := map[string]struct {
		stream    chan string
		inBetween string
		want      []string
	}{
		"Should return an empty Stream for nil input Stream": {
			stream:    nil,
			inBetween: " - ",
			want:      []string{},
		},
		"Should return an empty Stream for empty input Stream": {
			stream: func() chan string {
				c := make(chan string)
				go func() {
					defer close(c)
				}()
				return c
			}(),
			inBetween: " - ",
			want:      []string{},
		},
		"Should return the original input Stream when it has a single value": {
			stream: func() chan string {
				c := make(chan string)
				go func() {
					defer close(c)
					c <- "four"
				}()
				return c
			}(),
			inBetween: " - ",
			want: []string{
				"four",
			},
		},
		"Should return the Set with given value interspersed": {
			stream: func() chan string {
				c := make(chan string)
				go func() {
					defer close(c)
					c <- "four"
					c <- "twelve"
					c <- "one"
					c <- "six"
					c <- "three"
				}()
				return c
			}(),
			inBetween: " - ",
			want: []string{
				"four",
				" - ",
				"twelve",
				" - ",
				"one",
				" - ",
				"six",
				" - ",
				"three"},
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[string]{
				stream: tc.stream,
			}
			out := s.Intersperse(tc.inBetween)
			got := []string{}
			for e := range out.stream {
				got = append(got, e)
			}
			if !assert.ElementsMatch(t, got, tc.want) {
				t.Errorf("Stream.Intersperse() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStream_GroupBy(t *testing.T) {
	tt := map[string]struct {
		stream     chan int
		classifier Function[int, Any]
		want       map[Any][]int
	}{
		"Should return empty map when iterator with nil stream": {
			stream: nil,
			classifier: func(i int) Any {
				return i & 1
			},
			want: map[Any][]int{},
		},
		"Should return empty map when empty stream": {
			stream: func() chan int {
				c := make(chan int)
				go func() {
					defer close(c)
				}()
				return c
			}(),
			classifier: func(i int) Any {
				return i & 1
			},
			want: map[Any][]int{},
		},
		"Should group by odd / even numbers": {
			stream: func() chan int {
				c := make(chan int)
				go func() {
					defer close(c)
					c <- 1
					c <- 2
					c <- 3
					c <- 4
				}()
				return c
			}(),
			classifier: func(i int) Any {
				return i & 1
			},
			want: map[Any][]int{
				0: {2, 4},
				1: {1, 3},
			},
		},
	}
	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			rp := Stream[int]{
				stream: tc.stream,
			}

			got := rp.GroupBy(tc.classifier)
			if !cmp.Equal(tc.want, got) {
				t.Error(cmp.Diff(tc.want, got))
			}
		})
	}
}

func TestStream_Count(t *testing.T) {
	tt := map[string]struct {
		stream chan int
		want   int
	}{
		"Should return 0 for a nil channel": {
			stream: nil,
			want:   0,
		},
		"Should return 0 for an empty closed channel": {
			stream: func() chan int {
				c := make(chan int)
				go func() {
					defer close(c)
				}()
				return c
			}(),
			want: 0,
		},
		"Should return 3 for a size 3 closed channel": {
			stream: func() chan int {
				c := make(chan int, 1)
				go func() {
					defer close(c)
					c <- 1
					c <- 2
					c <- 1
				}()
				return c
			}(),
			want: 3,
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[int]{
				stream: tc.stream,
			}
			if got := s.Count(); got != tc.want {
				t.Errorf("Stream.Count() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStream_AnyMatch(t *testing.T) {
	dataGenerator := func() chan any {
		c := make(chan any, 2)
		go func() {
			defer close(c)
			c <- "a"
			c <- false
			c <- "b"
			c <- -17
			c <- "c"
		}()
		return c
	}

	tt := map[string]struct {
		stream    chan any
		predicate Predicate[any]
		want      bool
	}{
		"Should not match any when channel is nil": {
			stream:    nil,
			predicate: True[any](),
			want:      false,
		},
		"Should not match any": {
			stream:    dataGenerator(),
			predicate: func(e any) bool { return e == "not in here" },
			want:      false,
		},
		"Should match any": {
			stream:    dataGenerator(),
			predicate: func(e any) bool { return e == "b" },
			want:      true,
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[any]{
				stream: tc.stream,
			}
			if got := s.AnyMatch(tc.predicate); got != tc.want {
				t.Errorf("Stream.AnyMatch() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStream_NoneMatch(t *testing.T) {
	dataGenerator := func() chan any {
		c := make(chan any, 2)
		go func() {
			defer close(c)
			c <- "a"
			c <- false
			c <- "b"
			c <- -17
			c <- "c"
		}()
		return c
	}

	tt := map[string]struct {
		stream    chan any
		predicate Predicate[any]
		want      bool
	}{
		"Should satisfy when channel is nil": {
			stream:    nil,
			predicate: True[any](),
			want:      true,
		},
		"Should satisfy": {
			stream:    dataGenerator(),
			predicate: func(e any) bool { return e == "not in here" },
			want:      true,
		},
		"Should not satisfy": {
			stream:    dataGenerator(),
			predicate: func(e any) bool { return e == "b" },
			want:      false,
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[any]{
				stream: tc.stream,
			}
			if got := s.NoneMatch(tc.predicate); got != tc.want {
				t.Errorf("Stream.NoneMatch() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStream_AllMatch(t *testing.T) {
	dataGenerator := func() chan any {
		c := make(chan any, 2)
		go func() {
			defer close(c)
			c <- "a"
			c <- false
			c <- "b"
			c <- -17
			c <- "c"
		}()
		return c
	}

	tt := map[string]struct {
		stream    chan any
		predicate Predicate[any]
		want      bool
	}{
		"Should not match all when channel is nil": {
			stream:    nil,
			predicate: True[any](),
			want:      false,
		},
		"Should match all": {
			stream:    dataGenerator(),
			predicate: func(e any) bool { return e != "not in here" },
			want:      true,
		},
		"Should not match all": {
			stream:    dataGenerator(),
			predicate: func(e any) bool { return e == "b" },
			want:      false,
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[any]{
				stream: tc.stream,
			}
			if got := s.AllMatch(tc.predicate); got != tc.want {
				t.Errorf("Stream.AllMatch() = %v, want %v", got, tc.want)
			}
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

func TestStream_ToSlice(t *testing.T) {
	tt := map[string]struct {
		stream chan int
		want   []int
	}{
		"Should produce an empty stream when in-stream is nil": {
			stream: nil,
			want:   []int{},
		},
		"Should produce an empty stream when in-stream is empty": {
			stream: func() chan int {
				c := make(chan int)
				go func() { defer close(c) }()
				return c
			}(),
			want: []int{},
		},
		"Should produce a flat stream when in-stream is not empty": {
			stream: func() chan int {
				c := make(chan int)
				go func() {
					defer close(c)
					c <- 1
					c <- 2
					c <- 3
				}()
				return c
			}(),
			want: []int{
				1,
				2,
				3,
			},
		},
	}

	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			s := Stream[int]{
				stream: tc.stream,
			}
			got := s.ToSlice()
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestStream_DistinctPanicsWhenNilChannel(t *testing.T) {
	assert.PanicsWithValue(t, PanicMissingChannel, func() { Stream[string]{stream: nil}.Distinct(func(string) uint32 { return 0 }) })
}

func TestStream_Distinct(t *testing.T) {
	data5 := []any{
		1,
		"two",
		true,
		4,
		"five",
	}

	data10 := []any{
		1,
		"five",
		"two",
		1,
		"five",
		"five",
		true,
		4,
		false,
		"five",
	}

	generateStream := func(data []any) chan any {
		c := make(chan any, 200)
		go func() {
			defer close(c)
			for _, e := range data {
				c <- e
			}
		}()
		return c
	}

	tt := map[string]struct {
		stream chan any
		want   []any
	}{
		"Should return same values when the in-stream has distinct values": {
			stream: generateStream(data5),
			want:   data5,
		},
		"Should return distinct values when the in-stream has repeat values": {
			stream: generateStream(data10),
			want: []any{
				1,
				"five",
				"two",
				true,
				4,
				false,
			},
		},
	}

	hashFn := func(t any) uint32 {
		switch tp := t.(type) {
		case bool:
			if tp {
				return 1
			}
			return 0
		case int:
			return crc32.ChecksumIEEE([]byte(strconv.Itoa(tp)))
		case string:
			return crc32.ChecksumIEEE([]byte(tp))
		default:
			panic("unknown type")
		}
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[any]{
				stream: tc.stream,
			}
			got := s.Distinct(hashFn).ToSlice()
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestStream_Peek(t *testing.T) {
	computeSumTotal := func(callCount, total *int) Consumer[int] {
		return func(value int) {
			*callCount++
			*total += value
		}
	}

	tt := []struct {
		name          string
		stream        chan int
		consumer      func(callCount, total *int) Consumer[int]
		want          []int
		wantTotal     int
		wantCallCount int
	}{
		{
			name:          "Should peek and return empty stream when nil in-stream",
			stream:        nil,
			consumer:      computeSumTotal,
			want:          []int{},
			wantTotal:     0,
			wantCallCount: 0,
		},
		{
			name: "Should peek and return empty stream when empty in-stream",
			stream: func() chan int {
				c := make(chan int)
				go func() {
					defer close(c)
				}()
				return c
			}(),
			consumer:      computeSumTotal,
			want:          []int{},
			wantTotal:     0,
			wantCallCount: 0,
		},
		{
			name: "Should peek and return stream when populated in-stream",
			stream: func() chan int {
				c := make(chan int)
				go func() {
					defer close(c)
					c <- 1
					c <- 2
					c <- 3
					c <- 5
					c <- 8
				}()
				return c
			}(),
			consumer: computeSumTotal,
			want: []int{
				1,
				2,
				3,
				5,
				8,
			},
			wantTotal:     19,
			wantCallCount: 5,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			callCount, total := 0, 0

			s := Stream[int]{
				stream: tc.stream,
			}

			got := s.Peek(tc.consumer(&callCount, &total))
			assert.EqualValues(t, tc.want, got.ToSlice())
			assert.Equal(t, tc.wantTotal, total)
			assert.Equal(t, tc.wantCallCount, callCount)
		})
	}
}

var float2int = func() Function[float32, Any] {
	return func(f float32) Any {
		return int(f)
	}
}()

var int2string = func() Function[int, Any] {
	return func(i int) Any {
		return fmt.Sprintf("%d", i)
	}
}()

var string2int = func() Function[string, Any] {
	return func(s string) Any {
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		return i
	}
}()

var functionTimesTwo = func() Function[int, Any] {
	return func(i int) Any {
		return 2 * i
	}
}()

var functionSlowTimesTwo = func() Function[int, Any] {
	return func(i int) Any {
		time.Sleep(50 * time.Millisecond)
		return 2 * i
	}
}()
