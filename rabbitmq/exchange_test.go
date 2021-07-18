package rabbitmq

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_exchange_Publish(t *testing.T) {
	r, err := CreatePublisher(url)
	if err != nil {
		logrus.WithError(err).Infof("Create publisher fail, url: %s", url)
	}

	defer r.Close()

	type exchange struct {
		name       string
		kind       ExchangeType
		routingKey string
	}
	type args struct {
		message interface{}
	}
	tests := []struct {
		name     string
		exchange exchange
		args     args
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "test publish defaul exchange",
			exchange: exchange{
				name:       "",
				kind:       ExchangeDefault,
				routingKey: "test.queue1",
			},
			args: args{message: "goalllllllll!!!!!!!"},
		},
		{
			name: "test publish fanout exchange",
			exchange: exchange{
				name:       "test.exchange.fanout",
				kind:       ExchangeFanout,
				routingKey: "",
			},
			args: args{message: "bắn nè...chíu chíu"},
		},
		{
			name: "test publish direct exchange",
			exchange: exchange{
				name:       "test.exchange.direct",
				kind:       ExchangeDirect,
				routingKey: "key.direct",
			},
			args: args{message: "nhận lấy nè"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := r.NewDurableExchange(tt.exchange.name, tt.exchange.kind, tt.exchange.routingKey)
			if err != nil {
				logrus.WithError(err).Infof("New durable exchange fail, %s %s %s", tt.exchange.name, tt.exchange.kind, tt.exchange.routingKey)
			}
			if err := e.Publish(context.Background(), tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
