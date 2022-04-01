package fuego

// Optional is a container object which may or may not contain a value (NO: nil is considered a non-value).
// IMPORTANT NOTE:
// Currently, Go 1.18 does not permit nil generic types.
// See: https://github.com/golang/go/issues/22729
//
// See IsPresent().
//
// Additional methods that depend on the presence or absence of a contained value are provided,
// such as OrElse() (returns a default value if no value is present) and IfPresent()
// (performs an action if a value is present).
//
// This is a value-based class; programmers should treat instances that are equal as interchangeable.
type Optional[T any] struct {
	value   T
	present bool
}

// IsPresent return true if a value is present.
// If no value is present, the object is considered empty and IsPresent() returns false.
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// Filter returns an Optional describing the value if a value is present, and the value matches the
// given predicate, otherwise returns an empty Optional.
func (o Optional[T]) Filter(p Predicate[T]) Optional[T] {
	if o.present && p(o.value) {
		return o
	}

	return OptionalEmpty[T]()
}

// IfPresent performs the given action with the value, if a value is present,
// otherwise performs the given empty-based action.
func (o Optional[T]) IfPresent(c Consumer[T]) {
	if o.present {
		c(o.value)
	}
}

// Get returns the value if present, otherwise panics.
func (o Optional[T]) Get() T {
	if o.present {
		return o.value
	}

	panic(PanicNoSuchElement)
}

// Or returns an Optional describing the value if present, otherwise returns an Optional produced by the
// supplying function.
func (o Optional[T]) Or(s Supplier[Optional[T]]) Optional[T] {
	if o.present {
		return o
	}

	return s()
}

// OrElse returns the value if present, otherwise returns other.
func (o Optional[T]) OrElse(other T) T {
	if o.present {
		return o.value
	}

	return other
}

// OrElseGet returns the value if present, otherwise returns the result produced by the supplying function.
func (o Optional[T]) OrElseGet(other Supplier[T]) T {
	if o.present {
		return o.value
	}

	return other()
}

// FlatMap returns the result of applying the given Optional-bearing mapping function to the value if
// present, otherwise returns an empty Optional.
func (o Optional[T]) FlatMap(f Function[T, Optional[Any]]) Optional[Any] {
	if o.present {
		return f(o.value)
	}

	return OptionalEmpty[Any]()
}

// Map returns an Optional describing (as if by ofNullable(T)) the result of applying the given mapping
// function to the value if present, otherwise returns an empty Optional.
// If the mapping function returns a nil result then this method returns an empty Optional.
func (o Optional[T]) Map(f Function[T, Any]) Optional[Any] {
	if o.present {
		val := f(o.value)
		if val == nil {
			return OptionalEmpty[Any]()
		}

		return OptionalOf(val)
	}

	return OptionalEmpty[Any]()
}

// OptionalOf returns an Optional describing the given (NO: non-nil) value.
// IMPORTANT NOTE:
// Currently, Go 1.18 does not permit nil generic types.
// See: https://github.com/golang/go/issues/22729
func OptionalOf[T any](val T) Optional[T] {
	return Optional[T]{
		value:   val,
		present: true,
	}
}

// OptionalEmpty returns an empty Optional instance. No value is present for this Optional.
func OptionalEmpty[T any]() Optional[T] {
	return Optional[T]{
		present: false,
	}
}
