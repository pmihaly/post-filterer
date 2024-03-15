package main

import (
	"errors"
	"fmt"
	"log"
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

type PostWeight struct {
	PostCategory string
	Ratio        float64
}

func NewPostWeight(category string, ratio float64) PostWeight {
	weight := PostWeight{
		PostCategory: category,
		Ratio:        ratio,
	}

	return weight
}

func (parentWeight PostWeight) AddChildren(orphans []PostWeight) []PostWeight {
	childWeights := []PostWeight{}

	for _, orphan := range orphans {
		category := fmt.Sprintf("%v/%v", parentWeight.PostCategory, orphan.PostCategory)
		ratio := orphan.Ratio * parentWeight.Ratio

		childWeights = append(childWeights, NewPostWeight(category, ratio))
	}

	return childWeights
}

type PostFilterer struct {
	PostWeights      []PostWeight
	ratiosByCategory map[string]float64
}

func NewPostFilterer(postWeights []PostWeight) (*PostFilterer, error) {
	if len(postWeights) == 0 {
		return nil, errors.New("PostFilterer has to have at least one post weight")
	}

	ratiosByCategory := make(map[string]float64)

	for _, weight := range postWeights {
		ratiosByCategory[weight.PostCategory] = weight.Ratio
	}

	filterer := &PostFilterer{
		PostWeights:      postWeights,
		ratiosByCategory: ratiosByCategory,
	}

	return filterer, nil
}

func (filterer *PostFilterer) dropIrrelevantPosts(posts []Post) []Post {
	relevantPosts := []Post{}

	for _, post := range posts {
		_, isRelevant := filterer.ratiosByCategory[post.Category]

		if isRelevant {
			relevantPosts = append(relevantPosts, post)
		}
	}

	return relevantPosts
}

func (filterer *PostFilterer) getPostCountsByCategory(posts []Post) map[string]int {
	postCountsByCategory := make(map[string]int)

	for category, ratio := range filterer.ratiosByCategory {
		postCountsByCategory[category] = int(ratio * float64(len(posts)))
	}

	return postCountsByCategory
}

func (filterer *PostFilterer) FilterPosts(posts []Post) []Post {
	relevantPosts := filterer.dropIrrelevantPosts(posts)

	if len(filterer.ratiosByCategory) == 1 {
		return relevantPosts
	}

	postCountsByCategory := filterer.getPostCountsByCategory(relevantPosts)

	filteredPosts := []Post{}

	for _, post := range posts {
		postCount, hasCount := postCountsByCategory[post.Category]

		if !hasCount {
			continue
		}

		if postCount == 1 {
			continue
		}

		filteredPosts = append(filteredPosts, post)
		postCountsByCategory[post.Category] = postCount - 1
	}

	return filteredPosts
}

func mergeWeights(arrays ...[]PostWeight) []PostWeight {
	var merged []PostWeight
	for _, arr := range arrays {
		merged = append(merged, arr...)
	}
	return merged
}

func main() {
	weights := mergeWeights(
		NewPostWeight("top", 0.2).AddChildren([]PostWeight{
			NewPostWeight("daily", 0.5),
			NewPostWeight("weekly", 0.3),
			NewPostWeight("monthly", 0.2),
		}),
		NewPostWeight("trending", 0.2).AddChildren([]PostWeight{
			NewPostWeight("city", 0.4),
			NewPostWeight("area", 0.3),
			NewPostWeight("country", 0.3),
		}), []PostWeight{NewPostWeight("promoted", 0.1)},
		NewPostWeight("following", 0.5).AddChildren([]PostWeight{
			NewPostWeight("immediate-follow", 0.6),
			NewPostWeight("follow-of-follow", 0.4),
		}),
	)

	filterer, err := NewPostFilterer(weights)

	if err != nil {
		log.Fatalf("creating filterer failed: %v", err)
		return
	}

	posts := []Post{
		NewPost("chatgpt-generated-corporate-bs", "top/daily"),
		NewPost("ricks-vacation", "following/immediate-follow"),
		NewPost("lisas-puppy", "trending/area"),
		NewPost("theprimagen-shilling-rust", "top/monthly"),
		NewPost("jonass-cold-take", "following/follow-of-follow"),
		NewPost("definetly-a-linux-iso-torrent", "trending/city"),
		NewPost("useless-infographic-definetly-not-stolen-from-linkedin", "promoted"),
		NewPost("coding-humor", "top/daily"),
		NewPost("family-picnic", "following/immediate-follow"),
		NewPost("local-concert-highlights", "trending/area"),
		NewPost("book-recommendations", "top/monthly"),
		NewPost("tech-gadget-review", "following/follow-of-follow"),
		NewPost("art-exhibition-in-city", "trending/city"),
		NewPost("limited-time-offer", "promoted"),
		NewPost("memes-of-the-day", "top/daily"),
		NewPost("weekend-getaway-ideas", "following/immediate-follow"),
		NewPost("community-event-update", "trending/area"),
		NewPost("music-album-release", "top/monthly"),
		NewPost("opinion-poll", "following/follow-of-follow"),
		NewPost("street-food-festival", "trending/city"),
		NewPost("exclusive-discount-code", "promoted"),
	}

	filteredPosts := filterer.FilterPosts(posts)

	log.Printf("%v", filteredPosts)
}
