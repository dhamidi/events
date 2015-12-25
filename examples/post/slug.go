package post

import (
	"encoding"
	"errors"
	"regexp"
)

var (
	SlugRegexp   = regexp.MustCompile(`[A-Za-z][-A-Za-z0-9]*`)
	ErrMalformed = errors.New("malformed")
)

type Slug string

func (self *Slug) UnmarshalText(src []byte) error {
	if !SlugRegexp.Match(src) {
		return ErrMalformed
	}

	*self = Slug(src)
	return nil
}

func (self Slug) MarshalText() ([]byte, error) {
	return []byte(self), nil
}

func (self Slug) String() string {
	return string(self)
}

var (
	_ encoding.TextUnmarshaler = new(Slug)
)
