package fuego

// TODO: (see Vavr) FlatMap, Map, GetOrElseThrow, Peek, Collect, When

// PanicNoSuchElement signifies that the requested element is not present.
const PanicNoSuchElement = "no such element"

// A Maybe is a maybe monad.
type Maybe struct {
	value   Entry
	isEmpty bool
}

// MaybeNone is a Maybe that does not have a value.
func MaybeNone() Maybe {
	return Maybe{
		value:   nil,
		isEmpty: true,
	}
}

// MaybeSome creates a new Maybe with the given value.
// Note: MaybeOf(nil) == None() whereas MaybeSome(nil) == MaybeSome(nil).
func MaybeSome(i Entry) Maybe {
	return Maybe{
		value:   i,
		isEmpty: false,
	}
}

// MaybeOf creates a new Maybe with the given value.
// If the value is nil then return None otherwise Some(value).
// Note: MaybeOf(nil) == None() whereas MaybeSome(nil) == MaybeSome(nil).
func MaybeOf(i Entry) Maybe {
	if i == nil || i.Equal(nil) {
		return MaybeNone()
	}
	return MaybeSome(i)
}

// IsEmpty returns true when this Maybe does not have
// a value.
func (m Maybe) IsEmpty() bool { return m.isEmpty }

// Get the value of this Maybe or panic if none exists.
func (m Maybe) Get() Entry {
	if m.IsEmpty() {
		panic(PanicNoSuchElement) // TODO: return MaybeNone??
	}
	return m.value
}

// GetOrElse gets the value of this Maybe or the given Entry if none exists.
func (m Maybe) GetOrElse(elseEntry Entry) Entry {
	if m.IsEmpty() {
		return elseEntry
	}
	return m.value
}

// OrElse returns this Maybe or the given Maybe if this Maybe is empty.
func (m Maybe) OrElse(elseMaybe Maybe) Maybe {
	if m.IsEmpty() {
		return elseMaybe
	}
	return m
}

// Filter returns MaybeSome(value) if this is a MaybeSome
// and the value satisfies the predicate otherwise returns
// MaybeNone.
func (m Maybe) Filter(predicate Predicate) Maybe {
	if m.IsEmpty() || predicate(m.Get()) {
		return m
	}
	return MaybeNone()
}
