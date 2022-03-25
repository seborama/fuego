package fuego

// PanicMissingChannel signifies that the Stream is missing a channel.
const PanicMissingChannel = "stream creation requires a channel"

// PanicNoSuchElement signifies that the requested element is not present.
// This is usually when the Stream is empty.
const PanicNoSuchElement = "no such element"
