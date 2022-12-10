package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert/types"
)

var Region = "ap-northeast-1"

func callLambda(ctx context.Context, event events.S3Event) (string, error) {
	//endpoint取得してから、clientを再作成する
	//cfg, _ := config.LoadDefaultConfig(context.TODO())
	//cfg.Region = Region
	//client := mediaconvert.NewFromConfig(cfg)
	//endpoints, _ := client.DescribeEndpoints(context.TODO(), &mediaconvert.DescribeEndpointsInput{})
	//fmt.Println(endpoints)
	var customResolver = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == mediaconvert.ServiceID && region == Region {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               os.Getenv("AWS_MEDIACONVERT_ENDPOINT"),
				SigningRegion:     Region,
				HostnameImmutable: true,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	client := mediaconvert.NewFromConfig(cfg) //本命のclient
	//====

	jobParams := createParams(ctx, event)
	output, err := client.CreateJob(ctx, jobParams)
	if err != nil {
		return "", err
	}
	//err := req.Send()
	res, _ := json.Marshal(output.Job)
	return string(res), nil
}

func handleRequest(ctx context.Context, event events.S3Event) (string, error) {
	// event
	//eventJson, _ := json.MarshalIndent(event, "", "  ")
	//log.Printf("EVENT: %s", eventJson)
	// AWS SDK call

	usage, err := callLambda(ctx, event)
	if err != nil {
		return "ERROR", err
	}
	return usage, nil
}
func main() {
	runtime.Start(handleRequest)
	//cfg, _ := config.LoadDefaultConfig(context.TODO())
	//client := mediaconvert.NewFromConfig(cfg)
	//endpoints, _ := client.DescribeEndpoints(context.TODO(), &mediaconvert.DescribeEndpointsInput{})
	//json, _ := json.Marshal(endpoints)
	//fmt.Println(string(json))
}

func createParams(ctx context.Context, event events.S3Event) *mediaconvert.CreateJobInput {
	//ジョブテンプレートに加えてIAMRole、inputのオーディオのarn、output先を指定
	s3 := event.Records[0].S3
	tmplName := "audio-to-hls"
	role := "arn:aws:iam::138767642386:role/service-role/MediaConvert_Default_Role"
	bucket := s3.Bucket.Name
	key := s3.Object.URLDecodedKey
	folder := filepath.Dir(key)
	input := "s3://" + bucket + "/" + key
	output := "s3://" + bucket + "/" + folder + "/"
	fmt.Println(input)
	fmt.Println(output)
	//==============

	jobSettings := types.JobSettings{}

	jobSettings.Inputs = []types.Input{{FileInput: &input}}
	jobSettings.OutputGroups = []types.OutputGroup{{OutputGroupSettings: &types.OutputGroupSettings{HlsGroupSettings: &types.HlsGroupSettings{Destination: &output}}}}
	params := &mediaconvert.CreateJobInput{}
	params.Role = &role
	params.JobTemplate = &tmplName
	params.Settings = &jobSettings

	//AWSコンソール上でjson出力できるのでそれを参考にparamsを埋める

	return params
}
