package events

type EventStore interface {
	LoadHistory(Aggregate) error
	Append(Event) error
}

type EventStoreInMemory struct {
	events      []Event
	byAggregate map[string][]Event
}

func NewEventStoreInMemory() *EventStoreInMemory {
	return &EventStoreInMemory{
		events:      []Event{},
		byAggregate: map[string][]Event{},
	}
}

func (self *EventStoreInMemory) LoadHistory(aggregate Aggregate) error {
	for _, event := range self.byAggregate[aggregate.AggregateId()] {
		event.Apply(aggregate)
	}

	return nil
}

func (self *EventStoreInMemory) Append(event Event) error {
	self.events = append(self.events, event)
	self.byAggregate[event.AggregateId()] = append(
		self.byAggregate[event.AggregateId()],
		event,
	)

	return nil
}
