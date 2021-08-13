package rabbitmq

type Producer interface {
	Publish(message Message) error
	PublishWithRetry(message Message, numOfRetries int64) error
}

func (r *rabbitmq) CreateProducer(options ...ChannOption) Producer {
	options = append([]ChannOption{initExchange()}, options...)
	return r.newChannel(options...)
}
