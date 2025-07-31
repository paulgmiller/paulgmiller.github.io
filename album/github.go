package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v65/github"
)

type GitHubClient struct {
	client *github.Client
	owner  string
	repo   string
}

func NewGitHubClient() *GitHubClient {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is required")
	}

	owner := os.Getenv("GITHUB_OWNER")
	if owner == "" {
		owner = "paulgmiller" // default
	}

	repo := os.Getenv("GITHUB_REPO")
	if repo == "" {
		repo = "paulgmiller.github.io" // default
	}

	client := github.NewClient(nil).WithAuthToken(token)

	return &GitHubClient{
		client: client,
		owner:  owner,
		repo:   repo,
	}
}

func (g *GitHubClient) CreateJekyllPost(ctx context.Context, title string, content string) error {
	// Generate filename with current date
	now := time.Now()
	filename := fmt.Sprintf("_posts/%s-%s.md", now.Format("2006-01-02"), slugify(title))

	// Create file content with Jekyll front matter
	fullContent := fmt.Sprintf(`---
layout: post
title: "%s"
date: %s
tags: [photos]
---

%s
`, title, now.Format("2006-01-02"), content)

	// Create the file in the repository
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(fmt.Sprintf("Add photo album: %s", title)),
		Content: []byte(fullContent),
		Branch:  github.String("master"),
	}

	_, _, err := g.client.Repositories.CreateFile(ctx, g.owner, g.repo, filename, opts)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	log.Printf("Successfully created Jekyll post: %s", filename)
	return nil
}

// Simple slugify function to convert title to URL-friendly format
func slugify(s string) string {
	// Basic implementation - replace spaces with dashes and remove special chars
	result := ""
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			result += string(r)
		} else if r == ' ' {
			result += "-"
		}
	}
	return result
}