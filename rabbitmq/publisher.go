package rabbitmq

type Producer interface {
	Publish(message Msg) error
	PublishWithRetry(message Msg, numOfRetries int64) error
}

func (r *rabbitmq) CreatePublisher(options ...ChannOption) Producer {
	options = append([]ChannOption{initExchange()}, options...)
	return r.newChannel(options...)
}
