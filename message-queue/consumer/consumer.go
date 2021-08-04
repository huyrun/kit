package consumer

import (
	"fmt"

	"github.com/huypher/kit/log"

	"github.com/huypher/kit/util"

	message_queue "github.com/huypher/kit/message-queue"
)

type Handler func(msg []byte) error

type Consumer interface {
	Start()
	Bind(ch message_queue.Queue)
}

type consumer struct {
	handler Handler
	channel message_queue.Queue
	name    string
	id      string
}

type Option func(*consumer)

func NewConsumer(opts ...Option) *consumer {
	c := &consumer{
		id: util.RandomString(10),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func Name(name string) Option {
	return func(c *consumer) {
		c.name = name
	}
}

func HandlerFunc(handler Handler) Option {
	return func(c *consumer) {
		c.handler = handler
	}
}

func (c *consumer) Start() {
	log.Infof("Consumer %s@%s starting...", c.name, c.id)

	if c.handler == nil {
		panic(fmt.Sprintf("Consumer %s@%s don't have handler is nil", c.name, c.id))
	}

	if c.channel == nil {
		panic(fmt.Sprintf("Consumer %s@%s don't listen a topic", c.name, c.id))
	}

	go func() {
		log.Infof("Consumer %s@%s started", c.name, c.id)
		for {
			msg, ok := <-c.channel
			if ok {
				err := c.handler(msg)
				if err != nil {
					log.Error(err).Infof("Handle msg failed %s", string(msg))
					go c.requeue(msg)
				}
			} else {
				break
			}
		}
	}()
}

func (c *consumer) Bind(ch message_queue.Queue) {
	c.channel = ch
}

func (c *consumer) requeue(msg []byte) {
	for {
		if len(c.channel) < message_queue.QueueSizeDefault {
			c.channel <- msg
			break
		}
	}
}
