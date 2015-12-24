package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

var (
	ErrNotAcceptable  = errors.New("not acceptable")
	ErrUnknownCommand = errors.New("unknown command")
)

type Application struct {
	commands   map[string]CommandConstructor
	EventStore EventStore
}

func NewApplication() *Application {
	return &Application{
		EventStore: NewEventStoreInMemory(),
		commands:   map[string]CommandConstructor{},
	}
}

func (self *Application) RegisterCommand(commandName string, makeCommand CommandConstructor) *Application {
	self.commands[commandName] = makeCommand
	return self
}

func (self *Application) HandleCommand(message Message) error {
	log.Printf("HandleCommand: %s", message)
	commandName := message.RoutingKey()
	command, err := self.NewCommand(commandName)
	if err != nil {
		return err
	}

	if err := self.Unmarshal(message, command); err != nil {
		return err
	}

	aggregate := command.Aggregate()
	if err := self.EventStore.LoadHistory(aggregate); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "DEBUG:\n%s\n", (func() []byte { data, _ := json.Marshal(aggregate); return data })())
	event, err := command.Execute()
	if err != nil {
		return err
	}

	if err := self.EventStore.Append(event); err != nil {
		return err
	}

	return nil

}

func (self *Application) NewCommand(commandName string) (Command, error) {
	constructor, found := self.commands[commandName]
	if !found {
		return nil, ErrUnknownCommand
	}

	return constructor(), nil
}

func (self *Application) Unmarshal(message Message, dest interface{}) error {
	contentType := message.ContentType()
	body := message.Body()
	codec, err := self.CodecFor(contentType)
	if err != nil {
		return err
	}

	return codec.Decode(body, dest)
}

func (self *Application) CodecFor(contentType string) (*Codec, error) {
	switch contentType {
	case "application/json":
		return JSONCodec, nil
	}

	return nil, ErrNotAcceptable
}
