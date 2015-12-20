package value

import "net/mail"

type Email struct {
	*mail.Address
}

var (
	NullEmail = &Email{
		Address: &mail.Address{
			Name:    "John Doe",
			Address: "john.doe@example.com",
		},
	}
)

func (self *Email) UnmarshalText(src []byte) error {
	addr, err := mail.ParseAddress(string(src))
	if err != nil {
		return err
	}

	self.Address = addr
	return nil
}

func (self Email) MarshalText() ([]byte, error) {
	return []byte(self.Address.String()), nil
}
