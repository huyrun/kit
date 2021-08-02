package rabbitmq

import (
	"reflect"
	"testing"

	"github.com/streadway/amqp"
)

func Test_createChannel(t *testing.T) {
	r := NewRabbitmq(Address(addr))

	type args struct {
		conn *amqp.Connection
	}
	tests := []struct {
		name string
		args args
		want *amqp.Channel
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				conn: r.connection,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createChannel(tt.args.conn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createChannel() = %v, want %v", got, tt.want)
			}
		})
	}
}
