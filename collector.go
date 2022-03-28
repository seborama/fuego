package fuego

// NOTICE:
// The code in this file was inspired by Java Collectors,
// Vavr and somewhat Scala.
// Many thanks to those great projects!

// TODO: evolve towards a go-style decorator pattern?

// A Collector is a mutable reduction operation, optionally transforming the accumulated result.
//
// Collectors can be combined to express complex operations in a concise manner.
//
// In other words, a collector allows creating custom actions on a Stream. **fuego** comes shipped with a
// number of methods such as `MapToInt`, `Head`, `LastN`, `Filter`, etc, and Collectors also provide a
// few additional methods. But what if you need something else?
// And it is straightforward or readable when combining the existing methods fuego offers?
// Enters `Collector`:
// implement you own requirement functionally! Focus on *what* needs to be done in your streams (and
// delegate the details of the *how* to the implementation of your `Collector`).
//
// Example:
//
//  strs := []string{
//      "a",
//      "bb",
//      "cc",
//      "ddd",
//  }
//
//  NewStreamFromSlice[string](strs, 1e3).
//      Collect(
//          GroupingBy(
//              stringLength,
//              Mapping(
//                  stringToUpper,
//                  ToEntryMap())))
//  // Result: map[1:[A] 2:[BB CC] 3:[DDD]]
//
// Type T: type of input elements to the reduction operation
// Type A: mutable accumulation type of the reduction operation (often hidden as an implementation detail)
// Type R: result type of the reduction operation
type Collector[T, A, R any] struct {
	supplier    Supplier[A]
	accumulator BiFunction[A, T, A]
	// combiner BiFunction[A]/BinaryOperator[A] // this is for joining parallel collectors
	finisher Function[A, R]
}

// NewCollector creates a new Collector.
func NewCollector[T, A, R any](supplier Supplier[A], accumulator BiFunction[A, T, A], finisher Function[A, R]) Collector[T, A, R] {
	if supplier == nil {
		panic(PanicCollectorMissingSupplier)
	}

	if accumulator == nil {
		panic(PanicCollectorMissingAccumulator)
	}

	if finisher == nil {
		panic(PanicCollectorMissingFinisher)
	}

	return Collector[T, A, R]{
		supplier:    supplier,
		accumulator: accumulator,
		finisher:    finisher,
	}
}

// type MutationCollector func(Function, Collector) Collector
// type Collecting func(MutationCollector) MutationCollector

// GroupingBy groups the elements of the downstream Collector
// by classifying them with the provided classifier function.
//
// Type T: the type of the input elements
// Type K: the type of the keys
// Type A: the intermediate accumulation type of the downstream collector
// Type D: the result type of the downstream reduction
//
// classifier: a classifier function mapping input elements to keys
// downstream: a Collector implementing the downstream reduction
func GroupingBy[T any, K comparable, A, D any](classifier Function[T, K], downstream Collector[T, A, D]) Collector[T, map[K]A, map[K]D] {
	supplier := func() map[K]A { return map[K]A{} }

	accumulator := func(supply map[K]A, element T) map[K]A {
		key := classifier(element)
		container, ok := supply[key]
		if !ok {
			container = downstream.supplier()
		}
		container = downstream.accumulator(container, element)
		supply[key] = container

		return supply
	}

	finisher := func(e map[K]A) map[K]D {
		if downstream.finisher == nil {
			panic(PanicCollectorMissingFinisher)
		}

		m := map[K]D{}
		for k, v := range e {
			m[k] = downstream.finisher(v)
		}
		return m
	}

	return NewCollector(supplier, accumulator, finisher)
}

// Mapping adapts a Collector with elements of type U to a collector with elements of type T.
func Mapping[T, U, A, R any](mapper Function[T, U], downstream Collector[U, A, R]) Collector[T, A, R] {
	supplier := downstream.supplier

	accumulator := func(supplier A, entry T) A {
		return downstream.accumulator(supplier, mapper(entry))
	}

	finisher := downstream.finisher

	return NewCollector(supplier, accumulator, finisher)
}

// // FlatMapping adapts the Entries a Collector accepts to another type
// // by applying a flat mapping function which maps input elements to a
// // `Stream`.
// func FlatMapping(mapper StreamFunction, collector Collector) Collector {
// 	supplier := collector.supplier

// 	accumulator := func(supplierA Entry, entry Entry) Entry {
// 		container := supplierA
// 		stream := mapper(entry)
// 		stream.ForEach(
// 			func(e Entry) {
// 				container = collector.accumulator(container, e)
// 			})
// 		return container
// 	}

// 	finisher := collector.finisher

// 	return NewCollector(supplier, accumulator, finisher)
// }

// Filtering filters the entries a Collector accepts to a subset that satisfy the given predicate.
func Filtering[T, A, R any](predicate Predicate[T], collector Collector[T, A, R]) Collector[T, A, R] {
	supplier := collector.supplier

	accumulator := func(supplier A, entry T) A {
		if predicate(entry) {
			return collector.accumulator(supplier, entry)
		}
		return supplier
	}

	finisher := collector.finisher

	return NewCollector(supplier, accumulator, finisher)
}

// // Reducing returns a collector that performs a reduction of
// // its input elements using the provided BiFunction.
// func Reducing(f2 BiFunction) Collector {
// 	supplier := func() Entry { // TODO: use chan Entry instead with a finisher that converts to slice?
// 		return Tuple2{E1: EntryBool(false), E2: nil}
// 	}

// 	accumulator := func(supplierA Entry, entry Entry) Entry {
// 		present := supplierA.(Tuple2).E1.(EntryBool)
// 		result := supplierA.(Tuple2).E2

// 		if present {
// 			result = f2(result, entry)
// 		} else {
// 			present = true
// 			result = entry
// 		}
// 		return Tuple2{E1: present, E2: result}
// 	}

// 	finisher := func(e Entry) Entry {
// 		return e.(Tuple2).E2
// 	}

// 	return NewCollector(supplier, accumulator, finisher)
// }

// // ToEntryMap returns a collector that accumulates the input
// // entries into an EntryMap.
// func ToEntryMap(keyMapper, valueMapper Function) Collector {
// 	supplier := func() Entry { // TODO: use chan Entry instead with a finisher that converts to map?
// 		return EntryMap{}
// 	}

// 	accumulator := func(supplier, entry Entry) Entry {
// 		key := keyMapper(entry)
// 		value := valueMapper(entry)
// 		return supplier.(EntryMap).Merge(key, value, func(v1, v2 Entry) Entry { panic(fmt.Sprintf("%s: '%v'", PanicDuplicateKey, key)) })
// 	}

// 	finisher := IdentityFinisher

// 	return NewCollector(supplier, accumulator, finisher)
// }

// // ToEntryMapWithKeyMerge returns a collector that accumulates the input
// // entries into an EntryMap. Duplicate keys are managed by mergeFunction.
// // See EntryMap.Merge() for details about the mergeFunction.
// func ToEntryMapWithKeyMerge(keyMapper, valueMapper Function, mergeFunction BiFunction) Collector {
// 	supplier := func() Entry { // TODO: use chan Entry instead with a finisher that converts to EntryMap?
// 		return EntryMap{}
// 	}

// 	accumulator := func(supplier, entry Entry) Entry {
// 		key := keyMapper(entry)
// 		value := valueMapper(entry)
// 		return supplier.(EntryMap).Merge(key, value, mergeFunction)
// 	}

// 	finisher := IdentityFinisher

// 	return NewCollector(supplier, accumulator, finisher)
// }

// ToSlice returns a collector that accumulates the input entries into a Go slice.
func ToSlice[T any]() Collector[T, []T, []T] {
	supplier := func() []T { // TODO: use chan A instead with a finisher that converts to []A]?
		return []T{}
	}

	accumulator := func(supplier []T, element T) []T {
		return append(supplier, element)
	}

	finisher := IdentityFinisher[[]T]

	return NewCollector(supplier, accumulator, finisher)
}

// IdentityFinisher is a basic finisher that returns the
// original value passed to it, unmodified.
func IdentityFinisher[T any](t T) T {
	return t
}

// Collect reduces and optionally mutates the stream with the supplied Collector.
//
// This is a continuous terminal operation and hence expects
// the producer to close the stream in order to complete.
func Collect[T, A, R any](s Stream[T], c Collector[T, A, R]) R {
	if s.stream == nil {
		panic(PanicMissingChannel)
	}

	result := c.supplier()
	for e := range s.stream {
		result = c.accumulator(result, e)
	}

	finishedResult := c.finisher(result)

	return finishedResult
}
