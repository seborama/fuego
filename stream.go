package fuego

import (
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
//
// Example:
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

// R is an alias for type `any`.
type R any

// Map returns a Stream consisting of the result of
// applying the given function to the elements of this stream.
//
// This function streams continuously until the in-stream is closed at
// which point the out-stream will be closed too.
func (s Stream[T]) Map(mapper Function[T, R]) Stream[R] {
	return NewConcurrentStream(s.orderlyConcurrentDo(mapper), s.concurrency)
}

// orderlyConcurrentDo executes a Function on the stream.
// Execution is concurrent and order is preserved.
// See note on method Map() about the lack of support for parameterised methods in Go.
func (s Stream[T]) orderlyConcurrentDo(fn Function[T, R]) chan R {
	outstream := make(chan R, cap(s.stream))

	go func() {
		defer close(outstream)

		if s.stream == nil {
			return
		}

		pipelineCh := make(chan chan R, s.concurrency)

		pipelineWriter := func(pipelineWCh chan chan R) {
			defer close(pipelineWCh)

			for val := range s.stream {
				resultCh := make(chan R, 1)
				pipelineWCh <- resultCh
				go func(resultCh chan<- R, val T) {
					defer close(resultCh)
					resultCh <- fn(val)
				}(resultCh, val)
			}
		}

		go func() {
			pipelineWriter(pipelineCh)
		}()

		pipelineReader := func(pipelineRCh chan chan R) {
			for resultCh := range pipelineRCh {
				outstream <- <-resultCh
			}
		}
		pipelineReader(pipelineCh)
	}()

	return outstream
}

// ForEach executes the given consumer function for each entry in this stream.
//
// This is a continuous terminal operation. It will only complete if the producer closes the stream.
//
// Example:
// total := 0
//
// computeSumTotal := func(value int) {
// 	total += value
// }
//
// NewStreamFromSlice([]int{1, 3, 2}, 0).
// 	ForEach(computeSumTotal)
//
// fmt.Println("total =", total) // total == 6
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

// SC is a typed Stream cast function from a non-parameterised Stream[R] to a parameterised Stream[U].
// SC receives a typed Stream[U].
//
// SC exists to address the current lack of support in Go for parameterised methods.
// See doc.go for more details.
//
// See C for A typed cast.
func SC[U any](from Stream[R], to Stream[U]) Stream[U] {
	toCh := make(chan U, from.concurrency)
	to.stream = toCh

	go func() {
		defer close(to.stream)

		for f := range from.stream {
			to.stream <- interface{}(f).(U)
		}
	}()

	return to
}

// C is a typed cast function from a non-parameterised Stream[R] to a parameterised type Stream[U].
// C receives a type U and creates a Stream[U].
//
// C exists to address the current lack of support in Go for parameterised methods.
// See doc.go for more details.
//
// See SC for A Stream cast.
func C[U any](from Stream[R], to U) Stream[U] {
	toCh := make(chan U, from.concurrency)

	toStream := NewConcurrentStream(toCh, from.concurrency)

	go func() {
		defer close(toStream.stream)

		for f := range from.stream {
			toStream.stream <- interface{}(f).(U)
		}
	}()

	return toStream
}

// Useful constants that represent native go types.
const (
	Bool       = false
	Int        = int(0)
	Int8       = int8(0)
	Int16      = int16(0)
	Int32      = int32(0)
	Int64      = int64(0)
	Uint       = uint(0)
	Uint8      = uint8(0)
	Uint16     = uint16(0)
	Uint32     = uint32(0)
	Uint64     = uint64(0)
	Float32    = float32(0)
	Float64    = float64(0)
	Complex64  = complex64(0)
	Complex128 = complex128(0)
	String     = ""
)

// Useful variables that represent native go types.
var (
	SBool       = []bool{}
	SInt        = []int{}
	SInt8       = []int8{}
	SInt16      = []int16{}
	SInt32      = []int32{}
	SInt64      = []int64{}
	SUint       = []uint{}
	SUint8      = []uint8{}
	SUint16     = []uint16{}
	SUint32     = []uint32{}
	SUint64     = []uint64{}
	SFloat32    = []float32{}
	SFloat64    = []float64{}
	SComplex64  = []complex64{}
	SComplex128 = []complex128{}
	SString     = []string{}
)

// Useful variables that represent native go types.
var (
	BoolPtr       = ptr(Bool)
	IntPtr        = ptr(Int)
	Int8Ptr       = ptr(Int8)
	Int16Ptr      = ptr(Int16)
	Int32Ptr      = ptr(Int32)
	Int64Ptr      = ptr(Int64)
	UintPtr       = ptr(Uint)
	Uint8Ptr      = ptr(Uint8)
	Uint16Ptr     = ptr(Uint16)
	Uint32Ptr     = ptr(Uint32)
	Uint64Ptr     = ptr(Uint64)
	Float32Ptr    = ptr(Float32)
	Float64Ptr    = ptr(Float64)
	Complex64Ptr  = ptr(Complex64)
	Complex128Ptr = ptr(Complex128)
	StringPtr     = ptr(String)
)

// Useful variables that represent native go types.
var (
	SBoolPtr       = []*bool{}
	SIntPtr        = []*int{}
	SInt8Ptr       = []*int8{}
	SInt16Ptr      = []*int16{}
	SInt32Ptr      = []*int32{}
	SInt64Ptr      = []*int64{}
	SUintPtr       = []*uint{}
	SUint8Ptr      = []*uint8{}
	SUint16Ptr     = []*uint16{}
	SUint32Ptr     = []*uint32{}
	SUint64Ptr     = []*uint64{}
	SFloat32Ptr    = []*float32{}
	SFloat64Ptr    = []*float64{}
	SComplex64Ptr  = []*complex64{}
	SComplex128Ptr = []*complex128{}
	SStringPtr     = []*string{}
)

func ptr[T any](t T) *T { return &t }
