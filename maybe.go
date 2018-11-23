package fuego

// A Maybe is a maybe monad
type Maybe struct {
	value   Entry
	isEmpty bool
}

func Maybe_None() Maybe {
	return Maybe{
		value:   nil,
		isEmpty: true,
	}
}

// Maybe_Some creates a new Maybe with the given value.
// Note: Maybe_Of(nil) == None() whereas Maybe_Some(nil) == Maybe_Some(nil).
func Maybe_Some(i Entry) Maybe {
	return Maybe{
		value:   i,
		isEmpty: false,
	}
}

// Maybe_Of creates a new Maybe with the given value.
// If the value is nil then return None otherwise Some(value).
// Note: Maybe_Of(nil) == None() whereas Maybe_Some(nil) == Maybe_Some(nil).
func Maybe_Of(i Entry) Maybe {
	if i == nil || i.Equal(nil) {
		return Maybe_None()
	}
	return Maybe_Some(i)
}

func (m Maybe) IsEmpty() bool { return m.isEmpty }

// Filter    \
// Map        > can we use Stream for those?
// FlatMap   /
// Iterator or Stream?

func (m Maybe) Get() Entry {
	if m.IsEmpty() {
		panic(PanicNoSuchElement)
	}
	return m.value
}

func (m Maybe) GetOrElse(e Entry) Entry {
	if m.IsEmpty() {
		return e
	}
	return m.value
}

func (m Maybe) OrElse(other Maybe) Maybe {
	if m.IsEmpty() {
		return other
	}
	return m
}
