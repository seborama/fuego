package fuego

// TODO: consider two types of streams: CStreams (channel based as shown here) and SStreams (slice based). The former allows for infinite streams and thinner memory usage within the CStream object but lacks performance when the operation requires to deal with the end of the steam (it has to consume all the elements of the steam sequentially). SStreams require the entire data to be stored internally from the onset. However,  slices are seekable and can read from the end or be consumed backwards easily.

// TODO list:
// Methods with ** require the Stream to be finite and closed (or use a Future, perhaps Future.Stream()?)
// **Distinct()
// **DropRight(uint64) - drops the last n elements of the Stream. Only meaningful if the Stream is closed.
// FilterNot(Predicate) - <=> to Filter(Not(Predicate))
// Peek(Consumer) - Like ForEach but returns Stream as it was at the point of Peek
// Limit(uint64) - Returns a Stream consisting of at most n elements.
// MapToString(ToStringFunction)
// FlatMap
// FlatMapToXXX (Int, Uint, etc) => is this the same as FlatMap().MapToXXX()?
// **Sorted(Comparator)
// TakeWhile - see DropXXX()?
// ToSlice
// Collect
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
func NewStreamFromSlice(slice []Entry, bufsize int) Stream {
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

// Map returns a slice of channel of Set consisting of the results of
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
			resultMap[k] = append(resultMap[k], val)
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

// Drop the first 'n' elements of this stream.
func (s Stream) Drop(n uint64) Stream {
	if n >= 1 && s.stream != nil {
		i := uint64(1)
		for range s.stream {
			if i >= n {
				break
			}
			i++
		}
	}
	return s
}

// DropWhile drops the first elements of this stream while the predicate
// is satisfied.
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
			outstream <- val // this one belongs!
			break
		}

		// flush the remainder to outstream
		for val := range s.stream {
			outstream <- val
		}
	}()

	return NewStream(outstream)
}

// Last returns the last element in this stream.
func (s Stream) Last() Entry {
	return s.LastN(1)[0]
}

// LastN returns the last n elements in this stream.
func (s Stream) LastN(n uint64) []Entry {
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

	result := []Entry{val}

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

// Head returns the first element in this stream.
func (s Stream) Head() Entry {
	return s.HeadN(1)[0]
}

// HeadN returns the first n elements in this stream.
func (s Stream) HeadN(n uint64) []Entry {
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

	result := []Entry{val}

	count := uint64(len(result))

	for val = range s.stream {
		result = append(result, val)
		if count++; count >= n {
			break
		}
	}

	return result
}

// EndsWith returns true when this stream ends
// with the supplied elements.
func (s Stream) EndsWith(slice []Entry) bool {
	defer func() {
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
func (s Stream) StartsWith(slice []Entry) bool {
	defer func() {
		_ = recover()
	}()

	startElements := s.HeadN(uint64(len(slice)))
	if len(startElements) != len(slice) {
		return false
	}

	for idx, el := range slice {
		if !el.Equal(startElements[idx]) {
			return false
		}
	}

	return true
}
