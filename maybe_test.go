package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaybe_Of_Nil(t *testing.T) {
	none := Maybe_Of(nil)
	assert.True(t, none.IsEmpty())
	assert.Equal(t, Maybe_None(), none)
}

func TestMaybe_Of_Value(t *testing.T) {
	some := Maybe_Of(EntryInt(997))
	assert.False(t, some.IsEmpty())
	assert.Equal(t, Maybe_Some(EntryInt(997)), some)
}

func TestMaybe_None_IsEmpty(t *testing.T) {
	none := Maybe_None()
	assert.True(t, none.IsEmpty())
}

func TestMaybe_None_Get(t *testing.T) {
	none := Maybe_None()
	assert.PanicsWithValue(t, PanicNoSuchElement, func() { none.Get() })
}

func TestMaybe_None_GetOrElse(t *testing.T) {
	none := Maybe_None()
	other := EntryInt(333)
	assert.EqualValues(t, other, none.GetOrElse(other))
}

func TestMaybe_None_OrElse(t *testing.T) {
	none := Maybe_None()
	other := Maybe_Some(EntryInt(333))
	assert.Equal(t, other, none.OrElse(other))
}

func TestMaybe_Some_IsEmpty(t *testing.T) {
	e := EntryInt(997)
	some := Maybe_Some(e)
	assert.False(t, some.IsEmpty())
}

func TestMaybe_Some_Get(t *testing.T) {
	e := EntryInt(997)
	some := Maybe_Some(e)
	assert.EqualValues(t, e, some.Get())
}

func TestMaybe_Some_GetOrElse(t *testing.T) {
	e := EntryInt(997)
	some := Maybe_Some(e)
	other := EntryInt(333)
	assert.EqualValues(t, e, some.GetOrElse(other))
}

func TestMaybe_Some_OrElse(t *testing.T) {
	e := EntryInt(997)
	some := Maybe_Some(e)
	other := Maybe_Some(EntryInt(333))
	assert.Equal(t, some, some.OrElse(other))
}

func TestMaybe_Some_IsEmptyWithNil(t *testing.T) {
	e := Entry(nil)
	some := Maybe_Some(e)
	assert.False(t, some.IsEmpty())
}

func TestMaybe_Some_GetWithNil(t *testing.T) {
	e := Entry(nil)
	some := Maybe_Some(e)
	assert.Nil(t, some.Get())
}

func TestMaybe_Some_GetOrElseWithNil(t *testing.T) {
	e := Entry(nil)
	some := Maybe_Some(e)
	other := EntryInt(333)
	assert.EqualValues(t, e, some.GetOrElse(other))
}

func TestMaybe_Some_OrElseWithNil(t *testing.T) {
	e := Entry(nil)
	some := Maybe_Some(e)
	other := Maybe_Some(EntryInt(333))
	assert.Equal(t, some, some.OrElse(other))
}
