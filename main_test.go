package main

import (
	"reflect"
	"testing"
)

func TestMixerCreationFailsWithNoPostGroups(t *testing.T) {
	mixer, err := NewPostMixer([]PostGroup{})

	if mixer != nil {
		t.Errorf("got %v, want %v", mixer, nil)
	}

	if err == nil {
		t.Errorf("got %v, want non-nil error", err)
	}

}

func TestBasic(t *testing.T) {
	var tests = []struct {
		name       string
		postGroups []PostGroup
		posts      []Post
		want       []Post
	}{
		{"no posts return empty array", []PostGroup{NewPostGroup("trending", 1.0)}, []Post{}, []Post{}},
		{"single group returns the input array", []PostGroup{NewPostGroup("trending", 1.0)}, []Post{
			NewPost("ricks-vacation", "trending"),
			NewPost("lisas-puppy", "trending"),
			NewPost("theprimagen-shilling-rust", "trending"),
		}, []Post{
			NewPost("ricks-vacation", "trending"),
			NewPost("lisas-puppy", "trending"),
			NewPost("theprimagen-shilling-rust", "trending"),
		}},
		{"groups filter out irrelevant categories", []PostGroup{NewPostGroup("trending", 0.4), NewPostGroup("following", 0.6)}, []Post{
			NewPost("ricks-vacation", "following"),
			NewPost("lisas-puppy", "trending"),
			NewPost("theprimagen-shilling-rust", "hot"),
		}, []Post{
			NewPost("ricks-vacation", "following"),
			NewPost("lisas-puppy", "trending"),
		}},
		{"mixing by 50-50", []PostGroup{NewPostGroup("trending", 0.5), NewPostGroup("following", 0.5)}, []Post{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mixer, err := NewPostMixer(tt.postGroups)

			if err != nil {
				t.Fatalf("failed to create mixer: %v", err)
			}

			result, err := mixer.MixPosts(tt.posts)

			if err != nil {
				t.Fatalf("failed to mix posts: %v", err)
			}

			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("got %v, want %v (lengths of %v and %v)", result, tt.want, len(result), len(tt.want))
			}
		})
	}
}
