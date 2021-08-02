package rabbitmq

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func Test_channel_Consume(t *testing.T) {
	r1 := NewRabbitmq(Address(addr))
	c := r1.CreatePublisher(
		ExchangeName("exchange_direct"),
		ExchangeKind(ExchangeDirect),
		RegisterMarshalFunc(json.Marshal),
	)

	r2 := NewRabbitmq(Address(addr))
	q := r2.CreateConsumer(
		QueueName("queue_test1"),
		PriorityQueue(200),
		RegisterHandlerFunc(func(bytes []byte) error {
			fmt.Println(string(bytes))
			time.Sleep(2 * time.Second)
			return nil
		}),
	)

	q.Bind("exchange_direct", "q1")

	messages := []Msg{}
	for i := 1; i <= 100; i++ {
		messages = append(messages, Msg{
			Headers:    Header{},
			Body:       i,
			RoutingKey: "q1",
			Priority:   i,
		})
	}

	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range messages {
				c.Publish(msg)
			}
		})
	}
	time.Sleep(5 * time.Second)
	q.Consume()
	for {
		continue
	}
}
