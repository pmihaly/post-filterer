package main

import (
	"errors"
)

type Post struct{}
type Category struct{}

type PostMixer struct {
	Categories []Category
}

func NewPostMixer(categories []Category) (*PostMixer, error) {
	if len(categories) == 0 {
		return nil, errors.New("PostMixer has to have at least one category")
	}

	mixer := &PostMixer{
		Categories: categories,
	}

	return mixer, nil
}

func (mixer *PostMixer) MixPosts(posts []Post) ([]Post, error) {
	return posts, nil
}

func main() {
	println("this works")
}
