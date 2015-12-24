package user

import (
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/dhamidi/events"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrLoginFailed = errors.New("login failed")
)

type LogIn struct {
	User *User

	Username Username
	Password string
}

func (self *LogIn) Aggregate() events.Aggregate {
	self.User = NewUser(self.Username)
	return self.User
}

func (self *LogIn) Execute() (events.Event, error) {
	if !self.User.SignedUp {
		return nil, ErrLoginFailed
	}

	if err := bcrypt.CompareHashAndPassword(self.User.PasswordHash, []byte(self.Password)); err != nil {
		return nil, ErrLoginFailed
	}

	randomBytes := make([]byte, 10)
	if _, err := rand.Read(randomBytes); err != nil {
		panic(err)
	}
	sessionId := fmt.Sprintf("%x", randomBytes)

	return &LoggedIn{
		SessionId: sessionId,
		Username:  self.User.Username.String(),
	}, nil
}

type LoggedIn struct {
	SessionId string
	Username  string
}

func (self *LoggedIn) EventName() string {
	return EventLoggedIn
}

func (self *LoggedIn) AggregateId() string {
	return self.Username
}

func (self *LoggedIn) Apply(events.Aggregate) error {
	return nil
}
