package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGitHubClientPublishPullRequest(t *testing.T) {
	var (
		createdRef string
		prHead     string
		fileBody   string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/paul/blog/git/ref/heads/master":
			_ = json.NewEncoder(w).Encode(map[string]any{"object": map[string]string{"sha": "base-sha"}})
		case r.Method == http.MethodPost && r.URL.Path == "/repos/paul/blog/git/refs":
			var req createRefRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			createdRef = req.Ref
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]string{"ref": req.Ref})
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/repos/paul/blog/contents/"):
			http.NotFound(w, r)
		case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/repos/paul/blog/contents/"):
			var req putContentRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			raw, _ := base64.StdEncoding.DecodeString(req.Content)
			fileBody = string(raw)
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{"commit": map[string]string{"sha": "commit-sha"}})
		case r.Method == http.MethodPost && r.URL.Path == "/repos/paul/blog/pulls":
			var req pullRequestRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			prHead = req.Head
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]string{"html_url": "https://github.com/paul/blog/pull/1"})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	}))
	defer server.Close()

	client := &githubClient{
		apiBaseURL:  server.URL,
		authorEmail: "paul@example.com",
		authorName:  "Paul",
		baseBranch:  "master",
		client:      server.Client(),
		owner:       "paul",
		repo:        "blog",
		token:       "token",
	}

	result, err := client.Publish(context.Background(), githubPublishRequest{
		BaseBranch:  "master",
		Body:        []byte("hello"),
		CommitBody:  "Draft hello",
		FilePath:    "content/posts/test.md",
		PRBody:      "body",
		PRTitle:     "title",
		PublishMode: publishModePR,
		Slug:        "test",
	})
	if err != nil {
		t.Fatalf("Publish returned error: %v", err)
	}
	if got, want := result.URL, "https://github.com/paul/blog/pull/1"; got != want {
		t.Fatalf("url = %q, want %q", got, want)
	}
	if got, want := createdRef, "refs/heads/"+prHead; got != want {
		t.Fatalf("created ref = %q, want %q", got, want)
	}
	if got, want := fileBody, "hello"; got != want {
		t.Fatalf("file body = %q, want %q", got, want)
	}
}

func TestEnsureUniquePathAddsSuffix(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/paul/blog/contents/content/posts/test.md":
			_ = json.NewEncoder(w).Encode(map[string]string{"sha": "existing"})
		case "/repos/paul/blog/contents/content/posts/test-01.md":
			http.NotFound(w, r)
		default:
			t.Fatalf("unexpected request: %s", r.URL.String())
		}
	}))
	defer server.Close()

	client := &githubClient{
		apiBaseURL: server.URL,
		client:     server.Client(),
		baseBranch: "master",
		owner:      "paul",
		repo:       "blog",
		token:      "token",
	}

	got, err := client.EnsureUniquePath(context.Background(), "content/posts/test.md")
	if err != nil {
		t.Fatalf("EnsureUniquePath returned error: %v", err)
	}
	if want := "content/posts/test-01.md"; got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
}
