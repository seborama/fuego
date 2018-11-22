package fuego_test

import (
	"hash/crc32"
	"testing"

	"github.com/raviqqe/hamt"
	ƒ "github.com/seborama/fuego"
	"github.com/stretchr/testify/assert"
)

// EntryString is a ƒ.Entry for 'int'.
type EntryInt int

// Hash returns a hash for 'i'.
func (i EntryInt) Hash() uint32 {
	return uint32(i)
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

// EntryString is a ƒ.Entry for 'string'.
type EntryString string

// Hash returns a hash for 'i'.
func (i EntryString) Hash() uint32 {
	return crc32.ChecksumIEEE([]byte(i))
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
func isEvenNumber(t ƒ.Entry) bool {
	k := (t.(ƒ.MapEntry)).K
	return int(k.(EntryInt))&1 == 0
}

// isOddNumber is a Predicate for odd numbers.
func isOddNumber(t ƒ.Entry) bool {
	v := t.(EntryInt)
	return v&1 == 0
}

// concatenateStringsBiFunc returns a fuego.BiFunction that
// joins 'i' and 'j' together with a '-' in between.
func concatenateStringsBiFunc(i, j ƒ.Entry) ƒ.Entry {
	iStr := i.(EntryString)
	jStr := j.(EntryString)
	return EntryString(iStr + "-" + jStr)
}

// timesTwo returns a fuego.Function than multiplies 'i' by 2.
func timesTwo() ƒ.Function {
	return func(i ƒ.Entry) ƒ.Entry {
		return EntryInt(2 * i.(EntryInt))
	}
}

// isEvenNumberFunction is a Function that returns true when 'i' is
// an even number.
func isEvenNumberFunction(i ƒ.Entry) ƒ.Entry {
	return ƒ.EntryBool(i.(EntryInt)&1 == 0)
}

// intGreaterThanPredicate is a Predicate for numbers greater
// than 'rhs'.
func intGreaterThanPredicate(rhs int) ƒ.Predicate {
	return func(lhs ƒ.Entry) bool {
		return int(lhs.(EntryInt)) > rhs
	}
}
