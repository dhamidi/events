package events

type MessageBus interface {
	HandleMessage(routingKey string, handler MessageHandler) MessageBus
}

func HandleMessageFunc(bus MessageBus, routingKey string, handler func(request Message) error) MessageBus {
	return bus.HandleMessage(routingKey, MessageHandlerFunc(handler))
}

type MessageHandlerFunc func(request Message) error

func (self MessageHandlerFunc) HandleMessage(request Message) error {
	return self(request)
}

type MessageHandler interface {
	HandleMessage(request Message) error
}

type Message interface {
	RoutingKey() string
	ContentType() string
	Body() []byte
	Header(string) string
	Acknowledge(event Event) error
	Reject(err error) error
}
