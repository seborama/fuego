package fuego

// TODO: consider two types of streams: CStreams (channel based as shown here) and SStreams (slice based). The former allows for infinite streams and thinner memory usage within the CStream object but lacks performance when the operation requires to deal with the end of the steam (it has to consume all the elements of the steam sequentially). SStreams require the entire data to be stored internally from the onset. However,  slices are seekable and can read from the end or be consumed backwards easily.

// TODO: a stream should probably be marked as invalid after most (or all?) operations on it because the channel will have likely changed state.

// TODO list:
// Methods with ** require the Stream to be finite and closed (or use a Future, perhaps Future.Stream()?)
// **Distinct()
// **DropRight(uint64) - drops the last n elements of the Stream. Only meaningful if the Stream is closed.
// FilterNot(Predicate) - <=> to Filter(Not(Predicate))
// Peek(Consumer) - Like ForEach but returns Stream as it was at the point of Peek
// Limit(uint64) - Returns a Stream consisting of at most n elements.
// MapToString(ToStringFunction)
// FindAny / FindFirst?
// FlatMapToXXX (Int, Uint, etc) => is this the same as FlatMap().MapToXXX()?
// **Sorted(Comparator)
// Contains
// ContainsAll
// Tail
// Fold / FoldLeft
// Find / FindLast
// Zip / Unzip
// Concat?
// Range(from, toExclusive)?
// RangeBy(from,toExclusive, step)?
// Unfold()
// Fold()?

// PanicMissingChannel signifies that the Stream is missing a channel.
const PanicMissingChannel = "stream creation requires a channel"

// Stream is a sequence of elements supporting sequential and
// (in the future?) parallel operations.
type Stream struct {
	stream chan Entry
}

// NewStream creates a new Stream.
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
// The slice data is published to the stream after which the
// stream is closed.
func NewStreamFromSlice(slice EntrySlice, bufsize int) Stream {
	c := make(chan Entry, bufsize)

	go func() {
		// TODO: add test to confirm the stream gets closed
		defer close(c) // slices have finite size: close stream after all data was read.
		for _, element := range slice {
			c <- element
		}
	}()

	return NewStream(c)
}

// TODO: implement NewStreamFromMap -> Stream of Keys / Stream of Values?

// Map returns a Stream consisting of the result of
// applying the given function to the elements of this stream.
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) Map(mapper Function) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() {
		defer close(outstream) // TODO: add test to confirm the stream gets closed
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
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) FlatMap(mapper StreamFunction) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() {
		defer close(outstream) // TODO: add test to confirm the stream gets closed
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
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) Filter(predicate Predicate) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() {
		defer close(outstream) // TODO: add test to confirm the stream gets closed
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
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
func (s Stream) ForEach(consumer Consumer) {
	if s.stream == nil {
		return
	}

	for val := range s.stream {
		consumer(val)
	}
}

// LeftReduce accumulates the elements of this Stream by
// applying the given function.
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
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
// See LeftReduce for more info.
func (s Stream) Reduce(f2 BiFunction) Entry {
	return s.LeftReduce(f2)
}

// Intersperse inserts an element between all elements of this Stream.
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) Intersperse(e Entry) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() {
		defer close(outstream) // TODO: add test to confirm the stream gets closed
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
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) MapToInt(toInt ToIntFunction) IntStream {
	outstream := make(chan EntryInt, cap(s.stream))

	go func() {
		defer close(outstream) // TODO: add test to confirm the stream gets closed
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
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) MapToFloat(toFloat ToFloatFunction) FloatStream {
	outstream := make(chan EntryFloat, cap(s.stream))

	go func() {
		defer close(outstream) // TODO: add test to confirm the stream gets closed
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
// This is a special case of a reduction and is equivalent to:
//   s.MapToInt(func(Entry) { return EntryInt(1) }).Sum()
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
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
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
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
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
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
func (s Stream) NoneMatch(p Predicate) bool {
	return !s.AnyMatch(p)
}

// Drop the first 'n' elements of this stream and returns a new stream.
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
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) DropWhile(p Predicate) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() {
		defer close(outstream) // TODO: add test to confirm the stream gets closed
		if s.stream == nil {
			return
		}

		// drop elements as required
		for val := range s.stream {
			if p(val) {
				continue
			}
			outstream <- val // this element belongs!
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
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) DropUntil(p Predicate) Stream {
	return s.DropWhile(p.Negate())
}

// Last returns the last Entry in this stream.
func (s Stream) Last() Entry {
	return s.LastN(1)[0]
}

// LastN returns a slice of the last n elements in this stream.
func (s Stream) LastN(n uint64) EntrySlice {
	if s.stream == nil {
		panic(PanicMissingChannel)
	}

	if n < 1 {
		panic(PanicNoSuchElement)
	}

	val, ok := <-s.stream
	if !ok {
		panic(PanicNoSuchElement)
	}

	result := EntrySlice{val}

	count := uint64(len(result))
	flushTrigger := uint64(100)
	if n > flushTrigger {
		flushTrigger = n
	}

	for val = range s.stream {
		result = append(result, val)
		if count++; count > flushTrigger {
			// this is simply to reduce the number of
			// slice resizing operations
			result = result[uint64(len(result))-n:]
			count = 0
		}
	}

	if uint64(len(result)) > n {
		return result[uint64(len(result))-n:]
	}
	return result
}

// Head returns the first Entry in this stream.
func (s Stream) Head() Entry {
	head := s.HeadN(1)
	if len(head) != 1 {
		panic(PanicNoSuchElement)
	}
	return head[0]
}

// HeadN returns a slice of the first n elements in this stream.
func (s Stream) HeadN(n uint64) EntrySlice {
	if s.stream == nil {
		panic(PanicMissingChannel)
	}

	return s.Take(n).Collect(
		NewCollector(
			func() Entry { return EntrySlice{} },
			func(e1, e2 Entry) Entry { return e1.(EntrySlice).Append(e2) },
			nil,
		)).(EntrySlice)
}

// EndsWith returns true when this stream ends
// with the supplied elements.
func (s Stream) EndsWith(slice EntrySlice) bool {
	defer func() {
		// TODO: this doesn't look great... Need to re-write LastN like HeadN as a collect of TakeRight (to be implemented)
		_ = recover()
	}()

	endElements := s.LastN(uint64(len(slice)))
	if len(endElements) != len(slice) {
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
// This function streams continuously until the 'n' elements are picked
// or the in-stream  is closed at which point the out-stream
// will be closed too.
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

// TakeWhile returns a stream of the first elements of this
// stream while the predicate is satisfied.
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) TakeWhile(p Predicate) Stream {
	if s.stream == nil {
		panic(PanicMissingChannel)
	}

	outstream := make(chan Entry, cap(s.stream))

	go func() {
		defer close(outstream) // TODO: add test to confirm the stream gets closed

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
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) TakeUntil(p Predicate) Stream {
	return s.TakeWhile(p.Negate())
}

// Collect reduces and optionally mutates the stream with
// the supplied Collector.
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
func (s Stream) Collect(c Collector) interface{} {
	if s.stream == nil {
		panic(PanicMissingChannel)
	}

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
// This is a special case of a reduction.
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
