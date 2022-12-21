package utils

import (
	"fmt"
	"testing"
	"time"
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
	fmt.Println(jwt)
	PrintStruct(parsedClaims)
	got := parsedClaims.UserId
	want := userId
	if got != want {
		t.Errorf("got %d, want %d .", got, want)
	}
}
