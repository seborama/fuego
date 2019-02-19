package fuego

// TODO: evolve towards a go-style decorator pattern?
type Collector struct {
	supplier    Getter     // TODO: rename to "Supplier" with a method of "supply"
	accumulator BiFunction // TODO: this should be a BiConsumer but is it against pure functional design?
	// combiner BiFunction/BiOperator // this is for joining paralle collectors
	finisher Function
}

func NewCollector(supplier Getter, accumulator BiFunction, finisher Function) Collector {
	return Collector{
		supplier:    supplier,
		accumulator: accumulator,
		finisher:    finisher,
	}
}

type CollectorMiddleware func(Function, Collector) Collector

func GroupingBy(classifier Function, collector Collector) Collector {
	supplier := collector.supplier

	accumulator := func(supplier Entry, entry Entry) Entry {
		return collector.accumulator(supplier, classifier(entry))
	}

	finisher := collector.finisher

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

func ToEntryMap() Collector {
	var supplier = func() Entry { // TODO: use chan Entry instead?
		return EntryMap{}
	}

	accumulator := func(supplier, entry Entry) Entry {
		return supplier.(EntryMap).Append(entry.(Tuple2))
	}

	finisher := func(e Entry) Entry {
		return e
	}

	return NewCollector(supplier, accumulator, finisher)
}

func ToEntrySlice() Collector {
	var supplier = func() Entry { // TODO: use chan Entry instead?
		return EntrySlice{}
	}

	accumulator := func(supplier, entry Entry) Entry {
		return supplier.(EntrySlice).Append(entry)
	}

	finisher := func(e Entry) Entry {
		return e
	}

	return NewCollector(supplier, accumulator, finisher)
}
