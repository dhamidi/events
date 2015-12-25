package post

import (
	"errors"
	"strings"
	"time"

	"github.com/dhamidi/events"
)

type Publish struct {
	Slug Slug
	post *Post
}

var (
	EventPublished     = "post.published"
	ErrTitleTooShort   = errors.New("title too short")
	ErrContentTooShort = errors.New("content too short")
)

func NewPublish() *Publish {
	return &Publish{}
}

func (self *Publish) Execute() (events.Event, error) {
	trimmedTitle := strings.TrimSpace(self.post.Title)
	trimmedContent := strings.TrimSpace(self.post.Content)
	if len(trimmedTitle) == 0 {
		return nil, ErrTitleTooShort
	}

	if len(trimmedContent) == 0 {
		return nil, ErrContentTooShort
	}

	return &Published{
		Slug:        self.Slug.String(),
		PublishedAt: Now(),
	}, nil
}

func (self *Publish) Aggregate() events.Aggregate {
	self.post = NewPost(self.Slug)
	return self.post
}

type Published struct {
	Slug        string
	PublishedAt time.Time
}

func (self *Published) EventName() string {
	return EventPublished
}

func (self *Published) AggregateId() string {
	return self.Slug
}

func (self *Published) Apply(events.Aggregate) error {
	return nil
}
