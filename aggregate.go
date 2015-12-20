package events

type Aggregate interface {
	AggregateId() string
}
