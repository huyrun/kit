package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Exchange interface {
	Publish(ctx context.Context, message interface{}) error
	PublishWithRetry(ctx context.Context, message interface{}, numOfRetries int64) error
}

type ExchangeInfo struct {
	Name       string
	Type       ExchangeType
	RoutingKey string
}

type exchange struct {
	name       string
	kind       ExchangeType
	routingKey string

	rabbitmq     *rabbitmq
	channel      *amqp.Channel
	errorChannel chan *amqp.Error
	closed       bool

	marshalFunc MarshalFunc
}

// open a channel
func (e *exchange) openChannel() error {
	if e.rabbitmq != nil {
		if e.rabbitmq.connection != nil {
			if !e.rabbitmq.connection.IsClosed() {
				for {
					channel, err := e.rabbitmq.connection.Channel()
					if err == nil {
						e.channel = channel

						e.closed = false
						e.errorChannel = make(chan *amqp.Error)
						e.channel.NotifyClose(e.errorChannel)
						go e.reopenChannel()
						return nil
					}

					logrus.WithError(err).Infof("Open exchange channel failed. Retrying in %s... ", retry_open_channel_after.String())
					time.Sleep(retry_open_channel_after)
				}
			}
		}
	}

	return nil
}

func (e *exchange) reopenChannel() {
	amqpErr := <-e.errorChannel

	if !e.closed {
		for e.rabbitmq.connection.IsClosed() {
			// Do not things until connection re-connected
		}

		time.Sleep(retry_open_channel_after)

		if amqpErr != nil {
			logrus.Infof("Reopen after channel closed by err: %v", amqpErr.Error())
		}

		log.Printf("Reopening channel for exchange...")
		err := e.openChannel()
		if err != nil {
			logrus.WithError(err).Infof("Reopen channel of exchange (%s) failed", e.name)
		}

		log.Printf("Recovery exchange...")
		if e.name != "" && e.kind != ExchangeDefault {
			err := e.declareDurableExchange()
			if err != nil {
				logrus.WithError(err).Infof("Redeclared exchange failed")
			}
		}
	}
}

func (e *exchange) reOpenChannelOnly() error {
	if !e.closed {
		e.openChannel()
	}

	return nil
}

func (e *exchange) declareExchange() error {
	err := e.channel.ExchangeDeclare(
		e.name,         // name
		string(e.kind), // kind
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)

	if err != nil {
		logrus.WithError(err).Infof("Declare exchange %s, kind %s failed", e.name, e.kind)
		return err
	}

	return nil
}

func (e *exchange) declareDurableExchange() error {
	err := e.channel.ExchangeDeclare(
		e.name,         // name
		string(e.kind), // kind
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)

	if err != nil {
		logrus.WithError(err).Infof("Declare exchange %s, kind %s failed", e.name, e.kind)
		return err
	}

	return nil
}

func (e *exchange) Close() {
	log.Println("Closing channel of exchange...")
	e.closed = true
	e.channel.Close()
}

func (e *exchange) Publish(ctx context.Context, message interface{}) error {
	msg, err := e.marshalFunc(message)
	if err != nil {
		logrus.WithError(err).Infof("Marshal message failed: %s", msg)
		return err
	}

	err = e.channel.Publish(
		e.name,       // exchange
		e.routingKey, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: 2,
			ContentType:  "application/json",
			Body:         msg,
		})
	if err != nil {
		logrus.WithError(err).Infof("Publish message to queue failed: %s", message)
		return err
	}

	return nil
}

// To retry forever, set numOfRetries a number which is less than 0
func (e *exchange) PublishWithRetry(ctx context.Context, message interface{}, numOfRetries int64) error {
	var b backoff.BackOff
	if numOfRetries < 0 {
		b = backoff.NewExponentialBackOff()
	} else {
		b = backoff.WithMaxRetries(backoff.NewExponentialBackOff(), uint64(numOfRetries))
	}

	err := backoff.RetryNotify(func() error {
		msg, err := e.marshalFunc(message)
		if err != nil {
			logrus.WithError(err).Infof("Marshal message failed: %s", message)
			return err
		}

		err = e.channel.Publish(
			e.name,       // exchange
			e.routingKey, // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				DeliveryMode: 2,
				ContentType:  "application/json",
				Body:         msg,
			})
		if err != nil {
			logrus.WithError(err).Infof("Publish message to queue failed: %s", message)
			return err
		}

		return nil
	}, b, func(err error, t time.Duration) {
		fmt.Printf("Rabbitmq publish fail err = %v, retry after %v, message: %s\n", err, t, message)
	})
	if err != nil {
		logrus.WithError(err).Infof("Publish message to queue failed: %s", message)
		return err
	}

	return nil
}
