package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptional_OptionalOf_IsPresent_True(t *testing.T) {
	o := OptionalOf(123)
	assert.True(t, o.IsPresent())
}

func TestOptional_OptionalOf_ZeroValue_IsPresent_True(t *testing.T) {
	o := OptionalOf(0)
	assert.True(t, o.IsPresent())
}

func TestOptional_OptionalOf_Nil_IsPresent_False(t *testing.T) {
	o := OptionalOf[*string](nil)
	assert.False(t, o.IsPresent())
}

func TestOptional_OptionalEmpty_IsPresent_False(t *testing.T) {
	o := OptionalEmpty[int]()
	assert.False(t, o.IsPresent())
}

func TestOptional_Filter_Present(t *testing.T) {
	o := OptionalOf(123).Filter(func(i int) bool { return i == 123 })
	assert.True(t, o.IsPresent())
}

func TestOptional_Filter_Empty(t *testing.T) {
	o := OptionalOf(123).Filter(func(i int) bool { return i == -1 })
	assert.False(t, o.IsPresent())
}

func TestOptional_IfPresent_Present(t *testing.T) {
	var got int
	OptionalOf(123).IfPresent(func(el int) {
		got = el
	})
	assert.Equal(t, 123, got)
}

func TestOptional_IfPresent_Empty(t *testing.T) {
	var called bool
	OptionalEmpty[int]().IfPresent(func(el int) {
		called = true
	})
	assert.False(t, called)
}

func TestOptional_Get_Present(t *testing.T) {
	got := OptionalOf(123).Get()
	assert.Equal(t, 123, got)
}

func TestOptional_Or_NotEmpty(t *testing.T) {
	got := OptionalOf(123).Or(func() Optional[int] { return OptionalOf(456) })
	assert.Equal(t, OptionalOf(123), got)
}

func TestOptional_Or_Empty(t *testing.T) {
	got := OptionalEmpty[int]().Or(func() Optional[int] { return OptionalOf(456) })
	assert.Equal(t, OptionalOf(456), got)
}

func TestOptional_OrElse_NotEmpty(t *testing.T) {
	got := OptionalOf(123).OrElse(456)
	assert.Equal(t, 123, got)
}

func TestOptional_OrElse_Empty(t *testing.T) {
	got := OptionalEmpty[int]().OrElse(456)
	assert.Equal(t, 456, got)
}

func TestOptional_OrElseGet_NotEmpty(t *testing.T) {
	got := OptionalOf(123).OrElseGet(func() int { return 456 })
	assert.Equal(t, 123, got)
}

func TestOptional_OrElseGet_Empty(t *testing.T) {
	got := OptionalEmpty[int]().OrElseGet(func() int { return 456 })
	assert.Equal(t, 456, got)
}

func TestOptional_FlatMap_NotEmpty(t *testing.T) {
	got := OptionalOf(123).FlatMap(func(t int) Optional[Any] { return OptionalOf[Any](t * 2) })
	assert.Equal(t, OptionalOf[Any](246), got)
}

func TestOptional_FlatMap_Empty(t *testing.T) {
	got := OptionalEmpty[int]().FlatMap(func(t int) Optional[Any] { return OptionalOf[Any](t * 2) })
	assert.Equal(t, OptionalEmpty[Any](), got)
}

func TestOptional_Map_NotEmpty_NotNil(t *testing.T) {
	got := OptionalOf(123).Map(func(t int) Any { return t * 2 })
	assert.Equal(t, OptionalOf[Any](246), got)
}

func TestOptional_Map_NotEmpty_Nil(t *testing.T) {
	got := OptionalOf(123).Map(func(t int) Any { return nil })
	assert.Equal(t, OptionalEmpty[Any](), got)
}

func TestOptional_Map_Empty(t *testing.T) {
	got := OptionalEmpty[int]().Map(func(t int) Any { return t * 2 })
	assert.Equal(t, OptionalEmpty[Any](), got)
}
