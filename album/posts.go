package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"mime"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type publishMode string

const (
	publishModeDirect publishMode = "direct"
	publishModePR     publishMode = "pr"
)

type uploadedImage struct {
	ContentType string
	Data        []byte
	FileName    string
}

type postRequest struct {
	Body        string
	Images      []uploadedImage
	Now         time.Time
	PublishMode publishMode
	Tags        []string
	Title       string
}

type renderedPost struct {
	Body      string
	Date      time.Time
	Slug      string
	Tags      []string
	Title     string
	UserTitle string
}

func (r postRequest) Validate() error {
	if strings.TrimSpace(r.Body) == "" && len(r.Images) == 0 {
		return errors.New("body or images required")
	}
	if r.Now.IsZero() {
		return errors.New("missing timestamp")
	}
	if r.PublishMode != publishModeDirect && r.PublishMode != publishModePR {
		return errors.New("publish_mode must be pr or direct")
	}
	return nil
}

func buildPost(req postRequest) (renderedPost, error) {
	if err := req.Validate(); err != nil {
		return renderedPost{}, err
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = generateTitle(req)
	}
	slug := slugify(title)
	if slug == "" {
		slug = "post"
	}

	return renderedPost{
		Body:      strings.TrimSpace(req.Body),
		Date:      req.Now,
		Slug:      slug,
		Tags:      req.Tags,
		Title:     title,
		UserTitle: strings.TrimSpace(req.Title),
	}, nil
}

func (p renderedPost) FilePath(now time.Time) string {
	return fmt.Sprintf("content/posts/%s-%s.md", now.Format("2006-01-02"), p.Slug)
}

func (p renderedPost) CommitMessage(mode publishMode) string {
	if mode == publishModeDirect {
		return fmt.Sprintf("Publish %s", p.Title)
	}
	return fmt.Sprintf("Draft %s", p.Title)
}

func (p renderedPost) PullRequestTitle() string {
	return "New post: " + p.Title
}

func (p renderedPost) PullRequestBody() string {
	return fmt.Sprintf("Create `%s` from the mobile publisher.", p.FilePath(p.Date))
}

func (p renderedPost) RenderMarkdown(imageURLs []string) string {
	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.WriteString("layout: post\n")
	buf.WriteString("title: " + yamlQuote(p.Title) + "\n")
	buf.WriteString("tags: [" + strings.Join(p.Tags, ", ") + "]\n")
	buf.WriteString("date: " + p.Date.Format(time.RFC3339) + "\n")
	buf.WriteString("---\n\n")

	if p.Body != "" {
		buf.WriteString(p.Body)
		buf.WriteString("\n")
	}
	for _, imageURL := range imageURLs {
		buf.WriteString("\n![](")
		buf.WriteString(imageURL)
		buf.WriteString(")\n")
	}
	return buf.String()
}

func parseTags(raw string) []string {
	splitter := regexp.MustCompile(`[,\n]`)
	var tags []string
	seen := map[string]bool{}
	for _, part := range splitter.Split(raw, -1) {
		tag := slugifyTag(part)
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true
		tags = append(tags, tag)
	}
	return tags
}

func parsePublishMode(value string) publishMode {
	if strings.EqualFold(strings.TrimSpace(value), string(publishModeDirect)) {
		return publishModeDirect
	}
	return publishModePR
}

func generateTitle(req postRequest) string {
	if body := compactWhitespace(req.Body); body != "" {
		words := strings.Fields(body)
		if len(words) > 8 {
			words = words[:8]
		}
		return strings.Join(words, " ")
	}
	return "Photo post " + req.Now.Format("2006-01-02 15:04")
}

func compactWhitespace(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacer := regexp.MustCompile(`[^a-z0-9]+`)
	value = replacer.ReplaceAllString(value, "-")
	return strings.Trim(value, "-")
}

func slugifyTag(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, " ", "-")
	replacer := regexp.MustCompile(`[^a-z0-9-]+`)
	value = replacer.ReplaceAllString(value, "")
	return strings.Trim(value, "-")
}

func yamlQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func newImageObjectName(now time.Time, originalName string, contentType string) (string, error) {
	randomSuffix, err := randomHex(6)
	if err != nil {
		return "", err
	}
	ext := strings.ToLower(filepath.Ext(originalName))
	if ext == "" {
		exts, _ := mime.ExtensionsByType(contentType)
		if len(exts) > 0 {
			ext = exts[0]
		}
	}
	if ext == "" {
		ext = ".jpg"
	}
	base := strings.TrimSuffix(filepath.Base(originalName), filepath.Ext(originalName))
	base = slugify(base)
	if base == "" {
		base = "photo"
	}
	return fmt.Sprintf("%s-%s-%s%s", now.Format("20060102-150405"), base, randomSuffix, ext), nil
}

func randomHex(bytesLen int) (string, error) {
	buf := make([]byte, bytesLen)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
