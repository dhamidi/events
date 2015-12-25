package post

type Post struct {
	Slug Slug
}

func NewPost(slug Slug) *Post {
	return &Post{
		Slug: slug,
	}
}

func (self *Post) AggregateId() string {
	return self.Slug.String()
}
