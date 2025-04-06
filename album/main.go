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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	REGEX        = `(https:\/\/lh3\.googleusercontent\.com\/\w{2}\/[a-zA-Z0-9\-_]{64,})`
	BUCKET_NAME  = "blogimages"
	ENDPOINT_URL = "https://222b2fcd50aae5b52660992fbfd93b11.r2.cloudflarestorage.com"
)

var (
	accessKeyID     = "67d604ab768283b886fa7e1d746a9dc9"
	secretAccessKey = os.Getenv("SECRET_ACCESS_KEY")
)

func NewS3Uploader() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(ENDPOINT_URL)
	})
}

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

func mirror(photoURLs []string, client *s3.Client) []string {
	var mirroredURLs []string

	for _, urlStr := range photoURLs {
		parsedURL, _ := url.Parse(urlStr)
		fileName := path.Base(parsedURL.Path)
		mirrorURL := fmt.Sprintf("https://images.northbriton.net/%s", fileName)
		// Check if the file already exists by making a HEAD request
		headResp, err := http.Head(mirrorURL)
		if err == nil && headResp.StatusCode == http.StatusOK {
			log.Printf("Already exists: %s", mirrorURL)
			mirroredURLs = append(mirroredURLs, mirrorURL)
			continue
		}
		resp, err := http.Get(urlStr + "=s0")
		if err != nil {
			log.Printf("Failed to download %s: %v", urlStr, err)
			continue
		}
		defer resp.Body.Close()

		buffer := new(bytes.Buffer)
		_, err = io.Copy(buffer, resp.Body)
		if err != nil {
			log.Printf("Failed to read response body for %s: %v", urlStr, err)
			continue
		}

		_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(BUCKET_NAME),
			Key:    aws.String(fileName),
			Body:   buffer,
			ACL:    types.ObjectCannedACLPublicRead,
		})
		if err != nil {
			log.Printf("Failed to upload %s: %v", urlStr, err)
			continue
		}

		mirroredURLs = append(mirroredURLs, mirrorURL)
		log.Printf("Uploaded %s → %s", urlStr, mirrorURL)
	}

	return mirroredURLs
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a Google Photos album URL")
	}

	albumURL := os.Args[1]
	if secretAccessKey == "" {
		log.Fatal("Please set the SECRET_ACCESS_KEY environment variable")
	}

	uploader := NewS3Uploader()
	photoURLs, err := getPhotoURLs(albumURL)
	if err != nil {
		log.Fatalf("Failed to retrieve photo URLs: %v", err)
	}

	mirroredURLs := mirror(photoURLs, uploader)

	fmt.Println("<div class=\"fotorama\" data-allowfullscreen=\"true\">")
	fmt.Printf("    <!--%s-->", albumURL)

	for _, url := range mirroredURLs {
		fmt.Printf("    <img src=\"https://images.northbriton.net/cdn-cgi/image/width=800/%s\" data-full=\"%s\">\n", url, url)
	}

	fmt.Println("</div>")
}
