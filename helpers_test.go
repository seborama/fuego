package fuego_test

import (
	"hash/crc32"
	"testing"

	"github.com/raviqqe/hamt"
	ƒ "github.com/seborama/fuego"
	"github.com/stretchr/testify/assert"
)

// EntryString is a hamt.Entry for 'int'.
type EntryInt int

// Hash returns a hash for 'i'.
func (i EntryInt) Hash() uint32 {
	return uint32(i)
}

// Value returns the inner value of this EntryInt.
// TODO Call this FlatMap instead?
func (i EntryInt) Value() EntryInt {
	return i
}

// Equal returns true if 'e' and 'i' are equal.
func (i EntryInt) Equal(e hamt.Entry) bool {
	j, ok := e.(EntryInt)

	if !ok {
		return false
	}

	return i == j
}

func TestEntryInt(t *testing.T) {
	t.Log(EntryInt(42))
}

func TestEntryIntKey(t *testing.T) {
	assert.Equal(t, uint32(42), EntryInt(42).Hash())
}

// EntryString is a hamt.Entry for 'string'.
type EntryString string

// Hash returns a hash for 'i'.
func (i EntryString) Hash() uint32 {
	return crc32.ChecksumIEEE([]byte(i))
}

// Value returns the inner value of this EntryInt.
// TODO Call this FlatMap instead?
func (i EntryString) Value() EntryString {
	return i
}

// Equal returns true if 'e' and 'i' are equal.
func (i EntryString) Equal(e hamt.Entry) bool {
	j, ok := e.(EntryString)

	if !ok {
		return false
	}

	return i == j
}

func TestEntryString(t *testing.T) {
	t.Log(EntryString("Hello World"))
}

func TestEntryStringKey(t *testing.T) {
	assert.Equal(t, uint32(0x4a17b156), EntryString("Hello World").Hash())
}

// isEvenNumber is a Predicate for even numbers.
func isEvenNumber(t interface{}) bool {
	k := (t.(ƒ.MapEntry)).K.(EntryInt)
	return k.Value()&1 == 0
}

// isOddNumber is a Predicate for odd numbers.
func isOddNumber(t interface{}) bool {
	v := t.(EntryInt)
	return v&1 == 0
}

// concatenateStringsBiFunc returns a fuego.BiFunction that
// joins 'i' and 'j' together with a '-' in between.
func concatenateStringsBiFunc(i, j interface{}) interface{} {
	iStr := i.(EntryString)
	jStr := j.(EntryString)
	return EntryString(iStr + "-" + jStr)
}

// timesTwo returns a fuego.Function than multiplies 'i' by 2.
func timesTwo() ƒ.Function {
	return func(i interface{}) interface{} {
		return (EntryInt(2 * i.(int)))
	}
}

// isEvenNumberFunction is a Function that returns true when 'i' is
// an even number.
func isEvenNumberFunction(i interface{}) interface{} {
	return i.(int)&1 == 0
}

// intGreaterThanPredicate is a Predicate for numbers greater than
// 'rhs'.
func intGreaterThanPredicate(rhs int) ƒ.Predicate {
	return func(lhs interface{}) bool {
		return lhs.(int) > rhs
	}
}
