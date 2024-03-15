package main

import (
	"errors"
)

type Post struct {
	Id       string
	Category string
}

func NewPost(id, category string) Post {
	post := Post{
		Id:       id,
		Category: category,
	}

	return post
}

type PostGroup struct {
	PostCategory string
	Ratio        float32
}

func NewPostGroup(category string, ratio float32) PostGroup {
	group := PostGroup{
		PostCategory: category,
		Ratio:        ratio,
	}

	return group
}

type PostMixer struct {
	PostGroups []PostGroup
	categories map[string]struct{}
}

func NewPostMixer(postGroups []PostGroup) (*PostMixer, error) {
	if len(postGroups) == 0 {
		return nil, errors.New("PostMixer has to have at least one post group")
	}

	categories := make(map[string]struct{})

	for _, group := range postGroups {
		categories[group.PostCategory] = struct{}{}
	}

	mixer := &PostMixer{
		PostGroups: postGroups,
		categories: categories,
	}

	return mixer, nil
}

func (mixer *PostMixer) MixPosts(posts []Post) ([]Post, error) {
	mixedPosts := []Post{}

	for _, post := range posts {
		_, shouldIncludePost := mixer.categories[post.Category]

		if shouldIncludePost {
			mixedPosts = append(mixedPosts, post)
		}

	}

	return mixedPosts, nil
}

func main() {
	println("this works")
}
