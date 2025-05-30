package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"

	"github.com/samber/lo"
)

const (
	REGEX = `(https:\/\/lh3\.googleusercontent\.com\/\w{2}\/[a-zA-Z0-9\-_]{64,})`
)

func getPhotoURLs(albumURL string) ([]string, error) {
	resp, err := http.Get(albumURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(REGEX)
	matches := re.FindAllString(string(body), -1)

	if len(matches) <= 1 {
		return nil, fmt.Errorf("no images found")
	}

	matches = lo.Uniq(matches[1 : len(matches)-1])

	return matches, nil
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

var header = `<div class="fotorama" data-allowfullscreen="true">
<!--%s-->
`

func output(urls []string, albumURL string, w io.Writer) error {
	if _, err := fmt.Fprintf(w, header, albumURL); err != nil {
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
	photos, err := getPhotoURLs(albumURL)
	if err != nil {
		http.Error(w, "failed to scrape: "+err.Error(), http.StatusInternalServerError)
		return
	}
	mirroredURLs, err := mirror(r.Context(), photos, u)
	if err != nil {
		http.Error(w, "failed to mirror: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	if err := output(mirroredURLs, albumURL, w); err != nil {
		http.Error(w, "failed to write: "+err.Error(), http.StatusInternalServerError)
		return
	}

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

	photoURLs, err := getPhotoURLs(albumURL)
	if err != nil {
		log.Fatalf("Failed to retrieve photo URLs: %v", err)
	}

	mirroredURLs, err := mirror(ctx, photoURLs, uploader)
	if err != nil {
		log.Fatalf("Failed to mirror photos: %v", err)
	}

	if err := output(mirroredURLs, albumURL, os.Stdout); err != nil {
		log.Fatalf("Failed to mirror photos: %v", err)
	}

}
