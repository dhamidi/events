package sessions

import (
	"github.com/dhamidi/events"
	"github.com/dhamidi/events/examples/user"
)

type SessionsInMemory struct {
	ActiveSessions map[string]bool
}

func NewInMemory() *SessionsInMemory {
	return &SessionsInMemory{
		ActiveSessions: map[string]bool{},
	}
}

func (self *SessionsInMemory) SubscribeTo() []string {
	return []string{
		user.EventLoggedIn,
	}
}

func (self *SessionsInMemory) HandleEvent(event events.Event) error {
	switch e := event.(type) {
	case *user.LoggedIn:
		return self.addSession(e)
	}
	return nil
}

func (self *SessionsInMemory) IsLoggedIn(sessionId string) bool {
	return self.ActiveSessions[sessionId]
}

func (self *SessionsInMemory) addSession(loggedIn *user.LoggedIn) error {
	self.ActiveSessions[loggedIn.SessionId] = true
	return nil
}
