package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeq_AsGo_Immutability(t *testing.T) {
	unit := NewSeq()
	unit = unit.Append(EntryInt(1))
	unit = unit.Append(EntryInt(2))
	unit = unit.Append(EntryInt(3))
	assert.Equal(t, 3, unit.Size())

	seq := unit.AsGo()
	assert.EqualValues(t, unit.seq, seq)

	seq[0] = EntryInt(1000)
	assert.NotEqual(t, unit.seq[0], seq[0])
}
