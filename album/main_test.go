package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestParseAlbum(t *testing.T) {
	image := func(char string) string {
		return "https://lh3.googleusercontent.com/pw/" + strings.Repeat(char, 64)
	}
	body := fmt.Sprintf(`<html><head><title>Gem &amp; Snow - Google Photos</title></head><body>%s %s %s %s</body></html>`,
		image("a"), image("b"), image("b"), image("c"))

	album, err := parseAlbum([]byte(body))
	if err != nil {
		t.Fatalf("parseAlbum returned an error: %v", err)
	}
	if album.title != "Gem & Snow" {
		t.Fatalf("title = %q, want %q", album.title, "Gem & Snow")
	}
	if len(album.photoURLs) != 1 || album.photoURLs[0] != image("b") {
		t.Fatalf("photoURLs = %#v, want the unique gallery image", album.photoURLs)
	}
}

func TestParseAlbumUsesDefaultTitle(t *testing.T) {
	image := func(char string) string {
		return "https://lh3.googleusercontent.com/pw/" + strings.Repeat(char, 64)
	}
	body := fmt.Sprintf(`<html><body>%s %s %s</body></html>`, image("a"), image("b"), image("c"))

	album, err := parseAlbum([]byte(body))
	if err != nil {
		t.Fatalf("parseAlbum returned an error: %v", err)
	}
	if album.title != defaultAlbumTitle {
		t.Fatalf("title = %q, want %q", album.title, defaultAlbumTitle)
	}
}

func TestOutputIncludesHugoFrontMatterAndGallery(t *testing.T) {
	var outputBuffer bytes.Buffer
	err := output(
		[]string{"https://images.northbriton.net/photo"},
		"https://photos.app.goo.gl/example",
		"Gem Lake",
		"2026-08-27",
		&outputBuffer,
	)
	if err != nil {
		t.Fatalf("output returned an error: %v", err)
	}

	want := `---
layout: post
title: Gem Lake
date: 2026-08-27
---

<div class="fotorama" data-allowfullscreen="true">
<!--https://photos.app.goo.gl/example-->
    <img src="https://images.northbriton.net/cdn-cgi/image/width=800/https://images.northbriton.net/photo" data-full="https://images.northbriton.net/photo">
</div>
`
	if outputBuffer.String() != want {
		t.Fatalf("output =\n%s\nwant:\n%s", outputBuffer.String(), want)
	}
}

func TestNewPostURL(t *testing.T) {
	content := "---\ntitle: Gem Lake\n---\n"
	newFileURL, err := url.Parse(newPostURL("Gem Lake!", "2026-08-27", content))
	if err != nil {
		t.Fatalf("generated URL is invalid: %v", err)
	}
	if got, want := newFileURL.Scheme+"://"+newFileURL.Host+newFileURL.Path, githubNewPostURL; got != want {
		t.Fatalf("URL base = %q, want %q", got, want)
	}
	if got, want := newFileURL.Query().Get("filename"), "2026-08-27-gem-lake.md"; got != want {
		t.Fatalf("filename = %q, want %q", got, want)
	}
	if got := newFileURL.Query().Get("value"); got != content {
		t.Fatalf("value = %q, want %q", got, content)
	}
}

func TestYAMLStringQuotesSpecialTitles(t *testing.T) {
	if got, want := yamlString("Trip: day one"), `"Trip: day one"`; got != want {
		t.Fatalf("yamlString = %q, want %q", got, want)
	}
	if got, want := yamlString("Box 60"), "Box 60"; got != want {
		t.Fatalf("yamlString = %q, want %q", got, want)
	}
}
