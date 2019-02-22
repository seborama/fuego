package fuego

import "testing"

func TestEntriesEqual(t *testing.T) {
	type args struct {
		e1 Entry
		e2 Entry
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should equal",
			args: args{
				e1: EntryString("hi"),
				e2: EntryString("hi"),
			},
			want: true,
		},
		{
			name: "Should not equal",
			args: args{
				e1: EntryInt(0),
				e2: EntryFloat(0),
			},
			want: false,
		},
		{
			name: "Should not equal when e1 is nil",
			args: args{
				e1: nil,
				e2: EntryFloat(0),
			},
			want: false,
		},
		{
			name: "Should not equal when e2 is nil",
			args: args{
				e1: EntryInt(0),
				e2: nil,
			},
			want: false,
		},
		{
			name: "Should equal when both nil",
			args: args{
				e1: nil,
				e2: nil,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EntriesEqual(tt.args.e1, tt.args.e2); got != tt.want {
				t.Errorf("EntriesEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
