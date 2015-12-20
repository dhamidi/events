package events

type Event interface {
	EventName() string
	AggregateId() string
	Apply(Aggregate) error
}
