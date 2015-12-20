package events

import "fmt"

type InternalError struct {
	Internal error
}

func NewInternalError(err error) *InternalError {
	return &InternalError{
		Internal: err,
	}
}

func (self *InternalError) Error() string {
	return fmt.Sprintf("internal error: %s", self.Internal)
}
