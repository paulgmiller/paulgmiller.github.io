package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
	_ "time/tzdata"
	"unicode"

	"github.com/samber/lo"
)

const (
	REGEX             = `(https:\/\/lh3\.googleusercontent\.com\/\w{2}\/[a-zA-Z0-9\-_]{64,})`
	defaultAlbumTitle = "Photo Album"
	githubNewPostURL  = "https://github.com/paulgmiller/paulgmiller.github.io/new/master/content/posts"
)

var titleRegex = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)

type album struct {
	title     string
	photoURLs []string
}

func getAlbum(albumURL string) (album, error) {
	resp, err := http.Get(albumURL)
	if err != nil {
		return album{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return album{}, fmt.Errorf("album request returned %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return album{}, err
	}
	return parseAlbum(body)
}

func parseAlbum(body []byte) (album, error) {
	title := defaultAlbumTitle
	if titleMatch := titleRegex.FindSubmatch(body); len(titleMatch) == 2 {
		title = strings.TrimSpace(html.UnescapeString(string(titleMatch[1])))
		title = strings.TrimSpace(strings.TrimSuffix(title, " - Google Photos"))
		if title == "" || title == "Google Photos" {
			title = defaultAlbumTitle
		}
	}

	re := regexp.MustCompile(REGEX)
	matches := re.FindAllString(string(body), -1)

	if len(matches) <= 1 {
		return album{}, fmt.Errorf("no images found")
	}

	matches = lo.Uniq(matches[1 : len(matches)-1])

	return album{title: title, photoURLs: matches}, nil
}

func getPhotoURLs(albumURL string) ([]string, error) {
	album, err := getAlbum(albumURL)
	if err != nil {
		return nil, err
	}
	return album.photoURLs, nil
}

type uploader interface {
	Put(ctx context.Context, fileName string, body io.Reader) error
}

func mirror(ctx context.Context, photoURLs []string, client uploader) ([]string, error) {
	errors := make(chan error)
	mirroredURLs := make(chan string)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, urlStr := range photoURLs {
		go func(urlStr string) {
			parsedURL, _ := url.Parse(urlStr)
			fileName := path.Base(parsedURL.Path)
			mirrorURL := fmt.Sprintf("https://images.northbriton.net/%s", fileName)
			// Check if the file already exists by making a HEAD request
			headResp, err := http.Head(mirrorURL)
			if err == nil && headResp.StatusCode == http.StatusOK {
				log.Printf("Already exists: %s", mirrorURL)
				mirroredURLs <- mirrorURL
				return
			}
			resp, err := http.Get(urlStr + "=s0")
			if err != nil {
				log.Printf("Failed to download %s: %v", urlStr, err)
				errors <- err
				return
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				errors <- err
				return
			}
			err = client.Put(ctx, fileName, bytes.NewReader(data))
			if err != nil {
				log.Printf("Failed to upload %s: %v", urlStr, err)
				errors <- err
				return
			}

			mirroredURLs <- mirrorURL
			log.Printf("Uploaded %s → %s", urlStr, mirrorURL)
		}(urlStr)
	}

	var result []string
	for range len(photoURLs) {
		select {
		case url := <-mirroredURLs:
			result = append(result, url)
		case err := <-errors:
			log.Printf("Error occurred: %v", err)
			return nil, err
		}
	}

	return result, nil
}

var galleryHeader = `<div class="fotorama" data-allowfullscreen="true">
<!--%s-->
`

func output(urls []string, albumURL, title, date string, w io.Writer) error {
	if _, err := fmt.Fprintf(w, "---\nlayout: post\ntitle: %s\ndate: %s\n---\n\n", yamlString(title), date); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, galleryHeader, albumURL); err != nil {
		return err
	}
	for _, url := range urls {
		if _, err := fmt.Fprintf(w, "    <img src=\"https://images.northbriton.net/cdn-cgi/image/width=800/%s\" data-full=\"%s\">\n", url, url); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(w, "</div>"); err != nil {
		return err
	}
	return nil
}

func yamlString(value string) string {
	value = strings.Join(strings.Fields(value), " ")
	if value != "" && !strings.ContainsAny(value, ":#{}[],&*!|>'\"%@`") && !strings.Contains("-?:", value[:1]) {
		return value
	}
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

func slugify(value string) string {
	var slug strings.Builder
	separator := false
	for _, r := range strings.ToLower(value) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if separator && slug.Len() > 0 {
				slug.WriteByte('-')
			}
			slug.WriteRune(r)
			separator = false
		} else {
			separator = true
		}
	}
	if slug.Len() == 0 {
		return "album"
	}
	return slug.String()
}

func newPostURL(title, date, content string) string {
	query := url.Values{}
	query.Set("filename", fmt.Sprintf("%s-%s.md", date, slugify(title)))
	query.Set("value", content)
	return githubNewPostURL + "?" + query.Encode()
}

func today() string {
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		location = time.Local
	}
	return time.Now().In(location).Format(time.DateOnly)
}

func serve(w http.ResponseWriter, r *http.Request, u uploader) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}
	albumURL := r.FormValue("album")
	if albumURL == "" {
		http.Error(w, "missign album", http.StatusBadRequest)
		return
	}
	album, err := getAlbum(albumURL)
	if err != nil {
		http.Error(w, "failed to scrape: "+err.Error(), http.StatusInternalServerError)
		return
	}
	mirroredURLs, err := mirror(r.Context(), album.photoURLs, u)
	if err != nil {
		http.Error(w, "failed to mirror: "+err.Error(), http.StatusInternalServerError)
		return
	}
	date := today()
	var content bytes.Buffer
	if err := output(mirroredURLs, albumURL, album.title, date, &content); err != nil {
		http.Error(w, "failed to write: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, newPostURL(album.title, date, content.String()), http.StatusSeeOther)
}

func main() {

	ctx := context.Background()
	uploader := NewS3Uploader(ctx)

	if len(os.Args) < 2 {
		log.Println(" listening for form encoded album url on port 8080")
		http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serve(w, r, uploader)
		}))
	}
	albumURL := os.Args[1]

	album, err := getAlbum(albumURL)
	if err != nil {
		log.Fatalf("Failed to retrieve photo URLs: %v", err)
	}

	mirroredURLs, err := mirror(ctx, album.photoURLs, uploader)
	if err != nil {
		log.Fatalf("Failed to mirror photos: %v", err)
	}

	if err := output(mirroredURLs, albumURL, album.title, today(), os.Stdout); err != nil {
		log.Fatalf("Failed to mirror photos: %v", err)
	}

}
