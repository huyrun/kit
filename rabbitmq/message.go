package rabbitmq

const (
	ms  = 1
	s   = 1000
	min = 1000 * 60
)

type header map[HeaderType]interface{}

type body interface{}

type marshalFunc func(interface{}) ([]byte, error)

type handlerFunc func([]byte) error

type Msg interface {
	InitHeader()
	Header() header
	Body() body
	RoutingKey() string
	Priority() int
}

func DelayMilisecond(m Msg, d int) {
	delay(m, d*ms)
}

func DelaySecond(m Msg, d int) {
	delay(m, d*s)
}

func DelayMinute(m Msg, d int) {
	delay(m, d*min)
}
