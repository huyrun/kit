package rabbitmq

import (
	"log"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

const (
	retry_connect_after      = 1 * time.Second // 1 sec
	retry_open_channel_after = 1 * time.Second
)

type MarshalFunc func(interface{}) ([]byte, error)

type ExchangeType string

const (
	ExchangeDefault ExchangeType = ""
	ExchangeDirect  ExchangeType = amqp.ExchangeDirect
	ExchangeFanout  ExchangeType = amqp.ExchangeFanout
	ExchangeTopic   ExchangeType = amqp.ExchangeTopic
	ExchangeHeaders ExchangeType = amqp.ExchangeHeaders
)

type Header map[string]interface{}

type Body interface{}

type Msg struct {
	Headers Header
	Body    Body
}

// keep rabbitmq connect
type rabbitmq struct {
	address         string
	connection      *amqp.Connection
	errorConnection chan *amqp.Error
	closed          bool

	exchangePool []*exchange
	queuePool    []*queue
}

// create a connect to rabbitmq server
func (r *rabbitmq) connect() {
	for {
		conn, err := amqp.Dial(r.address)

		if err == nil {
			r.connection = conn
			r.errorConnection = make(chan *amqp.Error)
			r.connection.NotifyClose(r.errorConnection)
			go r.reconnector()
			log.Println("Rabbitmq connection established!")
			return
		}

		logrus.WithError(err).Infof("Connection to rabbitmq failed. Retrying in %s... ", retry_connect_after.String())
		time.Sleep(retry_connect_after)
	}
}

func (r *rabbitmq) reconnector() {
	amqpErr := <-r.errorConnection
	if !r.closed {
		if amqpErr != nil {
			logrus.Infof("Reconnecting after connection closed by err: %v", amqpErr.Error())
		}

		// recovery connection
		log.Printf("Reconnecting...")
		r.connect()
	}
}

func (r *rabbitmq) Close() {
	log.Println("Closing rabbitmq connection")
	r.closed = true
	r.connection.Close()
}
