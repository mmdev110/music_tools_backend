package main

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	_ "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func (app *App) Status() (*Response, error) {
	backendStatus := "stopped"
	pingSuccess, err := ping(app.backendEndpoint)
	if err != nil {
		return nil, err
	}
	if pingSuccess {
		backendStatus = "running"
	}

	dbStatus, err := app.getDBStatus()
	if err != nil {
		return nil, err
	}
	res := Response{
		BackendStatus: backendStatus,
		DBStatus:      dbStatus,
	}
	return &res, nil
}

func (app *App) getDBStatus() (string, error) {
	rdsClient := rds.NewFromConfig(*app.awsConfig)
	output, err := rdsClient.DescribeDBInstances(context.Background(), &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(app.dbName),
	})
	if err != nil {
		return "", err
	}
	if len(output.DBInstances) == 0 {
		return "", errors.New("DB instance not found")
	}
	instance := output.DBInstances[0]
	//fmt.Printf("DB instance %v has database %v.\n", *instance.DBInstanceIdentifier,*instance.DBInstanceStatus)
	return *instance.DBInstanceStatus, nil
}

func (app *App) findLogWithinThreshold(min time.Duration) (bool, error) {
	//aws logs tail --since={threshold}m log_group_nameしてレコード無ければfalse

	client := cloudwatchlogs.NewFromConfig(*app.awsConfig)
	start := time.Now().Add(-min * time.Minute)
	output, err := client.FilterLogEvents(context.Background(), &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName:  aws.String(app.logGroupName),
		StartTime:     aws.Int64(start.Unix() * 1000),
		FilterPattern: aws.String("\"/refresh\""),
	})
	if err != nil {
		return false, err
	}
	if len(output.Events) == 0 {
		return false, nil
	}
	//found := output.Events[0]
	//fmt.Printf("message: %v", *found.Message)
	//fmt.Printf("timestamp: %v", *found.Timestamp)
	return true, nil
}
