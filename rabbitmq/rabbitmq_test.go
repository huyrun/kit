package rabbitmq

import (
	"reflect"
	"testing"
)

const (
	addr = "amqp://guest:guest@localhost:5672"
)

func TestNewRabbitmq(t *testing.T) {
	type args struct {
		options []Option
	}
	tests := []struct {
		name string
		args args
		want *rabbitmq
	}{
		// TODO: Add test cases.
		{
			name: "test create rabbitmq instance",
			args: args{
				options: []Option{
					Address(addr),
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRabbitmq(tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRabbitmq() = %v, want %v", got, tt.want)
			}
			for {
				continue
			}
		})
	}
}

func Test_rabbitmq_newChannel(t *testing.T) {
	r := NewRabbitmq(Address(addr))

	tests := []struct {
		name string
		want *channel
	}{
		// TODO: Add test cases.
		{
			name: "channel 1",
			want: nil,
		},
		{
			name: "channel 2",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.newChannel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newChannel() = %v, want %v", got, tt.want)
			}
		})
	}
	for {
		continue
	}
}
