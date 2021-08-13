package rabbitmq

import (
	"encoding/json"
	"testing"
)

type msg struct {
	header     map[string]interface{}
	body       interface{}
	routingKey string
	priority   int
}

func (m *msg) MessageHeaderInit()        { m.header = make(Header) }
func (m *msg) MessageHeader() Header     { return m.header }
func (m *msg) MessageBody() Body         { return m.body }
func (m *msg) MessageRoutingKey() string { return m.routingKey }
func (m *msg) MessagePriority() int      { return m.priority }

func Test_channel_Publish(t *testing.T) {
	r := NewRabbitmq(Address(addr))
	c := r.CreateProducer(
		ExchangeName("exchange_fanout"),
		ExchangeKind(ExchangeFanout),
		RegisterMarshalFunc(json.Marshal),
	)

	messages := []Message{}
	for i := 0; i < 100; i++ {
		msg := &msg{
			header:     Header{},
			body:       i,
			routingKey: "",
		}
		messages = append(messages, msg)
	}

	type args struct {
		messages []Message
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
