package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

const (
	retry_connect_after        = 1 * time.Second // 1 sec
	retry_create_channel_after = 1 * time.Second
)

type Rabbitmq interface {
	Close()
}

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

type marshalFunc func(interface{}) ([]byte, error)

type handlerFunc func([]byte) error

type Msg struct {
	Headers    Header
	Body       Body
	RoutingKey string
	Priority   int
}

// keep rabbitmq connect
type rabbitmq struct {
	address         string
	connection      *amqp.Connection
	errorConnection chan *amqp.Error
	channels        []*channel
}

type Option func(*rabbitmq)

func NewRabbitmq(options ...Option) *rabbitmq {
	r := new(rabbitmq)

	for _, o := range options {
		o(r)
	}

	r.connect()
	go r.channelsFaultTolerance()

	return r
}

func Address(addr string) Option {
	return func(r *rabbitmq) {
		r.address = addr
	}
}

// create a connect to rabbitmq server
func (r *rabbitmq) connect() {
	for {
		conn, err := amqp.Dial(r.address)
		if err != nil {
			logrus.WithError(err).Infof("Connection to rabbitmq failed. Retrying in %s... ", retry_connect_after.String())
			time.Sleep(retry_connect_after)
			continue
		}

		r.connection = conn
		r.errorConnection = make(chan *amqp.Error)
		r.connection.NotifyClose(r.errorConnection)
		go r.reconnector()
		log.Println("Rabbitmq connection established!")

		return
	}
}

func (r *rabbitmq) reconnector() {
	amqpErr := <-r.errorConnection
	if amqpErr != nil {
		logrus.Infof("Reconnecting after connection closed by err: %v", amqpErr.Error())
	}

	// recovery connection
	r.connection.Close()
	log.Printf("Reconnecting...")
	r.connect()
}

func (r *rabbitmq) Close() {
	log.Println("Closing rabbitmq connection")
	for _, c := range r.channels {
		c.c.Close()
	}
	r.connection.Close()
}

func (r *rabbitmq) newChannel(options ...ChannOption) *channel {
	channel := &channel{
		cError: make(chan *amqp.Error),
	}

	for _, o := range options {
		o(channel)
	}

	channel.c = createChannel(r.connection)
	channel.c.NotifyClose(channel.cError)

	if channel.exchange != nil && channel.exchange.kind != ExchangeDefault {
		err := channel.createExchange()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	if channel.queue != nil {
		err := channel.createQueue()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	r.channels = append(r.channels, channel)
	return channel
}

func (r *rabbitmq) channelsFaultTolerance() {
	for {
		for idx := range r.channels {
			c := r.channels[idx]
			select {
			case amqpErr := <-r.channels[idx].cError:
				if amqpErr != nil {
					for r.connection.IsClosed() {
						continue
					}
					logrus.Infof("Re-create after channel closed by err: %v", amqpErr.Error())
					r.resetup(c)
				}
			default:
				continue
			}
		}
	}
}

func (r *rabbitmq) resetup(channel *channel) {
	channel.c = createChannel(r.connection)
	channel.cError = make(chan *amqp.Error)
	channel.c.NotifyClose(channel.cError)

	if channel.queue != nil {
		channel.Consume()
	}
}
