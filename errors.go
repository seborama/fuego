package fuego

// PanicMissingChannel signifies that the Stream is missing a channel.
const PanicMissingChannel = "stream creation requires a channel"

// PanicNoSuchElement signifies that the requested element is not present.
// This is usually when the Stream is empty.
const PanicNoSuchElement = "no such element"

// PanicCollectorMissingSupplier signifies that the Supplier of a Collector was not provided.
const PanicCollectorMissingSupplier = "collector missing supplier"

// PanicCollectorMissingAccumulator signifies that the accumulator of a Collector was not provided.
const PanicCollectorMissingAccumulator = "collector missing accumulator"

// PanicCollectorMissingFinisher signifies that the Finisher of a Collector was not provided.
const PanicCollectorMissingFinisher = "collector missing finisher"
