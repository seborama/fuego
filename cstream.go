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
	ps := CStream{}
	for _, channel := range channels {
		ps.streams = append(ps.streams, NewStream(channel))
	}
	return ps
}

// AddStreamFromChannels adds Streams derived from the supplied channels to this CStream.
func (ps CStream) AddStreamFromChannels(channels []chan Entry) CStream {
	for _, channel := range channels {
		ps.streams = append(ps.streams, NewStream(channel))
	}
	// TODO rather than modify in-place, make a copy of ps return the copy? Note: is it worth it because channels will copy by reference anyway?
	return ps
}

// AddStreamFromSlices adds Streams derived from the supplied slices to this CStream.
func (ps CStream) AddStreamFromSlices(slices []EntrySlice, bufsize int) CStream {
	for _, slice := range slices {
		ps.streams = append(ps.streams, NewStreamFromSlice(slice, bufsize))
	}
	// TODO rather than modify in-place, make a copy of ps return the copy? Note: is it worth it because channels will copy by reference anyway?
	return ps
}

// AddStreams adds Streams to this CStream.
func (ps CStream) AddStreams(streams []Stream) CStream {
	for _, stream := range streams {
		ps.streams = append(ps.streams, stream)
	}
	// TODO rather than modify in-place, make a copy of ps and return the copy? Note: is it worth it because channels will copy by reference anyway?
	return ps
}

// ForEach is the concurrent equivalent of Stream.ForEach.
// See Stream.ForEach for further information.
func (ps CStream) ForEach(consumer Consumer) {
	var wg sync.WaitGroup

	for _, stream := range ps.streams {
		wg.Add(1)
		go func(stream Stream) {
			defer wg.Done()
			stream.ForEach(consumer)
		}(stream)
	}

	wg.Wait()
}
