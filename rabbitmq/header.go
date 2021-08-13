package rabbitmq

const (
	XDelayHeader = "x-delay"
)

func setHeader(msg Message, headers ...header) {
	if msg.MessageHeader() == nil {
		return
	}

	for _, header := range headers {
		for k, v := range header {
			msg.MessageHeader()[k] = v
		}
	}
}

func delHeader(msg Message, headerType string) {
	if msg.MessageHeader() == nil {
		return
	}

	delete(msg.MessageHeader(), headerType)
}
