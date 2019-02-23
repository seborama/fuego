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

func TestEntryString_Equal(t *testing.T) {
	type args struct {
		e Entry
	}
	tests := []struct {
		name string
		es   EntryString
		args args
		want bool
	}{
		{
			name: "Should return true when string is empty and other string is nil",
			es:   EntryString(""),
			args: args{e: nil},
			want: false,
		},
		{
			name: "Should return true when string is not empty and other string is nil",
			es:   EntryString("hi"),
			args: args{e: nil},
			want: false,
		},
		{
			name: "Should return true when both strings are empty",
			es:   EntryString(""),
			args: args{e: EntryString("")},
			want: true,
		},
		{
			name: "Should return false when one string is not empty",
			es:   EntryString(""),
			args: args{e: EntryString("one")},
			want: false,
		},
		{
			name: "Should return false when other string is not empty",
			es:   EntryString("one"),
			args: args{e: EntryString("")},
			want: false,
		},
		{
			name: "Should return true when both strings are non empty and equal",
			es:   EntryString("one"),
			args: args{e: EntryString("one")},
			want: true,
		},
		{
			name: "Should return false when both strings are non empty and differ",
			es:   EntryString("one"),
			args: args{e: EntryString("two")},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.es.Equal(tt.args.e); got != tt.want {
				t.Errorf("EntryString.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntryString_MapToEntryBytes(t *testing.T) {
	type args struct {
		bufsize int
	}
	tests := []struct {
		name string
		es   EntryString
		args args
		want Stream
	}{
		{
			name: "Should return empty stream for empty input EntryString",
			es:   EntryString(""),
			args: args{bufsize: 0},
			want: NewStreamFromSlice(EntrySlice{}, 0),
		},
		{
			name: "Should return byte stream for non-empty input EntryString",
			es:   EntryString("xyzå©ç"),
			args: args{bufsize: 0},
			want: NewStreamFromSlice(
				func() EntrySlice {
					expected := EntrySlice{}
					for _, val := range []byte("xyzå©ç") {
						expected = append(expected, EntryByte(val))
					}
					return expected
				}(), 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.es.MapToEntryBytes(tt.args.bufsize)
			assert.Equal(t, tt.want.ToSlice(), got.ToSlice())
		})
	}
}
