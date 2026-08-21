package handler

import (
	"testing"

	"github.com/gbadopt/gbadopt/internal/model"
)

func TestPostHelpersNoAlias(t *testing.T) {
	check := func(name string, fn func([]model.CommunityPost) []model.CommunityPost) {
		posts := []model.CommunityPost{
			{ID: 1, Images: "[]", PostType: "lost_notice", Status: "published"},
			{ID: 2, Images: `["a.jpg"]`, PostType: "story", Status: "published"},
			{ID: 3, Images: "[]", PostType: "story", Status: "draft"},
			{ID: 4, Images: `["b.jpg"]`, PostType: "story", Status: "published"},
		}
		_ = fn(posts)
		for i := 0; i < 4; i++ {
			if posts[i].ID != uint(i+1) {
				t.Fatalf("%s corrupted posts[%d]: %+v", name, i, posts)
			}
		}
	}
	check("keepFeatured", keepFeatured)
	check("keepStory", keepStory)
	check("keepPublished", keepPublished)
}
