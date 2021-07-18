package message_queue

const (
	QueueSizeDefault = 10000
)

type Queue chan []byte
