package fuego

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func functionTimesTwo() Function {
	return func(i Entry) Entry {
		num := int(i.(EntryInt))
		return EntryInt(2 * num)
	}
}

func entryIntEqualsTo(number EntryInt) Function {
	return func(subject Entry) Entry {
		subjectEntryInt, ok := subject.(EntryInt)
		if !ok {
			return EntryBool(false)
		}

		return EntryBool(number.Equal(subjectEntryInt))
	}
}

func TestStream_Map(t *testing.T) {
	type fields struct {
		stream chan Entry
	}
	type args struct {
		mapper Function
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Entry
	}{
		{
			name:   "Should return a Stream of nil",
			fields: fields{stream: nil},
			args: args{
				mapper: functionTimesTwo(),
			},
			want: nil,
		},
		{
			name: "Should return a Stream of doubled integers",
			fields: fields{
				stream: func() chan Entry {
					c := make(chan Entry, 1e3)
					c <- EntryInt(1)
					c <- EntryInt(3)
					c <- EntryInt(2)
					close(c)
					return c
				}()},
			args: args{
				mapper: functionTimesTwo(),
			},
			want: []Entry{
				EntryInt(2),
				EntryInt(6),
				EntryInt(4)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Stream{
				stream: tt.fields.stream,
			}

			var got []Entry
			if gotStream := s.Map(tt.args.mapper).stream; gotStream != nil {
				got = []Entry{}
				for val := range gotStream {
					got = append(got, val)
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stream.Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStream_Filter(t *testing.T) {
	type fields struct {
		stream chan Entry
	}
	type args struct {
		predicate Predicate
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Entry
	}{
		{
			name:   "Should return a Stream of nil",
			fields: fields{stream: nil},
			args: args{
				predicate: intGreaterThanPredicate(5),
			},
			want: nil,
		},
		{
			name: "Should give produce filtered values as per predicate",
			fields: fields{
				stream: func() chan Entry {
					c := make(chan Entry, 1e3)
					c <- EntryInt(17)
					c <- EntryInt(8)
					c <- EntryInt(2)
					close(c)
					return c
				}()},
			args: args{
				predicate: intGreaterThanPredicate(5),
			},
			want: []Entry{
				EntryInt(17),
				EntryInt(8)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Stream{
				stream: tt.fields.stream,
			}

			var got []Entry
			if gotStream := s.Filter(tt.args.predicate).stream; gotStream != nil {
				got = []Entry{}
				for val := range gotStream {
					got = append(got, val)
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stream.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStream_ForEach(t *testing.T) {
	var callCount, total int
	computeSumTotal := func(value Entry) {
		callCount++
		total += int(value.(EntryInt))
	}

	type fields struct {
		stream chan Entry
	}
	type args struct {
		consumer Consumer
	}
	type want struct {
		total, count int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "Should return a Stream of nil",
			fields: fields{stream: nil},
			args: args{
				consumer: computeSumTotal,
			},
			want: want{
				count: 0,
				total: 0,
			},
		},
		{
			name: "Should give produce filtered values as per predicate",
			fields: fields{
				stream: func() chan Entry {
					c := make(chan Entry, 1e3)
					c <- EntryInt(4)
					c <- EntryInt(1)
					c <- EntryInt(3)
					close(c)
					return c
				}()},
			args: args{
				consumer: computeSumTotal,
			},
			want: want{
				count: 3,
				total: 8,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount, total = 0, 0
			s := Stream{
				stream: tt.fields.stream,
			}
			s.ForEach(tt.args.consumer)
			assert.Equal(t, tt.want.count, callCount)
			assert.Equal(t, tt.want.total, total)
		})
	}
}

// func TestStream_LeftReduce(t *testing.T) {
// 	concatenateStringsBiFunc := func(i, j Entry) Entry {
// 		iStr := i.(EntryString)
// 		jStr := j.(EntryString)
// 		return iStr + "-" + jStr
// 	}

// 	type fields struct {
// 		iterator Iterator
// 	}
// 	type args struct {
// 		f2 BiFunction
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   interface{}
// 	}{
// 		{
// 			name: "Should return nil for a nil Stream",
// 			fields: fields{
// 				iterator: NewSetIterator(nil),
// 			},
// 			args: args{f2: concatenateStringsBiFunc},
// 			want: nil,
// 		},
// 		{
// 			name: "Should return nil for an empty Stream",
// 			fields: fields{
// 				iterator: SetIterator{set: NewHamtSet()},
// 			},
// 			args: args{f2: concatenateStringsBiFunc},
// 			want: nil,
// 		},
// 		{
// 			name: "Should return reduction of Set of single element",
// 			fields: fields{
// 				iterator: NewSetIterator(NewHamtSet().
// 					Insert(EntryString("three"))),
// 			},
// 			args: args{f2: concatenateStringsBiFunc},
// 			want: EntryString("three"),
// 		},
// 		{
// 			name: "Should return reduction of Set",
// 			fields: fields{
// 				iterator: NewSetIterator(NewHamtSet().
// 					Insert(EntryString("four")).
// 					Insert(EntryString("twelve")).
// 					Insert(EntryString("one")).
// 					Insert(EntryString("six")).
// 					Insert(EntryString("three"))),
// 			},
// 			args: args{f2: concatenateStringsBiFunc},
// 			want: EntryString("one-three-twelve-six-four"),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			rp := Stream{
// 				iterator: tt.fields.iterator,
// 			}
// 			if gotReduce := rp.Reduce(tt.args.f2); !assert.Exactly(t, tt.want, gotReduce) {
// 				return
// 			}

// 			if gotLeftReduce := rp.LeftReduce(tt.args.f2); !assert.Exactly(t, tt.want, gotLeftReduce) {
// 				return
// 			}
// 		})
// 	}
// }

// func TestStream_Intersperse(t *testing.T) {
// 	type fields struct {
// 		iterator Iterator
// 	}
// 	type args struct {
// 		e Entry
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   Stream
// 	}{
// 		{
// 			name:   "Should return a Stream of nil for nil iterator",
// 			fields: fields{iterator: nil},
// 			args: args{
// 				e: EntryString(" - "),
// 			},
// 			want: NewStream(
// 				NewSliceIterator([]Entry{})),
// 		},
// 		{
// 			name:   "Should return a Stream of nil for empty Set",
// 			fields: fields{iterator: SetIterator{set: NewHamtSet()}},
// 			args: args{
// 				e: EntryString(" - "),
// 			},
// 			want: NewStream(
// 				NewSliceIterator([]Entry{})),
// 		},
// 		{
// 			name: "Should return the original Set when it has a single value",
// 			fields: fields{
// 				iterator: NewSetIterator(NewHamtSet().
// 					Insert(EntryString("four")))},
// 			args: args{
// 				e: EntryString(" - "),
// 			},
// 			want: NewStream(
// 				NewSliceIterator([]Entry{
// 					EntryString("four")})),
// 		},
// 		{
// 			name: "Should return the Set with given value interspersed",
// 			fields: fields{
// 				iterator: NewSetIterator(NewOrderedSet().
// 					Insert(EntryString("four")).
// 					Insert(EntryString("twelve")).
// 					Insert(EntryString("one")).
// 					Insert(EntryString("six")).
// 					Insert(EntryString("three"))),
// 			},
// 			args: args{
// 				e: EntryString(" - "),
// 			},
// 			want: NewStream(
// 				NewSliceIterator([]Entry{
// 					EntryString("four"),
// 					EntryString(" - "),
// 					EntryString("twelve"),
// 					EntryString(" - "),
// 					EntryString("one"),
// 					EntryString(" - "),
// 					EntryString("six"),
// 					EntryString(" - "),
// 					EntryString("three")})),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			rp := Stream{
// 				iterator: tt.fields.iterator,
// 			}
// 			if got := rp.Intersperse(tt.args.e); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Stream.Intersperse() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestStream_GroupBy(t *testing.T) {
// 	type fields struct {
// 		iterator Iterator
// 	}
// 	type args struct {
// 		classifier Function
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   Map
// 	}{
// 		{
// 			name: "Should return empty map when nil iterator",
// 			fields: fields{
// 				iterator: nil,
// 			},
// 			args: args{
// 				classifier: func(i Entry) Entry {
// 					return i.(EntryInt) & 1
// 				},
// 			},
// 			want: NewOrderedMap(),
// 		},
// 		{
// 			name: "Should return empty map when iterator with nil Set",
// 			fields: fields{
// 				iterator: NewSetIterator(nil),
// 			},
// 			args: args{
// 				classifier: func(i Entry) Entry {
// 					return i.(EntryInt) & 1
// 				},
// 			},
// 			want: NewOrderedMap(),
// 		},
// 		{
// 			name: "Should return empty map when empty Set",
// 			fields: fields{
// 				iterator: NewSetIterator(NewOrderedSet()),
// 			},
// 			args: args{
// 				classifier: func(i Entry) Entry {
// 					return i.(EntryInt) & 1
// 				},
// 			},
// 			want: NewOrderedMap(),
// 		},
// 		{
// 			name: "Should group by odd / even numbers",
// 			fields: fields{
// 				iterator: NewSetIterator(NewOrderedSet().
// 					Insert(EntryInt(1)).
// 					Insert(EntryInt(2)).
// 					Insert(EntryInt(3)).
// 					Insert(EntryInt(4))),
// 			},
// 			args: args{
// 				classifier: func(i Entry) Entry {
// 					return i.(EntryInt) & 1
// 				},
// 			},
// 			want: NewOrderedMap().
// 				Insert(EntryInt(0), NewOrderedSet().
// 					Insert(EntryInt(2)).
// 					Insert(EntryInt(4))).
// 				Insert(EntryInt(1), NewOrderedSet().
// 					Insert(EntryInt(1)).
// 					Insert(EntryInt(3))),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			rp := Stream{
// 				iterator: tt.fields.iterator,
// 			}
// 			got := rp.GroupBy(tt.args.classifier)
// 			elementsMatch(t, tt.want, got)
// 		})
// 	}
// }

// func elementsMatch(t *testing.T, mapA, mapB Map) {
// 	keysA, valuesA := splitKeysValues(mapA)
// 	keysB, valuesB := splitKeysValues(mapB)
// 	assert.ElementsMatch(t, keysA, keysB, "keys differ")
// 	assert.ElementsMatch(t, valuesA, valuesB, "values differ")
// }

// func splitKeysValues(m Map) (keys, values []Entry) {
// 	m.EntrySet().Stream().ForEach(func(e Entry) {
// 		keys = append(keys, e.(MapEntry).K)
// 		e.(MapEntry).V.(OrderedSet).Stream().ForEach(func(e Entry) {
// 			values = append(values, e)
// 		})
// 	})
// 	return keys, values
// }

// func TestStream_GroupBy_IteratorResets(t *testing.T) {
// 	it := NewSetIterator(NewOrderedSet().
// 		Insert(EntryInt(1)).
// 		Insert(EntryInt(2)).
// 		Insert(EntryInt(3)).
// 		Insert(EntryInt(4)))

// 	rp := Stream{iterator: it}

// 	expected := NewOrderedMap().
// 		Insert(EntryInt(1), NewOrderedSet().
// 			Insert(EntryInt(1)).
// 			Insert(EntryInt(3))).
// 		Insert(EntryInt(0), NewOrderedSet().
// 			Insert(EntryInt(2)).
// 			Insert(EntryInt(4)))

// 	res1 := rp.GroupBy(func(i Entry) Entry {
// 		return i.(EntryInt) & 1
// 	})
// 	elementsMatch(t, res1, expected)

// 	res2 := rp.GroupBy(func(i Entry) Entry {
// 		return i.(EntryInt) & 1
// 	})
// 	elementsMatch(t, res2, expected)
// }

func TestNewStream(t *testing.T) {
	type args struct {
		s chan Entry
	}
	emptyChannel := make(chan Entry)
	tests := []struct {
		name string
		args args
		want Stream
	}{
		{
			name: "Should create a Stream with a nil channel",
			args: args{s: nil},
			want: Stream{stream: nil},
		},
		{
			name: "Should create a Stream with a channel",
			args: args{s: emptyChannel},
			want: Stream{stream: emptyChannel},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStream(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStream() = %v, want %v", got, tt.want)
			}
		})
	}
}
