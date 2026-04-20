package main

import (
	"strings"
	"testing"
	"time"
)

func TestBuildPostAutoTitleAndMarkdown(t *testing.T) {
	now := time.Date(2026, 3, 6, 8, 15, 0, 0, time.FixedZone("PST", -8*3600))
	post, err := buildPost(postRequest{
		Body:        "Whistler was cold and sunny and much better than expected.",
		Now:         now,
		PublishMode: publishModePR,
		Tags:        []string{"family", "photos"},
	})
	if err != nil {
		t.Fatalf("buildPost returned error: %v", err)
	}

	if got, want := post.Title, "Whistler was cold and sunny and much better"; got != want {
		t.Fatalf("title = %q, want %q", got, want)
	}
	if got, want := post.FilePath(now), "content/posts/2026-03-06-whistler-was-cold-and-sunny-and-much-better.md"; got != want {
		t.Fatalf("file path = %q, want %q", got, want)
	}

	markdown := post.RenderMarkdown([]string{"https://images.northbriton.net/test.jpg"})
	for _, fragment := range []string{
		"layout: post",
		"title: 'Whistler was cold and sunny and much better'",
		"tags: [family, photos]",
		"![](https://images.northbriton.net/test.jpg)",
	} {
		if !strings.Contains(markdown, fragment) {
			t.Fatalf("markdown missing %q\n%s", fragment, markdown)
		}
	}
}

func TestParseTagsAndFallbackTitle(t *testing.T) {
	tags := parseTags("Family, photos\nfamily, Hiking ")
	if got, want := strings.Join(tags, ","), "family,photos,hiking"; got != want {
		t.Fatalf("tags = %q, want %q", got, want)
	}

	now := time.Date(2026, 3, 6, 9, 4, 0, 0, time.UTC)
	post, err := buildPost(postRequest{
		Images:      []uploadedImage{{FileName: "IMG_1.JPG", ContentType: "image/jpeg", Data: []byte("abc")}},
		Now:         now,
		PublishMode: publishModeDirect,
	})
	if err != nil {
		t.Fatalf("buildPost returned error: %v", err)
	}
	if got, want := post.Title, "Photo post 2026-03-06 09:04"; got != want {
		t.Fatalf("title = %q, want %q", got, want)
	}
}

func TestNewImageObjectNamePreservesExtension(t *testing.T) {
	name, err := newImageObjectName(time.Date(2026, 3, 6, 9, 4, 5, 0, time.UTC), "IMG_1001.PNG", "image/png")
	if err != nil {
		t.Fatalf("newImageObjectName returned error: %v", err)
	}
	if !strings.HasPrefix(name, "20260306-090405-img-1001-") {
		t.Fatalf("unexpected prefix: %q", name)
	}
	if !strings.HasSuffix(name, ".png") {
		t.Fatalf("unexpected suffix: %q", name)
	}
}
