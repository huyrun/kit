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
	Publish(ctx context.Context, message Msg) error
	PublishWithRetry(ctx context.Context, message Msg, numOfRetries int64) error
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
	if e.rabbitmq == nil {
		return nil
	}

	if e.rabbitmq.connection == nil {
		return nil
	}

	if e.rabbitmq.connection.IsClosed() {
		return nil
	}

	for {
		channel, err := e.rabbitmq.connection.Channel()
		if err != nil {
			logrus.WithError(err).Infof("Open exchange channel failed. Retrying in %s... ", retry_open_channel_after.String())
			time.Sleep(retry_open_channel_after)
			continue
		}

		e.channel = channel

		e.closed = false
		e.errorChannel = make(chan *amqp.Error)
		e.channel.NotifyClose(e.errorChannel)
		go e.reopenChannel()
		return nil
	}
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

func (e *exchange) Publish(ctx context.Context, message Msg) error {
	msg, err := e.marshalFunc(message.Body)
	if err != nil {
		logrus.WithError(err).Infof("Marshal message failed: %v", message.Body)
		return err
	}

	headers := amqp.Table{}
	for k, v := range message.Headers {
		headers[k] = v
	}

	err = e.channel.Publish(
		e.name,       // exchange
		e.routingKey, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Headers:      headers,
			DeliveryMode: 2,
			ContentType:  "application/json",
			Body:         msg,
		})
	if err != nil {
		logrus.WithError(err).Infof("Publish message to queue failed: %v", message.Body)
		return err
	}

	return nil
}

// To retry forever, set numOfRetries a number which is less than 0
func (e *exchange) PublishWithRetry(ctx context.Context, message Msg, numOfRetries int64) error {
	var b backoff.BackOff
	if numOfRetries < 0 {
		b = backoff.NewExponentialBackOff()
	} else {
		b = backoff.WithMaxRetries(backoff.NewExponentialBackOff(), uint64(numOfRetries))
	}

	err := backoff.RetryNotify(func() error {
		return e.Publish(ctx, message)
	}, b, func(err error, t time.Duration) {
		fmt.Printf("Rabbitmq publish fail err = %v, retry after %v, message: %s\n", err, t, message)
	})
	if err != nil {
		logrus.WithError(err).Infof("Publish message to queue failed: %s", message)
		return err
	}

	return nil
}
