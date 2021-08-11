package rabbitmq

import (
	"encoding/json"
	"testing"
)

func Test_channel_Publish(t *testing.T) {
	r := NewRabbitmq(Address(addr))
	c := r.CreateProducer(
		ExchangeName("exchange_fanout"),
		ExchangeKind(ExchangeFanout),
		RegisterMarshalFunc(json.Marshal),
	)

	messages := []Msg{}
	for i := 0; i < 100; i++ {
		messages = append(messages, Msg{
			headers:    header{},
			Body:       i,
			RoutingKey: "",
		})
	}

	type args struct {
		messages []Msg
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test publish msg",
			args: args{
				messages: messages,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.args.messages {
				if err := c.Publish(msg); (err != nil) != tt.wantErr {
					t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
	for {
		continue
	}
}
