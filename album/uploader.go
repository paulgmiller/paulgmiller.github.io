package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	ENDPOINT_URL = "https://222b2fcd50aae5b52660992fbfd93b11.r2.cloudflarestorage.com"
	BUCKET_NAME  = "blogimages"
)

var (
	accessKeyID = "67d604ab768283b886fa7e1d746a9dc9"
)

type s3uploader struct {
	client *s3.Client
}

func NewS3Uploader(ctx context.Context) *s3uploader {
	secretAccessKey := os.Getenv("SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		log.Fatal("Please set the SECRET_ACCESS_KEY environment variable")
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(ENDPOINT_URL)
	})
	return &s3uploader{client: client}
}

func (s *s3uploader) Put(ctx context.Context, fileName string, body io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(fileName),
		Body:   body,
		ACL:    types.ObjectCannedACLPublicRead,
	})
	return err
}
