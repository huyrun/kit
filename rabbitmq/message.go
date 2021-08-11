package rabbitmq

const (
	ms  = 1
	s   = 1000
	min = 1000 * 60
)

type header map[string]interface{}

type body interface{}

type marshalFunc func(interface{}) ([]byte, error)

type handlerFunc func([]byte) error

type Msg struct {
	headers header

	Body       body
	RoutingKey string
	Priority   int
}

func (m *Msg) delay(d int) {
	if m.headers == nil {
		m.headers = make(header)
	}

	m.headers["x-delay"] = d
}

func (m *Msg) DelayMilisecond(d int) {
	m.delay(d * ms)
}

func (m *Msg) DelaySecond(d int) {
	m.delay(d * s)
}

func (m *Msg) DelayMinute(d int) {
	m.delay(d * min)
}
