package utils

import (
	"testing"
)

func TestRandomString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				n: 10,
			},
			want: "",
		},
		{
			name: "test2",
			args: args{
				n: 10,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomString(tt.args.n); got != tt.want {
				t.Errorf("RandomString() = %v, want %v", got, tt.want)
			}
		})
	}
}
