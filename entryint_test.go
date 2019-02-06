package fuego_test

import (
	"testing"

	ƒ "github.com/seborama/fuego"
	"github.com/stretchr/testify/assert"
)

func TestEntryInt(t *testing.T) {
	t.Log(ƒ.EntryInt(42))
}

func TestEntryIntHash(t *testing.T) {
	assert.Equal(t, uint32(0xffffffd6), ƒ.EntryInt(-42).Hash())
	assert.Equal(t, uint32(0), ƒ.EntryInt(0).Hash())
	assert.Equal(t, uint32(42), ƒ.EntryInt(42).Hash())
}

type EntryFake struct{}

func (ef EntryFake) Hash() uint32       { return 2 }
func (ef EntryFake) Equal(ƒ.Entry) bool { return true }

func TestEntryIntEqual(t *testing.T) {
	assert.True(t, ƒ.EntryInt(0).Equal(ƒ.EntryInt(0)))
	assert.True(t, ƒ.EntryInt(42).Equal(ƒ.EntryInt(42)))
	assert.True(t, ƒ.EntryInt(-42).Equal(ƒ.EntryInt(-42)))
	assert.False(t, ƒ.EntryInt(42).Equal(ƒ.EntryInt(-42)))
	assert.False(t, ƒ.EntryInt(2).Equal(EntryFake{}))
}
