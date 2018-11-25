package fuego

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO weak tests: the order of elements in HamtMap/HamtSet is not guaranteed

func TestHamtMap_Insert(t *testing.T) {
	type fields struct {
		myMap Map
	}
	type args struct {
		k Entry
		v interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Map
	}{
		{
			name: "Should Insert entries into Map",
			fields: fields{
				myMap: NewHamtMap().
					Insert(EntryInt(1), "one"),
			},
			args: args{
				k: EntryInt(5),
				v: "five",
			},
			want: NewHamtMap().
				Insert(EntryInt(1), "one").
				Insert(EntryInt(5), "five"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.myMap.Insert(tt.args.k, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHamtMap_Delete(t *testing.T) {
	type fields struct {
		myMap Map
	}
	type args struct {
		k Entry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Map
	}{
		{
			name: "Should Insert entries into Map",
			fields: fields{
				myMap: NewHamtMap().
					Insert(EntryInt(1), "one").
					Insert(EntryInt(5), "five").
					Insert(EntryInt(2), "two"),
			},
			args: args{
				k: EntryInt(5),
			},
			want: NewHamtMap().
				Insert(EntryInt(1), "one").
				Insert(EntryInt(2), "two"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.myMap.Delete(tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHamtMap_Size(t *testing.T) {
	m := NewHamtMap().
		Insert(EntryInt(1), "one").
		Insert(EntryInt(5), "five").
		Insert(EntryInt(2), "two")

	assert.Equal(t, 3, m.Size())
}

func TestHamtMap_Merge(t *testing.T) {
	type fields struct {
		myMap Map
	}
	type args struct {
		t Map
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Map
	}{
		{
			name: "Should merge two excluding maps",
			fields: fields{
				myMap: NewHamtMap().
					Insert(EntryInt(7), "seven").
					Insert(EntryInt(2), "two"),
			},
			args: args{
				t: NewHamtMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(9), "nine"),
			},
			want: NewHamtMap().
				Insert(EntryInt(7), "seven").
				Insert(EntryInt(2), "two").
				Insert(EntryInt(3), "three").
				Insert(EntryInt(9), "nine"),
		},
		{
			name: "Should merge two overlapping sets",
			fields: fields{
				myMap: NewHamtMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(1), "one"),
			},
			args: args{
				t: NewHamtMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(1), "one"),
			},
			want: NewHamtMap().
				Insert(EntryInt(3), "three").
				Insert(EntryInt(1), "one"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.myMap.Merge(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map.Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHamtMap_Get(t *testing.T) {
	type fields struct {
		myMap Map
	}
	type args struct {
		k Entry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "Should not find missing item",
			fields: fields{
				myMap: NewHamtMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(1), "one"),
			},
			args: args{
				EntryInt(999),
			},
			want: EntryNone{},
		},
		{
			name: "Should find existing item",
			fields: fields{
				myMap: NewHamtMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(1), "one"),
			},
			args: args{
				EntryInt(3),
			},
			want: "three",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.myMap.Get(tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHamtMap_Has(t *testing.T) {
	type fields struct {
		myMap Map
	}
	type args struct {
		k Entry
		v interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Should not have missing item",
			fields: fields{
				myMap: NewHamtMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(1), "one"),
			},
			args: args{
				k: EntryInt(999),
				v: "sos",
			},
			want: false,
		},
		{
			name: "Should have existing item",
			fields: fields{
				myMap: NewHamtMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(1), "one"),
			},
			args: args{
				k: EntryInt(3),
				v: "three",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.myMap.Has(tt.args.k, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHamtMap_HasKey(t *testing.T) {
	type fields struct {
		myMap Map
	}
	type args struct {
		k Entry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Should not have missing item",
			fields: fields{
				myMap: NewHamtMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(1), "one"),
			},
			args: args{
				k: EntryInt(999),
			},
			want: false,
		},
		{
			name: "Should have existing item",
			fields: fields{
				myMap: NewHamtMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(1), "one"),
			},
			args: args{
				k: EntryInt(3),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.myMap.HasKey(tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map.HasKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHamtMap_HasValue(t *testing.T) {
	type fields struct {
		myMap Map
	}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Should not have missing item",
			fields: fields{
				myMap: NewHamtMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(1), "one"),
			},
			args: args{
				v: "sos",
			},
			want: false,
		},
		{
			name: "Should have existing item",
			fields: fields{
				myMap: NewHamtMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(1), "one"),
			},
			args: args{
				v: "three",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.myMap.HasValue(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map.TestMap_HasValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHamtMap_FirstRest(t *testing.T) {
	type fields struct {
		myMap Map
	}
	tests := []struct {
		name   string
		fields fields
		want   Entry
		want1  interface{}
		want2  Map
	}{
		{
			name: "Should have existing item",
			fields: fields{
				myMap: NewHamtMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(11), "eleven").
					Insert(EntryInt(12), "twelve").
					Insert(EntryInt(-23), "minus twenty three").
					Insert(EntryInt(7), "seven"),
			},
			want:  EntryInt(3),
			want1: "three",
			want2: NewHamtMap().
				Insert(EntryInt(11), "eleven").
				Insert(EntryInt(12), "twelve").
				Insert(EntryInt(-23), "minus twenty three").
				Insert(EntryInt(7), "seven"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := tt.fields.myMap.FirstRest()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map.FirstRest() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Map.FirstRest() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("Map.FirstRest() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestHamtMap_EntrySet(t *testing.T) {
	type fields struct {
		myMap Map
	}
	tests := []struct {
		name   string
		fields fields
		want   Set
	}{
		{
			name: "Should return EntrySet",
			fields: fields{
				myMap: NewHamtMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(11), "eleven").
					Insert(EntryInt(12), "twelve").
					Insert(EntryInt(-23), "minus twenty three").
					Insert(EntryInt(7), "seven"),
			},
			want: NewHamtSet().
				Insert(NewMapEntry(EntryInt(3), "three")).
				Insert(NewMapEntry(EntryInt(11), "eleven")).
				Insert(NewMapEntry(EntryInt(12), "twelve")).
				Insert(NewMapEntry(EntryInt(-23), "minus twenty three")).
				Insert(NewMapEntry(EntryInt(7), "seven")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.myMap.EntrySet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map.EntrySet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHamtMap_KeySet(t *testing.T) {
	type fields struct {
		myMap Map
	}
	tests := []struct {
		name   string
		fields fields
		want   Set
	}{
		{
			name: "Should return an empty key set",
			fields: fields{
				myMap: HamtMap{},
			},
			want: NewHamtSet(),
		},
		{
			name: "Should return the key set",
			fields: fields{
				myMap: NewHamtMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(11), "eleven").
					Insert(EntryInt(12), "twelve").
					Insert(EntryInt(-23), "minus twenty three").
					Insert(EntryInt(7), "seven"),
			},
			want: NewHamtSet().
				Insert(EntryInt(3)).
				Insert(EntryInt(11)).
				Insert(EntryInt(12)).
				Insert(EntryInt(-23)).
				Insert(EntryInt(7)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.fields.myMap.KeySet()
			if got := tt.fields.myMap.KeySet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map.KeySet() = %v, want %v", got, tt.want)
			}
		})
	}
}
