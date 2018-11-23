package fuego

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedMap_Insert(t *testing.T) {
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
			name: "Should Insert entry into Map",
			fields: fields{
				myMap: NewOrderedMap().
					Insert(EntryInt(1), "one"),
			},
			args: args{
				k: EntryInt(5),
				v: "five",
			},
			want: NewOrderedMap().
				Insert(EntryInt(1), "one").
				Insert(EntryInt(5), "five"),
		},
		{
			name: "Should replace duplicate entry on Insert into Map",
			fields: fields{
				myMap: NewOrderedMap().
					Insert(EntryInt(1), "one").
					Insert(EntryInt(5), "five"),
			},
			args: args{
				k: EntryInt(5),
				v: "cinq",
			},
			want: NewOrderedMap().
				Insert(EntryInt(1), "one").
				Insert(EntryInt(5), "cinq"),
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

func TestOrderedMap_Delete(t *testing.T) {
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
			name: "Should Delete entry from Map",
			fields: fields{
				myMap: NewOrderedMap().
					Insert(EntryInt(1), "one").
					Insert(EntryInt(5), "five").
					Insert(EntryInt(2), "two"),
			},
			args: args{
				k: EntryInt(5),
			},
			want: NewOrderedMap().
				Insert(EntryInt(1), "one").
				Insert(EntryInt(2), "two"),
		},
		{
			name: "Should not Delete entry from Map",
			fields: fields{
				myMap: NewOrderedMap().
					Insert(EntryInt(1), "one").
					Insert(EntryInt(5), "five").
					Insert(EntryInt(2), "two"),
			},
			args: args{
				k: EntryInt(0101010),
			},
			want: NewOrderedMap().
				Insert(EntryInt(1), "one").
				Insert(EntryInt(5), "five").
				Insert(EntryInt(2), "two"),
		},
		{
			name: "Should remove first entry in Map",
			fields: fields{
				myMap: NewOrderedMap().
					Insert(EntryInt(1), "one").
					Insert(EntryInt(33), "thirty-three").
					Insert(EntryInt(5), "five"),
			},
			args: args{
				k: EntryInt(1),
			},
			want: NewOrderedMap().
				Insert(EntryInt(33), "thirty-three").
				Insert(EntryInt(5), "five"),
		},
		{
			name: "Should remove last entry in Map",
			fields: fields{
				myMap: NewOrderedMap().
					Insert(EntryInt(1), "one").
					Insert(EntryInt(33), "thirty-three").
					Insert(EntryInt(5), "five"),
			},
			args: args{
				k: EntryInt(5),
			},
			want: NewOrderedMap().
				Insert(EntryInt(1), "one").
				Insert(EntryInt(33), "thirty-three"),
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

func TestOrderedMap_Size(t *testing.T) {
	m := NewOrderedMap().
		Insert(EntryInt(1), "one").
		Insert(EntryInt(5), "five").
		Insert(EntryInt(2), "two")

	assert.Equal(t, 3, m.Size())
}

func TestOrderedMap_Merge(t *testing.T) {
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
				myMap: NewOrderedMap().
					Insert(EntryInt(7), "seven").
					Insert(EntryInt(2), "two"),
			},
			args: args{
				t: NewOrderedMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(9), "nine"),
			},
			want: NewOrderedMap().
				Insert(EntryInt(7), "seven").
				Insert(EntryInt(2), "two").
				Insert(EntryInt(3), "three").
				Insert(EntryInt(9), "nine"),
		},
		{
			name: "Should merge two overlapping sets",
			fields: fields{
				myMap: NewOrderedMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(1), "one"),
			},
			args: args{
				t: NewOrderedMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(1), "one"),
			},
			want: NewOrderedMap().
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

func TestOrderedMap_Get(t *testing.T) {
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
				myMap: NewOrderedMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(1), "one"),
			},
			args: args{
				EntryInt(999),
			},
			want: notFound{},
		},
		{
			name: "Should find existing item",
			fields: fields{
				myMap: NewOrderedMap().
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

func TestOrderedMap_Has(t *testing.T) {
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
				myMap: NewOrderedMap().
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
				myMap: NewOrderedMap().
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

func TestOrderedMap_HasKey(t *testing.T) {
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
				myMap: NewOrderedMap().
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
				myMap: NewOrderedMap().
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

func TestOrderedMap_HasValue(t *testing.T) {
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
				myMap: NewOrderedMap().
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
				myMap: NewOrderedMap().
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

func TestOrderedMap_FirstRest(t *testing.T) {
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
				myMap: NewOrderedMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(11), "eleven").
					Insert(EntryInt(12), "twelve").
					Insert(EntryInt(-23), "minus twenty three").
					Insert(EntryInt(7), "seven"),
			},
			want:  EntryInt(3),
			want1: "three",
			want2: NewOrderedMap().
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

func TestOrderedMap_EntrySet(t *testing.T) {
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
				myMap: NewOrderedMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(11), "eleven").
					Insert(EntryInt(12), "twelve").
					Insert(EntryInt(-23), "minus twenty three").
					Insert(EntryInt(7), "seven"),
			},
			want: NewOrderedSet().
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

func TestOrderedMap_KeySet(t *testing.T) {
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
				myMap: OrderedMap{},
			},
			want: NewOrderedSet(),
		},
		{
			name: "Should return the key set",
			fields: fields{
				myMap: NewOrderedMap().
					Insert(EntryInt(3), "three").
					Insert(EntryInt(11), "eleven").
					Insert(EntryInt(12), "twelve").
					Insert(EntryInt(-23), "minus twenty three").
					Insert(EntryInt(7), "seven"),
			},
			want: NewOrderedSet().
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
