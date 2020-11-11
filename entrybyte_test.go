package fuego_test

import (
	"testing"

	ƒ "github.com/seborama/fuego/v10"
	"github.com/stretchr/testify/assert"
)

func TestEntryByte(t *testing.T) {
	t.Log(ƒ.EntryByte(42))
}

func TestEntryByteHash(t *testing.T) {
	assert.Equal(t, uint32(0), ƒ.EntryByte(0).Hash())
	assert.Equal(t, uint32(42), ƒ.EntryByte(42).Hash())
}

func TestEntryByteEqual(t *testing.T) {
	assert.True(t, ƒ.EntryByte(0).Equal(ƒ.EntryByte(0)))
	assert.True(t, ƒ.EntryByte(42).Equal(ƒ.EntryByte(42)))
	assert.False(t, ƒ.EntryByte(2).Equal(EntryFake{}))
}
