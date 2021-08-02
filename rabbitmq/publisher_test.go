package rabbitmq

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_rabbitmq_CreatePublisher(t *testing.T) {
	r := NewRabbitmq(Address(addr))

	type args struct {
		options []ChannOption
	}
	tests := []struct {
		name string
		args args
		want Producer
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				options: []ChannOption{
					ExchangeName("exchange_fanout"),
					ExchangeKind(ExchangeFanout),
					RegisterMarshalFunc(json.Marshal),
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.CreatePublisher(tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreatePublisher() = %v, want %v", got, tt.want)
			}
		})
	}
	for {
		continue
	}
}
