package fuego

import (
	"fmt"
	"hash/crc32"
	"reflect"
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
	tt := map[string]struct {
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

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[int]{
				stream: tc.stream,
			}

			var got []int
			if gotStream := s.Filter(tc.predicate).stream; gotStream != nil {
				for val := range gotStream {
					got = append(got, val)
				}
			}

			assert.Equal(t, tc.want, got)
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
	tt := map[string]struct {
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

	for name, tc := range tt {
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

func TestStream_Drop(t *testing.T) {
	data1 := []any{
		1,
	}

	data := []any{
		"a",
		false,
		"b",
		-17,
		"c",
	}

	dataGenerator := func(slice []any) chan any {
		c := make(chan any, 2)
		go func() {
			defer close(c)
			for _, val := range slice {
				c <- val
			}
		}()
		return c
	}

	tt := map[string]struct {
		stream chan any
		n      uint64
		want   []any
	}{
		"Should return an empty stream when nil channel": {
			stream: nil,
			n:      1,
			want:   []any{},
		},
		"Should not change the stream if n < 1": {
			stream: dataGenerator(data),
			n:      0,
			want:   data,
		},
		"Should drop all elements when n > number of elements in the stream": {
			stream: dataGenerator(data),
			n:      uint64(len(data) + 10),
			want:   []any{},
		},
		"Should drop the first n elements when n < number of elements in the stream": {
			stream: dataGenerator(data),
			n:      2,
			want:   data[2:],
		},
		"Should drop the sole element in the stream": {
			stream: dataGenerator(data1),
			n:      1,
			want:   []any{},
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[any]{stream: tc.stream}
			gotStream := s.Drop(tc.n)
			if tc.want == nil {
				assert.Nil(t, gotStream.stream)
				return
			}
			got := []any{}
			for val := range gotStream.stream {
				got = append(got, val)
			}
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestStream_DropWhile(t *testing.T) {
	data := []any{
		"a",
		false,
		"b",
		-17,
		"c",
	}

	dataGenerator := func() chan any {
		c := make(chan any, 2)
		go func() {
			defer close(c)
			for _, val := range data {
				c <- val
			}
		}()
		return c
	}

	tt := map[string]struct {
		stream chan any
		p      Predicate[any]
		want   []any
	}{
		"Should return empty out-stream when nil in-stream": {
			stream: nil,
			p:      True[any](),
			want:   []any{},
		},
		"Should not change the stream if predicate never satisfies": {
			stream: dataGenerator(),
			p:      False[any](),
			want:   data,
		},
		"Should drop the first few elements that satisfy the predicate": {
			stream: dataGenerator(),
			p: func(e any) bool {
				return e == "a" || e == false
			},
			want: data[2:],
		},
		"Should drop all elements when the predicate always satisfies": {
			stream: dataGenerator(),
			p:      True[any](),
			want:   []any{},
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[any]{stream: tc.stream}
			gotStream := s.DropWhile(tc.p)
			got := []any{}
			for val := range gotStream.stream {
				got = append(got, val)
			}
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestStream_DropUntil(t *testing.T) {
	data := []any{
		"a",
		false,
		"b",
		-17,
		"c",
	}

	dataGenerator := func() chan any {
		c := make(chan any, 2)
		go func() {
			defer close(c)
			for _, val := range data {
				c <- val
			}
		}()
		return c
	}

	tt := map[string]struct {
		stream chan any
		p      Predicate[any]
		want   []any
	}{
		"Should return empty out-stream when nil in-stream": {
			stream: nil,
			p:      True[any](),
			want:   []any{},
		},
		"Should return empty stream if predicate never satisfies": {
			stream: dataGenerator(),
			p:      False[any](),
			want:   []any{},
		},
		"Should drop the first few elements that satisfy the predicate": {
			stream: dataGenerator(),
			p: func(e any) bool {
				return e == "b"
			},
			want: []any{
				"b",
				-17,
				"c",
			},
		},
		"Should not drop any element when the predicate always satisfies": {
			stream: dataGenerator(),
			p:      True[any](),
			want:   data,
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[any]{stream: tc.stream}
			gotStream := s.DropUntil(tc.p)
			got := []any{}
			for val := range gotStream.stream {
				got = append(got, val)
			}
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestStream_LastX_PanicsWhenNilChannel(t *testing.T) {
	assert.PanicsWithValue(t, PanicMissingChannel, func() { Stream[int]{stream: nil}.LastN(1) })
	assert.PanicsWithValue(t, PanicMissingChannel, func() { Stream[int]{stream: nil}.Last() })
}

func TestStream_LastX_PanicsWhenEmptyChannel(t *testing.T) {
	emptyStream := func() chan int {
		c := make(chan int)
		go func() {
			defer close(c)
		}()
		return c
	}

	assert.PanicsWithValue(t, PanicNoSuchElement, func() { NewStream(emptyStream()).LastN(1) })
	assert.PanicsWithValue(t, PanicNoSuchElement, func() { NewStream(emptyStream()).Last() })
}

func TestStream_LastNWithInvalidArgumentPanics(t *testing.T) {
	tt := map[string]struct {
		n         uint64
		wantPanic bool
	}{
		"Should panic when N is less than 1": {
			n:         0,
			wantPanic: true,
		},
		"Should not panic when N is 1": {
			n:         1,
			wantPanic: false,
		},
	}

	populatedStream := func() chan string {
		c := make(chan string)
		go func() {
			defer close(c)
			c <- "joy"
		}()
		return c
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[string]{
				stream: populatedStream(),
			}
			if tc.wantPanic {
				assert.PanicsWithValue(t, PanicNoSuchElement, func() { s.LastN(tc.n) })
			} else {
				assert.NotPanics(t, func() { s.LastN(tc.n) })
			}
		})
	}
}

func TestStream_LastN(t *testing.T) {
	data := []any{
		"one",
		"two",
		"thee",
		"four",
		"five",
	}

	largeData := []any{}
	for i := 1; i < 1e5; i++ {
		largeData = append(largeData, i)
	}

	populatedStream := func(slice []any) chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			for _, val := range slice {
				c <- val
			}
		}()
		return c
	}

	tt := map[string]struct {
		stream chan any
		n      uint64
		want   []any
	}{
		"Should return the last n elements when the stream is longer": {
			stream: populatedStream(data),
			n:      2,
			want:   data[3:],
		},
		"Should return all the elements when the stream is shorter": {
			stream: populatedStream(data),
			n:      uint64(len(data) + 10),
			want:   data,
		},
		"Should return the last n elements when the stream is significantly larger": {
			stream: populatedStream(largeData),
			n:      2e3,
			want:   largeData[len(largeData)-2e3:],
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[any]{
				stream: tc.stream,
			}
			got := s.LastN(tc.n)
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestStream_HeadX_PanicsWhenNilChannel(t *testing.T) {
	assert.PanicsWithValue(t, PanicMissingChannel, func() { Stream[any]{stream: nil}.HeadN(1) })
	assert.PanicsWithValue(t, PanicMissingChannel, func() { Stream[any]{stream: nil}.Head() })
}

func TestStream_Head_PanicsWhenEmptyChannel(t *testing.T) {
	emptyStream := func() chan any {
		c := make(chan any)
		go func() {
			defer close(c)
		}()
		return c
	}

	assert.PanicsWithValue(t, PanicNoSuchElement, func() { NewStream(emptyStream()).Head() })
}

func TestStream_Head(t *testing.T) {
	data1 := []any{
		"one",
	}

	data5 := []any{
		"one",
		"two",
		"thee",
		"four",
		"five",
	}

	generateStream := func(slice []any) chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			for _, val := range slice {
				c <- val
			}
		}()
		return c
	}

	tt := map[string]struct {
		stream chan any
		want   any
	}{
		"Should return first element when the stream has one element": {
			stream: generateStream(data1),
			want:   "one",
		},
		"Should return first element when the stream has multiple elements": {
			stream: generateStream(data5),
			want:   "one",
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[any]{
				stream: tc.stream,
			}
			if got := s.Head(); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Stream.Head() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStream_HeadN(t *testing.T) {
	data0 := []any{}

	data := []any{
		"one",
		"two",
		"thee",
		"four",
		"five",
	}

	largeData := []any{}
	for i := 1; i < 1e5; i++ {
		largeData = append(largeData, i)
	}

	generateStream := func(slice []any) chan any {
		c := make(chan any)
		go func() {
			defer close(c)
			for _, val := range slice {
				c <- val
			}
		}()
		return c
	}

	tt := map[string]struct {
		stream chan any
		n      uint64
		want   []any
	}{
		"Should return empty slice when the stream is empty": {
			stream: generateStream(data0),
			n:      2,
			want:   []any{},
		},
		"Should return the first n elements when the stream is longer": {
			stream: generateStream(data),
			n:      2,
			want:   data[:2],
		},
		"Should return all the elements when the stream is shorter": {
			stream: generateStream(data),
			n:      2e3,
			want:   data,
		},
		"Should return the first n elements when the stream is significantly larger": {
			stream: generateStream(largeData),
			n:      2e3,
			want:   largeData[:2000],
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[any]{
				stream: tc.stream,
			}
			got := s.HeadN(tc.n)
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestStream_TakeX_LimitPanicWithNilChannel(t *testing.T) {
	s := Stream[any]{stream: nil}
	assert.PanicsWithValue(t, PanicMissingChannel, func() { s.Limit(1) })
	assert.PanicsWithValue(t, PanicMissingChannel, func() { s.TakeUntil(False[any]()) })
	assert.PanicsWithValue(t, PanicMissingChannel, func() { s.TakeWhile(True[any]()) })
}

func TestStream_Take_Limit(t *testing.T) {
	data1 := []any{
		1,
	}

	data := []any{
		"a",
		false,
		"b",
		-17,
		"c",
	}

	dataGenerator := func(slice []any) chan any {
		c := make(chan any, 2)
		go func() {
			defer close(c)
			for _, val := range slice {
				c <- val
			}
		}()
		return c
	}

	tt := map[string]struct {
		stream chan any
		n      uint64
		want   []any
	}{
		"Should return empty stream when n < 1": {
			stream: dataGenerator(data),
			n:      0,
			want:   []any{},
		},
		"Should return all elements when n > number of elements in the stream": {
			stream: dataGenerator(data),
			n:      uint64(len(data) + 10),
			want:   data,
		},
		"Should return the first n elements when n < number of elements in the stream": {
			stream: dataGenerator(data),
			n:      2,
			want:   data[:2],
		},
		"Should return the sole element in the stream": {
			stream: dataGenerator(data1),
			n:      1,
			want:   data1,
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[any]{stream: tc.stream}
			gotStream := s.Limit(tc.n)
			if tc.want == nil {
				assert.Nil(t, gotStream.stream)
				return
			}
			got := []any{}
			for val := range gotStream.stream {
				got = append(got, val)
			}
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestStream_TakeWhile(t *testing.T) {
	data := []any{
		"a",
		false,
		"b",
		-17,
		"c",
	}

	dataGenerator := func() chan any {
		c := make(chan any, 2)
		go func() {
			defer close(c)
			for _, val := range data {
				c <- val
			}
		}()
		return c
	}

	tt := map[string]struct {
		stream chan any
		p      Predicate[any]
		want   []any
	}{
		"Should return empty stream if predicate never satisfies": {
			stream: dataGenerator(),
			p:      False[any](),
			want:   []any{},
		},
		"Should take the first few elements that satisfy the predicate": {
			stream: dataGenerator(),
			p: func(e any) bool {
				return e == "a" || e == false
			},
			want: data[:2],
		},
		"Should take all elements when the predicate always satisfies": {
			stream: dataGenerator(),
			p:      True[any](),
			want:   data,
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[any]{stream: tc.stream}
			gotStream := s.TakeWhile(tc.p)
			got := []any{}
			for val := range gotStream.stream {
				got = append(got, val)
			}
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestStream_TakeUntil(t *testing.T) {
	data := []any{
		"a",
		false,
		"b",
		-17,
		"c",
	}

	dataGenerator := func() chan any {
		c := make(chan any, 2)
		go func() {
			defer close(c)
			for _, val := range data {
				c <- val
			}
		}()
		return c
	}

	tt := map[string]struct {
		stream chan any
		p      Predicate[any]
		want   []any
	}{
		"Should return whole stream if predicate never satisfies": {
			stream: dataGenerator(),
			p:      False[any](),
			want:   data,
		},
		"Should take the first few elements until predicate satisfies": {
			stream: dataGenerator(),
			p: func(e any) bool {
				return e == "b"
			},
			want: []any{
				"a",
				false,
			},
		},
		"Should return empty stream when the predicate always satisfies": {
			stream: dataGenerator(),
			p:      True[any](),
			want:   []any{},
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[any]{stream: tc.stream}
			gotStream := s.TakeUntil(tc.p)
			got := []any{}
			for val := range gotStream.stream {
				got = append(got, val)
			}
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestStream_StartsWithPanicsWithNilChannel(t *testing.T) {
	s := Stream[any]{stream: nil}
	assert.PanicsWithValue(t, PanicMissingChannel, func() { s.StartsWith([]any{}) })
}

func TestStream_StartsWith(t *testing.T) {
	data0 := []any{}
	data1 := []any{16}
	data4 := []any{
		true,
		1,
		4,
		"two",
	}
	data5 := []any{
		1,
		"two",
		true,
		4,
		"five",
	}

	generateStream := func(data []any) chan any {
		c := make(chan any)
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
		slice  []any
		want   bool
	}{
		"Should not match with an empty stream": {
			stream: generateStream(data0),
			slice:  data1,
			want:   false,
		},
		"Should not match with an empty slice: case with an empty stream ": {
			stream: generateStream(data0),
			slice:  data0,
			want:   false,
		},
		"Should not match with an empty slice: case with a populated stream ": {
			stream: generateStream(data1),
			slice:  data0,
			want:   false,
		},
		"Should not match when stream size is less than slice to compare even when the elements match": {
			stream: generateStream(data5[:3]),
			slice:  data5,
			want:   false,
		},
		"Should match when stream size and data matches slice to compare": {
			stream: generateStream(data5),
			slice:  data5,
			want:   true,
		},
		"Should match when stream starts with the slice to compare": {
			stream: generateStream(data5),
			slice:  data5[:3],
			want:   true,
		},
		"Should not match when stream does not end with slice to compare": {
			stream: generateStream(data5),
			slice:  data4,
			want:   false,
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[any]{
				stream: tc.stream,
			}
			assert.Equal(t, tc.want, s.StartsWith(tc.slice))
		})
	}
}

func TestStream_EndsWith(t *testing.T) {
	data0 := []any{}
	data1 := []any{16}
	data4 := []any{
		true,
		1,
		4,
		"two",
	}
	data5 := []any{
		1,
		"two",
		true,
		4,
		"five",
	}

	generateStream := func(data []any) chan any {
		c := make(chan any)
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
		slice  []any
		want   bool
	}{
		"Should not match with a nil channel": {
			stream: nil,
			slice:  data1,
			want:   false,
		},
		"Should not match with an empty stream and with input": {
			stream: generateStream(data0),
			slice:  data1,
			want:   false,
		},
		"Should not match with an empty stream and with no input": {
			stream: generateStream(data0),
			slice:  data0,
			want:   false,
		},
		"Should not match when stream size is less than slice to compare even when the elements match": {
			stream: generateStream(data5[:3]),
			slice:  data5,
			want:   false,
		},
		"Should match when stream size and data match slice to compare": {
			stream: generateStream(data5),
			slice:  data5,
			want:   true,
		},
		"Should match when the stream ends with the slice to compare": {
			stream: generateStream(data5),
			slice:  data5[3:],
			want:   true,
		},
		"Should not match when stream does not end with slice to compare": {
			stream: generateStream(data5),
			slice:  data4,
			want:   false,
		},
		"Should not match when stream is shorter than slice to compare": {
			stream: generateStream(data5[:2]),
			slice:  data5,
			want:   false,
		},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			s := Stream[any]{
				stream: tc.stream,
			}
			assert.Equal(t, tc.want, s.EndsWith(tc.slice))
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

	tt := map[string]struct {
		stream        chan int
		consumer      func(callCount, total *int) Consumer[int]
		want          []int
		wantTotal     int
		wantCallCount int
	}{
		"Should peek and return empty stream when nil in-stream": {
			stream:        nil,
			consumer:      computeSumTotal,
			want:          []int{},
			wantTotal:     0,
			wantCallCount: 0,
		},
		"Should peek and return empty stream when empty in-stream": {
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
		"Should peek and return stream when populated in-stream": {
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

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
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
