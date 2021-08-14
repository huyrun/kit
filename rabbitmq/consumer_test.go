package rabbitmq

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_rabbitmq_CreateConsumer(t *testing.T) {
	r := NewRabbitmq(Address(addr))

	type args struct {
		options []ChannOption
	}
	tests := []struct {
		name string
		args args
		want *consumer
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				options: []ChannOption{
					QueueName("queue_test"),
					RegisterHandlerFunc(func(bytes []byte) error {
						fmt.Println(string(bytes))
						return nil
					}),
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.CreateConsumer(tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateConsumer() = %v, want %v", got, tt.want)
			}
		})
	}
}
