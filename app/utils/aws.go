package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/app/conf"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sestypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

var Client *s3.Client

func GenerateSignedUrl(path string, method string, seconds time.Duration) (string, error) {
	cfg, loadErr := config.LoadDefaultConfig(context.TODO())
	if loadErr != nil {
		log.Fatalf("failed to load configuration, %v", loadErr)
	}
	//fmt.Printf("@@GenerateSignedUrl: %s %s\n", path, method)
	Client = s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(Client)
	bucketName := conf.AWS_BUCKET_NAME
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

func SendEmail(to, title, body string) {
	fmt.Println("email")
	cfg, loadErr := config.LoadDefaultConfig(context.TODO())
	if loadErr != nil {
		log.Fatalf("failed to load configuration, %v", loadErr)
	}
	client := sesv2.NewFromConfig(cfg)

	from := conf.SUPPORT_EMAIL
	input := &sesv2.SendEmailInput{
		FromEmailAddress: &from,
		Destination: &sestypes.Destination{
			ToAddresses: []string{to},
		},
		Content: &sestypes.EmailContent{
			Simple: &sestypes.Message{
				Body: &sestypes.Body{
					Text: &sestypes.Content{
						Data: &body,
					},
				},
				Subject: &sestypes.Content{
					Data: &title,
				},
			},
		},
	}
	result, err := client.SendEmail(context.TODO(), input)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.MessageId)
}
