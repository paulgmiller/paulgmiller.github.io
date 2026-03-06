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
	googlePhotosRegex = `(https:\/\/lh3\.googleusercontent\.com\/\w{2}\/[a-zA-Z0-9\-_]{64,})`
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

	re := regexp.MustCompile(googlePhotosRegex)
	matches := re.FindAllString(string(body), -1)
	if len(matches) <= 1 {
		return nil, fmt.Errorf("no images found")
	}

	matches = lo.Uniq(matches[1 : len(matches)-1])
	return matches, nil
}

type uploader interface {
	Put(ctx context.Context, fileName string, contentType string, body io.Reader) error
}

func mirror(ctx context.Context, photoURLs []string, client uploader, publicBaseURL string) ([]string, error) {
	errors := make(chan error)
	mirroredURLs := make(chan string)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, urlStr := range photoURLs {
		go func(urlStr string) {
			parsedURL, _ := url.Parse(urlStr)
			fileName := path.Base(parsedURL.Path)
			mirrorURL := fmt.Sprintf("%s/%s", publicBaseURL, fileName)
			headResp, err := http.Head(mirrorURL)
			if err == nil && headResp.StatusCode == http.StatusOK {
				log.Printf("already exists: %s", mirrorURL)
				mirroredURLs <- mirrorURL
				return
			}

			resp, err := http.Get(urlStr + "=s0")
			if err != nil {
				log.Printf("failed to download %s: %v", urlStr, err)
				errors <- err
				return
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				errors <- err
				return
			}

			contentType := resp.Header.Get("Content-Type")
			if contentType == "" {
				contentType = http.DetectContentType(data)
			}
			if err := client.Put(ctx, fileName, contentType, bytes.NewReader(data)); err != nil {
				log.Printf("failed to upload %s: %v", urlStr, err)
				errors <- err
				return
			}

			mirroredURLs <- mirrorURL
			log.Printf("uploaded %s -> %s", urlStr, mirrorURL)
		}(urlStr)
	}

	var result []string
	for range len(photoURLs) {
		select {
		case url := <-mirroredURLs:
			result = append(result, url)
		case err := <-errors:
			log.Printf("error occurred: %v", err)
			return nil, err
		}
	}

	return result, nil
}

var legacyHeader = `<div class="fotorama" data-allowfullscreen="true">
<!--%s-->
`

func output(urls []string, albumURL string, imageTransformURL string, w io.Writer) error {
	if _, err := fmt.Fprintf(w, legacyHeader, albumURL); err != nil {
		return err
	}
	for _, imageURL := range urls {
		if _, err := fmt.Fprintf(w, "    <img src=\"%s/%s\" data-full=\"%s\">\n", imageTransformURL, imageURL, imageURL); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(w, "</div>")
	return err
}

func serveMirror(w http.ResponseWriter, r *http.Request, cfg Config, u uploader) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}
	albumURL := r.FormValue("album")
	if albumURL == "" {
		http.Error(w, "missing album", http.StatusBadRequest)
		return
	}

	photos, err := getPhotoURLs(albumURL)
	if err != nil {
		http.Error(w, "failed to scrape: "+err.Error(), http.StatusInternalServerError)
		return
	}
	mirroredURLs, err := mirror(r.Context(), photos, u, cfg.PublicImageBaseURL)
	if err != nil {
		http.Error(w, "failed to mirror: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if err := output(mirroredURLs, albumURL, cfg.ImageTransformBaseURL, w); err != nil {
		http.Error(w, "failed to write: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	cfg := loadConfig()
	u := NewS3Uploader(context.Background())

	if len(os.Args) >= 2 {
		albumURL := os.Args[1]
		photoURLs, err := getPhotoURLs(albumURL)
		if err != nil {
			log.Fatalf("failed to retrieve photo URLs: %v", err)
		}

		mirroredURLs, err := mirror(context.Background(), photoURLs, u, cfg.PublicImageBaseURL)
		if err != nil {
			log.Fatalf("failed to mirror photos: %v", err)
		}

		if err := output(mirroredURLs, albumURL, cfg.ImageTransformBaseURL, os.Stdout); err != nil {
			log.Fatalf("failed to write output: %v", err)
		}
		return
	}

	log.Printf("listening on :%s with base path %q", cfg.Port, cfg.BasePath)
	if err := http.ListenAndServe(":"+cfg.Port, newServer(cfg, u)); err != nil {
		log.Fatal(err)
	}
}
