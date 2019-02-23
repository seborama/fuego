package fuego

import (
	"hash/crc32"
	"strings"
)

// EntryString is an Entry for 'string'.
type EntryString string

// Hash returns a hash for 'i'.
func (es EntryString) Hash() uint32 {
	return crc32.ChecksumIEEE([]byte(es))
}

// Equal returns true if 'e' and 'i' are equal.
func (es EntryString) Equal(e Entry) bool {
	es2, ok := e.(EntryString)

	if !ok {
		return false
	}

	return es == es2
}

// ToUpper transforms this EntryString to upper case.
func (es EntryString) ToUpper() EntryString {
	return EntryString(
		strings.ToUpper(string(es)),
	)
}

// ToLower transforms this EntryString to lower case.
func (es EntryString) ToLower() EntryString {
	return EntryString(
		strings.ToLower(string(es)),
	)
}

// Len transforms this EntryString to lower case.
func (es EntryString) Len() EntryInt {
	return EntryInt(len(es))
}

// MapToEntryBytes transforms the bytes of this EntryString to
// a Stream of EntryBytes.
func (es EntryString) MapToEntryBytes(bufsize int) Stream {
	ebs := []Entry{}
	for _, b := range []byte(es) {
		ebs = append(ebs, EntryByte(b))
	}
	return NewStreamFromSlice(ebs, bufsize)
}
