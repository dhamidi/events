package events

type EventStore interface {
	LoadHistory(Aggregate) error
	Append(Event) error
	Subscribe(subscriber EventSubscriber) error
}

type EventStoreInMemory struct {
	events      []Event
	byAggregate map[string][]Event
	listeners   map[string][]EventHandler
}

func NewEventStoreInMemory() *EventStoreInMemory {
	return &EventStoreInMemory{
		events:      []Event{},
		byAggregate: map[string][]Event{},
		listeners:   map[string][]EventHandler{},
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

	return self.notify(event)
}

func (self *EventStoreInMemory) notify(event Event) error {
	eventName := event.EventName()

	for _, listener := range self.listeners[eventName] {
		listener.HandleEvent(event)
	}

	return nil
}

func (self *EventStoreInMemory) Subscribe(subscriber EventSubscriber) error {
	eventNames := subscriber.SubscribeTo()
	for _, eventName := range eventNames {
		self.listeners[eventName] = append(
			self.listeners[eventName],
			subscriber,
		)
	}

	return nil
}
