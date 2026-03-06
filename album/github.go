package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const githubAPIBaseURL = "https://api.github.com"

type githubClient struct {
	apiBaseURL  string
	authorEmail string
	authorName  string
	baseBranch  string
	client      *http.Client
	owner       string
	repo        string
	token       string
}

type githubPublishRequest struct {
	BaseBranch  string
	Body        []byte
	CommitBody  string
	FilePath    string
	PRBody      string
	PRTitle     string
	PublishMode publishMode
	Slug        string
}

type publishResult struct {
	Message string
	URL     string
}

type gitRefResponse struct {
	Object struct {
		SHA string `json:"sha"`
	} `json:"object"`
}

type createRefRequest struct {
	Ref string `json:"ref"`
	SHA string `json:"sha"`
}

type putContentRequest struct {
	Branch    string          `json:"branch"`
	Committer githubCommitter `json:"committer"`
	Content   string          `json:"content"`
	Message   string          `json:"message"`
	SHA       string          `json:"sha,omitempty"`
}

type putContentResponse struct {
	Commit struct {
		SHA string `json:"sha"`
	} `json:"commit"`
}

type githubCommitter struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type pullRequestRequest struct {
	Base  string `json:"base"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Title string `json:"title"`
}

type pullRequestResponse struct {
	HTMLURL string `json:"html_url"`
}

type githubContentResponse struct {
	SHA string `json:"sha"`
}

func newGitHubClient(cfg Config) *githubClient {
	return &githubClient{
		apiBaseURL:  githubAPIBaseURL,
		authorEmail: cfg.GitHubAuthorEmail,
		authorName:  cfg.GitHubAuthorName,
		baseBranch:  cfg.GitHubBaseBranch,
		client:      &http.Client{Timeout: 30 * time.Second},
		owner:       cfg.GitHubOwner,
		repo:        cfg.GitHubRepo,
		token:       cfg.GitHubToken,
	}
}

func (g *githubClient) Publish(ctx context.Context, req githubPublishRequest) (publishResult, error) {
	baseBranch := req.BaseBranch
	if baseBranch == "" {
		baseBranch = g.baseBranch
	}
	if g.token == "" || g.owner == "" || g.repo == "" {
		return publishResult{}, errors.New("github is not configured")
	}

	if req.PublishMode == publishModeDirect {
		sha, err := g.putFile(ctx, baseBranch, req.FilePath, req.CommitBody, req.Body)
		if err != nil {
			return publishResult{}, err
		}
		return publishResult{
			Message: "Published directly to " + baseBranch,
			URL:     fmt.Sprintf("https://github.com/%s/%s/commit/%s", g.owner, g.repo, sha),
		}, nil
	}

	headSHA, err := g.getBranchSHA(ctx, baseBranch)
	if err != nil {
		return publishResult{}, err
	}
	branchName := fmt.Sprintf("mobile-post-%s-%s", time.Now().UTC().Format("20060102-150405"), req.Slug)
	if err := g.createBranch(ctx, branchName, headSHA); err != nil {
		return publishResult{}, err
	}
	if _, err := g.putFile(ctx, branchName, req.FilePath, req.CommitBody, req.Body); err != nil {
		return publishResult{}, err
	}
	prURL, err := g.createPullRequest(ctx, branchName, baseBranch, req.PRTitle, req.PRBody)
	if err != nil {
		return publishResult{}, err
	}
	return publishResult{
		Message: "Opened pull request",
		URL:     prURL,
	}, nil
}

func (g *githubClient) EnsureUniquePath(ctx context.Context, requestedPath string) (string, error) {
	if g.token == "" || g.owner == "" || g.repo == "" {
		return requestedPath, nil
	}

	ext := pathExt(requestedPath)
	base := strings.TrimSuffix(requestedPath, ext)
	candidate := requestedPath
	for i := 0; i < 20; i++ {
		exists, err := g.fileExists(ctx, g.baseBranch, candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%02d%s", base, i+1, ext)
	}
	return "", fmt.Errorf("could not find unique path for %s", requestedPath)
}

func (g *githubClient) getBranchSHA(ctx context.Context, branch string) (string, error) {
	var response gitRefResponse
	url := fmt.Sprintf("%s/repos/%s/%s/git/ref/heads/%s", g.apiBaseURL, g.owner, g.repo, url.PathEscape(branch))
	if err := g.doJSON(ctx, http.MethodGet, url, nil, &response, http.StatusOK); err != nil {
		return "", err
	}
	if response.Object.SHA == "" {
		return "", errors.New("missing branch sha")
	}
	return response.Object.SHA, nil
}

func (g *githubClient) createBranch(ctx context.Context, branch string, sha string) error {
	request := createRefRequest{
		Ref: "refs/heads/" + branch,
		SHA: sha,
	}
	url := fmt.Sprintf("%s/repos/%s/%s/git/refs", g.apiBaseURL, g.owner, g.repo)
	return g.doJSON(ctx, http.MethodPost, url, request, nil, http.StatusCreated)
}

func (g *githubClient) putFile(ctx context.Context, branch string, filePath string, message string, body []byte) (string, error) {
	request := putContentRequest{
		Branch:  branch,
		Content: base64.StdEncoding.EncodeToString(body),
		Committer: githubCommitter{
			Email: g.authorEmail,
			Name:  g.authorName,
		},
		Message: message,
	}

	existing, err := g.getFileSHA(ctx, branch, filePath)
	if err != nil {
		var httpErr *githubHTTPError
		if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusNotFound {
			return "", err
		}
	}
	if existing != "" {
		request.SHA = existing
	}

	var response putContentResponse
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", g.apiBaseURL, g.owner, g.repo, githubPathEscape(filePath))
	if err := g.doJSON(ctx, http.MethodPut, url, request, &response, http.StatusCreated, http.StatusOK); err != nil {
		return "", err
	}
	if response.Commit.SHA == "" {
		return "", errors.New("missing commit sha")
	}
	return response.Commit.SHA, nil
}

func (g *githubClient) createPullRequest(ctx context.Context, head string, base string, title string, body string) (string, error) {
	request := pullRequestRequest{
		Base:  base,
		Body:  body,
		Head:  head,
		Title: title,
	}
	var response pullRequestResponse
	url := fmt.Sprintf("%s/repos/%s/%s/pulls", g.apiBaseURL, g.owner, g.repo)
	if err := g.doJSON(ctx, http.MethodPost, url, request, &response, http.StatusCreated); err != nil {
		return "", err
	}
	if response.HTMLURL == "" {
		return "", errors.New("missing pull request url")
	}
	return response.HTMLURL, nil
}

func (g *githubClient) fileExists(ctx context.Context, branch string, filePath string) (bool, error) {
	sha, err := g.getFileSHA(ctx, branch, filePath)
	if err != nil {
		var httpErr *githubHTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return sha != "", nil
}

func (g *githubClient) getFileSHA(ctx context.Context, branch string, filePath string) (string, error) {
	var response githubContentResponse
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s?ref=%s", g.apiBaseURL, g.owner, g.repo, githubPathEscape(filePath), url.QueryEscape(branch))
	if err := g.doJSON(ctx, http.MethodGet, url, nil, &response, http.StatusOK); err != nil {
		var httpErr *githubHTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
			return "", err
		}
		return "", err
	}
	return response.SHA, nil
}

type githubHTTPError struct {
	Body       string
	StatusCode int
}

func (e *githubHTTPError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("github request failed with %d", e.StatusCode)
	}
	return fmt.Sprintf("github request failed with %d: %s", e.StatusCode, e.Body)
}

func (g *githubClient) doJSON(ctx context.Context, method string, endpoint string, requestBody any, responseBody any, expectedStatuses ...int) error {
	var body io.Reader
	if requestBody != nil {
		data, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	for _, expected := range expectedStatuses {
		if resp.StatusCode == expected {
			if responseBody == nil {
				return nil
			}
			return json.NewDecoder(resp.Body).Decode(responseBody)
		}
	}

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return &githubHTTPError{
		Body:       strings.TrimSpace(string(bodyBytes)),
		StatusCode: resp.StatusCode,
	}
}

func githubPathEscape(value string) string {
	return strings.ReplaceAll(url.PathEscape(value), "%2F", "/")
}

func pathExt(filePath string) string {
	lastSlash := strings.LastIndex(filePath, "/")
	lastDot := strings.LastIndex(filePath, ".")
	if lastDot <= lastSlash {
		return ""
	}
	return filePath[lastDot:]
}
