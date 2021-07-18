package rabbitmq

import (
	"context"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	delay_nack = 500
)

type Queue interface {
	Consume(ctx context.Context, handler handlerFunc, key string)
	BindToExchange(exchangeInfo ExchangeInfo) error
	Close()
}

type handlerFunc func(context.Context, []byte) error

//type Tracer interface {
//	Transaction(ctx context.Context, name string) (context.Context, func(err error))
//	Segment(ctx context.Context, name string) (context.Context, func(err error))
//}

type queue struct {
	ctx  context.Context
	name string
	tag  string

	rabbitmq     *rabbitmq
	channel      *amqp.Channel
	errorChannel chan *amqp.Error
	closed       bool

	handler handlerFunc
	//tracer  Tracer

	binding map[string]ExchangeInfo
}

func (q *queue) Consume(ctx context.Context, handler handlerFunc, key string) {
	q.ctx = ctx
	q.tag = key
	deliveries, err := q.registerQueueConsumer()
	if err != nil {
		logrus.WithError(err).Infof("Register queue consumer failed")
	}

	q.handler = handler

	q.executeMessageConsumer(deliveries)
}

func (q *queue) BindToExchange(exchangeInfo ExchangeInfo) error {
	err := q.channel.QueueBind(
		q.name,                  // queue name
		exchangeInfo.RoutingKey, // routing key
		exchangeInfo.Name,       // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	q.binding[exchangeInfo.Name] = exchangeInfo

	return nil
}

func (q *queue) Close() {
	log.Println("Closing channel of queue...")
	q.closed = true
	q.channel.Close()
}

// open a channel
func (q *queue) openChannel() error {
	if q.rabbitmq == nil {
		return nil
	}

	if q.rabbitmq.connection == nil {
		return nil
	}

	if q.rabbitmq.connection.IsClosed() {
		return nil
	}

	for {
		channel, err := q.rabbitmq.connection.Channel()
		if err != nil {
			logrus.WithError(err).Infof("Open queue channel failed. Retrying in %s... ", retry_open_channel_after.String())
			time.Sleep(retry_open_channel_after)
			continue
		}

		q.channel = channel

		q.closed = false
		q.errorChannel = make(chan *amqp.Error)
		q.channel.NotifyClose(q.errorChannel)
		go q.reopenChannel()
		return nil
	}
}

func (q *queue) reopenChannel() {
	amqpErr := <-q.errorChannel

	if !q.closed {
		for q.rabbitmq.connection.IsClosed() {
			// Do not things until connection re-connected
		}

		time.Sleep(retry_open_channel_after)

		if amqpErr != nil {
			logrus.Infof("Reopen after channel closed by err: %v", amqpErr.Error())
		}

		log.Printf("Reopening channel for queue...")
		err := q.openChannel()
		if err != nil {
			logrus.WithError(err).Infof("Reopen channel of queue %s failed", q.name)
		}

		log.Printf("Recovery queue...")
		err = q.declareDurableQueue()
		if err != nil {
			logrus.WithError(err).Infof("Redeclared queue %s failed", q.name)
		}

		log.Printf("Re-bind queue to exchange...")
		for _, e := range q.binding {
			err = q.bind(e)
			if err != nil {
				logrus.WithError(err).Infof("Bind queue %s to exchange %s with routingKey %s failed", q.name, e.Name, e.RoutingKey)
			}
		}

		log.Printf("Recovery consumer...")
		deliveries, err := q.registerQueueConsumer()
		if err != nil {
			logrus.WithError(err).Infof("Register queue consumer failed")
		}

		log.Printf("Recovery execute consumer...")
		q.executeMessageConsumer(deliveries)

	}
}

// for connection call to reopen channel when lost rabbibtqm connection
func (q *queue) reOpenChannelOnly() error {
	if !q.closed {
		err := q.openChannel()
		if err != nil {
			logrus.WithError(err).Infof("Open channel failed")
			return err
		}
	}

	return nil
}

func (q *queue) declareDurableQueue() error {
	_, err := q.channel.QueueDeclare(
		q.name, // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)

	if err != nil {
		logrus.WithError(err).Infof("Queue declaration failed: %s", q.name)
		return err
	}

	err = q.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	return nil
}

func (q *queue) registerQueueConsumer() (<-chan amqp.Delivery, error) {
	msgs, err := q.channel.Consume(
		q.name, // queue
		q.tag,  // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logrus.WithError(err).Infof("Consuming messages from queue %s failed", q.name)
	}

	return msgs, err
}

func (q *queue) executeMessageConsumer(deliveries <-chan amqp.Delivery) {
	for {
		for delivery := range deliveries {
			logrus.Infof("Rabbitmq consume value: %v", string(delivery.Body))
			err := q.handler(q.ctx, delivery.Body)
			if err != nil {
				logrus.WithError(err).Infof("handle message error at queue (%s) %s", q.name, string(delivery.Body))
				time.Sleep(delay_nack * time.Millisecond)
				delivery.Nack(false, true)
			} else {
				delivery.Ack(false)
			}
		}
	}
}

func (q *queue) bind(exchangeInfo ExchangeInfo) error {
	return q.channel.QueueBind(
		q.name,                  // queue name
		exchangeInfo.RoutingKey, // routing key
		exchangeInfo.Name,       // exchange
		false,
		nil,
	)
}
