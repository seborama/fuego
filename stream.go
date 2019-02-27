package fuego

import "fmt"

// TODO: consider two types of streams: CStreams (channel based as shown here) and SStreams (slice based). The former allows for infinite streams and thinner memory usage within the CStream object but lacks performance when the operation requires to deal with the end of the steam (it has to consume all the elements of the steam sequentially). SStreams require the entire data to be stored internally from the onset. However,  slices are seekable and can read from the end or be consumed backwards easily.

// TODO: a stream should probably be marked as invalid after most (or all?) operations on it because the channel will have likely changed state.

// TODO list:
// Methods with ** require the Stream to be finite and closed (or use a Future, perhaps Future.Stream()?)
// **DropRight(uint64) - drops the last n elements of the Stream. Only meaningful if the Stream is closed.
// FilterNot(Predicate) - <=> to Filter(Not(Predicate))
// MapToString(ToStringFunction)
// FlatMapToXXX (Int, Uint, etc) => is this the same as FlatMap().MapToXXX()?
// **Sorted(Comparator)
// Contains
// ContainsAll
// Tail
// Zip / Unzip
// Concat?
// Range(from, toExclusive)?
// RangeBy(from,toExclusive, step)?
// Unfold()
// Fold()?
// Fold / FoldLeft
// TODO: implement NewStreamFromMap -> Stream of Keys / Stream of Values?

// PanicMissingChannel signifies that the Stream is missing a channel.
const PanicMissingChannel = "stream creation requires a channel"

// Stream is a sequence of elements supporting sequential and
// (in the future?) parallel operations.
//
// A Stream is a wrapper over a Go channel ('nil' channels are prohibited).
//
// NOTE
//
// Concurrent streams are challenging to implement owing to
// ordering issues in parallel processing. At the moment, the view
// is that the most sensible approach is to delegate control to users.
// Multiple fuego streams can be created and data distributed
// across as desired. This empowers users of fuego to implement the
// desired behaviour of their pipelines.
//
// Creation
//
// When providing a Go channel to create a Stream, beware that until you
// close the channel, the Stream's internal Go function that processes
// the data on the channel will remain active. It will block until
// either new data is produced or the channel is closed by the producer.
// When a producer forgets to close the channel, the Go function will stray.
//
// Streams created from a slice do not suffer from this issue because
// they are closed when the slice content is fully pushed to the Stream.
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryInt(1),
//      ƒ.EntryInt(2),
//      ƒ.EntryInt(3),
//  }, 1e3)
//  // or if you already have a channel of Entry:
//  c := make(chan ƒ.Entry) // you could add a buffer size as a second arg, if desired
//  go func() {
//      defer close(c)
//      c <- ƒ.EntryString("one")
//      c <- ƒ.EntryString("two")
//      c <- ƒ.EntryString("three")
//      // c <- ...
//  }()
//  NewStream(c)
type Stream struct {
	stream           chan Entry
	concurrencyLevel uint
}

// NewStream creates a new Stream.
//
// This function leaves the provided channel is the same state
// of openness.
func NewStream(c chan Entry) Stream {
	if c == nil {
		panic(PanicMissingChannel)
	}
	return Stream{
		stream: c,
	}
}

// NewStreamFromSlice creates a new Stream from a Go slice.
//
// The slice data is published to the stream after which the
// stream is closed.
func NewStreamFromSlice(slice EntrySlice, bufsize int) Stream {
	c := make(chan Entry, bufsize)

	go func() {
		defer close(c) // slices have finite size: close stream after all data was read.
		for _, element := range slice {
			c <- element
		}
	}()

	return NewStream(c)
}

// Map returns a Stream consisting of the result of
// applying the given function to the elements of this stream.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) Map(mapper Function) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() {
		defer close(outstream)
		if s.stream == nil {
			return
		}
		for val := range s.stream {
			outstream <- mapper(val)
		}
	}()

	return Stream{
		stream: outstream,
	}
}

// FlatMap takes a StreamFunction to flatten the entries
// in this stream and produce a new stream.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
//
// Example
//
//  a := EntrySlice{EntryInt(1), EntryInt(2), EntryInt(3)}
//  b := EntrySlice{EntryInt(4), EntryInt(5)}
//  c := EntrySlice{EntryInt(6), EntryInt(7), EntryInt(8)}
//
//  sliceOfEntrySlicesOfEntryInts := EntrySlice{a, b, c}
//
//  fmt.Printf("Before flattening: %+v\n", sliceOfEntrySlicesOfEntryInts)
//
//  sliceOfEntryInts := NewStreamFromSlice(sliceOfEntrySlicesOfEntryInts, 0).
//      FlatMap(FlattenEntrySliceToEntry(0)).
//      ToSlice()
//
//  fmt.Printf("After flattening: %+v\n", sliceOfEntryInts)
//
//  // Output:
//  // Before flattening: [[1 2 3] [4 5] [6 7 8]]
//  // After flattening: [1 2 3 4 5 6 7 8]
func (s Stream) FlatMap(mapper StreamFunction) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() {
		defer close(outstream)
		if s.stream == nil {
			return
		}
		for val := range s.stream {
			mapper(val).ForEach(func(e Entry) {
				outstream <- e
			})
		}
	}()

	return Stream{
		stream: outstream,
	}
}

// Filter returns a stream consisting of the elements of this stream that
// match the given predicate.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
//
// Example
//
//  s := ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryInt(1),
//      ƒ.EntryInt(2),
//      ƒ.EntryInt(3),
//  }, 0)
//
//  s.Filter(
//          FunctionPredicate(entryIntEqualsTo(ƒ.EntryInt(1))).
//              Or(
//                  FunctionPredicate(entryIntEqualsTo(ƒ.EntryInt(3)))),
//  )
//  // Result: []ƒ.EntryInt{1,3}
func (s Stream) Filter(predicate Predicate) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() {
		defer close(outstream)
		if s.stream == nil {
			return
		}
		for val := range s.stream {
			if predicate(val) {
				outstream <- val
			}
		}
	}()

	return Stream{
		stream: outstream,
	}
}

// ForEach executes the given function for each entry in this stream.
//
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
//
// Example
//
//  total := 0
//
//  computeSumTotal := func(value ƒ.Entry) {
//      total += int(value.(ƒ.EntryInt).Value())
//  }
//
//  s := ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryInt(1),
//      ƒ.EntryInt(2),
//      ƒ.EntryInt(3),
//  }, 0).
//      ForEach(calculateSumTotal)
//  // Result: total == 6
func (s Stream) ForEach(consumer Consumer) {
	if s.stream == nil {
		return
	}

	for val := range s.stream {
		consumer(val)
	}
}

// Peek is akin to ForEach but returns the Stream.
//
// This is useful e.g. for debugging.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) Peek(consumer Consumer) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() {
		defer close(outstream)
		s.ForEach(func(e Entry) {
			consumer(e)
			outstream <- e
		})
	}()

	return Stream{
		stream: outstream,
	}
}

// LeftReduce accumulates the elements of this Stream by
// applying the given function.
//
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("four"),
//      ƒ.EntryString("twelve)",
//      ƒ.EntryString("one"),
//      ƒ.EntryString("six"),
//      ƒ.EntryString("three"),
//  }, 1e3).
//      Reduce(concatenateStringsBiFunc)
//  // Result: ƒ.EntryString("four-twelve-one-six-three")
func (s Stream) LeftReduce(f2 BiFunction) Entry {
	if s.stream == nil {
		return nil // TODO: return Maybe
	}

	res := <-s.stream
	for val := range s.stream {
		res = f2(res, val)
	}

	return res
}

// Reduce is an alias for LeftReduce.
//
// See LeftReduce for more info.
func (s Stream) Reduce(f2 BiFunction) Entry {
	return s.LeftReduce(f2)
}

// Intersperse inserts an element between all elements of this Stream.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("four"),
//  }, 1e3).
//      Intersperse(ƒ.EntryString(" - "))
//  // Result: "three - two - four"
func (s Stream) Intersperse(e Entry) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() {
		defer close(outstream)
		if s.stream == nil {
			return
		}
		if val := <-s.stream; val != nil {
			outstream <- val
		}
		for val := range s.stream {
			outstream <- e
			outstream <- val
		}
	}()

	return Stream{
		stream: outstream,
	}
}

// GroupBy groups the elements of this Stream by classifying them.
//
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
func (s Stream) GroupBy(classifier Function) EntryMap {
	resultMap := EntryMap{}

	if s.stream != nil {
		for val := range s.stream {
			k := classifier(val)
			if resultMap[k] == nil {
				resultMap[k] = EntrySlice{}
			}
			resultMap[k] = append(resultMap[k].(EntrySlice), val)
			// TODO?: resultMap = resultMap.Append(Tuple2{k, val})
		}
	}

	return resultMap
}

// MapToInt produces an EntryInt stream.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) MapToInt(toInt ToIntFunction) IntStream {
	outstream := make(chan EntryInt, cap(s.stream))

	go func() {
		defer close(outstream)
		if s.stream == nil {
			return
		}
		for val := range s.stream {
			outstream <- toInt(val)
		}
	}()

	return NewIntStream(outstream)
}

// MapToFloat produces an EntryFloat stream.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) MapToFloat(toFloat ToFloatFunction) FloatStream {
	outstream := make(chan EntryFloat, cap(s.stream))

	go func() {
		defer close(outstream)
		if s.stream == nil {
			return
		}
		for val := range s.stream {
			outstream <- toFloat(val)
		}
	}()

	return NewFloatStream(outstream)
}

// Count the number of elements in the stream.
//
// This is a special case of a reduction and is equivalent to:
//   s.MapToInt(func(Entry) { return EntryInt(1) }).Sum()
//
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
func (s Stream) Count() int {
	if s.stream == nil {
		return 0
	}

	count := 0
	for range s.stream {
		count++
	}

	return count
}

// AllMatch returns whether all of the elements in the stream
// satisfy the predicate.
//
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("fourth"),
//  }, 1e3).
//      AllMatch(func(e ƒ.Entry) bool {
//          return strings.Contains(string(e.(ƒ.EntryString)), "t")
//      })
//  // Result: true
func (s Stream) AllMatch(p Predicate) bool {
	if s.stream == nil {
		return false
	}

	for val := range s.stream {
		if !p(val) {
			return false
		}
	}

	return true
}

// AnyMatch returns whether any of the elements in the stream
// satisfies the predicate.
//
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("four"),
//  }, 1e3).
//      AnyMatch(func(e ƒ.Entry) bool {
//          return e.Equal(ƒ.EntryString("three"))
//      })
//  // Result: true
func (s Stream) AnyMatch(p Predicate) bool {
	if s.stream == nil {
		return false
	}

	for val := range s.stream {
		if p(val) {
			return true
		}
	}

	return false
}

// NoneMatch returns whether none of the elements in the stream
// satisfies the predicate. It is the opposite of AnyMatch.
//
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("four"),
//  }, 1e3).
//      NoneMatch(func(e ƒ.Entry) bool { return e.Equal(ƒ.EntryString("nothing like this")) })
//  // Result: true
func (s Stream) NoneMatch(p Predicate) bool {
	return !s.AnyMatch(p)
}

// Drop the first 'n' elements of this stream and returns a new stream.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("fourth"),
//  }, 1e3).
//      Drop(2)
//  // Result: Stream of ƒ.EntryString("fourth")
func (s Stream) Drop(n uint64) Stream {
	return s.DropWhile(func() func(e Entry) bool {
		count := uint64(0)
		return func(e Entry) bool {
			count++
			return count <= n
		}
	}())
}

// DropWhile drops the first elements of this stream while the predicate
// is satisfied and returns a new stream.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("fourth"),
//  }, 1e3).
//      DropWhile(func(e ƒ.Entry) bool {
//          return e.Equal(ƒ.EntryString("three"))
//      })
//  // Result: Stream of ƒ.EntryString("two") and ƒ.EntryString("fourth")
func (s Stream) DropWhile(p Predicate) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() {
		defer close(outstream)
		if s.stream == nil {
			return
		}

		// drop elements as required
		for val := range s.stream {
			if p(val) {
				continue
			}
			outstream <- val
			break
		}

		// flush the remainder to outstream
		for val := range s.stream {
			outstream <- val
		}
	}()

	return NewStream(outstream)
}

// DropUntil drops the first elements of this stream until the predicate
// is satisfied and returns a new stream.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("fourth"),
//  }, 1e3).
//      DropUntil(func(e ƒ.Entry) bool {
//          return e.Equal(ƒ.EntryString("fourth"))
//      })
//  // Result: Stream of ƒ.EntryString("three") and ƒ.EntryString("two")
func (s Stream) DropUntil(p Predicate) Stream {
	return s.DropWhile(p.Negate())
}

// Last returns the last Entry in this stream.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("fourth"),
//  }, 1e3).
//      Last()
//  // Result: ƒ.EntryString("fourth")
func (s Stream) Last() Entry {
	return s.LastN(1)[0]
}

// LastN returns a slice of the last n elements in this stream.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("fourth"),
//  }, 1e3).
//      LastN(2)
//  // Result: []ƒ.Entry{ƒ.EntryString("two"), ƒ.EntryString("fourth")}
func (s Stream) LastN(n uint64) EntrySlice {
	s.panicIfNilChannel()

	if n < 1 {
		panic(PanicNoSuchElement)
	}

	val, ok := <-s.stream
	if !ok {
		panic(PanicNoSuchElement)
	}

	result := EntrySlice{val}

	count := uint64(result.Len())
	flushTrigger := uint64(100)
	if n > flushTrigger {
		flushTrigger = n
	}

	for val = range s.stream {
		result = append(result, val)
		if count++; count > flushTrigger {
			// this is simply to reduce the number of
			// slice resizing operations
			result = result[uint64(result.Len())-n:]
			count = 0
		}
	}

	if uint64(result.Len()) > n {
		return result[uint64(result.Len())-n:]
	}
	return result
}

// Head returns the first Entry in this stream.
//
// This function only consumes at most one element from the stream.
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("fourth"),
//  }, 1e3).
//      Head()
//  // Result: ƒ.EntryString("three")
func (s Stream) Head() Entry {
	head := s.HeadN(1)
	if head.Len() != 1 {
		panic(PanicNoSuchElement)
	}
	return head[0]
}

// HeadN returns a slice of the first n elements in this stream.
//
// This function only consumes at most 'n' elements from the stream.
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("fourth"),
//  }, 1e3).
//      HeadN(2)
//  // Result: []ƒ.Entry{ƒ.EntryString("three"), ƒ.EntryString("two")}
func (s Stream) HeadN(n uint64) EntrySlice {
	return s.Take(n).Collect(
		NewCollector(
			func() Entry { return EntrySlice{} },
			func(e1, e2 Entry) Entry { return e1.(EntrySlice).Append(e2) },
			nil,
		)).(EntrySlice)
}

// EndsWith returns true when this stream ends
// with the supplied elements.
//
// This is a potentially expensive method since it has
// to consume all the elements in the Stream.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("fourth"),
//  }, 1e3).
//      EndsWith([]ƒ.Entry{ƒ.EntryString("two"), ƒ.EntryString("fourth")})
//  // Result: true
func (s Stream) EndsWith(slice EntrySlice) bool {
	if slice.Len() == 0 {
		return false
	}

	endElements := func() EntrySlice {
		defer func() {
			// TODO: this doesn't look great... Need to re-write LastN like HeadN as a collect of TakeRight (to be implemented)
			_ = recover()
		}()

		return s.LastN(uint64(slice.Len()))
	}()

	if endElements.Len() != slice.Len() {
		return false
	}

	for idx, el := range slice {
		if !el.Equal(endElements[idx]) {
			return false
		}
	}

	return true
}

// StartsWith returns true when this stream starts
// with the elements in the supplied slice.
//
// This function only consume as much data from the stream as
// is necessary to prove (or disprove) it starts with the supplied
// slice data.
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("fourth"),
//  }, 1e3).
//      StartsWith([]ƒ.Entry{ƒ.EntryString("three"), ƒ.EntryString("two")})
//  // Result: true
func (s Stream) StartsWith(slice EntrySlice) bool {
	startElements := s.HeadN(uint64(slice.Len()))
	if slice.Len() == 0 || startElements.Len() != slice.Len() {
		return false
	}

	for idx, el := range slice {
		if !el.Equal(startElements[idx]) {
			return false
		}
	}

	return true
}

// Take returns a stream of the first 'n' elements of this stream.
//
// This function streams continuously until the 'n' elements are picked
// or the in-stream  is closed at which point the out-stream
// will be closed too.
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("fourth"),
//  }, 1e3).
//      Take(2)
//  // Result: Stream of []ƒ.Entry{ƒ.EntryString("three"), ƒ.EntryString("two")}
func (s Stream) Take(n uint64) Stream {
	counterIsLessThanOrEqualTo := func(maxCount uint64) Predicate {
		counter := uint64(0)
		return func(t Entry) bool {
			counter++
			return counter <= maxCount
		}
	}
	return s.TakeWhile(counterIsLessThanOrEqualTo(n))
}

// Limit is a synonym for Take.
func (s Stream) Limit(n uint64) Stream {
	return s.Take(n)
}

// TakeWhile returns a stream of the first elements of this
// stream while the predicate is satisfied.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("fourth"),
//  }, 1e3).
//      TakeWhile(func(e ƒ.Entry) bool {
//          return strings.HasPrefix(string(e.(ƒ.EntryString)), "t")
//      })
//  // Result: Stream of []ƒ.Entry{ƒ.EntryString("three"), ƒ.EntryString("two")}
func (s Stream) TakeWhile(p Predicate) Stream {
	s.panicIfNilChannel()

	outstream := make(chan Entry, cap(s.stream))

	go func() {
		defer close(outstream)

		for val := range s.stream {
			if !p(val) {
				return
			}
			outstream <- val
		}
	}()

	return NewStream(outstream)
}

// TakeUntil returns a stream of the first elements
// of this stream until the predicate is satisfied.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
//
// Example
//
//  ƒ.NewStreamFromSlice([]ƒ.Entry{
//      ƒ.EntryString("three"),
//      ƒ.EntryString("two"),
//      ƒ.EntryString("fourth"),
//  }, 1e3).
//      TakeUntil(func(e ƒ.Entry) bool {
//          return e.Equal(ƒ.EntryString("fourth"))
//      })
//  // Result: Stream of []ƒ.Entry{ƒ.EntryString("three"), ƒ.EntryString("two")}
func (s Stream) TakeUntil(p Predicate) Stream {
	return s.TakeWhile(p.Negate())
}

// Collect reduces and optionally mutates the stream with
// the supplied Collector.
//
// It should be noted that this method returns an `interface{}`
// which enables it to return `Entry` as well as any other Go types.
//
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
//
// Example
//
//  strs := EntrySlice{
//  	EntryString("a"),
//  	EntryString("bb"),
//  	EntryString("cc"),
//  	EntryString("ddd"),
//  }
//
//  NewStreamFromSlice(strs, 1e3).
//      Collect(
//          GroupingBy(
//              stringLength,
//              Mapping(
//                  stringToUpper,
//                  Filtering(
//                      stringLengthGreaterThan(1),
//                      ToEntrySlice()))))
//  // Result: map[1:[] 2:[BB CC] 3:[DDD]]
func (s Stream) Collect(c Collector) interface{} {
	s.panicIfNilChannel()

	result := c.supplier()
	for e := range s.stream {
		result = c.accumulator(result, e)
	}
	if c.finisher != nil {
		result = c.finisher(result)
	}
	return result
}

// ToSlice extracts the elements of the stream into
// an EntrySlice.
//
// This is a special case of a reduction.
//
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
func (s Stream) ToSlice() EntrySlice {
	result := EntrySlice{}
	if s.stream != nil {
		for val := range s.stream {
			result = append(result, val)
		}
	}
	return result
}

// Distinct returns a stream of the distinct elements of
// this stream.
//
// This operation is costly both in time and in memory. It is
// strongly recommended to use buffered channels for this operation.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) Distinct() Stream {
	s.panicIfNilChannel()

	outstream := make(chan Entry, cap(s.stream))

	go func() {
		defer close(outstream)

		unique := map[string]struct{}{}

		for val := range s.stream {
			hash := val.Hash()
			uniqueHash := fmt.Sprintf("%T%d", val, hash)[6:] // remove "fuego." prefix
			if _, isset := unique[uniqueHash]; !isset {
				unique[uniqueHash] = struct{}{}
				outstream <- val
			}
		}
	}()

	return NewStream(outstream)
}

// panicIfNilChannel panics if s.stream is nil.
func (s Stream) panicIfNilChannel() {
	if s.stream == nil {
		panic(PanicMissingChannel)
	}
}
