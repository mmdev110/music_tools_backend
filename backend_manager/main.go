package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/app/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// github actionsに送るrequest(のbody)
type WorkflowBody struct {
	EventType     string  `json:"event_type"`
	ClientPayload Payload `json:"client_payload"`
}
type Payload struct {
	Action string `json:"action"` //"Apply" or "Destroy"
}

// lambdaが受け取るrequest
type Event struct {
	Action string `json:"action"` //"apply" or "destroy" or "status"
	Force  bool   `json:"force"`
}

// lambdaが返すresponse
type Response struct {
	BackendStatus string `json:"backend_status"`
	DBStatus      string `json:"db_status"`
}

const (
	ENDPOINT_GITHUB   = "https://api.github.com/repos/mmdev110/music_tools_infra/dispatches"
	EVENT_TYPE_GITHUB = "backend_manager"
	ENDPOINT_BACKEND  = "http://backend.music-tools.ys-dev.net/_chk"
	DB_NAME           = "music-tools-prod-db"
	LOG_GROUP_NAME    = "/music_tools/prod/backend"
	REGION            = "ap-northeast-1"
)

type App struct {
	githubEndpoint     string
	githubEventType    string
	backendEndpoint    string
	githubToken        string
	dbName             string
	logGroupName       string
	forceExec          bool
	sendRequest        bool
	awsConfig          *aws.Config
	accessThresholdMin int // n分アクセスなければdestroyさせる
}

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	token := os.Getenv("github_token")
	if token == "" {
		return sendResponse(nil, errors.New("github_token not set"))
	}
	event := Event{}
	if err := json.Unmarshal([]byte(req.Body), &event); err != nil {
		return sendResponse(nil, err)
	}
	fmt.Printf("%+v\n", event)

	action := event.Action

	cfg := configureAWS(ctx)
	app := App{
		githubEndpoint:     ENDPOINT_GITHUB,
		githubEventType:    EVENT_TYPE_GITHUB,
		backendEndpoint:    ENDPOINT_BACKEND,
		githubToken:        token,
		dbName:             DB_NAME,
		logGroupName:       LOG_GROUP_NAME,
		forceExec:          event.Force,
		sendRequest:        true,
		awsConfig:          cfg,
		accessThresholdMin: 60,
	}

	if action == "status" {
		result, err := app.Status()
		return sendResponse(result, err)

	} else if action == "apply" {
		result, err := app.Apply()
		return sendResponse(result, err)

	} else if action == "destroy" {
		result, err := app.Destroy()
		return sendResponse(result, err)
	}
	return sendResponse(nil, fmt.Errorf("invalid action: %s", action))
}

func configureAWS(ctx context.Context) *aws.Config {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(REGION))
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	return &cfg

}
func sendResponse(body interface{}, err error) (events.APIGatewayProxyResponse, error) {
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	} else {
		bd, err := utils.ToJSON(body)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
		}
		return events.APIGatewayProxyResponse{Body: bd, StatusCode: 200}, nil
	}
}

func main() {
	//fmt.Println("start")
	lambda.Start(handleRequest)
	//bo, err := ping(http.MethodGet, "http://backend.music-tools.ys-dev.net/_chk")
	//handleRequest(context.Background())
}
