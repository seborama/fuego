package fuego

//go:generate ./bin/maptoXXX

import (
	"fmt"

	"go.uber.org/zap"
)

// Stream is a sequence of elements supporting sequential and
// (in specific circumstances) parallel operations.
//
// A Stream is a wrapper over a Go channel ('nil' channels are prohibited).
//
// NOTE:
//
// Concurrent streams are challenging to implement owing to
// ordering issues in parallel processing. At the moment, the view
// is that the most sensible approach is to delegate control to users.
// Multiple fuego streams can be created and data distributed
// across as desired. This empowers users of fuego to implement the
// desired behaviour of their pipelines.
//
// As of v8.0.0, fuego offers ordered concurrency for some linear
// operations such as Map().
//
// Creation
//
// When providing a Go channel to create a Stream, the Stream's remains open unless the channel
// is closed by the producer.
// Should the producer not close the channel unintentionally, the Go function will stray.
//
// Streams created from a slice are bounded since the slice has finite content.
type Stream[T any] struct {
	stream      chan T
	concurrency int
}

// NewStream creates a new Stream.
//
// This function does not close the provided channel.
func NewStream[T any](c chan T) Stream[T] {
	return NewConcurrentStream(c, 0)
}

// NewConcurrentStream creates a new Stream with a degree of concurrency of n.
func NewConcurrentStream[T any](c chan T, n int) Stream[T] {
	return Stream[T]{
		stream:      c,
		concurrency: n,
	}
}

// NewStreamFromSlice creates a new Stream from a Go slice.
//
// The slice data is published to the stream after which the stream is closed.
func NewStreamFromSlice[T any](slice []T, bufsize int) Stream[T] {
	c := make(chan T, bufsize)

	go func() {
		defer close(c) // slices have finite size: close stream after all data was read.

		for _, element := range slice {
			c <- element
		}
	}()

	return NewStream(c)
}

// Concurrency returns the stream's concurrency level (i.e. parallelism).
func (s Stream[T]) Concurrency() int {
	return s.concurrency
}

// Concurrent sets the level of concurrency for this Stream.
//
// This is used for concurrent methods such as Stream.Map.
//
// Note that to switch off concurrency, you should provide n = 0.
// With n = 1, concurrency is internal whereby the Stream writer
// will not block on writing a single element (i.e. buffered
// channel of 1). This already provides significant processing gains.
//
// Performance:
//
// Channels are inherently expensive to use owing to their internal
// mutex lock.
//
// Benefits will ONLY be observed when the execution has a degree
// of latency (at the very least, several dozens of nanoseconds).
// The higher the latency, the better the gains from concurrency
// (even on a single CPU core).
//
// If latency is too low or next to none, using concurrency will
// likely be slower than without, particularly when no CPU core is
// available.
func (s Stream[T]) Concurrent(n int) Stream[T] {
	// This is not accurate but improves performance (by avoiding the
	// creation of a new channel and iterating through this one).
	// It should be safe.
	return NewConcurrentStream(s.stream, n)
}

// Any is an alias for type `any`.
type Any any

// Map returns a Stream consisting of the result of
// applying the given function to the elements of this stream.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream[T]) Map(mapper Function[T, Any]) Stream[Any] {
	return NewConcurrentStream(orderlyConcurrentDo(s, mapper), s.concurrency)
}

// orderlyConcurrentDo executes a Function on the stream.
// Execution is concurrent and order is preserved.
// See note on method Map() about the lack of support for parameterised methods in Go.
func orderlyConcurrentDo[T, U any](s Stream[T], fn Function[T, U]) chan U {
	outstream := make(chan U, cap(s.stream))

	go func() {
		defer close(outstream)

		if s.stream == nil {
			return
		}

		pipelineCh := make(chan chan U, s.concurrency)

		pipelineWriter := func(pipelineWCh chan chan U) {
			defer close(pipelineWCh)

			for val := range s.stream {
				resultCh := make(chan U, 1)
				pipelineWCh <- resultCh

				go func(resultCh chan<- U, val T) {
					defer close(resultCh)
					resultCh <- fn(val)
				}(resultCh, val)
			}
		}

		go func() {
			pipelineWriter(pipelineCh)
		}()

		pipelineReader := func(pipelineRCh chan chan U) {
			for resultCh := range pipelineRCh {
				outstream <- <-resultCh
			}
		}
		pipelineReader(pipelineCh)
	}()

	return outstream
}

// FlatMap takes a StreamFunction to flatten the entries
// in this stream and produce a new stream.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream[T]) FlatMap(mapper StreamFunction[T, Any]) Stream[Any] {
	return NewConcurrentStream(orderlyConcurrentDoStream(s, mapper), s.concurrency)
}

// orderlyConcurrentDoStream executes a StreamFunction on the stream.
// Execution is concurrent and order is preserved.
func orderlyConcurrentDoStream[T, U any](s Stream[T], streamfn StreamFunction[T, U]) chan U {
	outstream := make(chan U, cap(s.stream))

	go func() {
		defer close(outstream)

		if s.stream == nil {
			return
		}

		pipelineCh := make(chan chan Stream[U], s.concurrency)

		pipelineWriter := func(pipelineWCh chan chan Stream[U]) {
			defer close(pipelineWCh)

			for val := range s.stream {
				resultCh := make(chan Stream[U], 1)
				pipelineWCh <- resultCh

				go func(resultCh chan<- Stream[U], val T) {
					defer close(resultCh)
					resultCh <- streamfn(val)
				}(resultCh, val)
			}
		}

		go func() {
			pipelineWriter(pipelineCh)
		}()

		pipelineReader := func(pipelineRCh chan chan Stream[U]) {
			for resultCh := range pipelineRCh {
				val := <-resultCh
				val.ForEach(func(e U) {
					outstream <- e
				})
			}
		}
		pipelineReader(pipelineCh)
	}()

	return outstream
}

// Filter returns a stream consisting of the elements of this stream that
// match the given predicate.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream[T]) Filter(predicate Predicate[T]) Stream[T] {
	outstream := make(chan T, cap(s.stream))

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

	return NewConcurrentStream(outstream, s.concurrency)
}

// LeftReduce accumulates the elements of this Stream by applying the given function.
//
// This is a continuous terminal operation. It will only complete if the producer closes the stream.
func (s Stream[T]) LeftReduce(f2 BiFunction[T, T, T]) T {
	if s.stream == nil {
		var t T
		return t // TODO: return Optional
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
func (s Stream[T]) Reduce(f2 BiFunction[T, T, T]) T {
	return s.LeftReduce(f2)
}

// Intersperse inserts an element between all elements of this Stream.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream[T]) Intersperse(e T) Stream[T] {
	outstream := make(chan T, cap(s.stream))

	go func() {
		defer close(outstream)
		if s.stream == nil {
			return
		}

		// this is to get around the inability to test generic types for nil in Go 1.18
		// nolint: wsl
		select {
		case val, ok := <-s.stream:
			if !ok {
				return
			}
			outstream <- val
		}

		for val := range s.stream {
			outstream <- e
			outstream <- val
		}
	}()

	return NewConcurrentStream(outstream, s.concurrency)
}

// GroupBy groups the elements of this Stream by classifying them.
//
// This is a continuous terminal operation and hence expects the producer to close the stream
// in order to complete.
func (s Stream[T]) GroupBy(classifier Function[T, Any]) map[Any][]T {
	resultMap := make(map[Any][]T)

	if s.stream != nil {
		for val := range s.stream {
			k := classifier(val)

			if resultMap[k] == nil {
				resultMap[k] = []T{}
			}

			resultMap[k] = append(resultMap[k], val)
		}
	}

	return resultMap
}

// ForEach executes the given consumer function for each entry in this stream.
//
// This is a continuous terminal operation. It will only complete if the producer closes the stream.
func (s Stream[T]) ForEach(c Consumer[T]) {
	if s.stream == nil {
		zap.L().Debug("empty stream")
		return
	}

	for val := range s.stream {
		zap.L().Debug("calling consumer", zap.Any("value", val))
		c(val)
	}
}

// Peek is akin to ForEach but returns the Stream.
//
// This is useful e.g. for debugging.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream[T]) Peek(consumer Consumer[T]) Stream[T] {
	outstream := make(chan T, cap(s.stream))

	go func() {
		defer close(outstream)
		s.ForEach(func(e T) {
			consumer(e)
			outstream <- e
		})
	}()

	return NewConcurrentStream(outstream, s.concurrency)
}

// ToSlice extracts the elements of the stream into a []T.
//
// This is a special case of a reduction.
//
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete (or
// it will block).
func (s Stream[T]) ToSlice() []T {
	result := []T{}

	if s.stream != nil {
		for val := range s.stream {
			result = append(result, val)
		}
	}

	return result
}

// Distinct returns a stream of the distinct elements of this stream.
// Distinctiveness is determined via the provided hashFn.
//
// This operation is costly both in time and in memory. It is
// strongly recommended to use buffered channels for this operation.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream[T]) Distinct(hashFn func(T) uint32) Stream[T] {
	if s.stream == nil {
		panic(PanicMissingChannel)
	}

	outstream := make(chan T, cap(s.stream))

	go func() {
		defer close(outstream)

		unique := map[string]struct{}{}

		for val := range s.stream {
			// hash is prefixed with the type in case T is an interface implemented by 2 or more types
			// that are present on the stream.
			uniqueHash := fmt.Sprintf("%T%d", val, hashFn(val))
			if _, ok := unique[uniqueHash]; !ok {
				unique[uniqueHash] = struct{}{}
				outstream <- val
			}
		}
	}()

	return NewConcurrentStream(outstream, s.concurrency)
}

// StreamAny returns this stream as a Stream[Any].
func (s Stream[T]) StreamAny() Stream[Any] {
	rCh := make(chan Any, cap(s.stream))

	r := NewConcurrentStream(rCh, s.concurrency)

	go func() {
		defer close(rCh)

		s.ForEach(func(el T) {
			r.stream <- el
		})
	}()

	return r
}
