package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryNone(t *testing.T) {
	t.Log(EntryNone{})
}

func TestEntryNoneHash(t *testing.T) {
	assert.Equal(t, uint32(0), EntryNone{}.Hash())
}

func TestEntryNoneEqual(t *testing.T) {
	assert.True(t, EntryNone{}.Equal(EntryNone{}))
	assert.False(t, EntryNone{}.Equal(EntryString("pop")))
}
