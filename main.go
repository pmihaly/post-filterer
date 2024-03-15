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

func (parentGroup PostGroup) AddChildren(orphans []PostGroup) []PostGroup {
	childGroups := []PostGroup{}

	for _, orphan := range orphans {
		category := fmt.Sprintf("%v/%v", parentGroup.PostCategory, orphan.PostCategory)
		ratio := orphan.Ratio * parentGroup.Ratio

		childGroups = append(childGroups, NewPostGroup(category, ratio))
	}

	return childGroups
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

func main() {
	mixingWeights := []PostGroup{}

	mixingWeights = append(mixingWeights, NewPostGroup("top", 0.2).AddChildren([]PostGroup{
		NewPostGroup("daily", 0.5),
		NewPostGroup("weekly", 0.3),
		NewPostGroup("monthly", 0.2),
	})...)
	mixingWeights = append(mixingWeights, NewPostGroup("trending", 0.2).AddChildren([]PostGroup{
		NewPostGroup("city", 0.4),
		NewPostGroup("area", 0.3),
		NewPostGroup("country", 0.3),
	})...)
	mixingWeights = append(mixingWeights, NewPostGroup("promoted", 0.1))
	mixingWeights = append(mixingWeights, NewPostGroup("following", 0.5).AddChildren([]PostGroup{
		NewPostGroup("immediate-follow", 0.6),
		NewPostGroup("follow-of-follow", 0.4),
	})...)

	mixer, err := NewPostMixer(mixingWeights)

	if err != nil {
		log.Fatalf("creating mixer failed: %v", err)
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
