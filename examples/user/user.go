package user

const (
	EventSignedUp = "user.signed-up"
	EventLoggedIn = "user.logged-in"
)

type User struct {
	Username     Username
	SignedUp     bool
	PasswordHash []byte
}

func NewUser(username Username) *User {
	return &User{
		Username: username,
	}
}

func (self *User) AggregateId() string {
	return self.Username.String()
}
