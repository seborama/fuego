package fuego

import "fmt"

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
// Type T: type of input elements to the reduction operation
// Type A: mutable accumulation type of the reduction operation (often hidden as an implementation detail)
// Type R: result type of the reduction operation.
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

// FlatMapping adapts the Entries a Collector accepts to another type by
// applying a flat mapping function which maps input elements to a `Stream`.
func FlatMapping[U, T, A, R any](mapper StreamFunction[T, U], collector Collector[U, A, R]) Collector[T, A, R] {
	supplier := collector.supplier

	accumulator := func(supplierA A, entry T) A {
		container := supplierA
		stream := mapper(entry)

		stream.ForEach(
			func(e U) {
				container = collector.accumulator(container, e)
			})

		return container
	}

	finisher := collector.finisher

	return NewCollector(supplier, accumulator, finisher)
}

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

// Reducing returns a collector that performs a reduction of
// its input elements using the provided BiFunction.
func Reducing[T any](f2 BiFunction[T, T, T]) Collector[T, Optional[T], T] {
	supplier := func() Optional[T] {
		return OptionalEmpty[T]()
	}

	accumulator := func(supplierA Optional[T], entry T) Optional[T] {
		result := entry

		supplierA.IfPresent(func(val T) { result = f2(val, entry) })

		return OptionalOf(result)
	}

	finisher := func(e Optional[T]) T {
		return e.Get()
	}

	return NewCollector(supplier, accumulator, finisher)
}

// ToSlice returns a collector that accumulates the input entries into a Go slice.
// Type T: type of the elements accumulated in the slice.
func ToSlice[T any]() Collector[T, []T, []T] {
	supplier := func() []T { // TODO: use chan A instead with a finisher that converts to []A?
		return []T{}
	}

	accumulator := func(supplier []T, element T) []T {
		return append(supplier, element)
	}

	finisher := IdentityFinisher[[]T]

	return NewCollector(supplier, accumulator, finisher)
}

// ToMap returns a collector that accumulates the input entries into a Go map.
// Type T: type from which the elements are accumulated in the map.
// Type K: type of the keys derived from T.
// Type V: type of the values derived from T.
func ToMap[T any, K comparable, V any](keyMapper Function[T, K], valueMapper Function[T, V]) Collector[T, map[K]V, map[K]V] {
	supplier := func() map[K]V { // TODO: use chan instead with a finisher that converts to map?
		return map[K]V{}
	}

	accumulator := func(supplier map[K]V, element T) map[K]V {
		key := keyMapper(element)
		value := valueMapper(element)

		if _, ok := supplier[key]; !ok {
			supplier[key] = value
			return supplier
		}

		panic(fmt.Sprintf("%s: '%v'", PanicDuplicateKey, key))
	}

	finisher := IdentityFinisher[map[K]V]

	return NewCollector(supplier, accumulator, finisher)
}

// ToMapWithMerge returns a collector that accumulates the input entries into a Go map.
// Key collision strategy is managed by mergeFn.
// Type T: type from which the elements are accumulated in the map.
// Type K: type of the keys derived from T.
// Type V: type of the values derived from T.
func ToMapWithMerge[T any, K comparable, V any](keyMapper Function[T, K], valueMapper Function[T, V], mergeFn BiFunction[V, V, V]) Collector[T, map[K]V, map[K]V] {
	supplier := func() map[K]V { // TODO: use chan instead with a finisher that converts to map?
		return map[K]V{}
	}

	accumulator := func(supplier map[K]V, element T) map[K]V {
		key := keyMapper(element)
		value := valueMapper(element)

		if _, ok := supplier[key]; !ok {
			supplier[key] = value
			return supplier
		}

		supplier[key] = mergeFn(supplier[key], value)

		return supplier
	}

	finisher := IdentityFinisher[map[K]V]

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
