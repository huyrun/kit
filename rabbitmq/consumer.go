package rabbitmq

// for consume message
type Consumer interface {
	Close()
	NewDurableQueue(queueName string) (Queue, error)
}

func CreateConsumer(addr string) (Consumer, error) {
	r := new(rabbitmq)
	r.address = addr
	r.connect()

	return r, nil
}

func (r *rabbitmq) NewDurableQueue(queueName string) (Queue, error) {
	q := new(queue)

	q.rabbitmq = r

	q.binding = make(map[string]ExchangeInfo)

	r.queuePool = append(r.queuePool, q) // for recovery if connection lost

	q.name = queueName

	err := q.openChannel()
	if err != nil {
		return q, err
	}

	err = q.declareDurableQueue()
	if err != nil {
		return q, err
	}

	return q, nil
}
