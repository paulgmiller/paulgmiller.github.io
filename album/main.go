package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

	return matches[1 : len(matches)-1], nil
}

func mirror(ctx context.Context, photoURLs []string, client *s3.Client) ([]string, error) {
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
			}
			defer resp.Body.Close()

			_, err = client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(BUCKET_NAME),
				Key:    aws.String(fileName),
				Body:   resp.Body,
				ACL:    types.ObjectCannedACLPublicRead,
			})
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
	for _ = range len(photoURLs) {
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

func output(urls []string, albumURL string, w io.Writer) error {
	if _, err := fmt.Fprintf(w, `<div class="fotorama" data-allowfullscreen="true">\n<!--%s-->\n`, albumURL); err != nil {
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

func serve(w http.ResponseWriter, r *http.Request) {
	album := r.URL.Query().Get("album")
	if album == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("need an album query arg"))
		return
	}
	photos, err := getPhotoURLs(album)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	mirror(r.Context(), photos, NewS3Uploader(r.Context()))

}

func main() {
	if secretAccessKey == "" {
		log.Fatal("Please set the SECRET_ACCESS_KEY environment variable")
	}

	if len(os.Args) < 2 {
		log.Fatal("Please provide a Google Photos album URL")
	}
	ctx := context.Background()
	uploader := NewS3Uploader(ctx)

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
