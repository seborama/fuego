package fuego_test

import (
	"testing"

	ƒ "github.com/seborama/fuego"
	"github.com/stretchr/testify/assert"
)

func TestEntryFloat(t *testing.T) {
	t.Log(ƒ.EntryFloat(3.1415926536))
}

func TestEntryFloatHash(t *testing.T) {
	assert.Equal(t, uint32(0xc0490fdb), ƒ.EntryFloat(-3.1415926536).Hash())
	assert.Equal(t, uint32(0), ƒ.EntryFloat(0).Hash())
	assert.Equal(t, uint32(0x40490fdb), ƒ.EntryFloat(3.1415926536).Hash())
}

func TestEntryFloatEqual(t *testing.T) {
	assert.True(t, ƒ.EntryFloat(0).Equal(ƒ.EntryFloat(0)))
	assert.True(t, ƒ.EntryFloat(3.1415926536).Equal(ƒ.EntryFloat(3.1415926536)))
	assert.True(t, ƒ.EntryFloat(-3.1415926536).Equal(ƒ.EntryFloat(-3.1415926536)))
	assert.False(t, ƒ.EntryFloat(3.1415926536).Equal(ƒ.EntryFloat(-3.1415926536)))
	assert.False(t, ƒ.EntryFloat(3.1415926536).Equal(EntryFake{}))
}
