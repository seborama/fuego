package fuego

import "fmt"

// NOTICE:
// The code in this file was inspired by Java Collectors,
// Vavr and somewhat Scala.
// Many thanks to those great projects!

// TODO: evolve towards a go-style decorator pattern?

// A Collector is a mutable reduction operation,
// optionally transforming the accumulated result.
//
// Collectors can be combined to express complex
// operations in a concise manner.
//
// In other words, a collector allows creating custom
// actions on a Stream. **fuego** comes shipped with a
// number of methods such as `MapToInt`, `Head`, `LastN`,
// `Filter`, etc, and Collectors also provide a few
// additional methods. But what if you need something else?
// And it is straighforward or readable when combining the
// existing methods fuego offers? Enters `Collector`:
// implement you own requirement functionally! Focus on
// *what* needs to be done in your streams (and
// delegate the details of the *how* to the implementation
// of your `Collector`).
//
// It should be noted that the `finisher` function is
// optional (i.e. it may acceptably be `nil`).
//
// Example
//
//  strs := EntrySlice{
//      EntryString("a"),
//      EntryString("bb"),
//      EntryString("cc"),
//      EntryString("ddd"),
//  }
//
//  NewStreamFromSlice(strs, 1e3).
//      Collect(
//          GroupingBy(
//              stringLength,
//              Mapping(
//                  stringToUpper,
//                  ToEntryMap())))
//  // Result: map[1:[A] 2:[BB CC] 3:[DDD]]
type Collector struct {
	supplier    Supplier
	accumulator BiFunction // TODO: this should be a BiConsumer but is it against pure functional design?
	// combiner BiFunction/BiOperator // this is for joining parallel collectors
	finisher Function
}

// NewCollector creates a new Collector.
func NewCollector(supplier Supplier, accumulator BiFunction, finisher Function) Collector {
	return Collector{
		supplier:    supplier,
		accumulator: accumulator,
		finisher:    finisher,
	}
}

// type MutationCollector func(Function, Collector) Collector
// type Collecting func(MutationCollector) MutationCollector

// GroupingBy groups the elements of the downstream Collector
// by classifying them with the provided classifier function.
func GroupingBy(classifier Function, downstream Collector) Collector {
	supplier := func() Entry { return EntryMap{} }

	accumulator := func(supplierA Entry, entry Entry) Entry {
		k := classifier(entry)
		container, ok := supplierA.(EntryMap)[k]
		if !ok {
			container = downstream.supplier()
		}
		container = downstream.accumulator(container, entry)
		supplierA.(EntryMap)[k] = container
		return supplierA
	}

	finisher := func(e Entry) Entry {
		if downstream.finisher == nil ||
			fmt.Sprintf("%p", downstream.finisher) == fmt.Sprintf("%p", IdentityFinisher) /* i.e. downstream.finisher is the IdentityFinisher */ {
			return IdentityFinisher(e)
		}

		m := supplier()
		for k, v := range e.(EntryMap) {
			m.(EntryMap)[k] = downstream.finisher(v)
		}
		return m
	}

	return NewCollector(supplier, accumulator, finisher)
}

// Mapping adapts the Entries a Collector accepts to another type.
func Mapping(mapper Function, collector Collector) Collector {
	supplier := collector.supplier

	accumulator := func(supplier Entry, entry Entry) Entry {
		return collector.accumulator(supplier, mapper(entry))
	}

	finisher := collector.finisher

	return NewCollector(supplier, accumulator, finisher)
}

// FlatMapping adapts the Entries a Collector accepts to another type
// by applying a flat mapping function which maps input elements to a
// `Stream`.
func FlatMapping(mapper StreamFunction, collector Collector) Collector {
	supplier := collector.supplier

	accumulator := func(supplierA Entry, entry Entry) Entry {
		container := supplierA
		stream := mapper(entry)
		stream.ForEach(
			func(e Entry) {
				container = collector.accumulator(container, e)
			})
		return container
	}

	finisher := collector.finisher

	return NewCollector(supplier, accumulator, finisher)
}

// Filtering adapts the Entries a Collector accepts to a subset
// that satisfies the given predicate.
func Filtering(predicate Predicate, collector Collector) Collector {
	supplier := collector.supplier

	accumulator := func(supplier Entry, entry Entry) Entry {
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
func Reducing(f2 BiFunction) Collector {
	supplier := func() Entry { // TODO: use chan Entry instead with a finisher that converts to slice?
		return Tuple2{E1: EntryBool(false), E2: nil}
	}

	accumulator := func(supplierA Entry, entry Entry) Entry {
		present := supplierA.(Tuple2).E1.(EntryBool)
		result := supplierA.(Tuple2).E2

		if present {
			result = f2(result, entry)
		} else {
			present = true
			result = entry
		}
		return Tuple2{E1: present, E2: result}
	}

	finisher := func(e Entry) Entry {
		return e.(Tuple2).E2
	}

	return NewCollector(supplier, accumulator, finisher)
}

// func ToEntryMap() Collector {
// 	var supplier = func() Entry { // TODO: use chan Entry instead with a finisher that converts to slice?
// 		return EntryMap{}
// 	}

// 	accumulator := func(supplier, entry Entry) Entry {
// 		log.Printf("DEBUG - ToEntryMap - %+v\n", entry)
// 		return supplier.(EntryMap).Append(entry.(Tuple2))
// 	}

// 	finisher := IdentityFinisher

// 	return NewCollector(supplier, accumulator, finisher)
// }

// ToEntrySlice returns a collector that accumulates the input
// entries into an EntrySlice.
func ToEntrySlice() Collector {
	supplier := func() Entry { // TODO: use chan Entry instead with a finisher that converts to slice?
		return EntrySlice{}
	}

	accumulator := func(supplier, entry Entry) Entry {
		return supplier.(EntrySlice).Append(entry)
	}

	finisher := IdentityFinisher

	return NewCollector(supplier, accumulator, finisher)
}

// IdentityFinisher is a basic finisher that returns the
// original value passed to it, unmodified.
func IdentityFinisher(e Entry) Entry {
	return e
}
