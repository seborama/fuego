package fuego

import (
	"hash/crc32"
	"testing"

	"github.com/stretchr/testify/assert"
)

// EntryString is an Entry for 'string'.
type EntryString string

// Hash returns a hash for 'i'.
func (i EntryString) Hash() uint32 {
	return crc32.ChecksumIEEE([]byte(i))
}

// Equal returns true if 'e' and 'i' are equal.
func (i EntryString) Equal(e Entry) bool {
	j, ok := e.(EntryString)

	if !ok {
		return false
	}

	return i == j
}

func TestEntryString(t *testing.T) {
	t.Log(EntryString("Hello World"))
}

func TestEntryStringHash(t *testing.T) {
	assert.Equal(t, uint32(0x4a17b156), EntryString("Hello World").Hash())
}
