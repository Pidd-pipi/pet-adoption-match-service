package service

import (
	"testing"

	"github.com/gbadopt/gbadopt/internal/model"
)

func TestPostFiltersNoAlias(t *testing.T) {
	check := func(name string, fn func([]model.CommunityPost) []model.CommunityPost) {
		posts := []model.CommunityPost{
			{ID: 1, Status: "draft", PostType: "story", Images: "[]"},
			{ID: 2, Status: "published", PostType: "story", Images: `["a.jpg"]`},
			{ID: 3, Status: "published", PostType: "lost_notice", Images: `["b.jpg"]`},
			{ID: 4, Status: "published", PostType: "story", Images: "[]"},
		}
		_ = fn(posts)
		for i := 0; i < 4; i++ {
			if posts[i].ID != uint(i+1) {
				t.Fatalf("%s corrupted posts[%d]: %+v", name, i, posts)
			}
		}
	}
	check("filterPublished", filterPublished)
	check("filterStory", filterStory)
	check("filterWithImages", filterWithImages)
}
