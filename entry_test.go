package fuego

import (
	"hash/crc32"
	"testing"

	"github.com/raviqqe/hamt"
	"github.com/stretchr/testify/assert"
)

// EntryInt is an Entry for 'int'.
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

// EntryString is an Entry for 'string'.
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
