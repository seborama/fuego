package fuego

import (
	"testing"
)

func TestMapEntry_Hash(t *testing.T) {
	type fields struct {
		k Entry
		v interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "Should return the hash of the Entry",
			fields: fields{
				k: EntryInt(31415),
				v: "hello world",
			},
			want: EntryInt(31415).Hash(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := MapEntry{
				K: tt.fields.k,
				V: tt.fields.v,
			}
			if got := me.Hash(); got != tt.want {
				t.Errorf("MapEntry.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapEntry_Equal(t *testing.T) {
	type fields struct {
		k Entry
		v interface{}
	}
	type args struct {
		e Entry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Should be deep equal",
			fields: fields{
				k: EntryInt(257),
				v: "test",
			},
			args: args{
				e: EntryInt(257),
			},
			want: true,
		},
		{
			name: "Should NOT be equal (keys differ)",
			fields: fields{
				k: EntryInt(123),
				v: "test",
			},
			args: args{
				e: EntryInt(257),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := MapEntry{
				K: tt.fields.k,
				V: tt.fields.v,
			}
			if got := me.Equal(tt.args.e); got != tt.want {
				t.Errorf("MapEntry.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapEntry_DeepEqual(t *testing.T) {
	type fields struct {
		k Entry
		v interface{}
	}
	type args struct {
		o MapEntry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Should be deep equal",
			fields: fields{
				k: EntryInt(257),
				v: "test",
			},
			args: args{
				o: MapEntry{
					K: EntryInt(257),
					V: "test",
				},
			},
			want: true,
		},
		{
			name: "Should NOT be deep equal (keys differ)",
			fields: fields{
				k: EntryInt(123),
				v: "test",
			},
			args: args{
				o: MapEntry{
					K: EntryInt(257),
					V: "test",
				},
			},
			want: false,
		},
		{
			name: "Should NOT be deep equal (values differ)",
			fields: fields{
				k: EntryInt(257),
				v: "not test",
			},
			args: args{
				o: MapEntry{
					K: EntryInt(257),
					V: "test",
				},
			},
			want: false,
		},
		{
			name: "Should NOT be deep equal (values differ in type)",
			fields: fields{
				k: EntryInt(257),
				v: "not test",
			},
			args: args{
				o: MapEntry{
					K: EntryInt(257),
					V: struct{ t string }{t: "test"},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := MapEntry{
				K: tt.fields.k,
				V: tt.fields.v,
			}
			if got := me.DeepEqual(tt.args.o); got != tt.want {
				t.Errorf("MapEntry.DeepEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapEntry_EqualMapEntry(t *testing.T) {
	type fields struct {
		K Entry
		V interface{}
	}
	type args struct {
		ome MapEntry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Should be equal",
			fields: fields{
				K: EntryString("Hello"),
				V: 1234,
			},
			args: args{
				ome: MapEntry{
					K: EntryString("Hello"),
					V: 1234,
				},
			},
			want: true,
		},
		{
			name: "Should differ by key",
			fields: fields{
				K: EntryString("World"),
				V: 1234,
			},
			args: args{
				ome: MapEntry{
					K: EntryString("Hello"),
					V: 1234,
				},
			},
			want: false,
		},
		{
			name: "Should equal when differing values but same keys",
			fields: fields{
				K: EntryString("Hello"),
				V: 9876,
			},
			args: args{
				ome: MapEntry{
					K: EntryString("Hello"),
					V: 1234,
				},
			},
			want: true,
		},
		{
			name: "Should differ by key and value",
			fields: fields{
				K: EntryString("World"),
				V: 9876,
			},
			args: args{
				ome: MapEntry{
					K: EntryString("Hello"),
					V: 1234,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := MapEntry{
				K: tt.fields.K,
				V: tt.fields.V,
			}
			if got := me.EqualMapEntry(tt.args.ome); got != tt.want {
				t.Errorf("MapEntry.EqualMapEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}
