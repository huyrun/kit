package rabbitmq

type Consumer interface {
	Consume()
	Bind(exchangeName, routingKey string) error
}

func (r *rabbitmq) CreateConsumer(options ...ChannOption) Consumer {
	options = append([]ChannOption{initQueue()}, options...)
	return r.newChannel(options...)
}
