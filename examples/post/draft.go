package post

import (
	"time"

	"github.com/dhamidi/events"
	"github.com/dhamidi/events/examples/user"
)

var Now = func() time.Time {
	return time.Now()
}

var (
	EventDrafted = "post.drafted"
)

type Draft struct {
	Author  user.Username
	Slug    Slug
	Title   string
	Content string

	post *Post
}

func NewDraft() *Draft {
	return &Draft{}
}

func (self *Draft) Execute() (events.Event, error) {
	return &Drafted{
		Author:    self.Author.String(),
		Slug:      self.Slug.String(),
		Title:     self.Title,
		Content:   self.Content,
		DraftedAt: Now(),
	}, nil
}

func (self *Draft) Aggregate() events.Aggregate {
	self.post = NewPost(self.Slug)
	return self.post
}

type Drafted struct {
	Author    string
	Slug      string
	Title     string
	Content   string
	DraftedAt time.Time
}

func (self *Drafted) EventName() string {
	return EventDrafted
}

func (self *Drafted) AggregateId() string {
	return self.Slug
}

func (self *Drafted) Apply(events.Aggregate) error {
	return nil
}
