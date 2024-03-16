package main

import (
	"reflect"
	"testing"
)

func TestFiltererCreationFailsWithNoPostWeights(t *testing.T) {
	filterer, err := NewPostFilterer([]PostWeight{})

	if filterer != nil {
		t.Errorf("got %v, want %v", filterer, nil)
	}

	if err == nil {
		t.Errorf("got %v, want non-nil error", err)
	}

}

func TestFlatHierarchy(t *testing.T) {
	var tests = []struct {
		name        string
		postWeights []PostWeight
		posts       []Post
		want        []Post
	}{
		{"no posts return empty array", []PostWeight{NewPostWeight("trending", 1.0)}, []Post{}, []Post{}},
		{"single weight returns the input array", []PostWeight{NewPostWeight("trending", 1.0)}, []Post{
			NewPost("ricks-vacation", "trending"),
			NewPost("lisas-puppy", "trending"),
			NewPost("theprimagen-shilling-rust", "trending"),
		}, []Post{
			NewPost("ricks-vacation", "trending"),
			NewPost("lisas-puppy", "trending"),
			NewPost("theprimagen-shilling-rust", "trending"),
		}},
		{"weights filter out irrelevant categories", []PostWeight{NewPostWeight("trending", 1.0), NewPostWeight("following", 1.0)}, []Post{
			NewPost("ricks-vacation", "following"),
			NewPost("lisas-puppy", "trending"),
			NewPost("theprimagen-shilling-rust", "hot"),
			NewPost("definetly-a-linux-iso-torrent", "trending"),
		}, []Post{
			NewPost("ricks-vacation", "following"),
			NewPost("lisas-puppy", "trending"),
			NewPost("definetly-a-linux-iso-torrent", "trending"),
		}},
		{"2 categories split by 50-50", []PostWeight{NewPostWeight("trending", 0.5), NewPostWeight("following", 0.5)}, []Post{
			NewPost("ricks-vacation", "following"),
			NewPost("lisas-puppy", "trending"),
			NewPost("jonass-cold-take", "following"),
			NewPost("definetly-a-linux-iso-torrent", "trending"),
			NewPost("this-shouldnt-be-making-into-the-feed-1", "following"),
			NewPost("this-shouldnt-be-making-into-the-feed-2", "following"),
			NewPost("this-shouldnt-be-making-into-the-feed-3", "following"),
		}, []Post{
			NewPost("ricks-vacation", "following"),
			NewPost("lisas-puppy", "trending"),
			NewPost("jonass-cold-take", "following"),
			NewPost("definetly-a-linux-iso-torrent", "trending"),
		}},
		{"filtering by weights should ignore irrelevant posts", []PostWeight{NewPostWeight("trending", 0.5), NewPostWeight("following", 0.5)}, []Post{
			NewPost("chatgpt-generated-corporate-bs", "hot"),
			NewPost("ricks-vacation", "following"),
			NewPost("lisas-puppy", "trending"),
			NewPost("theprimagen-shilling-rust", "hot"),
			NewPost("jonass-cold-take", "following"),
			NewPost("definetly-a-linux-iso-torrent", "trending"),
			NewPost("this-shouldnt-be-making-into-the-feed-1", "following"),
			NewPost("useless-infographic-definetly-not-stolen-from-linkedin", "hot"),
			NewPost("this-shouldnt-be-making-into-the-feed-2", "following"),
			NewPost("this-shouldnt-be-making-into-the-feed-3", "following"),
		}, []Post{
			NewPost("ricks-vacation", "following"),
			NewPost("lisas-puppy", "trending"),
			NewPost("jonass-cold-take", "following"),
			NewPost("definetly-a-linux-iso-torrent", "trending"),
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filterer, err := NewPostFilterer(tt.postWeights)

			if err != nil {
				t.Fatalf("failed to create filterer: %v", err)
			}

			result := filterer.FilterPosts(tt.posts)

			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("got %v, want %v (lengths of %v and %v)", result, tt.want, len(result), len(tt.want))
			}
		})
	}
}

func TestAddingChildren(t *testing.T) {
	deepHierarchy := []PostWeight{
		NewPostWeight("trending", 0.5),
	}

	deepHierarchy = append(deepHierarchy,
		NewPostWeight("following", 0.5).AddChildren([]PostWeight{NewPostWeight("immediate-follow", 0.6), NewPostWeight("follow-of-follow", 0.4)})...,
	)

	flatHierarchy := []PostWeight{
		NewPostWeight("trending", 0.5),
		NewPostWeight("following/immediate-follow", 0.3),
		NewPostWeight("following/follow-of-follow", 0.2),
	}

	if !reflect.DeepEqual(deepHierarchy, flatHierarchy) {
		t.Errorf("got %v, want %v", deepHierarchy, flatHierarchy)
	}
}
