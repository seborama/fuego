package fuego

// TODO: consider two types of streams: CStreams (channel based as shown here) and SStreams (slice based). The former allows for infinite streams and thinner memory usage within the CStream object but lacks performance when the operation requires to deal with the end of the steam (it has to consume all the elements of the steam sequentially). SStreams require the entire data to be stored internally from the onset. However,  slices are seekable and can read from the end or be consumed backwards easily.

// Stream is a sequence of elements supporting sequential and
// (in the future?) parallel operations.
type Stream struct {
	stream chan Entry
}

// NewStream creates a new Stream.
func NewStream(c chan Entry) Stream {
	return Stream{
		stream: c,
	}
}

// NewStreamFromSlice creates a new Stream from a Go slice.
func NewStreamFromSlice(s []Entry) Stream {
	c := make(chan Entry, 1e3)
	defer close(c)
	for _, element := range s {
		c <- element
	}

	return NewStream(c)
}

// TODO: implement NewStreamFromMap?

// Map returns a slice of channel of Set consisting of the results of
// applying the given function to the elements of this Set
func (s Stream) Map(mapper Function) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() { // TODO: introduce a cut-off to prevent the go func from straying
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

// Filter returns a stream consisting of the elements of this stream that
// match the given predicate.
func (s Stream) Filter(predicate Predicate) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() { // TODO: introduce a cut-off to prevent the go func from straying
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
func (s Stream) Reduce(f2 BiFunction) Entry {
	return s.LeftReduce(f2)
}

// Intersperse inserts an element between all elements of this Stream.
func (s Stream) Intersperse(e Entry) Stream {
	outstream := make(chan Entry, cap(s.stream))

	go func() { // TODO: introduce a cut-off to prevent the go func from straying
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
// MapToInt8(ToIntFunction) \
// MapToInt16(ToIntFunction)  \
// MapToInt32(ToIntFunction)   \
// MapToInt64(ToIntFunction)    \
// MapToUint8(ToUintFunction)     \
// MapToUint16(ToUintFunction)      \ All these exist as syntactic sugar to allow for things like x.Sum(), x.Max(), x.Substr(), etc
// MapToUint32(ToUintFunction)      /
// MapToUint64(ToUintFunction)     /
// MapToFloat32(ToFloatFunction)   /
// MapToFloat64(ToFloatFunction)  /
// MapToString(ToStringFunction) /
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
func (s Stream) MapToInt(toInt ToIntFunction) IntStream {
	outstream := make(chan EntryInt, cap(s.stream))

	go func() { // TODO: introduce a cut-off to prevent the go func from straying
		defer close(outstream)
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
func (s Stream) Count() int {
	if s.stream == nil {
		return 0
	}

	func() {
		defer func() { _ = recover() }()
		s.Close()
	}()

	count := 0
	for range s.stream {
		count++
	}

	return count
}

// Close the stream.
// Panics if s.stream is nil or already closed.
func (s Stream) Close() {
	close(s.stream)
}
