package rabbitmq

type HeaderType string

func (h HeaderType) String() string {
	return string(h)
}

const (
	XDelayHeader HeaderType = "x-delay"
)

func setHeader(msg Msg, headers ...header) {
	if msg.Header() == nil {
		return
	}

	for _, header := range headers {
		for k, v := range header {
			msg.Header()[k] = v
		}
	}
}

func delay(m Msg, d int) {
	setHeader(m, header{
		XDelayHeader: d,
	})
}
