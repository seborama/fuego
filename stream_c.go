package fuego

// This file is really an extension of stream.go.
//
// It contains the "concurrent" methods of Stream.
//
// The reason for splitting is two-fold:
// - convenient separation
// - help reduce the growth of stream.go

import "sync"

// PanicInvalidConcurrencyLevel signifies that the Stream is missing a channel.
const PanicInvalidConcurrencyLevel = "stream concurrency must be 1 or more"

// Concurrent sets the level of concurrency for this Stream.
// This is used for concurrent methods such as ForEachC.
func (s Stream) Concurrent(n uint) {
	s.concurrencyLevel = n
	s.panicIfInvalidConcurrency()
}

// ForEachC is a concurrent wrapper of ForEach.
// The level of concurrency is set by the last call made to method Concurrent.
// See 'ForEach' for full details.
// See 'concurrentDo' for limitations.
func (s Stream) ForEachC(consumer Consumer) {
	s.concurrentDo(func() { s.ForEach(consumer) })
}

// concurrentDo is a wrapper for concurrently running methods
// of Stream that are not concurrent.
// The level of concurrency is set by 'Concurrent'.
// This function concurrently consumes from the stream's channel.
// Consumption is ordered by the stream's channel but output is
// unordered: a slow consumer will be "out-raced" by faster consumers.
// Channels are inherently expensive to use owing to their internal
// mutex lock.
// Benefits will ONLY be observed when the execution of the function
// passed as an argument has a degree of latency (at the very least,
// several dozens of nanoseconds). The higher the latency, the better
// the gains from concurrency (even on a single CPU core).
// If latency is too low or none, using concurrency will likely be
// slower than without, particularly when no CPU core is available.
func (s Stream) concurrentDo(f func()) {
	s.panicIfInvalidConcurrency()

	var wg sync.WaitGroup

	for i := uint(1); i <= s.concurrencyLevel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	wg.Wait()
}

// panicIfInvalidConcurrency panics if the concurrency level
// is not valid.
func (s Stream) panicIfInvalidConcurrency() {
	if s.concurrencyLevel < 2 {
		panic(PanicInvalidConcurrencyLevel)
	}
}
