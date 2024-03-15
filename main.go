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

func (mixer *PostMixer) filterIrrelevantPosts(posts []Post) []Post {
	filteredPosts := []Post{}

	for _, post := range posts {
		_, isRelevant := mixer.ratiosByCategory[post.Category]

		if isRelevant {
			filteredPosts = append(filteredPosts, post)
		}
	}

	return filteredPosts
}

func (mixer *PostMixer) getPostCountsByCategory(posts []Post) map[string]int {
	postCountsByCategory := make(map[string]int)

	for category, ratio := range mixer.ratiosByCategory {
		postCountsByCategory[category] = int(ratio * float64(len(posts)))
	}

	return postCountsByCategory
}

func (mixer *PostMixer) MixPosts(posts []Post) ([]Post, error) {
	filteredPosts := mixer.filterIrrelevantPosts(posts)

	if len(mixer.ratiosByCategory) == 1 {
		return filteredPosts, nil
	}

	postCountsByCategory := mixer.getPostCountsByCategory(filteredPosts)

	mixedPosts := []Post{}

	for _, post := range posts {
		postCount, hasCount := postCountsByCategory[post.Category]

		if !hasCount {
			continue
		}

		if postCount == 1 {
			continue
		}

		mixedPosts = append(mixedPosts, post)
		postCountsByCategory[post.Category] = postCount - 1
	}

	return mixedPosts, nil
}

func main() {
	println("this works")
}
