package utils

import (
	"testing"
	"time"

	"example.com/app/testutil"
)

func TestJwt(t *testing.T) {
	userId := uint(100)
	jwt, err := GenerateJwt(userId, "access", 24*time.Hour)
	if err != nil {
		t.Fatalf("error found at GenerateJwt: %v", err)
	}

	parsedClaims, err2 := ParseJwt(jwt)
	if err != nil {
		t.Fatalf("error found at ParseJwt: %v", err2)
	}
	testutil.Checker(t, "user_id", parsedClaims.UserId, userId)
}
