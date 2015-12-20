package events

type Message interface {
	RoutingKey() string
	ContentType() string
	Body() []byte
}
