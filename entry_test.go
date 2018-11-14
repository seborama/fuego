package fuego

import (
	"testing"

	"github.com/raviqqe/hamt"
	"github.com/stretchr/testify/assert"
)

type EntryInt int

func (i EntryInt) Hash() uint32 {
	return uint32(i)
}

// TODO Call this FlatMap instead?
func (i EntryInt) Value() EntryInt {
	return i
}

func (i EntryInt) Equal(e hamt.Entry) bool {
	j, ok := e.(EntryInt)

	if !ok {
		return false
	}

	return i == j
}

func TestEntry(t *testing.T) {
	t.Log(hamt.Entry(EntryInt(42)))
}

func TestEntryKey(t *testing.T) {
	assert.Equal(t, uint32(42), hamt.Entry(EntryInt(42)).Hash())
}
