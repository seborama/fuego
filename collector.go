package fuego

// TODO: evolve towards a go-style decorator pattern?
type Collector struct {
	supplier    Getter
	accumulator BiFunction // TODO: this should be a BiConsumer but is it against pure functional design?
	// combiner BiFunction // this is for joining paralle collectors
	finisher Function
}

func NewCollector(supplier Getter, accumulator BiFunction, finisher Function) Collector {
	return Collector{
		supplier:    supplier,
		accumulator: accumulator,
		finisher:    finisher,
	}
}
