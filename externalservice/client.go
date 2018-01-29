package externalservice

import (
	"errors"
)

// Post is the data structure representing the data sent and received from the
// external service
type Post struct {
	ID int `json:"id"` // the primary key

	Title       string `json:"title" form:"title"`
	Description string `json:"description,omitempty" form:"description"`
}

// Client represents the client interface to the external service
type Client interface {
	GET(id int) (*Post, error)
	POST(id int, post *Post) (*Post, error)
}

// ClientMock implements Client interfacte
type ClientMock struct {
	Posts map[int]*Post
}

func (c *ClientMock) GET(id int) (*Post, error) {
	if c.Posts[id] == nil {
		return nil, errors.New("Post not found")
	}
	return c.Posts[id], nil
}

func (c *ClientMock) POST(id int, post *Post) (*Post, error) {
	if c.Posts[id] != nil {
		return nil, errors.New("Post id is already called")
	}
	post.ID = id
	c.Posts[id] = post
	return post, nil
}
