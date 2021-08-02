package rabbitmq

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type ChannOption func(*channel)

type channel struct {
	c      *amqp.Channel
	cError chan *amqp.Error

	exchange *exchange
	queue    *queue
}

func createChannel(conn *amqp.Connection) *amqp.Channel {
	if conn == nil {
		panic("Rabbitmq is not connected")
	}

	for {
		chann, err := conn.Channel()
		if err != nil {
			logrus.WithError(err).Infof("Open channel failed. Retrying in %s... ", retry_create_channel_after.String())
			time.Sleep(retry_create_channel_after)
			continue
		}

		return chann
	}
}
