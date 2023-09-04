package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

var Key = "AKIASATZSF4JJV5MKS4K"
var Secret = "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"

func Test_Handler(t *testing.T) {
	ctx := context.TODO()
	opt := config.WithRegion("ap-northeast-1")
	opt = config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(Key, Secret, ""))
	cfg, err := config.LoadDefaultConfig(ctx,
		opt,
	)
	if err != nil {
		t.Fatalf("failed to load configuration, %v", err)
	}
	app := App{
		githubEndpoint:  ENDPOINT_GITHUB,
		githubEventType: EVENT_TYPE_GITHUB,
		backendEndpoint: ENDPOINT_BACKEND,
		githubToken:     "",
		dbName:          "music-tools-prod-db",
		logGroupName:    "/music_tools/prod/backend",
		forceExec:       false,
		sendRequest:     true,
		awsConfig:       &cfg,
	}
	status, err := app.getDBStatus()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(status)

	found, err := app.findLogWithinThreshold(time.Duration(20))
	if err != nil {
		t.Error(err)
	}
	fmt.Println("log: ", found)
}
