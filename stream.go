package fuego

// TODO: consider two types of streams: CStreams (channel based as shown here) and SStreams (slice based). The former allows for infinite streams and thinner memory usage within the CStream object but lacks performance when the operation requires to deal with the end of the steam (it has to consume all the elements of the steam sequentially). SStreams require the entire data to be stored internally from the onset. However,  slices are seekable and can read from the end or be consumed backwards easily.

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
func NewStreamFromSlice(slice []Entry) Stream {
	c := make(chan Entry, 1e3)

	go func() {
		// TODO: add test to confirm the stream gets closed
		defer close(c) // slices have finite size: close stream after all data was read.
		for _, element := range slice {
			c <- element
		}
	}()

	return NewStream(c)
}

// TODO: implement NewStreamFromMap?

// Map returns a slice of channel of Set consisting of the results of
// applying the given function to the elements of this stream.
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) Map(mapper Function) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() { // TODO: introduce a cut-off to prevent the go func from straying
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

	go func() { // TODO: introduce a cut-off to prevent the go func from straying
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

	go func() { // TODO: introduce a cut-off to prevent the go func from straying
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

// TODO list:
// Methods with ** require the Stream to be finite and closed (or use a Future, perhaps Future.Stream()?)
// **Distinct()
// Drop(uint64) - drops the first n elements of the Stream.
// **DropRight(uint64) - drops the last n elements of the Stream. Only meaningful if the Stream is closed.
// DropWhile(Predicate) - Returns the a Stream representing the longest suffix of this iterable whose first element does not satisfy the predicate
// FilterNot(Predicate) - <=> to Filter(Not(Predicate))
// **EndsWith([]Entry) - Tests whether this Stream ends with the []Entry
// Peek(Consumer) - Like ForEach but returns Stream as it was at the point of Peek
// Limit(uint64) - Returns a Stream consisting of at most n elements.
// MapToString(ToStringFunction)
// FlatMap
// FlatMapToXXX (Int, Uint, etc) => is this the same as FlatMap().MapToXXX()?
// **Sorted(Comparator)
// Skip
// TakeWhile - see DropXXX()
// ToSlice
// Collect
// Contains
// ContainsAll
// Head
// Last
// Tail
// Count (as a map reduction operation - is that different to the Count() already implemented?)
// Fold / FoldLeft
// Find / FindLast
// Zip / Unzip
// Concat?
// Range(from, toExclusive)?
// RangeBy(from,toExclusive, step)?
// Unfold()
// Fold()?

// MapToInt produces an EntryInt stream.
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream) MapToInt(toInt ToIntFunction) IntStream {
	outstream := make(chan EntryInt, cap(s.stream))

	go func() { // TODO: introduce a cut-off to prevent the go func from straying
		defer close(outstream) // TODO: add test to confirm the stream gets closed
		if s.stream == nil {
			return
		}
		for val := range s.stream {
			outstream <- toInt(val)
		}
	}()

	return IntStream{
		stream: outstream,
	}
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

// Close the stream and returns true if success.
func (s Stream) Close() bool {
	closed := false
	func() {
		defer func() { _ = recover() }()
		close(s.stream)
		closed = true
	}()
	return closed
}
