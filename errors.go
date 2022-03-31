package fuego

// PanicMissingChannel signifies that the Stream is missing a channel.
const PanicMissingChannel = "stream creation requires a channel"

// PanicNoSuchElement signifies that the requested element is not present.
// Examples: when the Stream is empty, or when an Optional does not have a value.
const PanicNoSuchElement = "no such element"

// PanicCollectorMissingSupplier signifies that the Supplier of a Collector was not provided.
const PanicCollectorMissingSupplier = "collector missing supplier"

// PanicCollectorMissingAccumulator signifies that the accumulator of a Collector was not provided.
const PanicCollectorMissingAccumulator = "collector missing accumulator"

// PanicCollectorMissingFinisher signifies that the Finisher of a Collector was not provided.
const PanicCollectorMissingFinisher = "collector missing finisher"

// PanicNilNotPermitted signifies that the `nil` value is not allowed in the context.
const PanicNilNotPermitted = "nil not permitted"

// PanicDuplicateKey signifies that an attempt was made to duplicate a key in a container (such as a map).
const PanicDuplicateKey = "duplicate key"
