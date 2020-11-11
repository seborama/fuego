package fuego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryMap(t *testing.T) {
	t.Log(EntryMap{EntryInt(0): EntrySlice{}})
}

func TestEntryMapHash(t *testing.T) {
	type fields struct {
		entrymap EntryMap
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name:   "Should return 0 for an empty entrymap",
			fields: fields{entrymap: EntryMap{}},
			want:   0,
		},
		{
			name: "Should return hash for a single-item entrymap",
			fields: fields{entrymap: EntryMap{
				EntryInt(1): EntrySlice{EntryString("one")},
			}},
			want: 0x7a6c8730,
		},
		{
			name: "Should return hash for a multi-item entrymap",
			fields: fields{entrymap: EntryMap{
				EntryInt(7):  EntrySlice{EntryString("seven"), EntryInt(7)},
				EntryInt(13): EntrySlice{EntryString("thirteen")},
				EntryInt(28): EntrySlice{EntryString("twenty eight")},
			}},
			want: 0xee7059c1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.fields.entrymap.Hash())
		})
	}
}

func TestEntryMapHas_DiffersWhenEntrySliceValuesAreSameButInDifferentOrder(t *testing.T) {
	slice1 := EntryMap{
		EntryInt(7): EntrySlice{EntryInt(7), EntryString("seven")},
	}
	slice2 := EntryMap{
		EntryInt(7): EntrySlice{EntryString("seven"), EntryInt(7)},
	}

	assert.NotEqual(t, slice2.Hash(), slice1.Hash())
}

func TestEntryMapEqual(t *testing.T) {
	type fields struct {
		map1 EntryMap
		map2 Entry
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Should return false for when comparee is not an EntryMap",
			fields: fields{
				map1: EntryMap{},
				map2: EntrySlice{}},
			want: false,
		},
		{
			name: "Should return true for two empty maps",
			fields: fields{
				map1: EntryMap{},
				map2: EntryMap{}},
			want: true,
		},
		{
			name: "Should return false for one empty entrymap and one non-empty entrymap",
			fields: fields{
				map1: EntryMap{},
				map2: EntryMap{
					EntryInt(7): EntrySlice{EntryString("seven")},
				}},
			want: false,
		},
		{
			name: "Should return true for two identical multi-item maps",
			fields: fields{
				map1: EntryMap{
					EntryInt(7):  EntrySlice{EntryString("seven"), EntryInt(7)},
					EntryInt(13): EntrySlice{EntryString("thirteen")},
					EntryInt(28): EntrySlice{EntryString("twenty eight")},
				},
				map2: EntryMap{
					EntryInt(7):  EntrySlice{EntryString("seven"), EntryInt(7)},
					EntryInt(13): EntrySlice{EntryString("thirteen")},
					EntryInt(28): EntrySlice{EntryString("twenty eight")},
				},
			},
			want: true,
		},
		{
			name: "Should return true for two maps with same items but in different key order",
			fields: fields{
				map1: EntryMap{
					EntryInt(7):  EntrySlice{EntryString("seven"), EntryInt(7)},
					EntryInt(13): EntrySlice{EntryString("thirteen")},
					EntryInt(28): EntrySlice{EntryString("twenty eight")},
				},
				map2: EntryMap{
					EntryInt(13): EntrySlice{EntryString("thirteen")},
					EntryInt(7):  EntrySlice{EntryString("seven"), EntryInt(7)},
					EntryInt(28): EntrySlice{EntryString("twenty eight")},
				},
			},
			want: true,
		},
		{
			name: "Should return true for two maps with same items but in different key and value order",
			fields: fields{
				map1: EntryMap{
					EntryInt(7):  EntrySlice{EntryInt(7), EntryString("seven")},
					EntryInt(13): EntrySlice{EntryString("thirteen")},
					EntryInt(28): EntrySlice{EntryString("twenty eight")},
				},
				map2: EntryMap{
					EntryInt(13): EntrySlice{EntryString("thirteen")},
					EntryInt(7):  EntrySlice{EntryString("seven"), EntryInt(7)},
					EntryInt(28): EntrySlice{EntryString("twenty eight")},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.fields.map1.Equal(tt.fields.map2))
		})
	}
}

func TestEntryMap_Stream(t *testing.T) {
	tests := []struct {
		name string
		em   EntryMap
		want EntrySlice
	}{
		{
			name: "Should create a Stream with a nil channel",
			em:   nil,
			want: nil,
		},
		{
			name: "Should create an empty Stream with an empty channel",
			em:   EntryMap{},
			want: EntrySlice{},
		},
		{
			name: "Should create a Stream with a populated channel",
			em: EntryMap{
				EntryInt(7):  EntrySlice{EntryString("seven"), EntryInt(7)},
				EntryInt(13): EntrySlice{EntryString("thirteen")},
				EntryInt(28): EntrySlice{EntryString("twenty eight")},
			},
			want: EntrySlice{
				Tuple2{
					EntryInt(7),
					EntrySlice{EntryString("seven"), EntryInt(7)},
				},
				Tuple2{
					EntryInt(13),
					EntrySlice{EntryString("thirteen")},
				},
				Tuple2{
					EntryInt(28),
					EntrySlice{EntryString("twenty eight")},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got EntrySlice
			if gotStream := tt.em.Stream(0).stream; gotStream != nil {
				got = EntrySlice{}
				for val := range gotStream {
					got = append(got, val)
				}
			}

			if !assert.ElementsMatch(t, got, tt.want) {
				t.Errorf("EntryMap.Stream() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntryMap_Len(t *testing.T) {
	tests := []struct {
		name string
		em   EntryMap
		want int
	}{
		{
			name: "Should return 0 for nil map",
			em:   nil,
			want: 0,
		},
		{
			name: "Should return 0 for empty map",
			em:   EntryMap{},
			want: 0,
		},
		{
			name: "Should return 1 for map of 1 Entry",
			em:   EntryMap{EntryInt(1): EntryInt(123)},
			want: 1,
		},
		{
			name: "Should return 3 for map of 3 Entries",
			em: EntryMap{
				EntryInt(1): EntrySlice{
					EntryInt(123),
					EntryInt(12),
				},
				EntryInt(2): EntryInt(3),
				EntryInt(3): EntryInt(4),
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.em.Len(); got != tt.want {
				t.Errorf("EntryMap.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntryMap_Has(t *testing.T) {
	m := EntryMap{
		EntryInt(1): EntryString("one"),
	}

	assert.True(t, m.HasKey(EntryInt(1)))
	assert.False(t, m.HasKey(EntryInt(2)))
}

func TestEntryMap_Merge(t *testing.T) {
	const (
		keyToAdd   = EntryInt(1)
		valueToAdd = EntryString("new_one")
	)

	tests := []struct {
		name string
		em   EntryMap
		want EntryMap
	}{
		{
			name: "empty map",
			em:   EntryMap{},
			want: EntryMap{
				EntryInt(1): EntryString("new_one"),
			},
		},
		{
			name: "new key",
			em: EntryMap{
				EntryInt(2): EntryString("two"),
			},
			want: EntryMap{
				EntryInt(1): EntryString("new_one"),
				EntryInt(2): EntryString("two"),
			},
		},
		{
			name: "existing key",
			em: EntryMap{
				EntryInt(1): EntryString("one"),
			},
			want: EntryMap{
				EntryInt(1): EntryString("new_one"),
			},
		},
	}

	biIdentity := func(e1, e2 Entry) Entry { return e2 }

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.em.Merge(keyToAdd, valueToAdd, biIdentity)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
