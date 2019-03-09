package fuego

import (
	"sync"
)

// CStream is a concurrent stream.
type CStream struct {
	streams []Stream
}

// NewCStream creates a new concurrent stream.
func NewCStream(channels []chan Entry) CStream {
	cs := CStream{}
	for _, channel := range channels {
		cs.streams = append(cs.streams, NewStream(channel))
	}
	return cs
}

// AddStreamFromChannels adds Streams derived from the supplied channels to this CStream.
func (cs CStream) AddStreamFromChannels(channels []chan Entry) CStream {
	for _, channel := range channels {
		cs.streams = append(cs.streams, NewStream(channel))
	}
	// TODO rather than modify in-place, make a copy of cs return the copy? Note: is it worth it because channels will copy by reference anyway?
	return cs
}

// AddStreamFromSlices adds Streams derived from the supplied slices to this CStream.
func (cs CStream) AddStreamFromSlices(slices []EntrySlice, bufsize int) CStream {
	for _, slice := range slices {
		cs.streams = append(cs.streams, NewStreamFromSlice(slice, bufsize))
	}
	// TODO rather than modify in-place, make a copy of cs return the copy? Note: is it worth it because channels will copy by reference anyway?
	return cs
}

// AddStreams adds Streams to this CStream.
func (cs CStream) AddStreams(streams []Stream) CStream {
	cs.streams = append(cs.streams, streams...)
	// TODO rather than modify in-place, make a copy of cs and return the copy? Note: is it worth it because channels will copy by reference anyway?
	return cs
}

// ForEach is the concurrent equivalent of Stream.ForEach.
// See Stream.ForEach for further information.
func (cs CStream) ForEach(consumer Consumer) {
	var wg sync.WaitGroup

	for _, stream := range cs.streams {
		wg.Add(1)
		go func(stream Stream) {
			defer wg.Done()
			stream.ForEach(consumer)
		}(stream)
	}

	wg.Wait()
}

// Filter is the concurrent equivalent of Stream.Filter.
// See Stream.Filter for further information.
func (cs CStream) Filter(predicate Predicate) CStream {
	cstream := make([]Stream, len(cs.streams))

	var wg sync.WaitGroup

	for idx, stream := range cs.streams {
		wg.Add(1)
		go func(idx int, stream Stream) {
			defer wg.Done()
			cstream[idx] = stream.Filter(predicate)
		}(idx, stream)
	}

	wg.Wait()

	return NewCStream(nil).AddStreams(cstream)
}
