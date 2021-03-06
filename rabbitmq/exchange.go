package rabbitmq

import (
	"time"

	"github.com/huypher/kit/log"

	"github.com/cenkalti/backoff"
	"github.com/streadway/amqp"
)

const (
	xDelayedMessage = "x-delayed-message"
	xDelayedType    = "x-delayed-type"
)

type exchange struct {
	name          string
	kind          ExchangeType
	args          amqp.Table
	isDelayedType bool
	marshalFunc   marshalFunc
}

func initExchange() ChannOption {
	return func(c *channel) {
		c.exchange = new(exchange)
		c.exchange.args = make(amqp.Table)
	}
}

func ExchangeName(name string) ChannOption {
	return func(c *channel) {
		c.exchange.name = name
	}
}

func ExchangeKind(kind ExchangeType) ChannOption {
	return func(c *channel) {
		c.exchange.kind = kind
	}
}

func RegisterMarshalFunc(marshalFunc marshalFunc) ChannOption {
	return func(c *channel) {
		c.exchange.marshalFunc = marshalFunc
	}
}

func ExchangeDelayedType() ChannOption {
	return func(c *channel) {
		c.exchange.isDelayedType = true
	}
}

func (c *channel) createExchange() error {
	e := c.exchange

	exchangeType := string(e.kind)
	if e.isDelayedType {
		exchangeType = xDelayedMessage
		e.args[xDelayedType] = string(e.kind)
	}

	err := c.c.ExchangeDeclare(
		e.name,
		exchangeType,
		true,
		true,
		false,
		false,
		e.args,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *channel) publish(message Message) error {
	if c.exchange == nil {
		panic("can not use this channel to publish msg")
	}

	e := c.exchange

	msgBody := message.MessageBody()
	msgHeader := message.MessageHeader()
	msgRoutingKey := message.MessageRoutingKey()
	msgPriority := message.MessagePriority()

	msg, err := e.marshalFunc(msgBody)
	if err != nil {
		log.Error(err).Infof("Marshal message failed: %v", msgBody)
		return err
	}

	headers := amqp.Table{}
	for k, v := range msgHeader {
		headers[k] = v
	}

	err = c.c.Publish(
		e.name,        // exchange
		msgRoutingKey, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			Headers:      headers,
			DeliveryMode: 2,
			ContentType:  "application/json",
			Body:         msg,
			Priority:     uint8(msgPriority),
		})
	if err != nil {
		log.Error(err).Infof("Publish message to exchange (%s) failed: %v", c.exchange.name, msgBody)
		return err
	}

	return nil
}

// To retry forever, set numOfRetries a number which is less than 0
func (c *channel) publishWithRetry(message Message, numOfRetries int64) error {
	var b backoff.BackOff
	if numOfRetries < 0 {
		b = backoff.NewExponentialBackOff()
	} else {
		b = backoff.WithMaxRetries(backoff.NewExponentialBackOff(), uint64(numOfRetries))
	}

	err := backoff.RetryNotify(func() error {
		return c.publish(message)
	}, b, func(err error, t time.Duration) {
		log.Infof("Rabbitmq publish fail err = %v, retry after %v, message: %v\n", err, t, message)
	})
	if err != nil {
		log.Error(err).Infof("Publish message to queue failed: %v", message)
		return err
	}

	return nil
}
