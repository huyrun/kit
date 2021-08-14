package rabbitmq

type producer struct {
	c *channel
}

func (r *rabbitmq) CreateProducer(options ...ChannOption) *producer {
	options = append([]ChannOption{initExchange()}, options...)
	return &producer{c: r.newChannel(options...)}
}

func (p *producer) Publish(msg Message) error {
	return p.c.publish(msg)
}

func (p *producer) PublishWithRetry(msg Message, numsOfRetries int64) error {
	return p.c.publishWithRetry(msg, numsOfRetries)
}
