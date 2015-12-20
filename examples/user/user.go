package user

const (
	EventSignedUp = "user.signed-up"
)

type User struct {
	Username Username
	SignedUp bool
}

func NewUser(username Username) *User {
	return &User{
		Username: username,
	}
}

func (self *User) AggregateId() string {
	return self.Username.String()
}
