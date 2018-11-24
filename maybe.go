package fuego

// A Maybe is a maybe monad.
type Maybe struct {
	value   Entry
	isEmpty bool
}

// A MaybeNone is a Maybe that does not have a value.
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

func (m Maybe) IsEmpty() bool { return m.isEmpty }

// Filter    \
// Map        > can we use Stream for those?
// FlatMap   /
// Iterator or Stream?

// Get the value of this Maybe or panic if none exists.
func (m Maybe) Get() Entry {
	if m.IsEmpty() {
		panic(PanicNoSuchElement)
	}
	return m.value
}

// GetOrElse gets the value of this Maybe or the given Entry if none exists.
func (m Maybe) GetOrElse(e Entry) Entry {
	if m.IsEmpty() {
		return e
	}
	return m.value
}

// GetOrElse returns this Maybe or the given Maybe if this Maybe is empty.
func (m Maybe) OrElse(other Maybe) Maybe {
	if m.IsEmpty() {
		return other
	}
	return m
}
