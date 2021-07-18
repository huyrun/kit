package rabbitmq

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

const (
	url = "amqp://guest:guest@localhost:5672"
)

func TestCreatePublisher(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		args    args
		want    Publisher
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "test create publisher",
			args:    args{addr: url},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreatePublisher(tt.args.addr)
			defer got.Close()
			//// for test reconnect rabbitmq
			//for {
			//	fmt.Println("running...")
			//	time.Sleep(2 * time.Second)
			//}
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePublisher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreatePublisher() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rabbitmq_NewExchange(t *testing.T) {
	type args struct {
		exchangeName string
		exchangeType ExchangeType
		routingKey   string
	}

	r, err := CreatePublisher(url)
	if err != nil {
		logrus.WithError(err).Infof("Create publihser fail: %s", url)
	}
	defer r.Close()

	tests := []struct {
		name    string
		args    args
		want    Exchange
		wantErr bool
	}{
		// TODO: Add test cases.
		//{
		//	name: "test new default exchange",
		//	args: args{
		//		exchangeName: "",
		//		exchangeType: ExchangeDefault,
		//		routingKey:   "test.exchange",
		//	},
		//	wantErr: true,
		//},
		//{
		//	name: "test new direct exchange",
		//	args: args{
		//		exchangeName: "test.exchange.direct",
		//		exchangeType: ExchangeDirect,
		//		routingKey:   "key.direct",
		//	},
		//	wantErr: true,
		//},
		//{
		//	name: "test new topic exchange",
		//	args: args{
		//		exchangeName: "test.exchange.topic",
		//		exchangeType: ExchangeTopic,
		//		routingKey:   "key.direct",
		//	},
		//	wantErr: true,
		//},
		{
			name: "test new fanout exchange",
			args: args{
				exchangeName: "test.exchange.fanout",
				exchangeType: ExchangeFanout,
				routingKey:   "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.NewDurableExchange(tt.args.exchangeName, tt.args.exchangeType, tt.args.routingKey)
			//// for test reopen channel by publish message when lost connection
			//for i := 0; i < 50; i++ {
			//	fmt.Println("running...", i)
			//	time.Sleep(2 * time.Second)
			//}
			//got.Publish("bằng bằng bằng")
			//for {
			//	fmt.Println("rerunning...")
			//	time.Sleep(2 * time.Second)
			//}
			if (err != nil) != tt.wantErr {
				t.Errorf("NewExchange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExchange() got = %v, want %v", got, tt.want)
			}
		})
	}
}
