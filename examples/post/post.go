package post

type Post struct {
	Slug    Slug
	Title   string
	Content string
}

func NewPost(slug Slug) *Post {
	return &Post{
		Slug: slug,
	}
}

func (self *Post) AggregateId() string {
	return self.Slug.String()
}
