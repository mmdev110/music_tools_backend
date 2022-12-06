package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/aws/aws-sdk-go-v2"
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

func GenerateSignedUrl(userId uint, fileName string, method string, seconds int64) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	Client = s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(Client)
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	key := getKey(userId, fileName)
	fmt.Println(key)
	if method == http.MethodGet {
		presignedGetRequest, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(15 * time.Minute)
		})
		if err != nil {
			log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n", bucketName, key, err)

		}
		return presignedGetRequest.URL, err
	} else if method == http.MethodPut {
		presignedPutRequest, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(15 * time.Minute)
		})
		if err != nil {
			log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n", bucketName, key, err)

		}
		return presignedPutRequest.URL, err
	}
	return "", fmt.Errorf("method %s not allowed", method)
}

func getKey(userId uint, fileName string) string {
	return strconv.Itoa(int(userId)) + "/" + fileName
}
