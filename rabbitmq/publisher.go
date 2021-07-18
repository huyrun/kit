package rabbitmq

import (
	"encoding/json"
	"errors"
)

// for publish message
type Publisher interface {
	Close()
	NewDurableExchange(exchangeName string, exchangeType ExchangeType, routingKey string) (Exchange, error)
}

func CreatePublisher(addr string) (Publisher, error) {
	r := new(rabbitmq)
	r.address = addr
	r.connect()

	return r, nil
}

func (r *rabbitmq) NewDurableExchange(exchangeName string, exchangeType ExchangeType, routingKey string) (Exchange, error) {
	e := new(exchange)

	if exchangeType != ExchangeDefault && exchangeName == "" {
		return nil, errors.New("Missing exchange name for non-default exchange")
	}
	if exchangeType == ExchangeDefault && routingKey == "" {
		return nil, errors.New("Missing routing key for default exchange")
	}

	r.exchangePool = append(r.exchangePool, e) // for recovery when connection lost

	e.name = exchangeName
	e.kind = exchangeType
	e.routingKey = routingKey
	e.rabbitmq = r
	e.marshalFunc = json.Marshal

	err := e.openChannel()
	if err != nil {
		return e, err
	}

	if e.name != "" && e.kind != ExchangeDefault {
		err = e.declareDurableExchange()
		if err != nil {
			return e, err
		}
	}

	return e, nil
}
