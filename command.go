package events

type Command interface {
	Execute() (Event, error)
	Aggregate() Aggregate
}

type CommandConstructor func() Command
