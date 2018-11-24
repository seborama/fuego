package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryBool(t *testing.T) {
	t.Log(EntryBool(true))
}

func TestEntryBoolHash(t *testing.T) {
	assert.Equal(t, uint32(1), EntryBool(true).Hash())
	assert.Equal(t, uint32(0), EntryBool(false).Hash())
}

func TestEntryBoolEqual(t *testing.T) {
	assert.True(t, EntryBool(true).Equal(EntryBool(true)))
	assert.True(t, EntryBool(false).Equal(EntryBool(false)))
	assert.False(t, EntryBool(true).Equal(EntryBool(false)))
	assert.False(t, EntryBool(false).Equal(EntryBool(true)))

	assert.False(t, EntryBool(true).Equal(EntryString("pop")))
}
