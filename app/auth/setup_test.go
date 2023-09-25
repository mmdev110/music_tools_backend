package auth

import (
	"fmt"
	"os"
	"testing"
)

var auth = Auth{}

func TestMain(m *testing.M) {
	fmt.Println("@@@@@@@@@")
	region := os.Getenv("AWS_REGION")
	userPoolID := os.Getenv("AWS_COGNITO_USER_POOL_ID")
	appClientID := os.Getenv("AWS_COGNITO_APP_CLIENT_ID")
	fmt.Println(userPoolID)
	if err := auth.Configure(region, userPoolID, appClientID); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}
