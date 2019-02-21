package fuego

// NOTICE:
// The code in this file was inspired by Java Collectors,
// Vavr and somewhat Scala.
// Many thanks to those great projects!

// TODO: evolve towards a go-style decorator pattern?
type Collector struct {
	supplier    Supplier
	accumulator BiFunction // TODO: this should be a BiConsumer but is it against pure functional design?
	// combiner BiFunction/BiOperator // this is for joining paralle collectors
	finisher Function
}

func NewCollector(supplier Supplier, accumulator BiFunction, finisher Function) Collector {
	return Collector{
		supplier:    supplier,
		accumulator: accumulator,
		finisher:    finisher,
	}
}

// type MutationCollector func(Function, Collector) Collector
// type Collecting func(MutationCollector) MutationCollector

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

	finisher := downstream.finisher

	return NewCollector(supplier, accumulator, finisher)
}

func Mapping(mapper Function, collector Collector) Collector {
	supplier := collector.supplier

	accumulator := func(supplier Entry, entry Entry) Entry {
		return collector.accumulator(supplier, mapper(entry))
	}

	finisher := collector.finisher

	return NewCollector(supplier, accumulator, finisher)
}

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

// func ToEntryMap() Collector {
// 	var supplier = func() Entry { // TODO: use chan Entry instead?
// 		return EntryMap{}
// 	}

// 	accumulator := func(supplier, entry Entry) Entry {
// 		log.Printf("DEBUG - ToEntryMap - %+v\n", entry)
// 		return supplier.(EntryMap).Append(entry.(Tuple2))
// 	}

// 	finisher := IdentityFinisher

// 	return NewCollector(supplier, accumulator, finisher)
// }

func ToEntrySlice() Collector {
	var supplier = func() Entry { // TODO: use chan Entry instead?
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
