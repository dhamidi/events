package events

import "errors"

var (
	ErrNotAcceptable = errors.New("not acceptable")
)

type Application struct{}

func (self *Application) HandleCommand(message Message) error {
	commandName := message.RoutingKey()
	command, err := self.NewCommand(commandName)
	if err != nil {
		return err
	}

	if err := self.Unmarshal(message, command); err != nil {
		return err
	}

	return nil

}

func (self *Application) NewCommand(commandName string) (Command, error) {
	return (Command)(nil), nil
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
