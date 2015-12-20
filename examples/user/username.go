package user

import (
	"errors"
	"strings"
)

type Username string

var (
	ErrEmpty = errors.New("empty")
)

func (self *Username) UnmarshalText(src []byte) error {
	name := strings.TrimSpace(string(src))
	if len(name) == 0 {
		return ErrEmpty
	}

	*self = (Username)(name)
	return nil
}

func (self Username) MarshalText() ([]byte, error) {
	return []byte(self), nil
}

func (self Username) String() string {
	return string(self)
}
