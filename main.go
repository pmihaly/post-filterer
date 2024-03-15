package main

import (
	"errors"
	"math"
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
	Ratio        float64
}

func NewPostGroup(category string, ratio float64) PostGroup {
	group := PostGroup{
		PostCategory: category,
		Ratio:        ratio,
	}

	return group
}

type PostMixer struct {
	PostGroups       []PostGroup
	ratiosByCategory map[string]float64
}

func NewPostMixer(postGroups []PostGroup) (*PostMixer, error) {
	if len(postGroups) == 0 {
		return nil, errors.New("PostMixer has to have at least one post group")
	}

	ratiosByCategory := make(map[string]float64)

	for _, group := range postGroups {
		ratiosByCategory[group.PostCategory] = group.Ratio
	}

	mixer := &PostMixer{
		PostGroups:       postGroups,
		ratiosByCategory: ratiosByCategory,
	}

	return mixer, nil
}

func (mixer *PostMixer) MixPosts(posts []Post) ([]Post, error) {
	postCountsByCategory := make(map[string]int)

	for category, ratio := range mixer.ratiosByCategory {
		postCountsByCategory[category] = int(math.Ceil(ratio * float64(len(posts))))
	}

	mixedPosts := []Post{}

	for index, post := range posts {
		maxPostCount, hasCount := postCountsByCategory[post.Category]

		if !hasCount {
			continue
		}

		if index >= maxPostCount {
			continue
		}

		mixedPosts = append(mixedPosts, post)
	}

	return mixedPosts, nil
}

func main() {
	println("this works")
}
