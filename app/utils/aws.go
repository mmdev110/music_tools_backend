package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var Client *s3.Client

func ConfigureAWS() {
	fmt.Println(os.Getenv("AWS_ACCESS_KEY_ID"))
	fmt.Println(os.Getenv("AWS_SECRET_ACCESS_KEY"))
	fmt.Println(os.Getenv("AWS_BUCKET_NAME"))

}

func GenerateSignedUrl(path string, method string, seconds time.Duration) (string, error) {
	cfg, loadErr := config.LoadDefaultConfig(context.TODO())
	if loadErr != nil {
		log.Fatalf("failed to load configuration, %v", loadErr)
	}
	fmt.Printf("@@GenerateSignedUrl: %s %s\n", path, method)
	Client = s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(Client)
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	var url string
	var err error
	if method == http.MethodGet {
		presignedGetRequest, err1 := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(path),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = seconds * time.Second
		})
		url = presignedGetRequest.URL
		err = err1
	} else if method == http.MethodPut {
		presignedPutRequest, err2 := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(path),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = seconds * time.Second
		})
		url = presignedPutRequest.URL
		err = err2
	}
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n", bucketName, path, err)
		return "", err
	}
	if url == "" {
		return "", fmt.Errorf("method %s not allowed", method)
	}
	return url, nil
}
