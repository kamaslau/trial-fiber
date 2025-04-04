// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Mutation struct {
}

type Post struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Content string  `json:"content"`
	Excerpt *string `json:"excerpt,omitempty"`
}

type Query struct {
}

type CreatePostInput struct {
	Name    string  `json:"name"`
	Content string  `json:"content"`
	Excerpt *string `json:"excerpt,omitempty"`
}
