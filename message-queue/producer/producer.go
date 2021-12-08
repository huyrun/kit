package producer

import (
	"encoding/json"
	"github.com/huypher/kit/utils"
	"sync"

	message_queue "github.com/huypher/kit/message-queue"

	"github.com/huypher/kit/message-queue/consumer"
)

type MarshalFunc func(interface{}) ([]byte, error)

type Producer interface {
	Subscribe(consumers []consumer.Consumer)
	Publish(msg interface{}) error
}

type producer struct {
	mu          sync.RWMutex
	id          string
	topic       string
	channels    []message_queue.Queue
	next        int
	marshalFunc MarshalFunc
}

type Option func(*producer)

func NewProducer(otps ...Option) *producer {
	p := &producer{
		id:          utils.RandomString(10),
		channels:    make([]message_queue.Queue, 0),
		next:        -1,
		marshalFunc: json.Marshal,
	}

	for _, opt := range otps {
		opt(p)
	}

	return p
}

func Topic(topic string) Option {
	return func(p *producer) {
		p.topic = topic
	}
}

func (ps *producer) Subscribe(consumers []consumer.Consumer) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, consumer := range consumers {
		ch := make(message_queue.Queue, message_queue.QueueSizeDefault)
		consumer.Bind(ch)
		ps.channels = append(ps.channels, ch)
	}
}

func (ps *producer) Publish(msg interface{}) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	bytesMsg, err := ps.marshalFunc(msg)
	if err != nil {
		return err
	}

	numsOfConsumer := len(ps.channels)
	if numsOfConsumer > 0 {
		ps.next = (ps.next + 1) % numsOfConsumer
		ps.channels[ps.next] <- bytesMsg
	}

	return nil
}
