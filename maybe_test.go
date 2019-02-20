package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaybeOf_Nil(t *testing.T) {
	none := MaybeOf(nil)
	assert.True(t, none.IsEmpty())
	assert.Exactly(t, MaybeNone(), none)
}

// entryNil is an Entry for an empty struct.
type entryNil struct{}

// Hash returns a hash for 'i'.
func (i entryNil) Hash() uint32 {
	return 0
}

// Equal returns true if 'e' and 'i' are equal.
func (i entryNil) Equal(e Entry) bool {
	return e == nil
}

func TestMaybeOf_EntryNil(t *testing.T) {
	none := MaybeOf(entryNil{})
	assert.True(t, none.IsEmpty())
	assert.Exactly(t, MaybeNone(), none)
}

func TestMaybeOf_Value(t *testing.T) {
	some := MaybeOf(EntryInt(997))
	assert.False(t, some.IsEmpty())
	assert.Exactly(t, MaybeSome(EntryInt(997)), some)
}

func TestMaybeNone_IsEmpty(t *testing.T) {
	none := MaybeNone()
	assert.True(t, none.IsEmpty())
}

func TestMaybeNone_Get(t *testing.T) {
	none := MaybeNone()
	assert.PanicsWithValue(t, PanicNoSuchElement, func() { none.Get() })
}

func TestMaybeNone_GetOrElse(t *testing.T) {
	none := MaybeNone()
	other := EntryInt(333)
	assert.Exactly(t, other, none.GetOrElse(other))
}

func TestMaybeNone_OrElse(t *testing.T) {
	none := MaybeNone()
	other := MaybeSome(EntryInt(333))
	assert.Exactly(t, other, none.OrElse(other))
}

func TestMaybeSome_IsEmpty(t *testing.T) {
	e := EntryInt(997)
	some := MaybeSome(e)
	assert.False(t, some.IsEmpty())
}

func TestMaybeSome_Get(t *testing.T) {
	e := EntryInt(997)
	some := MaybeSome(e)
	assert.Exactly(t, e, some.Get())
}

func TestMaybeSome_GetOrElse(t *testing.T) {
	e := EntryInt(997)
	some := MaybeSome(e)
	other := EntryInt(333)
	assert.Exactly(t, e, some.GetOrElse(other))
}

func TestMaybeSome_OrElse(t *testing.T) {
	e := EntryInt(997)
	some := MaybeSome(e)
	other := MaybeSome(EntryInt(333))
	assert.Exactly(t, some, some.OrElse(other))
}

func TestMaybeSome_IsEmptyWithNil(t *testing.T) {
	e := Entry(nil)
	some := MaybeSome(e)
	assert.False(t, some.IsEmpty())
}

func TestMaybeSome_GetWithNil(t *testing.T) {
	e := Entry(nil)
	some := MaybeSome(e)
	assert.Nil(t, some.Get())
}

func TestMaybeSome_GetOrElseWithNil(t *testing.T) {
	e := Entry(nil)
	some := MaybeSome(e)
	other := EntryInt(333)
	assert.Exactly(t, e, some.GetOrElse(other))
}

func TestMaybeSome_OrElseWithNil(t *testing.T) {
	e := Entry(nil)
	some := MaybeSome(e)
	other := MaybeSome(EntryInt(333))
	assert.Exactly(t, some, some.OrElse(other))
}
