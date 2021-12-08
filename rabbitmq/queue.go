package rabbitmq

import (
	"context"
	"fmt"
	"github.com/huypher/kit/utils"
	"time"

	"github.com/huypher/kit/log"

	"github.com/streadway/amqp"
)

const (
	delay_nack = 500

	xMaxPriority = "x-max-priority"
)

type queue struct {
	ctx  context.Context
	name string
	args amqp.Table

	handlerFunc handlerFunc
}

func initQueue() ChannOption {
	return func(c *channel) {
		c.queue = new(queue)
		c.queue.args = make(amqp.Table)
	}
}

func QueueName(name string) ChannOption {
	return func(c *channel) {
		c.queue.name = name
	}
}

func PriorityQueue(priorityRange int) ChannOption {
	return func(c *channel) {
		c.queue.args[xMaxPriority] = priorityRange
	}
}

func RegisterHandlerFunc(handlerFunc handlerFunc) ChannOption {
	return func(c *channel) {
		c.queue.handlerFunc = handlerFunc
	}
}

func (c *channel) consume() {
	go c.executeMessageConsumer(c.registerQueueConsumer())
}

func (c *channel) registerQueueConsumer() <-chan amqp.Delivery {
	q := c.queue
	if q == nil {
		panic("queue is nil")
	}

	msgs, err := c.c.Consume(
		q.name,                // queue
		genConsumerID(q.name), // consumer
		false,                 // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	if err != nil {
		panic(fmt.Sprintf("Consuming messages from queue %s failed, err: %s", q.name, err.Error()))
	}

	return msgs
}

func (c *channel) executeMessageConsumer(deliveries <-chan amqp.Delivery) {
	q := c.queue
	if q == nil {
		panic("queue is nil")
	}

	for {
		for delivery := range deliveries {
			err := q.handlerFunc(delivery.Body)
			if err != nil {
				log.Error(err).Infof("handle message error at queue (%s) %s", q.name, string(delivery.Body))
				time.Sleep(delay_nack * time.Millisecond)
				delivery.Nack(false, true)
				continue
			}
			delivery.Ack(false)
		}
	}
}

func (c *channel) createQueue() error {
	q := c.queue
	if q == nil {
		panic("queue is nil")
	}

	_, err := c.c.QueueDeclare(
		q.name, // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		q.args, // arguments
	)

	if err != nil {
		log.Error(err).Infof("Queue declaration failed: %s", q.name)
		return err
	}

	err = c.c.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	return nil
}

func (c *channel) bind(exchangeName, routingKey string) error {
	if c.queue == nil {
		panic("queue is nil")
	}

	return c.c.QueueBind(
		c.queue.name, // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
}

func genConsumerID(queueName string) string {
	return fmt.Sprintf("%s-%s", queueName, utils.RandomString(10))
}
