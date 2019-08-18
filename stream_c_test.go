package fuego

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStream_Concurrent_panicIfInvalidConcurrency(t *testing.T) {
	type fields struct {
		stream           chan Entry
		concurrencyLevel int
	}
	type args struct {
		n int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantPanic bool
	}{
		{
			name: "Should panic if concurrency level < 0",
			fields: fields{
				stream: make(chan Entry, 10),
			},
			args:      args{-1},
			wantPanic: true,
		},
		{
			name: "Should accept concurrency level 0 (the minimum - it is mapped to 1)",
			fields: fields{
				stream: make(chan Entry, 10),
			},
			args:      args{0},
			wantPanic: false,
		},
		{
			name: "Should accept high concurrency level",
			fields: fields{
				stream: make(chan Entry, 100),
			},
			args:      args{50},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Stream{
				stream:           tt.fields.stream,
				concurrencyLevel: tt.fields.concurrencyLevel,
			}
			if tt.wantPanic {
				assert.PanicsWithValue(t, PanicInvalidConcurrencyLevel, func() { s.Concurrent(tt.args.n) })
			} else {
				assert.NotPanics(t, func() { s.Concurrent(tt.args.n) })
			}
		})
	}
}

func TestStream_ForEachC(t *testing.T) {
	var consumeCount int64
	produceCount := int64(5000)
	concurrencyLevel := 1000

	channel := make(chan Entry, concurrencyLevel)
	go func() {
		defer close(channel)
		for i := int64(1); i <= produceCount; i++ {
			channel <- EntryInt(i)
		}
	}()

	s := Stream{
		stream:           channel,
		concurrencyLevel: concurrencyLevel,
	}
	s.ForEachC(func(e Entry) {
		atomic.AddInt64(&consumeCount, 1)
		time.Sleep(10 * time.Microsecond)
	})

	assert.Equal(t, produceCount, atomic.LoadInt64(&consumeCount))
}

func TestStream_concurrentDo_PanicsWhenInvalidConcurrency(t *testing.T) {
	s := Stream{
		stream:           nil,
		concurrencyLevel: -1,
	}

	assert.PanicsWithValue(t, PanicInvalidConcurrencyLevel, func() { s.concurrentDo(func() {}) })
}

func TestStream_concurrentDo(t *testing.T) {
	var concurrencyCount, consumeCount int64
	produceCount := int64(5000)
	concurrencyLevel := int64(1000)

	channel := make(chan Entry, concurrencyLevel)
	go func() {
		defer close(channel)
		for i := int64(1); i <= produceCount; i++ {
			channel <- EntryInt(i)
		}
	}()

	f := func() {
		atomic.AddInt64(&concurrencyCount, 1)
		for range channel {
			atomic.AddInt64(&consumeCount, 1)
		}
		time.Sleep(10 * time.Microsecond)
	}

	s := Stream{
		stream:           channel,
		concurrencyLevel: int(concurrencyLevel),
	}

	s.concurrentDo(f)
	assert.Equal(t, produceCount, atomic.LoadInt64(&consumeCount))
	assert.Equal(t, concurrencyLevel, atomic.LoadInt64(&concurrencyCount))
}
