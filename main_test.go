package main

import (
	"reflect"
	"testing"
)

func TestMixerCreationFailsWithNoCategories(t *testing.T) {
	mixer, err := NewPostMixer([]Category{})

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
		categories []Category
		posts      []Post
		want       []Post
	}{
		{"no posts", []Category{{}}, []Post{}, []Post{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mixer, err := NewPostMixer(tt.categories)

			if err != nil {
				t.Fatalf("failed to create mixer: %v", err)
			}

			result, err := mixer.MixPosts(tt.posts)

			if err != nil {
				t.Fatalf("failed to mix posts: %v", err)
			}

			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("got %t, want %t", result, tt.want)
			}
		})
	}
}
