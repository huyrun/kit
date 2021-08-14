package rabbitmq

type consumer struct {
	c *channel
}

func (r *rabbitmq) CreateConsumer(options ...ChannOption) *consumer {
	options = append([]ChannOption{initQueue()}, options...)
	return &consumer{c: r.newChannel(options...)}
}

func (c *consumer) Consume() {
	c.c.consume()
}

func (c *consumer) Bind(exchangeName, routingKey string) error {
	return c.c.bind(exchangeName, routingKey)
}
