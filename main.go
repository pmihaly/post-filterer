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

type PostMixer struct {
	PostWeights      []PostWeight
	ratiosByCategory map[string]float64
}

func NewPostMixer(postWeights []PostWeight) (*PostMixer, error) {
	if len(postWeights) == 0 {
		return nil, errors.New("PostMixer has to have at least one post weight")
	}

	ratiosByCategory := make(map[string]float64)

	for _, weight := range postWeights {
		ratiosByCategory[weight.PostCategory] = weight.Ratio
	}

	mixer := &PostMixer{
		PostWeights:      postWeights,
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

func (mixer *PostMixer) MixPosts(posts []Post) []Post {
	filteredPosts := mixer.filterIrrelevantPosts(posts)

	if len(mixer.ratiosByCategory) == 1 {
		return filteredPosts
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

	return mixedPosts
}

func mergeWeights(arrays ...[]PostWeight) []PostWeight {
	var merged []PostWeight
	for _, arr := range arrays {
		merged = append(merged, arr...)
	}
	return merged
}

func main() {
	mixingWeights := mergeWeights(
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

	mixer, err := NewPostMixer(mixingWeights)

	if err != nil {
		log.Fatalf("creating mixer failed: %v", err)
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

	mixedPosts := mixer.MixPosts(posts)

	log.Printf("%v", mixedPosts)
}
