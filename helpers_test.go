package fuego_test

import (
	"hash/crc32"
	"testing"

	ƒ "github.com/seborama/fuego"
	"github.com/stretchr/testify/assert"
)

// EntryString is a ƒ.Entry for 'string'.
type EntryString string

// Hash returns a hash for 'i'.
func (i EntryString) Hash() uint32 {
	return crc32.ChecksumIEEE([]byte(i))
}

// Equal returns true if 'e' and 'i' are equal.
func (i EntryString) Equal(e ƒ.Entry) bool {
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

// concatenateStringsBiFunc returns a ƒ.BiFunction that
// joins 'i' and 'j' together with a '-' in between.
func concatenateStringsBiFunc(i, j ƒ.Entry) ƒ.Entry {
	iStr := i.(EntryString)
	jStr := j.(EntryString)
	return iStr + "-" + jStr
}

// timesTwo returns a ƒ.Function than multiplies 'i' by 2.
func timesTwo() ƒ.Function {
	return func(i ƒ.Entry) ƒ.Entry {
		return 2 * i.(ƒ.EntryInt)
	}
}

// isEvenNumberFunction is a Function that returns true when 'i' is
// an even number.
func isEvenNumberFunction(i ƒ.Entry) ƒ.Entry {
	return ƒ.EntryBool(i.(ƒ.EntryInt)&1 == 0)
}

// intGreaterThanPredicate is a Predicate for numbers greater
// than 'rhs'.
func intGreaterThanPredicate(rhs int) ƒ.Predicate {
	return func(lhs ƒ.Entry) bool {
		return int(lhs.(ƒ.EntryInt)) > rhs
	}
}
