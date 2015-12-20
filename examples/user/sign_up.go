package user

import (
	"github.com/dhamidi/events"
	"github.com/dhamidi/events/value"
	"golang.org/x/crypto/bcrypt"
)

type SignUp struct {
	Username Username
	Password string
	Email    value.Email

	User *User

	Crypt func(string) ([]byte, error)
}

func NewSignUp() *SignUp {
	return &SignUp{
		Crypt: cryptWithBcrypt,
	}
}

func (self *SignUp) Aggregate() events.Aggregate {
	self.User = NewUser(self.Username)
	return self.User
}

func cryptWithBcrypt(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (self *SignUp) Execute() (*SignedUp, error) {
	hashedPassword, err := self.Crypt(self.Password)
	if err != nil {
		return nil, events.NewInternalError(err)
	}

	return &SignedUp{
		Username:     self.Username.String(),
		PasswordHash: hashedPassword,
		Email:        self.Email.String(),
	}, nil
}

type SignedUp struct {
	Username     string
	PasswordHash []byte
	Email        string
}

func (self *SignedUp) AggregateId() string {
	return self.Username
}

func (self *SignedUp) EventName() string {
	return EventSignedUp
}

func (self *SignedUp) Apply(to events.Aggregate) error {
	user, ok := to.(*User)
	if !ok {
		return nil
	}
	user.SignedUp = true
	return nil
}
