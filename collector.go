package fuego

import "log"

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

func GroupingBy(classifier Function, collector Collector) Collector {
	// resultMap := EntryMap{}

	supplier := collector.supplier

	accumulator := func(supplier Entry, entry Entry) Entry {
		log.Printf("********** DEBUG - GroupingBy - %+v = %+v\n", entry, classifier(entry))
		k := classifier(entry)
		container := collector.accumulator(supplier, entry)
		log.Printf("DEBUG - GroupingBy - k=%+v - container=%+v\n", k, container)
		return container
	}

	finisher := collector.finisher

	return NewCollector(supplier, accumulator, finisher)
}

func Mapping(mapper Function, collector Collector) Collector {
	supplier := collector.supplier

	accumulator := func(supplier Entry, entry Entry) Entry {
		log.Printf("DEBUG - Mapping - %+v = %+v\n", entry, mapper(entry))
		return collector.accumulator(supplier, mapper(entry))
	}

	finisher := collector.finisher

	return NewCollector(supplier, accumulator, finisher)
}

// entryNone is an Entry for an empty struct.
type entryNone struct{}

// Hash returns a hash for 'i'.
func (i entryNone) Hash() uint32 {
	return 0
}

// Equal returns true if 'e' and 'i' are equal.
func (i entryNone) Equal(e Entry) bool {
	return e == entryNone{}
}

func Filtering(predicate Predicate, collector Collector) Collector {
	supplier := collector.supplier

	accumulator := func(supplier Entry, entry Entry) Entry {
		log.Printf("DEBUG - Filtering - %+v = %+v\n", entry, predicate(entry))
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
		log.Printf("DEBUG - ToEntryMap - %+v\n", entry)
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
		log.Printf("DEBUG - ToEntrySlice - %+v\n", entry)
		return supplier.(EntrySlice).Append(entry)
	}

	finisher := func(e Entry) Entry {
		return e
	}

	return NewCollector(supplier, accumulator, finisher)
}
