package events

type Message interface {
	RoutingKey() string
	ContentType() string
	Body() []byte
	Cookie(string) string
	Acknowledge(event Event) error
	Reject(err error) error
}
