package fuego

type Stream[T any] struct {
	stream           chan T
	concurrencyLevel int
}

// NewConcurrentStream creates a new Stream with a degree
// of concurrency of n.
func NewConcurrentStream[T any](c chan T, n int) Stream[T] {
	s := Stream[T]{
		stream:           c,
		concurrencyLevel: n,
	}

	return s
}

type R any

// Go does not support parameterised methods:
// https://go.googlesource.com/proposal/+/master/design/43651-type-parameters.md#no-parameterized-methods
//
// This is not possible:
// func (s Stream[T]) Map[R any](mapper Function[T, R]) Stream[R] {...}
//                        ^^^^^^
//
// One option would be to make `Map` a function rather than a method but constructs would be chained right-to-left
// rather than left-to-right, which I think is awkward.
// Example: "Map(Map(stream,f1),f2)" instead of "stream.Map(f1).Map(f2)".
func (s Stream[T]) Map(mapper Function[T, R]) Stream[R] {
	return NewConcurrentStream(s.orderlyConcurrentDo(mapper), s.concurrencyLevel)
}

func To[T any, R any](from Stream[T]) Stream[R] {
	outstream := make(chan R)
	return Stream[R]{
		stream:           outstream,
		concurrencyLevel: from.concurrencyLevel,
	}
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

		pipelineCh := make(chan chan R, s.concurrencyLevel)

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
