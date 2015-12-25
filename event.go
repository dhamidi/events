package events

type Event interface {
	EventName() string
	AggregateId() string
	Apply(Aggregate) error
}

type EventHandler interface {
	HandleEvent(event Event) error
}

type EventSubscriber interface {
	EventHandler
	SubscribeTo() []string
}
