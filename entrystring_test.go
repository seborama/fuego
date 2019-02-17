package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryString(t *testing.T) {
	t.Log(EntryString("Hello World"))
}

func TestEntryStringHash(t *testing.T) {
	assert.Equal(t, uint32(0x4a17b156), EntryString("Hello World").Hash())
}

func TestEntryString_ToUpper(t *testing.T) {
	tests := []struct {
		name string
		s    EntryString
		want EntryString
	}{
		{
			name: "Should return empty output for empty input",
			s:    EntryString(""),
			want: EntryString(""),
		},
		{
			name: "Should return upper cased output for input",
			s:    EntryString("aQw8 S.=!9s"),
			want: EntryString("AQW8 S.=!9S"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ToUpper(); got != tt.want {
				t.Errorf("EntryString.ToUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntryString_ToLower(t *testing.T) {
	tests := []struct {
		name string
		s    EntryString
		want EntryString
	}{
		{
			name: "Should return empty output for empty input",
			s:    EntryString(""),
			want: EntryString(""),
		},
		{
			name: "Should return lower cased output for input",
			s:    EntryString("aQw8 S.=!9s"),
			want: EntryString("aqw8 s.=!9s"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ToLower(); got != tt.want {
				t.Errorf("EntryString.ToLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntryString_Len(t *testing.T) {
	tests := []struct {
		name string
		es   EntryString
		want EntryInt
	}{
		{
			name: "Should return 0 for empty string",
			es:   EntryString(""),
			want: 0,
		},
		{
			name: "Should return 22 for utf8 string",
			es:   EntryString("éñ!qèøƒ-œ;”’"),
			want: 22,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.es.Len(); got != tt.want {
				t.Errorf("EntryString.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}
